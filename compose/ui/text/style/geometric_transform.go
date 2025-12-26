package style

import (
	"fmt"

	"github.com/zodimo/go-compose/pkg/floatutils"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
)

var TextGeometricTransformUnspecified = &TextGeometricTransform{
	ScaleX: floatutils.Float32Unspecified,
	SkewX:  floatutils.Float32Unspecified,
}

var TextGeometricTransformNone = &TextGeometricTransform{
	ScaleX: 1,
	SkewX:  0,
}

// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/text/style/TextGeometricTransform.kt;drc=4970f6e96cdb06089723da0ab8ec93ae3f067c7a;l=33

type TextGeometricTransform struct {
	ScaleX float32
	SkewX  float32
}

func LerpGeometricTransform(
	start *TextGeometricTransform,
	stop *TextGeometricTransform,
	fraction float32,
) TextGeometricTransform {
	if start == nil || stop == nil {
		panic("TextGeometricTransform must be specified")
	}
	if fraction == 0 {
		return *start
	}
	if fraction == 1 {
		return *stop
	}
	return TextGeometricTransform{
		ScaleX: lerp.Between32(start.ScaleX, stop.ScaleX, fraction),
		SkewX:  lerp.Between32(start.SkewX, stop.SkewX, fraction),
	}
}

func StringTextGeometricTransform(gt *TextGeometricTransform) string {
	if !IsSpecifiedTextGeometricTransform(gt) {
		return "TextGeometricTransformUnspecified"
	}
	return fmt.Sprintf("TextGeometricTransform(scaleX=%f, skewX=%f)", gt.ScaleX, gt.SkewX)
}

func IsSpecifiedTextGeometricTransform(gt *TextGeometricTransform) bool {
	return gt != nil && gt != TextGeometricTransformUnspecified
}

func TakeOrElseTextGeometricTransform(gt *TextGeometricTransform, defaultGT *TextGeometricTransform) *TextGeometricTransform {
	if !IsSpecifiedTextGeometricTransform(gt) {
		return defaultGT
	}
	return gt
}

// Identity (2 ns)
func SameTextGeometricTransform(a, b *TextGeometricTransform) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == TextGeometricTransformUnspecified
	}
	if b == nil {
		return a == TextGeometricTransformUnspecified
	}
	return a == b
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualTextGeometricTransform(a, b *TextGeometricTransform) bool {

	a = CoalesceTextGeometricTransform(a, TextGeometricTransformUnspecified)
	b = CoalesceTextGeometricTransform(b, TextGeometricTransformUnspecified)

	return a.ScaleX == b.ScaleX &&
		a.SkewX == b.SkewX
}

func EqualTextGeometricTransform(a, b *TextGeometricTransform) bool {
	if !SameTextGeometricTransform(a, b) {
		return SemanticEqualTextGeometricTransform(a, b)
	}
	return true
}

func MergeTextGeometricTransform(a, b *TextGeometricTransform) *TextGeometricTransform {
	a = CoalesceTextGeometricTransform(a, TextGeometricTransformUnspecified)
	b = CoalesceTextGeometricTransform(b, TextGeometricTransformUnspecified)

	if a == TextGeometricTransformUnspecified {
		return b
	}
	if b == TextGeometricTransformUnspecified {
		return a
	}

	// Both are custom: allocate new merged style
	return &TextGeometricTransform{
		ScaleX: floatutils.TakeOrElse(b.ScaleX, a.ScaleX),
		SkewX:  floatutils.TakeOrElse(b.SkewX, a.SkewX),
	}
}

func CoalesceTextGeometricTransform(ptr, def *TextGeometricTransform) *TextGeometricTransform {
	if ptr == nil {
		return def
	}
	return ptr
}
