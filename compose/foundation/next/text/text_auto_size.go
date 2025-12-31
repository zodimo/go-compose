package text

import (
	"math"

	"github.com/zodimo/go-compose/compose/ui/next/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// TextAutoSizeDefaults contains defaults for TextAutoSize APIs.
var TextAutoSizeDefaults = TextAutoSizeValues{
	MinFontSize: uiTextUnitSp(12),
	MaxFontSize: uiTextUnitSp(112),
	StepSize:    uiTextUnitSp(0.25),
}

type TextAutoSizeValues struct {
	MinFontSize uiTextUnit
	MaxFontSize uiTextUnit
	StepSize    uiTextUnit
}

// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/TextAutoSize.kt;drc=4970f6e96cdb06089723da0ab8ec93ae3f067c7a;l=35
type TextAutoSize interface {
	GetFontSize(scope TextAutoSizeLayoutScope, constraints uiConstraints, text uiAnnotatedString) uiTextUnit
}

type TextAutoSizeLayoutScope interface {
	PerformLayout(constraints uiConstraints, text uiAnnotatedString, fontSize uiTextUnit) uiTextLayoutResult
}

// NewStepBasedTextAutoSize automatically sizes the text with the biggest font size that fits the available space.
func NewStepBasedTextAutoSize(
	minFontSize uiTextUnit,
	maxFontSize uiTextUnit,
	stepSize uiTextUnit,
) TextAutoSize {
	if minFontSize.IsUnspecified() {
		minFontSize = TextAutoSizeDefaults.MinFontSize
	}
	if maxFontSize.IsUnspecified() {
		maxFontSize = TextAutoSizeDefaults.MaxFontSize
	}
	if stepSize.IsUnspecified() {
		stepSize = TextAutoSizeDefaults.StepSize
	}
	s := &stepBasedTextAutoSize{
		minFontSize: minFontSize,
		maxFontSize: maxFontSize,
		stepSize:    stepSize,
	}
	s.init()
	return s
}

type stepBasedTextAutoSize struct {
	minFontSize uiTextUnit
	maxFontSize uiTextUnit
	stepSize    uiTextUnit
}

func (s *stepBasedTextAutoSize) init() {
	// Checks for validity of AutoSize instance
	// In Go we handle this in the constructor or lazily?
	// The Kotlin init block throws exceptions.
	// We'll panic for now to match strictness or just log/ignore?
	// Kotlin throws IllegalArgumentException.
	if s.minFontSize.IsUnspecified() {
		panic("AutoSize.StepBased: TextUnit.Unspecified is not a valid value for minFontSize. Try using other values e.g. 10.sp")
	}
	if s.maxFontSize.IsUnspecified() {
		panic("AutoSize.StepBased: TextUnit.Unspecified is not a valid value for maxFontSize. Try using other values e.g. 100.sp")
	}
	if s.stepSize.IsUnspecified() {
		panic("AutoSize.StepBased: TextUnit.Unspecified is not a valid value for stepSize. Try using other values e.g. 0.25.sp")
	}

	if s.minFontSize.Type() == s.maxFontSize.Type() && s.minFontSize.Value() > s.maxFontSize.Value() {
		s.minFontSize = s.maxFontSize
	}

	// check if stepSize is too small
	if s.stepSize.Type() == unit.TextUnitTypeSp && s.stepSize.Value() < 0.0001 {
		panic("AutoSize.StepBased: stepSize must be greater than or equal to 0.0001f.sp")
	}

	if s.minFontSize.Value() < 0 {
		panic("AutoSize.StepBased: minFontSize must not be negative")
	}
	if s.maxFontSize.Value() < 0 {
		panic("AutoSize.StepBased: maxFontSize must not be negative")
	}
}

func (s *stepBasedTextAutoSize) GetFontSize(scope TextAutoSizeLayoutScope, constraints uiConstraints, text uiAnnotatedString) uiTextUnit {
	stepSize := float64(s.stepSize.Value())
	smallest := float64(s.minFontSize.Value())
	largest := float64(s.maxFontSize.Value())
	min := smallest
	max := largest

	current := (min + max) / 2

	for (max - min) >= stepSize {
		layoutResult := scope.PerformLayout(constraints, text, uiTextUnitSp(float32(current)))
		if didOverflow(layoutResult) {
			max = current
		} else {
			min = current
		}
		current = (min + max) / 2
	}

	// used size minus minFontSize must be divisible by stepSize
	current = (math.Floor((min-smallest)/stepSize)*stepSize + smallest)

	// We have found a size that fits, but we can still try one step up
	if (current + stepSize) <= largest {
		layoutResult := scope.PerformLayout(constraints, text, uiTextUnitSp(float32(current+stepSize)))
		if !didOverflow(layoutResult) {
			current += stepSize
		}
	}

	return uiTextUnitSp(float32(current))
}

func didOverflow(result uiTextLayoutResult) bool {
	switch result.LayoutInput().Overflow {
	// case style.OverFlowClip, style.OverFlowVisible:
	// 	return result.DidOverflowBounds()
	case style.OverFlowStartEllipsis, style.OverFlowMiddleEllipsis, style.OverFlowEllipsis:
		return didOverflowByEllipsize(result)
	default:
		panic("TextOverflow type not supported")
	}
}

func didOverflowByEllipsize(result uiTextLayoutResult) bool {
	lineCount := result.LineCount()
	if lineCount == 0 {
		return false
	}
	// Text only gets start- or middle-ellipsized if it is single line, so we can check if
	// the first line is ellipsized to cover all single-line ellipsis overflow
	if lineCount == 1 {
		return result.IsLineEllipsized(0)
	}

	switch result.LayoutInput().Overflow {
	// If the text is not single line but start or middle ellipsis has been set, fall
	// back to the behavior for TextOverflow.Clip
	// case style.OverFlowStartEllipsis, style.OverFlowMiddleEllipsis:
	// 	return result.DidOverflowBounds()
	// TextOverflow.Ellipsis is supported for multiline text and happens at the end of
	// the text, so we only need to check the last line
	case style.OverFlowEllipsis:
		return result.IsLineEllipsized(lineCount - 1)
	default:
		return false
	}
}

func (s *stepBasedTextAutoSize) Equals(other interface{}) bool {
	if otherS, ok := other.(*stepBasedTextAutoSize); ok {
		return s.minFontSize.Equals(otherS.minFontSize) &&
			s.maxFontSize.Equals(otherS.maxFontSize) &&
			s.stepSize.Equals(otherS.stepSize)
	}
	return false
}
