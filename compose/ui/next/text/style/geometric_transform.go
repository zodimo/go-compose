package style

import (
	"fmt"

	"github.com/zodimo/go-compose/pkg/floatutils"
	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// textGeometricTransformUnspecified is the sentinel for TextGeometricTransform. It
// is deliberately package-private: the only public way to obtain it is
// TextGeometricTransformCompanion.Unspecified().
var textGeometricTransformUnspecified = &TextGeometricTransform{
	ScaleX: floatutils.Float32Unspecified,
	SkewX:  floatutils.Float32Unspecified,
}

var TextGeometricTransformNone = &TextGeometricTransform{
	ScaleX: 1,
	SkewX:  0,
}

// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/next/text/style/TextGeometricTransform.kt;drc=4970f6e96cdb06089723da0ab8ec93ae3f067c7a;l=33

type TextGeometricTransform struct {
	ScaleX float32
	SkewX  float32
}

func LerpGeometricTransform(
	start *TextGeometricTransform,
	stop *TextGeometricTransform,
	fraction float32,
) *TextGeometricTransform {
	if start == nil || stop == nil {
		panic("TextGeometricTransform must be specified")
	}
	if fraction == 0 {
		return start
	}
	if fraction == 1 {
		return stop
	}
	return &TextGeometricTransform{
		ScaleX: lerp.Between32(start.ScaleX, stop.ScaleX, fraction),
		SkewX:  lerp.Between32(start.SkewX, stop.SkewX, fraction),
	}
}

func StringTextGeometricTransform(gt *TextGeometricTransform) string {
	if !isSpecifiedTextGeometricTransform(gt) {
		return "TextGeometricTransformUnspecified"
	}
	return fmt.Sprintf("TextGeometricTransform(scaleX=%f, skewX=%f)", gt.ScaleX, gt.SkewX)
}

func IsSpecifiedTextGeometricTransform(gt *TextGeometricTransform) bool {
	return isSpecifiedTextGeometricTransform(gt)
}

func TakeOrElseTextGeometricTransform(gt *TextGeometricTransform, defaultGT *TextGeometricTransform) *TextGeometricTransform {
	return takeOrElseTextGeometricTransform(gt, defaultGT)
}

// Identity (2 ns)
func SameTextGeometricTransform(a, b *TextGeometricTransform) bool {
	return sameTextGeometricTransform(a, b)
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualTextGeometricTransform(a, b *TextGeometricTransform) bool {
	return semanticEqualTextGeometricTransform(a, b)
}

func EqualTextGeometricTransform(a, b *TextGeometricTransform) bool {
	return equalTextGeometricTransform(a, b)
}

func MergeTextGeometricTransform(a, b *TextGeometricTransform) *TextGeometricTransform {
	return mergeTextGeometricTransform(a, b)
}

func CoalesceTextGeometricTransform(ptr, def *TextGeometricTransform) *TextGeometricTransform {
	return coalesceTextGeometricTransform(ptr, def)
}

// === Sentinel operations ===

func isSpecifiedTextGeometricTransform(gt *TextGeometricTransform) bool {
	return gt != nil && gt != textGeometricTransformUnspecified
}

func isUnspecifiedTextGeometricTransform(gt *TextGeometricTransform) bool {
	return !isSpecifiedTextGeometricTransform(gt)
}

func takeOrElseTextGeometricTransform(gt, defaultGT *TextGeometricTransform) *TextGeometricTransform {
	if !isSpecifiedTextGeometricTransform(gt) {
		return defaultGT
	}
	return gt
}

func coalesceTextGeometricTransform(ptr, def *TextGeometricTransform) *TextGeometricTransform {
	if ptr == nil {
		return def
	}
	return ptr
}

func mergeTextGeometricTransform(a, b *TextGeometricTransform) *TextGeometricTransform {
	a = coalesceTextGeometricTransform(a, textGeometricTransformUnspecified)
	b = coalesceTextGeometricTransform(b, textGeometricTransformUnspecified)

	if a == textGeometricTransformUnspecified {
		return b
	}
	if b == textGeometricTransformUnspecified {
		return a
	}

	// Both are custom: allocate new merged style. b wins where specified.
	return &TextGeometricTransform{
		ScaleX: floatutils.TakeOrElse(b.ScaleX, a.ScaleX),
		SkewX:  floatutils.TakeOrElse(b.SkewX, a.SkewX),
	}
}

func stringTextGeometricTransform(gt *TextGeometricTransform) string {
	return StringTextGeometricTransform(gt)
}

func sameTextGeometricTransform(a, b *TextGeometricTransform) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == textGeometricTransformUnspecified
	}
	if b == nil {
		return a == textGeometricTransformUnspecified
	}
	return a == b
}

func semanticEqualTextGeometricTransform(a, b *TextGeometricTransform) bool {

	a = coalesceTextGeometricTransform(a, textGeometricTransformUnspecified)
	b = coalesceTextGeometricTransform(b, textGeometricTransformUnspecified)

	if a == textGeometricTransformUnspecified && b == textGeometricTransformUnspecified {
		return true
	}

	return floatutils.Float32Equals(a.ScaleX, b.ScaleX, floatutils.Float32EqualityThreshold) &&
		floatutils.Float32Equals(a.SkewX, b.SkewX, floatutils.Float32EqualityThreshold)
}

func equalTextGeometricTransform(a, b *TextGeometricTransform) bool {
	if !sameTextGeometricTransform(a, b) {
		return semanticEqualTextGeometricTransform(a, b)
	}
	return true
}

// copyTextGeometricTransform returns a copy of gt; the source is never mutated.
func copyTextGeometricTransform(gt *TextGeometricTransform, options ...func(*TextGeometricTransform)) *TextGeometricTransform {
	cp := &TextGeometricTransform{}
	if isSpecifiedTextGeometricTransform(gt) {
		cp.ScaleX = gt.ScaleX
		cp.SkewX = gt.SkewX
	}
	for _, option := range options {
		option(cp)
	}
	return cp
}

// TextGeometricTransformCompanion is the single exported sentinel-family symbol for TextGeometricTransform.
var TextGeometricTransformCompanion = companions.NewStructCompanion(
	func() *TextGeometricTransform { return textGeometricTransformUnspecified },
	stringTextGeometricTransform,
	isSpecifiedTextGeometricTransform,
	isUnspecifiedTextGeometricTransform,
	takeOrElseTextGeometricTransform,
	mergeTextGeometricTransform,
	sameTextGeometricTransform,
	semanticEqualTextGeometricTransform,
	equalTextGeometricTransform,
	coalesceTextGeometricTransform,
	copyTextGeometricTransform,
)
