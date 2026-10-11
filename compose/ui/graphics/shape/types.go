package shape

import (
	"fmt"
	"image"

	"gioui.org/op"
	"gioui.org/op/clip"
	gioUnit "gioui.org/unit"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// https://developer.android.com/reference/kotlin/androidx/compose/ui/graphics/Shape

type ShapeOption func(Shape)

type Shape interface {
	CreateOutline(size image.Point, metric gioUnit.Metric) Outline
	mergeShape(other Shape) Shape
	sameShape(other Shape) bool
	semanticEqualShape(other Shape) bool
	copyShape(options ...ShapeOption) Shape
	stringShape() string
}

type Outline interface {
	Push(ops *op.Ops) clip.Stack
	Op(ops *op.Ops) clip.Op
	Path(ops *op.Ops) clip.PathSpec
}

func IsSpecifiedShape(s Shape) bool {
	return s != nil && s != shapeUnspecified
}

func TakeOrElseShape(s, defaultStyle Shape) Shape {
	if s == nil || s == shapeUnspecified {
		return defaultStyle
	}
	return s
}

func CoalesceShape(ptr, def Shape) Shape {
	if ptr == nil {
		return def
	}
	return ptr
}

func MergeShape(a, b Shape) Shape {
	a = CoalesceShape(a, shapeUnspecified)
	b = CoalesceShape(b, shapeUnspecified)

	if a == shapeUnspecified {
		return b
	}
	if b == shapeUnspecified {
		return a
	}

	return a.mergeShape(b)
}

func SameShape(a, b Shape) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == shapeUnspecified
	}
	if b == nil {
		return a == shapeUnspecified
	}
	return a.sameShape(b)
}

func SemanticEqualShape(a, b Shape) bool {

	a = CoalesceShape(a, shapeUnspecified)
	b = CoalesceShape(b, shapeUnspecified)

	return a.semanticEqualShape(b)
}

func EqualShape(a, b Shape) bool {
	if !SameShape(a, b) {
		return SemanticEqualShape(a, b)
	}
	return true
}

// === Sentinel operations ===

func isSpecifiedShape(s Shape) bool {
	return s != nil && s != shapeUnspecified
}

func isUnspecifiedShape(s Shape) bool {
	return !isSpecifiedShape(s)
}

func takeOrElseShape(s, def Shape) Shape {
	if s == nil || s == shapeUnspecified {
		return def
	}
	return s
}

func coalesceShape(ptr, def Shape) Shape {
	if ptr == nil {
		return def
	}
	return ptr
}

func mergeShape(a, b Shape) Shape {
	return MergeShape(a, b)
}

func sameShape(a, b Shape) bool {
	return SameShape(a, b)
}

func semanticEqualShape(a, b Shape) bool {
	return SemanticEqualShape(a, b)
}

func equalShape(a, b Shape) bool {
	return EqualShape(a, b)
}

func stringShape(s Shape) string {
	if !isSpecifiedShape(s) {
		return "Shape{Unspecified}"
	}
	return s.stringShape()
}

// ShapeCompanion is the single exported sentinel-family symbol for Shape.
var ShapeCompanion = companions.NewStructCompanion(
	func() Shape { return shapeUnspecified },
	stringShape,
	isSpecifiedShape,
	isUnspecifiedShape,
	takeOrElseShape,
	mergeShape,
	sameShape,
	semanticEqualShape,
	equalShape,
	coalesceShape,
	copyShape,
)

func CopyShape(s Shape, options ...ShapeOption) Shape {
	switch shape := s.(type) {
	case *circleShape:
		copy := *shape
		for _, option := range options {
			option(&copy)
		}
		return &copy
	case *rectangleShape:
		copy := *shape
		for _, option := range options {
			option(&copy)
		}
		return &copy
	case *RoundedCornerShape:
		copy := *shape
		for _, option := range options {
			option(&copy)
		}
		return &copy
	case *CutCornerShape:
		copy := *shape
		for _, option := range options {
			option(&copy)
		}
		return &copy

	default:
		panic(fmt.Sprintf("CopyShape: unknown shape type %s", s.stringShape()))
	}
}

func StringShape(shape Shape) string {
	return shape.stringShape()
}

// copyShape copies s and applies the given ShapeOptions on top; the source is
// never mutated. It is the companion's Copy operation.
func copyShape(s Shape, options ...func(Shape)) Shape {
	opts := make([]ShapeOption, len(options))
	for i, option := range options {
		opts[i] = ShapeOption(option)
	}
	return CopyShape(s, opts...)
}
