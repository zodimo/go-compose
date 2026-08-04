package textfield

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"time"

	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/textinput"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
	"github.com/zodimo/go-compose/pkg/sentinel"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	gioUnit "gioui.org/unit"
	"gioui.org/widget"
	gioMaterial "gioui.org/widget/material"
	"github.com/zodimo/go-compose/compose/material"
)

const Material3OutlinedTextFieldNodeID = "Material3OutlinedTextField"

// Outlined implements the Outlined Material Design 3 text field.
// It uses a custom widget implementation adapted from gio-x.
func Outlined(
	value string,
	onValueChange func(string),
	options ...TextFieldOption,
) Composable {

	opts := DefaultTextFieldOptions()
	for _, opt := range options {
		opt(&opts)
	}

	return func(c Composer) Composer {

		opts.Colors = ResolveTextFieldColors(c, opts.Colors)
		opts.SupportingText = sentinel.TakeOrElseString(opts.SupportingText, "")
		opts.Label = sentinel.TakeOrElseString(opts.Label, "")

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
				Editor: textinput.NewEditor(opts.SingleLine, opts.OnSubmit != nil, 0),
			}
		})
		outWidget := widgetVal.Get().(*OutlinedTextFieldWidget)

		// State tracker for synchronization
		trackerState := c.State(fmt.Sprintf("%d/%s/tracker/s%v", key, path, opts.SingleLine), func() any {
			return &TextFieldStateTracker{LastValue: ""}
		})
		tracker := trackerState.Get().(*TextFieldStateTracker)

		// Update static properties
		outWidget.Editor.SetSingleLine(opts.SingleLine)
		outWidget.Editor.SetSubmit(opts.OnSubmit != nil)
		outWidget.Editor.SetMask(opts.Mask)
		outWidget.Helper = opts.SupportingText
		outWidget.Editor.SetReadOnly(opts.ReadOnly)
		outWidget.SetError(opts.IsError, opts.SupportingText)

		c.StartBlock(Material3OutlinedTextFieldNodeID)
		c.Modifier(func(m ui.Modifier) ui.Modifier {
			return m.Then(opts.Modifier)
		})

		// Compose slots
		if opts.LeadingIcon != nil {
			c.WithComposable(opts.LeadingIcon)
		}
		if opts.TrailingIcon != nil {
			c.WithComposable(opts.TrailingIcon)
		}

		// Constructor
		gioTh := material.GioThemeForEngine(c).(*gioMaterial.Theme)
		th := textinput.NewTheme(gioTh)

		c.SetWidgetConstructor(outlinedTextFieldWidgetConstructor(outWidget, value, opts, handlerWrapper, onSubmitWrapper, tracker, th))

		return c.EndBlock()
	}
}

