package size

import (
	"go-compose-dev/internal/layoutnode"
	"image"
)

// Alignment is an interface that calculates the position of a child
// inside a parent container.
type Alignment interface {
	Align(size image.Point, space image.Point, layoutDirection layoutnode.LayoutDirection) image.Point
}

// BiasAlignment implements Alignment using a horizontal and vertical bias.
// bias is between -1 (start/top) and 1 (end/bottom).
type BiasAlignment struct {
	HorizontalBias float32
	VerticalBias   float32
}

func (b BiasAlignment) Align(size image.Point, space image.Point, _ layoutnode.LayoutDirection) image.Point {
	// Formula: (space - size) * (1 + bias) / 2
	// But usually bias is 0..1 or -1..1.
	// Compose uses -1..1.
	// (space - size) / 2 * (1 + bias)

	x := float32(space.X-size.X) / 2 * (1 + b.HorizontalBias)
	y := float32(space.Y-size.Y) / 2 * (1 + b.VerticalBias)
	return image.Point{X: int(x), Y: int(y)}
}

var (
	TopStart     = BiasAlignment{-1, -1}
	TopCenter    = BiasAlignment{0, -1}
	TopEnd       = BiasAlignment{1, -1}
	CenterStart  = BiasAlignment{-1, 0}
	Center       = BiasAlignment{0, 0}
	CenterEnd    = BiasAlignment{1, 0}
	BottomStart  = BiasAlignment{-1, 1}
	BottomCenter = BiasAlignment{0, 1}
	BottomEnd    = BiasAlignment{1, 1}
)
