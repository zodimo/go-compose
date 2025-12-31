package text

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

type TextStyleInterface interface {
	Alpha() float32
	Background() graphics.Color
	Color() graphics.Color
	FontFamily() font.FontFamily
	FontFeatureSettings() string
	FontSize() unit.TextUnit
	FontStyle() font.FontStyle
	FontSynthesis() *font.FontSynthesis
	FontWeight() font.FontWeight
	LetterSpacing() unit.TextUnit
	LineBreak() style.LineBreak
	LineHeight() unit.TextUnit
	TextAlign() style.TextAlign
	TextDecoration() *style.TextDecoration
	TextDirection() style.TextDirection
	ToString() string
}