func outlinedTextFieldWidgetConstructor(
	w *OutlinedTextFieldWidget,
	value string,
	opts TextFieldOptions,
	handler *HandlerWrapper,
	onSubmitHandler *OnSubmitWrapper,
	tracker *TextFieldStateTracker,
	theme *textinput.Theme,
) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {

		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			if !opts.Enabled {
				d := gtx.ToGio().Disabled()
				gtx = layoutnode.NewLayoutContext(&d)
			}
			// Map children to slots (Must be done here, after WrapChildren)
			children := node.Children()
			childIdx := 0

			w.Prefix = nil
			if opts.LeadingIcon != nil && childIdx < len(children) {
				child := children[childIdx]
				if coord, ok := child.(layoutnode.NodeCoordinator); ok {
					w.Prefix = func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
						return coord.Layout(gtx)
					}
				}
				childIdx++
			}

			w.Suffix = nil
			if opts.TrailingIcon != nil && childIdx < len(children) {
				child := children[childIdx]
				if coord, ok := child.(layoutnode.NodeCoordinator); ok {
					w.Suffix = func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
						return coord.Layout(gtx)
					}
				}
				childIdx++
			}

			// 1. Sync external value changes
			if value != tracker.LastValue {
				if w.Editor.Text() != value {
					w.Editor.SetText(value)
					if tracker.HasPendingCaret {
						maxPos := len(value)
						start, end := tracker.CaretStart, tracker.CaretEnd
						if start > maxPos {
							start = maxPos
						}
						if end > maxPos {
							end = maxPos
						}
						w.Editor.SetCaret(start, end)
						tracker.HasPendingCaret = false
					}
				}
				tracker.LastValue = value
			}

			// Check for submit events
			hasSubmit := w.Editor.ProcessEvents(gtx)
			if hasSubmit {
				if onSubmitHandler != nil && onSubmitHandler.Func != nil {
					onSubmitHandler.Func()
				}
			}

			// Check for text changes
			currentText := w.Editor.Text()
			if currentText != value {
				start, end := w.Editor.Selection()
				tracker.CaretStart = start
				tracker.CaretEnd = end
				tracker.HasPendingCaret = true

				if handler.Func != nil {
					handler.Func(currentText)
				}
				w.Editor.SetText(value)
				maxPos := len(value)
				if start > maxPos {
					start = maxPos
				}
				if end > maxPos {
					end = maxPos
				}
				w.Editor.SetCaret(start, end)
			}

			w.Colors = opts.Colors

			return w.Layout(gtx, theme, opts.Label)
		}
	})
}

// --- Adapted from gio-x/component/text_field.go ---

type OutlinedTextFieldWidget struct {
	Editor *textinput.Editor
	click  textinput.Click

	// Config
	Helper    string
	CharLimit uint
	Prefix    layoutnode.GioLayoutWidget
	Suffix    layoutnode.GioLayoutWidget
	Colors    TextFieldColors

	// Animation state
	state
	label  label
	border border
	helper helper
	anim   *Progress

	editorInset struct {
		Top, Right, Bottom, Left unit.Dp
	}

	errored bool
}

type helper struct {
	Color color.NRGBA
	Text  string
}

type label struct {
	TextSize gioUnit.Sp
	Inset    layout.Inset
	Smallest layout.Dimensions
}

type border struct {
	Thickness gioUnit.Dp
	Color     color.NRGBA
}

type state int

const (
	inactive state = iota
	hovered
	activated
	focused
)

// IsActive if input is in an active state (Active, Focused or Errored).
func (in *OutlinedTextFieldWidget) IsActive() bool {
	return in.state >= activated
}

// IsErrored if input is in an errored state.
// Typically this is when the validator has returned an error message.
func (in *OutlinedTextFieldWidget) IsErrored() bool {
	return in.errored
}

// SetError puts the input into an errored state with the specified error text.
func (in *OutlinedTextFieldWidget) SetError(isError bool, err string) {
	in.errored = isError
	in.helper.Text = err
}

// ClearError clears any errored status.
func (in *OutlinedTextFieldWidget) ClearError() {
	in.errored = false
	in.helper.Text = in.Helper
}

// Clear the input text and reset any error status.
func (in *OutlinedTextFieldWidget) Clear() {
	in.Editor.SetText("")
	in.ClearError()
}

// TextTooLong returns whether the current editor text exceeds the set character
// limit.
func (in *OutlinedTextFieldWidget) TextTooLong() bool {
	return !(in.CharLimit == 0 || uint(len(in.Editor.Text())) < in.CharLimit)
}

