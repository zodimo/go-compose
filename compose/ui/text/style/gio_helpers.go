package style

import (
	gioText "gioui.org/text"
)

func FromGioTextAlign(gioTextAlign gioText.Alignment) TextAlign {
	return TextAlign(gioTextAlign)
}

// textAlignToGioAlign converts compose TextAlign to gio text.Alignment.
func TextAlignToGioTextAlignment(textAlign TextAlign) gioText.Alignment {
	switch textAlign {
	case TextAlignLeft:
		return gioText.Start
	case TextAlignCenter:
		return gioText.Middle
	case TextAlignRight:
		return gioText.End
	case TextAlignJustify:
		// Gio doesn't support justify, fall back to Start
		return gioText.Start
	default:
		return gioText.Start
	}
}

func LineBreakToGioWrapPolicy(l LineBreak) gioText.WrapPolicy {
	linebreak := LineBreakCompanion.TakeOrElse(l, LineBreakParagraph)
	return gioText.WrapPolicy(linebreak)
}

func GioWrapPolicyToLineBreak(gioWrapPolicy gioText.WrapPolicy) LineBreak {
	return LineBreak(gioWrapPolicy)
}
