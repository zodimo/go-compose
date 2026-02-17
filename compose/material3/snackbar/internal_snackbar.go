// // SPDX-License-Identifier: Unlicense OR MIT

package snackbar

// import (
// 	"fmt"
// 	"image"
// 	"log"

// 	"gioui.org/f32"
// 	"gioui.org/layout"
// 	"gioui.org/op"
// 	"gioui.org/op/clip"
// 	"gioui.org/op/paint"
// 	"gioui.org/text"
// 	"gioui.org/unit"
// 	"gioui.org/widget"
// 	"git.sr.ht/~schnwalter/gio-mw/token"
// 	"git.sr.ht/~schnwalter/gio-mw/wdk"
// 	"git.sr.ht/~schnwalter/gio-mw/wdk/block"
// 	"git.sr.ht/~schnwalter/gio-mw/widget/overlay"
// )

// var (
// 	defaultSpacing         = unit.Dp(8)
// 	paddingStart           = unit.Dp(16)
// 	paddingEnd             = unit.Dp(16)
// 	paddingRightWithAction = unit.Dp(8)
// 	paddingRightWithIcon   = unit.Dp(0)
// 	paddingIcon            = unit.Dp(12)
// 	recommendedMaxLength   = 68
// )

// type Style struct {
// 	sTheme         *Theme
// 	supportingText string
// 	overlayItemId  int64

// 	// Action button support
// 	actionLabel     string
// 	actionClickable *widget.Clickable
// 	onAction        func()

// 	// Dismiss icon support
// 	withDismiss      bool
// 	dismissClickable *widget.Clickable
// 	onDismiss        func()
// }

// // ActionClicked returns true if the action button was clicked during this frame.
// func (s *Style) ActionClicked(gtx layout.Context) bool {
// 	if s.actionClickable == nil {
// 		return false
// 	}
// 	return s.actionClickable.Clicked(gtx)
// }

// // DismissClicked returns true if the dismiss icon was clicked during this frame.
// func (s *Style) DismissClicked(gtx layout.Context) bool {
// 	if s.dismissClickable == nil {
// 		return false
// 	}
// 	return s.dismissClickable.Clicked(gtx)
// }

// func (s *Style) AsOverlayItem() *overlay.Item {
// 	overlayItem := overlay.NewItem(s.Layout, block.GravityBottomCenter)
// 	s.overlayItemId = overlayItem.GetId()
// 	return overlayItem
// }

// func (s *Style) WithTheme(theme *Theme) *Style {
// 	s.sTheme = theme
// 	return s
// }

// func (s *Style) Layout(gtx layout.Context) layout.Dimensions {
// 	s.sTheme = BuildTheme(gtx)
// 	if s.supportingText == "" {
// 		return layout.Dimensions{}
// 	}

// 	// Process click events and invoke callbacks
// 	if s.actionClickable != nil && s.actionClickable.Clicked(gtx) && s.onAction != nil {
// 		s.onAction()
// 	}
// 	if s.dismissClickable != nil && s.dismissClickable.Clicked(gtx) && s.onDismiss != nil {
// 		s.onDismiss()
// 	}

// 	return block.Padding{
// 		Bottom: unit.Dp(16),
// 		Start:  unit.Dp(16),
// 		End:    unit.Dp(16),
// 	}.Layout(gtx, s.widgetLayout)
// }

// func (s *Style) widgetLayout(gtx layout.Context) layout.Dimensions {
// 	macroOp := op.Record(gtx.Ops)
// 	if gtx.Constraints.Min.Y < gtx.Dp(s.sTheme.EnabledContainerOneLineHeight) {
// 		gtx.Constraints.Min.Y = gtx.Dp(s.sTheme.EnabledContainerOneLineHeight)
// 	}
// 	contentDim := s.childrenLayout(gtx)
// 	contentCallOp := macroOp.Stop()

// 	// Build the snackbar shape.
// 	baseBox := wdk.Box{
// 		Shape:    wdk.FromCornerShapesToken(gtx, s.sTheme.EnabledContainerShape),
// 		EndPoint: contentDim.Size,
// 	}

// 	// Draw snackbar elevation shadow.
// 	sElevation := wdk.Elevation{
// 		Level:       s.sTheme.EnabledContainerElevation,
// 		ShadowColor: s.sTheme.EnabledContainerShadowColor,
// 	}
// 	sElevation.Layout(gtx, baseBox)

// 	// Draw snackbar background and content.
// 	clipStack := baseBox.Outline(gtx).Push(gtx.Ops)
// 	paint.Fill(gtx.Ops, s.sTheme.EnabledContainerColor.AsNRGBA())
// 	contentCallOp.Add(gtx.Ops)
// 	clipStack.Pop()

// 	return contentDim
// }

// func (s *Style) childrenLayout(gtx layout.Context) layout.Dimensions {
// 	c := block.Container{
// 		Gravity: block.GravityMiddleStart,
// 	}

// 	hasAction := s.actionLabel != "" && s.actionClickable != nil
// 	hasDismiss := s.withDismiss && s.dismissClickable != nil

