package textfield

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/zodimo/go-compose/internal/layoutnode"

	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/input"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

const Material3OutlinedTextFieldNodeID = "Material3OutlinedTextField"

// Outlined implements the Outlined Material Design 3 text field.
// It uses a custom widget implementation adapted from gio-x.
func Outlined(
	value string,
	onValueChange func(string),
	label string,
	options ...TextFieldOption,
) Composable {
	return func(c Composer) Composer {
		opts := DefaultTextFieldOptions()
		for _, opt := range options {
			opt(&opts)
		}

		key := c.GenerateID()
		path := c.GetPath()

		// Handler wrapper
		handlerWrapperState := c.State(fmt.Sprintf("%d/%s/handler_wrapper", key, path), func() any {
			return &HandlerWrapper{Func: onValueChange}
		})
		handlerWrapper := handlerWrapperState.Get().(*HandlerWrapper)
		handlerWrapper.Func = onValueChange

		// OnSubmit wrapper
		var onSubmitWrapper *OnSubmitWrapper
		if opts.OnSubmit != nil {
			onSubmitWrapperState := c.State(fmt.Sprintf("%d/%s/onsubmit_wrapper", key, path), func() any {
				return &OnSubmitWrapper{Func: opts.OnSubmit}
			})
			onSubmitWrapper = onSubmitWrapperState.Get().(*OnSubmitWrapper)
			onSubmitWrapper.Func = opts.OnSubmit
		}

		// Custom Outlined Widget State
		widgetStatePath := fmt.Sprintf("%d/%s/outlined_widget/s%v", key, path, opts.SingleLine)
		widgetVal := c.State(widgetStatePath, func() any {
			return &OutlinedTextFieldWidget{
				Editor: widget.Editor{
					SingleLine: opts.SingleLine,
					Submit:     opts.OnSubmit != nil,
				},
			}
		})
		outWidget := widgetVal.Get().(*OutlinedTextFieldWidget)

		// State tracker for synchronization
		trackerState := c.State(fmt.Sprintf("%d/%s/tracker/s%v", key, path, opts.SingleLine), func() any {
			return &TextFieldStateTracker{LastValue: ""}
		})
		tracker := trackerState.Get().(*TextFieldStateTracker)

		// Update static properties
		outWidget.Editor.SingleLine = opts.SingleLine
		outWidget.Editor.Submit = opts.OnSubmit != nil

		c.StartBlock(Material3OutlinedTextFieldNodeID)
		c.Modifier(func(m Modifier) Modifier {
			return m.Then(opts.Modifier)
		})

		// Constructor
		c.SetWidgetConstructor(outlinedTextFieldWidgetConstructor(outWidget, value, label, opts, handlerWrapper, onSubmitWrapper, tracker))

		return c.EndBlock()
	}
}

func outlinedTextFieldWidgetConstructor(
	w *OutlinedTextFieldWidget,
	value string,
	label string,
	opts TextFieldOptions,
	handler *HandlerWrapper,
	onSubmitHandler *OnSubmitWrapper,
	tracker *TextFieldStateTracker,
) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			// 1. Sync External Change
			if value != tracker.LastValue {
				if w.Editor.Text() != value {
					w.Editor.SetText(value)
				}
				tracker.LastValue = value
			}

			// 2. Events & Layout
			th := material.NewTheme() // Fallback TODO: real theme

			// Check for submit events
			for {
				ev, ok := w.Editor.Update(gtx)
				if !ok {
					break
				}
				if _, ok := ev.(widget.SubmitEvent); ok {
					if onSubmitHandler != nil && onSubmitHandler.Func != nil {
						onSubmitHandler.Func()
					}
				}
			}

			// Check for text changes
			currentText := w.Editor.Text()
			if currentText != value {
				if handler.Func != nil {
					handler.Func(currentText)
				}
			}

			return w.Layout(gtx, th, label)
		}
	})
}

// --- Adapted from gio-x/component/text_field.go ---

type OutlinedTextFieldWidget struct {
	widget.Editor
	click gesture.Click

	// Animation state
	state
	label  label
	border border
	anim   *Progress
}

type label struct {
	TextSize unit.Sp
	Inset    layout.Inset
	Smallest layout.Dimensions
}

type border struct {
	Thickness unit.Dp
	Color     color.NRGBA
}

type state int

const (
	inactive state = iota
	hovered
	activated
	focused
)

