package layout

import (
	"math"

	"github.com/zodimo/go-compose/internal/layoutnode"
)

// ContentScale defines how to scale the source content to fit the destination space.
type ContentScale interface {
	// Scale returns the scale factor to apply to the source to fit the destination.
	Scale(srcSize, dstSize layoutnode.LayoutDimensions) ScaleFactor
}

type ScaleFactor struct {
	ScaleX float32
	ScaleY float32
}

type contentScaleFunc func(srcSize, dstSize layoutnode.LayoutDimensions) ScaleFactor

func (f contentScaleFunc) Scale(srcSize, dstSize layoutnode.LayoutDimensions) ScaleFactor {
	return f(srcSize, dstSize)
}

var (
	// ContentScaleNone maintains the source size.
	ContentScaleNone ContentScale = contentScaleFunc(func(srcSize, dstSize layoutnode.LayoutDimensions) ScaleFactor {
		return ScaleFactor{ScaleX: 1.0, ScaleY: 1.0}
	})

	// ContentScaleFillBounds scales the source to fill the destination bounds, potentially changing the aspect ratio.
	ContentScaleFillBounds ContentScale = contentScaleFunc(func(srcSize, dstSize layoutnode.LayoutDimensions) ScaleFactor {
		return computeFillScale(srcSize, dstSize)
	})

	// ContentScaleFillWidth scales the source to fill the destination width, maintaining the aspect ratio.
	ContentScaleFillWidth ContentScale = contentScaleFunc(func(srcSize, dstSize layoutnode.LayoutDimensions) ScaleFactor {
		return computeFillWidth(srcSize, dstSize)
	})

	// ContentScaleFillHeight scales the source to fill the destination height, maintaining the aspect ratio.
	ContentScaleFillHeight ContentScale = contentScaleFunc(func(srcSize, dstSize layoutnode.LayoutDimensions) ScaleFactor {
		return computeFillHeight(srcSize, dstSize)
	})

	// ContentScaleFit scales the source to fit within the destination bounds, maintaining the aspect ratio.
	// The source will be completely visible.
	ContentScaleFit ContentScale = contentScaleFunc(func(srcSize, dstSize layoutnode.LayoutDimensions) ScaleFactor {
		return computeFitScale(srcSize, dstSize)
	})

	// ContentScaleCrop scales the source to fill the destination bounds, maintaining the aspect ratio.
	// The source may be cropped.
	ContentScaleCrop ContentScale = contentScaleFunc(func(srcSize, dstSize layoutnode.LayoutDimensions) ScaleFactor {
		return computeCropScale(srcSize, dstSize)
	})

	// ContentScaleInside scales the source to fit within the destination bounds only if the source is larger than the destination.
	// Otherwise, it maintains the source size.
	ContentScaleInside ContentScale = contentScaleFunc(func(srcSize, dstSize layoutnode.LayoutDimensions) ScaleFactor {
		if srcSize.Size.X <= dstSize.Size.X && srcSize.Size.Y <= dstSize.Size.Y {
			return ScaleFactor{ScaleX: 1.0, ScaleY: 1.0}
		}
		return computeFitScale(srcSize, dstSize)
	})
)

func computeFillScale(srcSize, dstSize layoutnode.LayoutDimensions) ScaleFactor {
	return ScaleFactor{
		ScaleX: float32(dstSize.Size.X) / float32(srcSize.Size.X),
		ScaleY: float32(dstSize.Size.Y) / float32(srcSize.Size.Y),
	}
}

func computeFillWidth(srcSize, dstSize layoutnode.LayoutDimensions) ScaleFactor {
	scale := float32(dstSize.Size.X) / float32(srcSize.Size.X)
	return ScaleFactor{ScaleX: scale, ScaleY: scale}
}

func computeFillHeight(srcSize, dstSize layoutnode.LayoutDimensions) ScaleFactor {
	scale := float32(dstSize.Size.Y) / float32(srcSize.Size.Y)
	return ScaleFactor{ScaleX: scale, ScaleY: scale}
}

func computeFitScale(srcSize, dstSize layoutnode.LayoutDimensions) ScaleFactor {
	widthScale := float32(dstSize.Size.X) / float32(srcSize.Size.X)
	heightScale := float32(dstSize.Size.Y) / float32(srcSize.Size.Y)
	scale := float32(math.Min(float64(widthScale), float64(heightScale)))
	return ScaleFactor{ScaleX: scale, ScaleY: scale}
}

func computeCropScale(srcSize, dstSize layoutnode.LayoutDimensions) ScaleFactor {
	widthScale := float32(dstSize.Size.X) / float32(srcSize.Size.X)
	heightScale := float32(dstSize.Size.Y) / float32(srcSize.Size.Y)
	scale := float32(math.Max(float64(widthScale), float64(heightScale)))
	return ScaleFactor{ScaleX: scale, ScaleY: scale}
}