// 	// Determine end padding based on what's displayed
// 	endPad := paddingEnd
// 	if hasAction {
// 		endPad = paddingRightWithAction
// 	}
// 	if hasDismiss {
// 		endPad = paddingRightWithIcon
// 	}

// 	return c.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
// 		return block.Padding{
// 			Start: paddingStart,
// 			End:   endPad,
// 		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
// 			var widgets []block.Segment

// 			tLength := len(s.supportingText)
// 			if tLength > recommendedMaxLength {
// 				// TODO: Use slog.
// 				warningText := fmt.Errorf("snackbar text too long, %v > %v", tLength, recommendedMaxLength)
// 				log.Println(warningText)
// 			}
// 			textWidget := func(gtx layout.Context) layout.Dimensions {
// 				gtx.Constraints.Min.Y = 0
// 				presentation := wdk.LabelStyle{
// 					Alignment: text.Middle,
// 					Color:     s.sTheme.EnabledSupportingTextColor,
// 					MaxLines:  2,
// 					Typestyle: token.TypestyleLabelMedium,
// 				}
// 				return wdk.LayoutLabel(gtx, presentation, s.supportingText)
// 			}

// 			// Use flex segment for text when action/dismiss is present so text fills remaining space
// 			if hasAction || hasDismiss {
// 				widgets = append(widgets, block.NewFlexSegment(textWidget))
// 			} else {
// 				widgets = append(widgets, block.NewSegment(textWidget))
// 			}

// 			// Add action button
// 			if hasAction {
// 				widgets = append(widgets, block.NewHorizontalSpacer(defaultSpacing))
// 				actionWidget := func(gtx layout.Context) layout.Dimensions {
// 					gtx.Constraints.Min.Y = 0
// 					return s.actionClickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
// 						presentation := wdk.LabelStyle{
// 							Alignment: text.Middle,
// 							Color:     s.sTheme.EnabledLabelColor,
// 							MaxLines:  1,
// 							Typestyle: token.TypestyleLabelMedium,
// 						}
// 						return block.Padding{
// 							Start: paddingRightWithAction,
// 							End:   paddingRightWithAction,
// 						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
// 							return wdk.LayoutLabel(gtx, presentation, s.actionLabel)
// 						})
// 					})
// 				}
// 				widgets = append(widgets, block.NewSegment(actionWidget).AlignMiddle())
// 			}

// 			// Add dismiss icon
// 			if hasDismiss {
// 				widgets = append(widgets, block.NewHorizontalSpacer(defaultSpacing))
// 				dismissWidget := func(gtx layout.Context) layout.Dimensions {
// 					gtx.Constraints.Min.Y = 0
// 					iconSize := gtx.Dp(s.sTheme.EnabledIconSize)
// 					return s.dismissClickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
// 						return block.Padding{
// 							Start:  paddingIcon,
// 							End:    paddingIcon,
// 							Top:    paddingIcon,
// 							Bottom: paddingIcon,
// 						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
// 							return layoutCloseIcon(gtx, iconSize, s.sTheme.EnabledIconColor)
// 						})
// 					})
// 				}
// 				widgets = append(widgets, block.NewSegment(dismissWidget).AlignMiddle())
// 			}

// 			return block.Line{
// 				Axis:     block.AxisHorizontal,
// 				Overflow: block.OverflowClip,
// 				Expand:   hasAction || hasDismiss,
// 			}.Layout(gtx, widgets...)
// 		})
// 	})
// }

// // layoutCloseIcon draws a simple X (close) icon.
// func layoutCloseIcon(gtx layout.Context, sizePx int, color token.MatColor) layout.Dimensions {
// 	sz := image.Pt(sizePx, sizePx)

// 	// Draw X using two lines
// 	nrgba := color.AsNRGBA()
// 	strokeWidth := float32(sizePx) / 12
// 	if strokeWidth < 1 {
// 		strokeWidth = 1
// 	}

// 	m := float32(sizePx) * 0.2 // margin

// 	// Line 1: top-left to bottom-right
// 	{
// 		var p clip.Path
// 		p.Begin(gtx.Ops)
// 		p.MoveTo(f32.Pt(m, m))
// 		p.LineTo(f32.Pt(float32(sizePx)-m, float32(sizePx)-m))
// 		spec := p.End()
// 		clipStack := clip.Stroke{Path: spec, Width: strokeWidth}.Op().Push(gtx.Ops)
// 		paint.ColorOp{Color: nrgba}.Add(gtx.Ops)
// 		paint.PaintOp{}.Add(gtx.Ops)
// 		clipStack.Pop()
// 	}

// 	// Line 2: top-right to bottom-left
// 	{
// 		var p clip.Path
// 		p.Begin(gtx.Ops)
// 		p.MoveTo(f32.Pt(float32(sizePx)-m, m))
// 		p.LineTo(f32.Pt(m, float32(sizePx)-m))
// 		spec := p.End()
// 		clipStack := clip.Stroke{Path: spec, Width: strokeWidth}.Op().Push(gtx.Ops)
// 		paint.ColorOp{Color: nrgba}.Add(gtx.Ops)
// 		paint.PaintOp{}.Add(gtx.Ops)
// 		clipStack.Pop()
// 	}

// 	return layout.Dimensions{Size: sz}
// }
