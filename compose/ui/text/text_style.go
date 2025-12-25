package text

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

type TextStyle struct {
	Color      graphics.Color
	FontSize   uiTextUnit //TextUnit.Unspecified,
	FontWeight font.FontWeight
	FontStyle  font.FontStyle
	// FontSynthesis          any
	FontFamily font.FontFamily
	// FontFeatureSettings    any
	LetterSpacing          uiTextUnit
	BaselineShift          style.BaselineShift
	TextGeometricTransform *uiTextGeometricTransform
	LocaleList             *uiLocaleList
	Background             graphics.Color
	TextDecoration         *uiTextDecoration
	Shadow                 *uiShadow
	TextAlign              style.TextAlign
	TextDirection          uiDirection
	LineHeight             uiTextUnit
	TextIndent             *uiIndent
}

func TextStyleDefaults() TextStyle {
	return TextStyle{
		Color:                  graphics.ColorUnspecified,
		FontSize:               unit.TextUnitUnspecified,
		FontWeight:             font.FontWeightUnspecified,
		FontStyle:              font.FontStyleUnspecified,
		FontFamily:             nil,
		LetterSpacing:          unit.TextUnitUnspecified,
		BaselineShift:          style.BaselineShiftUnspecified,
		TextGeometricTransform: nil,
		LocaleList:             nil,
		Background:             graphics.ColorUnspecified,
		TextDecoration:         nil,
		Shadow:                 nil,
		TextAlign:              style.TextAlignUnspecified,
		TextDirection:          style.TextDirectionUnspecified,
		LineHeight:             unit.TextUnitUnspecified,
		TextIndent:             nil,
	}
}
