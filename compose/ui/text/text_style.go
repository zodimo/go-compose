package text

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

var _ TextStyleInterface = (*TextStyle)(nil)

type TextStyle struct {
	font       font.Font
	color      graphics.Color
	background graphics.Color
	lineHeight unit.TextUnit
	lineBreak  style.LineBreak
}

func (ts TextStyle) Alpha() float32 {
	return ts.color.Alpha()
}
func (ts TextStyle) Background() graphics.Color {
	return ts.background
}
func (ts TextStyle) Color() graphics.Color {
	return ts.color
}
func (ts TextStyle) FontFamily() font.FontFamily {
	panic("FontFamily not implemented")
}
func (ts TextStyle) FontFeatureSettings() string {
	panic("FontFeatureSettings not implemented")
}
func (ts TextStyle) FontSize() unit.TextUnit {
	panic("FontSizenot implemented")
}
func (ts TextStyle) FontStyle() font.FontStyle {
	panic("FontStyle not implemented")
}
func (ts TextStyle) FontSynthesis() *font.FontSynthesis {
	panic("FontSynthesis not implemented")
}
func (ts TextStyle) FontWeight() font.FontWeight {
	panic("FontWeightnot implemented")
}
func (ts TextStyle) LetterSpacing() unit.TextUnit {
	panic("LetterSpacing not implemented")
}

func (ts TextStyle) LineBreak() style.LineBreak {
	return ts.lineBreak
}
func (ts TextStyle) LineHeight() unit.TextUnit {
	return ts.lineHeight
}
func (ts TextStyle) TextAlign() style.TextAlign {
	panic("TextAlign not implemented")
}
func (ts TextStyle) TextDecoration() *style.TextDecoration {
	panic("TextDecoration not implemented")
}
func (ts TextStyle) TextDirection() style.TextDirection {
	panic("TextDirection not implemented")
}
func (ts TextStyle) ToString() string {
	panic("ToString not implemented")
}

func TextStyleResolveDefaults(ts *TextStyle, direction unit.LayoutDirection) *TextStyle {
	return &TextStyle{}
}
