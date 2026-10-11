package text

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/next/text/style"
	"github.com/zodimo/go-compose/compose/ui/next/text/unit"
)

var DefaultColor = graphics.ColorBlack
var DefaultColorForegroundStyle = style.TextForegroundStyleFromColor(DefaultColor)

var DefaultFontSize = unit.Sp(14)
var DefaultLetterSpacing = unit.Sp(0)
var DefaultBackgroundColor = graphics.ColorTransparent

var DefaultTextAlign = style.TextAlignStart
var DefaultLineHeight = unit.TextUnitCompanion.Unspecified() // TODO: Should this stay unspecified ?
var DefaultLineBreak = style.LineBreakParagraph
var DefaultTextIndent = style.TextIndentNone
var DefaultHyphens = style.HyphensNone
var DefaultTextMotion = style.TextMotionStatic
