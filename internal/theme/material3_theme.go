package theme

import (
	"gioui.org/layout"
	"git.sr.ht/~schnwalter/gio-mw/defaults"
	"git.sr.ht/~schnwalter/gio-mw/defaults/schemes"
	"git.sr.ht/~schnwalter/gio-mw/token"
)

func defaultMaterial3Theme(gtx layout.Context) *token.Theme {
	scheme := schemes.SchemeBaselineLight()
	return defaults.NewTheme(gtx, scheme)
}
