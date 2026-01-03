package textfield

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/pkg/sentinel"

	"gioui.org/gesture"
	"gioui.org/io/input"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	gioUnit "gioui.org/unit"
	"gioui.org/widget"
	gioMaterial "gioui.org/widget/material"
	"github.com/zodimo/go-compose/compose/material"
)

const Material3FilledTextFieldNodeID = "Material3FilledTextField"

// Filled implements the Filled Material Design 3 text field.
// It uses a custom widget implementation.
func Filled(
	value string,
	onValueChange func(string),
	options ...TextFieldOption,
) Composable {

	opts := DefaultTextFieldOptions()
	for _, opt := range options {
		opt(&opts)
	}

	return func(c Composer) Composer {
		theme := material.Theme(c)

		opts.Colors = ResolveTextFieldColors(c, opts.Colors)
		opts.SupportingText = sentinel.TakeOrElseString(opts.SupportingText, "")

		key := c.GenerateID()
		path := c.GetPath()

		// Handler wrappers
		handlerWrapperState := c.State(fmt.Sprintf("%d/%s/handler_wrapper", key, path), func() any {
			return &HandlerWrapper{Func: onValueChange}
		})
		handlerWrapper := handlerWrapperState.Get().(*HandlerWrapper)
		handlerWrapper.Func = onValueChange

		var onSubmitWrapper *OnSubmitWrapper
		if opts.OnSubmit != nil {
			onSubmitWrapperState := c.State(fmt.Sprintf("%d/%s/onsubmit_wrapper", key, path), func() any {
				return &OnSubmitWrapper{Func: opts.OnSubmit}
			})
			onSubmitWrapper = onSubmitWrapperState.Get().(*OnSubmitWrapper)
			onSubmitWrapper.Func = opts.OnSubmit
		}

		// Widget state
		widgetStatePath := fmt.Sprintf("%d/%s/filled_widget/s%v", key, path, opts.SingleLine)
		widgetVal := c.State(widgetStatePath, func() any {
			return &FilledTextFieldWidget{
				Editor: widget.Editor{
					SingleLine: opts.SingleLine,
					Submit:     opts.OnSubmit != nil,
					Mask:       opts.Mask,
				},
			}
		})
		w := widgetVal.Get().(*FilledTextFieldWidget)

		// Tracker
		trackerState := c.State(fmt.Sprintf("%d/%s/tracker/s%v", key, path, opts.SingleLine), func() any {
			return &TextFieldStateTracker{LastValue: ""}
		})
		tracker := trackerState.Get().(*TextFieldStateTracker)

		// Update properties
		w.Editor.SingleLine = opts.SingleLine
		w.Editor.Submit = opts.OnSubmit != nil
		w.Editor.Mask = opts.Mask
		w.Helper = opts.SupportingText
		w.SetError(opts.IsError, opts.SupportingText)

		c.StartBlock(Material3FilledTextFieldNodeID)
		c.Modifier(func(m Modifier) Modifier {
			return m.Then(opts.Modifier)
		})

		// Compose slots
		if opts.LeadingIcon != nil {
			c.WithComposable(opts.LeadingIcon)
		}
		if opts.TrailingIcon != nil {
			c.WithComposable(opts.TrailingIcon)
		}

		c.SetWidgetConstructor(filledTextFieldWidgetConstructor(w, value, opts, handlerWrapper, onSubmitWrapper, tracker, theme))

		return c.EndBlock()
	}
}

func filledTextFieldWidgetConstructor(
	w *FilledTextFieldWidget,
	value string,
	opts TextFieldOptions,
	handler *HandlerWrapper,
	onSubmitHandler *OnSubmitWrapper,
	tracker *TextFieldStateTracker,
	theme material.ThemeInterface,
) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			// Map children to slots (Must be done here, after WrapChildren)
			children := node.Children()
			childIdx := 0

			w.Prefix = nil
			if opts.LeadingIcon != nil && childIdx < len(children) {
				child := children[childIdx]
				if coord, ok := child.(layoutnode.NodeCoordinator); ok {
					w.Prefix = func(gtx layout.Context) layout.Dimensions {
						return coord.Layout(gtx)
					}
				}
				childIdx++
			}

			w.Suffix = nil
			if opts.TrailingIcon != nil && childIdx < len(children) {
				child := children[childIdx]
				if coord, ok := child.(layoutnode.NodeCoordinator); ok {
					w.Suffix = func(gtx layout.Context) layout.Dimensions {
						return coord.Layout(gtx)
					}
				}
				childIdx++
			}

			// 1. Sync
			if value != tracker.LastValue {
				if w.Editor.Text() != value {
					w.Editor.SetText(value)
				}
				tracker.LastValue = value
			}

			// 2. Events
			th := theme.GioMaterialTheme()
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

			// 3. Change Detection
			currentText := w.Editor.Text()
			if currentText != value {
				if handler.Func != nil {
					handler.Func(currentText)
				}
			}

			w.Colors = opts.Colors

			// 4. Layout
			return w.Layout(gtx, th, opts.Label)
		}
	})
}

// FilledTextFieldWidget is the local implementation of the Filled text field.
type FilledTextFieldWidget struct {
	widget.Editor
	click gesture.Click

	Helper string
	Colors TextFieldColors
	Prefix layout.Widget
	Suffix layout.Widget

	state
	label  label
	border border
	helper helper
	anim   *Progress

	errored bool
}

