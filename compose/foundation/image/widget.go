package image

import (
	"image"

	"github.com/zodimo/go-compose/compose/ui/layout"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/modifiers/size"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// ImageWidget is a widget that displays an image with scaling, alignment, and opacity.
type ImageWidget struct {
	// img is the source image to display. A paint.ImageOp is built from it
	// at draw time via engine-bound code.
	img image.Image
	// ContentScale specifies how to scale the image to the constraints.
	ContentScale layout.ContentScale
	// Alignment specifies where to position the image within the constraints.
	Alignment size.Alignment
	// Alpha applies opacity to the image (0.0 to 1.0).
	Alpha float32
}

func (im ImageWidget) Layout(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
	src := paint.NewImageOp(im.img)
	srcSize := src.Size()
	dstSize := gtx.ToGio().Constraints.Max

	srcDim := layoutnode.LayoutDimensions{Size: srcSize}
	dstDim := layoutnode.LayoutDimensions{Size: dstSize}

	// Calculate scale factor using ContentScale
	scaleFactor := im.ContentScale.Scale(srcDim, dstDim)

	// Calculate scaled size
	scaledW := int(float32(srcSize.X) * scaleFactor.ScaleX)
	scaledH := int(float32(srcSize.Y) * scaleFactor.ScaleY)
	scaledSize := image.Point{X: scaledW, Y: scaledH}

	// Constrain to reported size
	reportedSize := gtx.ToGio().Constraints.Constrain(scaledSize)

	// Calculate alignment offset
	offset := im.Alignment.Align(scaledSize, reportedSize, layoutnode.LayoutDirectionLTR)

	// Clip to reported bounds
	defer clip.Rect{Max: reportedSize}.Push(gtx.ToGio().Ops).Pop()

	// Apply transform (scale + offset)
	trans := f32.NewAffine2D(
		scaleFactor.ScaleX, 0, float32(offset.X),
		0, scaleFactor.ScaleY, float32(offset.Y),
	)
	defer op.Affine(trans).Push(gtx.ToGio().Ops).Pop()

	// Apply alpha
	if im.Alpha < 1.0 {
		defer paint.PushOpacity(gtx.ToGio().Ops, im.Alpha).Pop()
	}

	// Draw image
	src.Add(gtx.ToGio().Ops)
	paint.PaintOp{}.Add(gtx.ToGio().Ops)

	return layoutnode.LayoutDimensions{Size: reportedSize}
}
