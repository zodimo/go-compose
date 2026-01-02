package button

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
)

// ButtonShapesUnspecified is the sentinel value for unspecified ButtonShapes
var ButtonShapesUnspecified = &ButtonShapes{
	Shape:        shape.ShapeUnspecified,
	PressedShape: shape.ShapeUnspecified,
}

// ButtonShapes holds shape configuration for button states
type ButtonShapes struct {
	Shape        shape.Shape
	PressedShape shape.Shape
}

// IsSpecifiedButtonShapes returns true if s is specified (not nil and not the sentinel)
func IsSpecifiedButtonShapes(s *ButtonShapes) bool {
	return s != nil && s != ButtonShapesUnspecified
}

// TakeOrElseButtonShapes returns s if specified, otherwise returns defaultShapes
func TakeOrElseButtonShapes(s, defaultShapes *ButtonShapes) *ButtonShapes {
	if s == nil || s == ButtonShapesUnspecified {
		return defaultShapes
	}
	return s
}

// MergeButtonShapes merges two ButtonShapes, preferring b's specified values over a's
func MergeButtonShapes(a, b *ButtonShapes) *ButtonShapes {
	a = CoalesceButtonShapes(a, ButtonShapesUnspecified)
	b = CoalesceButtonShapes(b, ButtonShapesUnspecified)

	if a == ButtonShapesUnspecified {
		return b
	}
	if b == ButtonShapesUnspecified {
		return a
	}

	// Both are custom: allocate new merged style
	return &ButtonShapes{
		Shape:        shape.TakeOrElseShape(b.Shape, a.Shape),
		PressedShape: shape.TakeOrElseShape(b.PressedShape, a.PressedShape),
	}
}

// StringButtonShapes returns a string representation of ButtonShapes
func StringButtonShapes(s *ButtonShapes) string {
	if !IsSpecifiedButtonShapes(s) {
		return "ButtonShapes{Unspecified}"
	}

	return fmt.Sprintf(
		"ButtonShapes{Shape: %s, PressedShape: %s}",
		shape.StringShape(s.Shape),
		shape.StringShape(s.PressedShape),
	)
}

// CoalesceButtonShapes returns ptr if not nil, otherwise returns def
func CoalesceButtonShapes(ptr, def *ButtonShapes) *ButtonShapes {
	if ptr == nil {
		return def
	}
	return ptr
}

// SameButtonShapes returns true if a and b are the same pointer or both unspecified
func SameButtonShapes(a, b *ButtonShapes) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == ButtonShapesUnspecified
	}
	if b == nil {
		return a == ButtonShapesUnspecified
	}
	return a == b
}

// SemanticEqualButtonShapes checks field-by-field equality
func SemanticEqualButtonShapes(a, b *ButtonShapes) bool {
	a = CoalesceButtonShapes(a, ButtonShapesUnspecified)
	b = CoalesceButtonShapes(b, ButtonShapesUnspecified)

	return shape.EqualShape(a.Shape, b.Shape) &&
		shape.EqualShape(a.PressedShape, b.PressedShape)
}

// EqualButtonShapes returns true if a and b are semantically equal
func EqualButtonShapes(a, b *ButtonShapes) bool {
	if SameButtonShapes(a, b) {
		return true
	}
	return SemanticEqualButtonShapes(a, b)
}

// ButtonShapesOption is a functional option for CopyButtonShapes
type ButtonShapesOption func(*ButtonShapes)

// WithButtonShape sets the shape option
func WithButtonShape(s shape.Shape) ButtonShapesOption {
	return func(o *ButtonShapes) {
		o.Shape = s
	}
}

// WithButtonPressedShape sets the pressed shape option
func WithButtonPressedShape(s shape.Shape) ButtonShapesOption {
	return func(o *ButtonShapes) {
		o.PressedShape = s
	}
}

// CopyButtonShapes creates a copy with optional modifications
func CopyButtonShapes(s *ButtonShapes, options ...ButtonShapesOption) *ButtonShapes {
	opt := *ButtonShapesUnspecified

	for _, option := range options {
		option(&opt)
	}

	s = CoalesceButtonShapes(s, ButtonShapesUnspecified)

	return &ButtonShapes{
		Shape:        shape.TakeOrElseShape(opt.Shape, s.Shape),
		PressedShape: shape.TakeOrElseShape(opt.PressedShape, s.PressedShape),
	}
}
