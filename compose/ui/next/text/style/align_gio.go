package style

import (
	gioText "gioui.org/text"
)

// textAlignToGioAlign converts compose TextAlign to gio text.Alignment.
func TextAlignToGioAlign(textAlign TextAlign) gioText.Alignment {
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