func (in *OutlinedTextFieldWidget) Layout(gtx layout.Context, th *material.Theme, hint string) layout.Dimensions {
	// Logic from gio-x Update + Layout
	in.update(gtx, th, hint)

	// Offset accounts for label height, which sticks above the border dimensions.
	defer op.Offset(image.Pt(0, in.label.Smallest.Size.Y/2)).Push(gtx.Ops).Pop()

	// Draw Label
	in.label.Inset.Layout(
		gtx,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Left:  unit.Dp(4),
				Right: unit.Dp(4),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, in.label.TextSize, hint)
				label.Color = in.border.Color
				return label.Layout(gtx)
			})
		})

	dims := layout.Flex{
		Axis: layout.Vertical,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(
				gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					cornerRadius := unit.Dp(4)
					dimsFunc := func(gtx layout.Context) layout.Dimensions {
						return layout.Dimensions{Size: image.Point{
							X: gtx.Constraints.Max.X,
							Y: gtx.Constraints.Min.Y,
						}}
					}
					border := widget.Border{
						Color:        in.border.Color,
						Width:        in.border.Thickness,
						CornerRadius: cornerRadius,
					}
					// Cutout logic
					if gtx.Source.Focused(&in.Editor) || in.Editor.Len() > 0 {
						visibleBorder := clip.Path{}
						visibleBorder.Begin(gtx.Ops)
						// Helper to make points clearer
						pt := func(x, y float32) f32.Point { return f32.Point{X: x, Y: y} }

						const buffer = 1000.0 // Draw way outside to avoid clipping corners

						// Start top-leftish (at label start)
						labelStartX := float32(gtx.Dp(in.label.Inset.Left))
						labelEndX := labelStartX + float32(in.label.Smallest.Size.X)
						labelEndY := float32(in.label.Smallest.Size.Y)

						// Trace the visible area (everything EXCEPT the label cutout)
						// We use a large bounding box method or exact path.
						// Current path: Start (0,0) -> Down -> Right -> Up -> Left(to LabelEnd) -> Down(cutout) -> Left -> Up -> Close.

						minY := float32(gtx.Constraints.Min.Y)
						maxX := float32(gtx.Constraints.Max.X)

						visibleBorder.MoveTo(pt(0, 0))
						visibleBorder.LineTo(pt(0, minY))    // Down to bottom-left
						visibleBorder.LineTo(pt(maxX, minY)) // Right to bottom-right
						visibleBorder.LineTo(pt(maxX, 0))    // Up to top-right
						visibleBorder.LineTo(pt(labelEndX, 0))
						visibleBorder.LineTo(pt(labelEndX, labelEndY))   // Dip down
						visibleBorder.LineTo(pt(labelStartX, labelEndY)) // Left across dip
						visibleBorder.LineTo(pt(labelStartX, 0))         // Up from dip
						visibleBorder.LineTo(pt(0, 0))                   // Back to start

						visibleBorder.Close()
						defer clip.Outline{
							Path: visibleBorder.End(),
						}.Op().Push(gtx.Ops).Pop()
					}
					return border.Layout(gtx, dimsFunc)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(12)).Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return layout.Flex{
								Axis:      layout.Horizontal,
								Alignment: layout.Middle,
							}.Layout(
								gtx,
								// Prefix would go here
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return material.Editor(th, &in.Editor, "").Layout(gtx)
								}),
								// Suffix would go here
							)
						},
					)
				}),
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					defer pointer.PassOp{}.Push(gtx.Ops).Pop()
					defer clip.Rect(image.Rectangle{
						Max: gtx.Constraints.Min,
					}).Push(gtx.Ops).Pop()
					in.click.Add(gtx.Ops)
					return layout.Dimensions{}
				}),
			)
		}),
		// Helper text would go here
	)
	return layout.Dimensions{
		Size: image.Point{
			X: dims.Size.X,
			Y: dims.Size.Y + in.label.Smallest.Size.Y/2,
		},
		Baseline: dims.Baseline,
	}
}

