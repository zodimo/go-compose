package text

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/text/unit"
)

var DefaultColor = graphics.ColorBlack

var DefaultFontSize = unit.Sp(14)
var DefaultLetterSpacing = unit.Sp(0)
var DefaultBackgroundColor = graphics.ColorTransparent
var DefaultTextDecoration = style.TextDecorationNone
var DefaultShadow = graphics.ShadowNone
var DefaultFontWeight = font.FontWeightNormal
var DefaultFontStyle = font.FontStyleNormal
var DefaultFontFamily = font.FontFamilyDefault

var DefaultTextAlign = style.TextAlignStart
var DefaultLineHeight = unit.TextUnitUnspecified // TODO: Should this stay unspecified ?
var DefaultLineBreak = style.LineBreakParagraph
