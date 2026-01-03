package textfield

import (
	"image"
	"image/color"
	"strconv"
	"time"

	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/input"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"

	gioUnit "gioui.org/unit"
)

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

type HandlerWrapper struct {
	Func func(string)
}

type OnSubmitWrapper struct {
	Func func()
}

type TextFieldStateTracker struct {
	LastValue string
}

// TextField implements the Material Design Text Field
// described here: https://material.io/components/text-fields
type TextFieldWidget struct {
	// Editor contains the edit buffer.
	Editor *widget.Editor
	// click detects when the mouse pointer clicks or hovers
	// within the textfield.
	click gesture.Click

	// Helper text to give additional context to a field.
	Helper string
	// CharLimit specifies the maximum number of characters the text input
	// will allow. Zero means "no limit".
	CharLimit uint
	// Prefix appears before the content of the text input.
	Prefix layout.Widget
	// Suffix appears after the content of the text input.
	Suffix layout.Widget

	Colors TextFieldColors

	// Animation state.
	state
	label  label
	border border
	helper helper
	anim   *Progress

	// errored tracks whether the input is in an errored state.
	// This is orthogonal to the other states: the input can be both errored
	// and inactive for example.
	errored bool
}

// IsActive if input is in an active state (Active, Focused or Errored).
func (in *TextFieldWidget) IsActive() bool {
	return in.state >= activated
}

// IsErrored if input is in an errored state.
// Typically this is when the validator has returned an error message.
func (in *TextFieldWidget) IsErrored() bool {
	return in.errored
}

// SetError puts the input into an errored state with the specified error text.
func (in *TextFieldWidget) SetError(isError bool, err string) {
	in.errored = isError
	in.helper.Text = err
}

// ClearError clears any errored status.
func (in *TextFieldWidget) ClearError() {
	in.errored = false
	in.helper.Text = in.Helper
}

// Clear the input text and reset any error status.
func (in *TextFieldWidget) Clear() {
	in.Editor.SetText("")
	in.ClearError()
}

// TextTooLong returns whether the current editor text exceeds the set character
// limit.
func (in *TextFieldWidget) TextTooLong() bool {
	return !(in.CharLimit == 0 || uint(len(in.Editor.Text())) < in.CharLimit)
}

func (in *TextFieldWidget) Layout(gtx layout.Context, th *material.Theme, hint string) layout.Dimensions {
	// Logic from gio-x Update + Layout
	in.update(gtx, th, hint)

	// Offset accounts for label height, which sticks above the border dimensions.
	defer op.Offset(image.Pt(0, in.label.Smallest.Size.Y/2)).Push(gtx.Ops).Pop()

	// Draw Label
	in.label.Inset.Layout(
		gtx,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Left:  gioUnit.Dp(4),
				Right: gioUnit.Dp(4),
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
					return layout.UniformInset(gioUnit.Dp(12)).Layout(
						gtx,
						func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return layout.Flex{
								Axis:      layout.Horizontal,
								Alignment: layout.Middle,
							}.Layout(
								gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if in.IsActive() && in.Prefix != nil {
										return in.Prefix(gtx)
									}
									return layout.Dimensions{}
								}),
								// Prefix would go here
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									// Resolve editor colors
									textColor := graphics.ColorToNRGBA(in.Colors.TextColor)
									if !gtx.Enabled() {
										textColor = graphics.ColorToNRGBA(in.Colors.DisabledTextColor)
									}
									selectionColor := graphics.ColorToNRGBA(in.Colors.SelectionColor)

									ed := material.Editor(th, in.Editor, "")
									ed.Color = textColor
									ed.SelectionColor = selectionColor

									return ed.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if in.IsActive() && in.Suffix != nil {
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
							helper := material.Label(th, gioUnit.Sp(12), in.helper.Text)
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
							count := material.Label(
								th,
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
	return layout.Dimensions{
		Size: image.Point{
			X: dims.Size.X,
			Y: dims.Size.Y + in.label.Smallest.Size.Y/2,
		},
		Baseline: dims.Baseline,
	}
}

func (in *TextFieldWidget) update(gtx layout.Context, th *material.Theme, hint string) {

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

	var (
		textNormal = th.TextSize
		textSmall  = th.TextSize * 0.8

		borderColor         = graphics.ColorToNRGBA(in.Colors.UnfocusedIndicatorColor)
		borderColorHovered  = graphics.ColorToNRGBA(in.Colors.HoveredIndicatorColor)
		borderColorActive   = graphics.ColorToNRGBA(in.Colors.FocusedIndicatorColor)
		borderColorError    = graphics.ColorToNRGBA(in.Colors.ErrorIndicatorColor)
		borderColorDisabled = graphics.ColorToNRGBA(in.Colors.DisabledIndicatorColor)

		borderThickness       = gioUnit.Dp(1)
		borderThicknessActive = gioUnit.Dp(2)

		helperColor         = graphics.ColorToNRGBA(in.Colors.SupportingTextColor)
		helperColorError    = graphics.ColorToNRGBA(in.Colors.ErrorSupportingTextColor)
		helperColorDisabled = graphics.ColorToNRGBA(in.Colors.DisabledSupportingTextColor)
	)

	if disabled {
		borderColor = borderColorDisabled
		borderColorHovered = borderColorDisabled
		borderColorActive = borderColorDisabled
		borderColorError = borderColorDisabled
		helperColor = helperColorDisabled
		helperColorError = helperColorDisabled
	}

	in.label.TextSize = gioUnit.Sp(lerp.Between32(float32(textSmall), float32(textNormal), 1.0-in.anim.Progress()))
	switch in.state {
	case inactive:
		in.border.Thickness = borderThickness
		in.border.Color = borderColor
		in.helper.Color = helperColor
	case hovered, activated:
		in.border.Thickness = borderThickness
		in.border.Color = borderColorHovered
		in.helper.Color = helperColor
	case focused:
		in.border.Thickness = borderThicknessActive
		in.border.Color = borderColorActive
		in.helper.Color = helperColor
	}

	if in.IsErrored() {
		in.border.Color = borderColorError
		in.helper.Color = helperColorError
	}

	// Calculate smallest label for cutout
	gtx.Constraints.Min.X = 0
	macro := op.Record(gtx.Ops)
	var spacing gioUnit.Dp
	if len(hint) > 0 {
		spacing = 4
	}
	in.label.Smallest = layout.Inset{
		Left:  spacing,
		Right: spacing,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		l := material.Label(th, textSmall, hint)
		l.Color = in.border.Color
		return l.Layout(gtx)
	})
	macro.Stop()

	labelTopInsetNormal := float32(in.label.Smallest.Size.Y) - float32(in.label.Smallest.Size.Y/4)
	topInsetDP := gioUnit.Dp(labelTopInsetNormal / gtx.Metric.PxPerDp)
	topInsetActiveDP := (topInsetDP / 2 * -1) - gioUnit.Dp(in.border.Thickness)
	in.label.Inset = layout.Inset{
		Top:  gioUnit.Dp(lerp.Between32(float32(topInsetDP), float32(topInsetActiveDP), in.anim.Progress())),
		Left: gioUnit.Dp(10),
	}
}
