package style

import gioText "gioui.org/text"

func TextAlignToGioTextAlignment(t TextAlign) gioText.Alignment {
	align := t.TakeOrElse(TextAlignStart)
	return gioText.Alignment(align)
}

func LineBreakToGioWrapPolicy(l LineBreak) gioText.WrapPolicy {
	linebreak := l.TakeOrElse(LineBreakParagraph)
	return gioText.WrapPolicy(linebreak)
}