func (in *OutlinedTextFieldWidget) update(gtx layout.Context, th *material.Theme, hint string) {
	disabled := gtx.Source == (input.Source{})
	for {
		ev, ok := in.click.Update(gtx.Source)
		if !ok {
			break
		}
		switch ev.Kind {
		case gesture.KindPress:
			gtx.Execute(key.FocusCmd{Tag: &in.Editor})
		}
	}
	in.state = inactive
	if in.click.Hovered() && !disabled {
		in.state = hovered
	}
	hasContents := in.Editor.Len() > 0
	if hasContents {
		in.state = activated
	}
	if gtx.Source.Focused(&in.Editor) && !disabled {
		in.state = focused
	}
	const (
		duration = time.Millisecond * 100
	)
	if in.anim == nil {
		in.anim = &Progress{}
	}
	if in.state == activated || hasContents {
		in.anim.Start(gtx.Now, Forward, 0)
	}
	if in.state == focused && !hasContents && !in.anim.Started() {
		in.anim.Start(gtx.Now, Forward, duration)
	}
	if in.state == inactive && !hasContents && in.anim.Finished() {
		in.anim.Start(gtx.Now, Reverse, duration)
	}
	if in.anim.Started() {
		gtx.Execute(op.InvalidateCmd{})
	}
	in.anim.Update(gtx.Now)

	// Styles - TODO: Use real theme colors
	var (
		textNormal         = th.TextSize
		textSmall          = th.TextSize * 0.8
		borderColor        = color.NRGBA{R: 120, G: 120, B: 120, A: 255} // Placeholder
		borderColorHovered = color.NRGBA{R: 0, G: 0, B: 0, A: 255}       // Placeholder
		borderColorActive  = th.Palette.ContrastBg

		borderThickness       = unit.Dp(1)
		borderThicknessActive = unit.Dp(2)
	)

	in.label.TextSize = unit.Sp(lerp(float32(textSmall), float32(textNormal), 1.0-in.anim.Progress()))
	switch in.state {
	case inactive:
		in.border.Thickness = borderThickness
		in.border.Color = borderColor
	case hovered, activated:
		in.border.Thickness = borderThickness
		in.border.Color = borderColorHovered
	case focused:
		in.border.Thickness = borderThicknessActive
		in.border.Color = borderColorActive
	}

	// Calculate smallest label for cutout
	gtx.Constraints.Min.X = 0
	macro := op.Record(gtx.Ops)
	var spacing unit.Dp
	if len(hint) > 0 {
		spacing = 4
	}
	in.label.Smallest = layout.Inset{
		Left:  spacing,
		Right: spacing,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.Label(th, textSmall, hint).Layout(gtx)
	})
	macro.Stop()

	labelTopInsetNormal := float32(in.label.Smallest.Size.Y) - float32(in.label.Smallest.Size.Y/4)
	topInsetDP := unit.Dp(labelTopInsetNormal / gtx.Metric.PxPerDp)
	topInsetActiveDP := (topInsetDP / 2 * -1) - unit.Dp(in.border.Thickness)
	in.label.Inset = layout.Inset{
		Top:  unit.Dp(lerp(float32(topInsetDP), float32(topInsetActiveDP), in.anim.Progress())),
		Left: unit.Dp(10),
	}
}

// Utils
type Progress struct {
	begin time.Time
	dir   Direction
	dur   time.Duration
	frac  float32
}

type Direction int

const (
	Forward Direction = iota
	Reverse
)

func (p *Progress) Start(now time.Time, dir Direction, dur time.Duration) {
	if p.dir == dir && (p.frac == 1.0 || p.frac == 0.0) && dur == p.dur {
		return
	}
	if p.dir != dir {
		p.frac = 1.0 - p.frac
		p.begin = now.Add(time.Duration(float32(dur)*(1-p.frac)) * -1)
	} else {
		p.begin = now.Add(time.Duration(float32(dur)*p.frac) * -1)
	}
	p.dir = dir
	p.dur = dur
}

func (p *Progress) Started() bool {
	return p.dur > 0 && (p.frac < 1.0 && p.frac > 0.0)
}

func (p *Progress) Finished() bool {
	return p.frac == 1.0
}

func (p *Progress) Update(now time.Time) {
	if p.dur == 0 {
		if p.dir == Forward {
			p.frac = 1.0
		} else {
			p.frac = 0.0
		}
		return
	}
	elapsed := now.Sub(p.begin)
	if elapsed > p.dur {
		p.frac = 1.0
		return
	}
	p.frac = float32(elapsed) / float32(p.dur)
}

func (p *Progress) Progress() float32 {
	if p.dir == Reverse {
		return 1.0 - p.frac
	}
	return p.frac
}

func lerp(start, end, progress float32) float32 {
	return start + (end-start)*progress
}
