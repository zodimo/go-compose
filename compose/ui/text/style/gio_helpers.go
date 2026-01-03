package style

import (
	gioText "gioui.org/text"
)

func FromGioTextAlign(gioTextAlign gioText.Alignment) TextAlign {
	return TextAlign(gioTextAlign)
}

func TextAlignToGioTextAlignment(t TextAlign) gioText.Alignment {
	align := t.TakeOrElse(TextAlignStart)
	return gioText.Alignment(align)
}

func LineBreakToGioWrapPolicy(l LineBreak) gioText.WrapPolicy {
	linebreak := l.TakeOrElse(LineBreakParagraph)
	return gioText.WrapPolicy(linebreak)
}

func GioWrapPolicyToLineBreak(gioWrapPolicy gioText.WrapPolicy) LineBreak {
	return LineBreak(gioWrapPolicy)
}