func (in *OutlinedTextFieldWidget) Layout(gtx layoutnode.LayoutContext, th *textinput.Theme, hint string) layoutnode.LayoutDimensions {
	// Logic from gio-x Update + Layout
	in.update(gtx, th, hint)

	g := *gtx.ToGio()

	// Offset accounts for label height, which sticks above the border dimensions.
	defer op.Offset(image.Pt(0, in.label.Smallest.Size.Y/2)).Push(g.Ops).Pop()

	// Draw Label
	in.label.Inset.Layout(
		g,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Left:  gioUnit.Dp(4),
				Right: gioUnit.Dp(4),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				label := gioMaterial.Label(th.T, in.label.TextSize, hint)
				label.Color = in.border.Color
				return label.Layout(gtx)
			})
		})

	dims := layout.Flex{
		Axis: layout.Vertical,
	}.Layout(
		g,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(
				gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					cornerRadius := gioUnit.Dp(4)
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
					if gtx.Source.Focused(&in.Editor.E) || in.Editor.Len() > 0 {
						visibleBorder := clip.Path{}
						visibleBorder.Begin(gtx.Ops)
						pt := func(x, y float32) f32.Point { return f32.Point{X: x, Y: y} }

						labelStartX := float32(gtx.Dp(in.label.Inset.Left))
						labelEndX := labelStartX + float32(in.label.Smallest.Size.X)
						labelEndY := float32(in.label.Smallest.Size.Y)

						minY := float32(gtx.Constraints.Min.Y)
						maxX := float32(gtx.Constraints.Max.X)

						visibleBorder.MoveTo(pt(0, 0))
						visibleBorder.LineTo(pt(0, minY))
						visibleBorder.LineTo(pt(maxX, minY))
						visibleBorder.LineTo(pt(maxX, 0))
						visibleBorder.LineTo(pt(labelEndX, 0))
						visibleBorder.LineTo(pt(labelEndX, labelEndY))
						visibleBorder.LineTo(pt(labelStartX, labelEndY))
						visibleBorder.LineTo(pt(labelStartX, 0))
						visibleBorder.LineTo(pt(0, 0))

						visibleBorder.Close()
						defer clip.Outline{
							Path: visibleBorder.End(),
						}.Op().Push(gtx.Ops).Pop()
					}
					return border.Layout(gtx, dimsFunc)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Left:  gioUnit.Dp(12),
						Right: gioUnit.Dp(12),
					}.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return layout.Flex{
								Axis:      layout.Horizontal,
								Alignment: layout.Middle,
							}.Layout(
								gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if in.Prefix != nil {
										return layoutnode.ToGioDimensions(in.Prefix(layoutnode.NewLayoutContext(&gtx)))
									}
									return layout.Dimensions{}
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									inset := layout.Inset{
										Top:    gioUnit.Dp(in.editorInset.Top),
										Right:  gioUnit.Dp(in.editorInset.Right),
										Bottom: gioUnit.Dp(in.editorInset.Bottom),
										Left:   gioUnit.Dp(in.editorInset.Left),
									}
									return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										textColor := graphics.ColorToNRGBA(in.Colors.TextColor)
										if !gtx.Enabled() {
											textColor = graphics.ColorToNRGBA(in.Colors.DisabledTextColor)
										}
										selectionColor := graphics.ColorToNRGBA(in.Colors.SelectionColor)

										ed := gioMaterial.Editor(th.T, &in.Editor.E, "")
										ed.Color = textColor
										ed.SelectionColor = selectionColor

										return ed.Layout(gtx)
									})
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if in.Suffix != nil {
										return layoutnode.ToGioDimensions(in.Suffix(layoutnode.NewLayoutContext(&gtx)))
									}
									return layout.Dimensions{}
								}),
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
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
				Spacing:   layout.SpaceBetween,
			}.Layout(
				gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if in.helper.Text == "" {
						return layout.Dimensions{}
					}
					return layout.Inset{
						Top:  gioUnit.Dp(4),
						Left: gioUnit.Dp(10),
					}.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							helper := gioMaterial.Label(th.T, gioUnit.Sp(12), in.helper.Text)
							helper.Color = in.helper.Color
							return helper.Layout(gtx)
						},
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if in.CharLimit == 0 {
						return layout.Dimensions{}
					}
					return layout.Inset{
						Top:   gioUnit.Dp(4),
						Right: gioUnit.Dp(10),
					}.Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							count := gioMaterial.Label(
								th.T,
								gioUnit.Sp(12),
								strconv.Itoa(in.Editor.Len())+"/"+strconv.Itoa(int(in.CharLimit)),
							)
							count.Color = in.helper.Color
							return count.Layout(gtx)
						},
					)
				}),
			)
		}),
	)
	return layoutnode.LayoutDimensions{
		Size: image.Point{
			X: dims.Size.X,
			Y: dims.Size.Y + in.label.Smallest.Size.Y/2,
		},
		Baseline: dims.Baseline,
	}
}

