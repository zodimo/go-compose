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
	FontSize() unit.TextUnit
	FontStyle() font.FontStyle
	FontWeight() font.FontWeight
	LetterSpacing() unit.TextUnit
	LineBreak() style.LineBreak
	LineHeight() unit.TextUnit
	TextAlign() style.TextAlign
	TextDecoration() *style.TextDecoration
	TextDirection() style.TextDirection
}

type SpanStyleInterface interface {
	Color() graphics.Color
	FontSize() unit.TextUnit
	FontWeight() font.FontWeight
	FontStyle() font.FontStyle
	FontFamily() font.FontFamily
	LetterSpacing() unit.TextUnit
	Background() graphics.Color
	TextDecoration() *style.TextDecoration
	Shadow() *graphics.Shadow
}

type ParaghStyleInterface interface {
	TextAlign() style.TextAlign
	TextDirection() style.TextDirection
	LineHeight() unit.TextUnit
	LineBreak() style.LineBreak
}