func (in *FilledTextFieldWidget) Layout(gtx layout.Context, th *gioMaterial.Theme, hint string) layout.Dimensions {
	in.update(gtx, th, hint)

	// Helper function to draw box
	drawBox := func(gtx layout.Context, size image.Point, color color.NRGBA) {
		defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, color)
	}

	dims := layout.Flex{
		Axis: layout.Vertical,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(
				gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					// Background
					bgColor := graphics.ColorToNRGBA(in.Colors.ContainerColor)
					if !gtx.Enabled() {
						bgColor = graphics.ColorToNRGBA(in.Colors.DisabledContainerColor)
					}
					drawBox(gtx, gtx.Constraints.Min, bgColor)

					// Active Indicator Line
					indicatorH := gtx.Dp(1)
					if in.state == focused || in.state == activated {
						indicatorH = gtx.Dp(2)
					}

					rect := image.Rectangle{
						Min: image.Point{0, gtx.Constraints.Min.Y - indicatorH},
						Max: gtx.Constraints.Min,
					}
					defer clip.Rect(rect).Push(gtx.Ops).Pop()
					paint.Fill(gtx.Ops, in.border.Color)

					return layout.Dimensions{Size: gtx.Constraints.Min}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Left:  gioUnit.Dp(16),
						Right: gioUnit.Dp(16),
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
										return in.Prefix(gtx)
									}
									return layout.Dimensions{}
								}),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return layout.Inset{
										Top:    gioUnit.Dp(16),
										Bottom: gioUnit.Dp(16),
									}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										textColor := graphics.ColorToNRGBA(in.Colors.TextColor)
										if !gtx.Enabled() {
											textColor = graphics.ColorToNRGBA(in.Colors.DisabledTextColor)
										}
										selectionColor := graphics.ColorToNRGBA(in.Colors.SelectionColor)

										ed := gioMaterial.Editor(th, &in.Editor, "")
										ed.Color = textColor
										ed.SelectionColor = selectionColor
										return ed.Layout(gtx)
									})
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if in.Suffix != nil {
										return in.Suffix(gtx)
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
			// Helper text
			if in.helper.Text == "" {
				return layout.Dimensions{}
			}
			return layout.Inset{
				Top:  gioUnit.Dp(4),
				Left: gioUnit.Dp(16),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				helper := gioMaterial.Label(th, gioUnit.Sp(12), in.helper.Text)
				helper.Color = in.helper.Color
				return helper.Layout(gtx)
			})
		}),
	)

	// Layout Label on top
	macro := op.Record(gtx.Ops)
	in.label.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		l := gioMaterial.Label(th, in.label.TextSize, hint)
		l.Color = in.border.Color
		if in.IsErrored() {
			l.Color = graphics.ColorToNRGBA(in.Colors.ErrorLabelColor)
		} else if in.state == focused {
			l.Color = graphics.ColorToNRGBA(in.Colors.FocusedLabelColor)
		} else {
			l.Color = graphics.ColorToNRGBA(in.Colors.LabelColor)
		}

		return l.Layout(gtx)
	})
	op.Defer(gtx.Ops, macro.Stop())

	return dims
}

func (in *FilledTextFieldWidget) update(gtx layout.Context, th *gioMaterial.Theme, hint string) {
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
	if in.Editor.Len() > 0 {
		in.state = activated
	}
	if gtx.Source.Focused(&in.Editor) && !disabled {
		in.state = focused
	}

	if in.anim == nil {
		in.anim = &Progress{}
	}
	// Animation logic
	if in.state == activated || in.Editor.Len() > 0 || (in.state == focused && in.Editor.Len() == 0) {
		in.anim.Start(gtx.Now, Forward, time.Millisecond*100)
	} else if in.state == inactive && in.Editor.Len() == 0 {
		in.anim.Start(gtx.Now, Reverse, time.Millisecond*100)
	}
	if in.anim.Started() {
		gtx.Execute(op.InvalidateCmd{})
	}
	in.anim.Update(gtx.Now)

	in.border.Color = graphics.ColorToNRGBA(in.Colors.UnfocusedIndicatorColor)
	in.helper.Color = graphics.ColorToNRGBA(in.Colors.SupportingTextColor)

	if in.state == focused {
		in.border.Color = graphics.ColorToNRGBA(in.Colors.FocusedIndicatorColor)
	}
	if in.IsErrored() {
		in.border.Color = graphics.ColorToNRGBA(in.Colors.ErrorIndicatorColor)
		in.helper.Color = graphics.ColorToNRGBA(in.Colors.ErrorSupportingTextColor)
	}

	textNormal := th.TextSize
	textSmall := th.TextSize * 0.75
	in.label.TextSize = gioUnit.Sp(lerp(float32(textNormal), float32(textSmall), in.anim.Progress()))

	startTop := float32(gtx.Dp(16))
	endTop := float32(gtx.Dp(8))
	in.label.Inset = layout.Inset{
		Top:  gioUnit.Dp(lerp(startTop, endTop, in.anim.Progress())),
		Left: gioUnit.Dp(16),
	}
}

func (in *FilledTextFieldWidget) IsErrored() bool {
	return in.errored
}

func (in *FilledTextFieldWidget) SetError(isError bool, err string) {
	in.errored = isError
	in.helper.Text = err
}