func (in *OutlinedTextFieldWidget) update(gtx layoutnode.LayoutContext, th *textinput.Theme, hint string) {
	disabled := textinput.Disabled(gtx)
	in.click.ProcessEvents(gtx, in.Editor.FocusTag())

	in.state = inactive
	if in.click.Hovered() && !disabled {
		in.state = hovered
	}
	hasContents := in.Editor.Len() > 0
	if hasContents {
		in.state = activated
	}
	if in.Editor.SourceFocused(gtx.ToGio().Source) && !disabled {
		in.state = focused
	}
	const (
		duration = time.Millisecond * 100
	)
	if in.anim == nil {
		in.anim = &Progress{}
	}
	if in.state == activated || hasContents {
		in.anim.Start(gtx.ToGio().Now, Forward, 0)
	}
	if in.state == focused && !hasContents && !in.anim.Started() {
		in.anim.Start(gtx.ToGio().Now, Forward, duration)
	}
	if in.state == inactive && !hasContents && in.anim.Finished() {
		in.anim.Start(gtx.ToGio().Now, Reverse, duration)
	}
	if in.anim.Started() {
		gtx.ToGio().Execute(op.InvalidateCmd{})
	}
	in.anim.Update(gtx.ToGio().Now)

	// Colors
	in.border.Color = graphics.ColorToNRGBA(in.Colors.UnfocusedIndicatorColor)
	in.helper.Color = graphics.ColorToNRGBA(in.Colors.SupportingTextColor)
	in.border.Thickness = gioUnit.Dp(1)

	if in.state == hovered {
		in.border.Color = graphics.ColorToNRGBA(in.Colors.HoveredIndicatorColor)
	} else if in.state == focused {
		in.border.Color = graphics.ColorToNRGBA(in.Colors.FocusedIndicatorColor)
		in.border.Thickness = gioUnit.Dp(2)
	}

	if disabled {
		in.border.Color = graphics.ColorToNRGBA(in.Colors.DisabledIndicatorColor)
		in.helper.Color = graphics.ColorToNRGBA(in.Colors.DisabledSupportingTextColor)
	}

	if in.IsErrored() {
		in.border.Color = graphics.ColorToNRGBA(in.Colors.ErrorIndicatorColor)
		in.helper.Color = graphics.ColorToNRGBA(in.Colors.ErrorSupportingTextColor)
	}

	// Label Logic
	g := *gtx.ToGio()
	textNormal := th.T.TextSize
	textSmall := th.T.TextSize * 0.8
	in.label.TextSize = gioUnit.Sp(lerp.Between32(float32(textSmall), float32(textNormal), 1.0-in.anim.Progress()))

	// Calculate smallest label for cutout
	g.Constraints.Min.X = 0
	macro := op.Record(g.Ops)
	var spacing gioUnit.Dp
	if len(hint) > 0 {
		spacing = 4
	}
	in.label.Smallest = layout.Inset{
		Left:  spacing,
		Right: spacing,
	}.Layout(g, func(gtx layout.Context) layout.Dimensions {
		l := gioMaterial.Label(th.T, textSmall, hint)
		l.Color = in.border.Color
		return l.Layout(gtx)
	})
	macro.Stop()

	// Calculate label position
	startTop := float32(g.Dp(16))
	endTop := float32(in.label.Smallest.Size.Y) / -2.0

	in.label.Inset = layout.Inset{
		Top:  gioUnit.Dp(lerp.Between32(startTop, endTop, in.anim.Progress())),
		Left: gioUnit.Dp(12),
	}

	// Editor Inset
	in.editorInset = struct {
		Top, Right, Bottom, Left unit.Dp
	}{
		Top:    unit.Dp(16),
		Bottom: unit.Dp(16),
	}
}
