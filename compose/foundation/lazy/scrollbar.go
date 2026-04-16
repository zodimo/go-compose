package lazy

import (
	"image"
	"image/color"
	"math"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

func clamp1(v float32) float32 {
	if v >= 1 {
		return 1
	} else if v <= 0 {
		return 0
	}
	return v
}

func fromListPosition(lp layout.Position, elements int, majorAxisSize int) (start, end float32) {
	lengthEstPx := float32(lp.Length)
	elementLenEstPx := lengthEstPx / float32(elements)

	listOffsetF := float32(lp.Offset)
	listOffsetL := float32(lp.OffsetLast)

	viewportStart := clamp1((float32(lp.First)*elementLenEstPx + listOffsetF) / lengthEstPx)
	viewportEnd := clamp1((float32(lp.First+lp.Count)*elementLenEstPx + listOffsetL) / lengthEstPx)
	viewportFraction := viewportEnd - viewportStart

	visiblePx := float32(majorAxisSize)
	visibleFraction := visiblePx / lengthEstPx

	err := visibleFraction - viewportFraction
	adjStart := viewportStart
	adjEnd := viewportEnd
	if viewportFraction < 1 {
		startShare := viewportStart / (1 - viewportFraction)
		endShare := (1 - viewportEnd) / (1 - viewportFraction)
		startErr := startShare * err
		endErr := endShare * err

		adjStart -= startErr
		adjEnd += endErr
	}
	return adjStart, adjEnd
}

func rangeIsScrollable(start, end float32) bool {
	return end-start < 1
}

type scrollbarStyle struct {
	scrollbar  *widget.Scrollbar
	trackPad   unit.Dp
	indicatorW unit.Dp
	indicatorR unit.Dp
	color      color.NRGBA
	hoverColor color.NRGBA
}

func layoutScrollbar(gtx layout.Context, list *widget.List, axis layout.Axis, itemCount int, dims layout.Dimensions) {
	majorAxisSize := axis.Convert(dims.Size).X
	start, end := fromListPosition(list.Position, itemCount, majorAxisSize)

	if !rangeIsScrollable(start, end) {
		return
	}

	sb := scrollbarStyle{
		scrollbar:  &list.Scrollbar,
		trackPad:   2,
		indicatorW: 6,
		indicatorR: 3,
		color:      color.NRGBA{R: 128, G: 128, B: 128, A: 150},
		hoverColor: color.NRGBA{R: 128, G: 128, B: 128, A: 200},
	}

	anchoring := layout.E
	if axis == layout.Horizontal {
		anchoring = layout.S
	}

	gtx.Constraints.Min = dims.Size
	gtx.Constraints.Max = dims.Size
	anchoring.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		list.Scrollbar.Update(gtx, axis, start, end)

		if list.Scrollbar.IndicatorHovered() {
			sb.color = sb.hoverColor
		}

		return sb.layout(gtx, axis, start, end)
	})

	if delta := list.ScrollDistance(); delta != 0 {
		list.List.ScrollBy(delta * float32(itemCount))
	}
}

func (s scrollbarStyle) layout(gtx layout.Context, axis layout.Axis, viewportStart, viewportEnd float32) layout.Dimensions {
	trackPad := gtx.Dp(s.trackPad)
	indicatorWidth := gtx.Dp(s.indicatorW)
	cornerRadius := gtx.Dp(s.indicatorR)

	convert := axis.Convert

	maxMajor := convert(gtx.Constraints.Max).X
	minMinor := indicatorWidth + trackPad + trackPad

	gtx.Constraints.Min = convert(image.Pt(maxMajor, minMinor))
	gtx.Constraints.Max = gtx.Constraints.Min

	inset := layout.Inset{
		Top:    s.trackPad,
		Bottom: s.trackPad,
		Left:   s.trackPad,
		Right:  s.trackPad,
	}
	if axis == layout.Horizontal {
		inset.Top, inset.Bottom, inset.Left, inset.Right = inset.Left, inset.Right, inset.Top, inset.Bottom
	}

	return layout.Background{}.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			area := image.Rectangle{Max: gtx.Constraints.Min}
			pointerArea := clip.Rect(area)
			defer pointerArea.Push(gtx.Ops).Pop()
			s.scrollbar.AddDrag(gtx.Ops)

			defer pointer.PassOp{}.Push(gtx.Ops).Pop()
			defer pointerArea.Push(gtx.Ops).Pop()
			s.scrollbar.AddTrack(gtx.Ops)

			return layout.Dimensions{Size: gtx.Constraints.Min}
		},
		func(gtx layout.Context) layout.Dimensions {
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = convert(gtx.Constraints.Min)
				gtx.Constraints.Max = convert(gtx.Constraints.Max)

				trackLen := gtx.Constraints.Min.X
				viewStart := int(math.Round(float64(viewportStart) * float64(trackLen)))
				viewEnd := int(math.Round(float64(viewportEnd) * float64(trackLen)))
				indicatorLen := viewEnd - viewStart
				minLen := gtx.Dp(18)
				if indicatorLen < minLen {
					indicatorLen = minLen
				}
				if viewStart+indicatorLen > trackLen {
					viewStart = trackLen - indicatorLen
				}
				if viewStart < 0 {
					viewStart = 0
				}

				indicatorDims := convert(image.Pt(indicatorLen, indicatorWidth))

				offset := convert(image.Pt(viewStart, 0))
				defer op.Offset(offset).Push(gtx.Ops).Pop()

				paint.FillShape(gtx.Ops, s.color, clip.RRect{
					Rect: image.Rectangle{
						Max: indicatorDims,
					},
					SW: cornerRadius,
					NW: cornerRadius,
					NE: cornerRadius,
					SE: cornerRadius,
				}.Op(gtx.Ops))

				area := clip.Rect(image.Rectangle{Max: indicatorDims})
				defer pointer.PassOp{}.Push(gtx.Ops).Pop()
				defer area.Push(gtx.Ops).Pop()
				s.scrollbar.AddIndicator(gtx.Ops)

				return layout.Dimensions{Size: convert(gtx.Constraints.Min)}
			})
		},
	)
}
