package material3

import (
	"gioui.org/layout"
	"git.sr.ht/~schnwalter/gio-mw/token"
	"git.sr.ht/~schnwalter/gio-mw/wdk"
)

// Helpers to auhment the token.Theme use with gio-mw widgets

type TokenColorSchemeOptions = func(scheme *token.Scheme)

func WithColorSchemeOptions(options ...TokenColorSchemeOptions) TokenColorSchemeOptions {
	return func(scheme *token.Scheme) {
		for _, option := range options {
			option(scheme)
		}
	}
}

func UpdateTokenTheme(gtx layout.Context, schemeOptions []TokenColorSchemeOptions) {
	theme := *wdk.GetMaterialTheme(gtx)
	scheme := *theme.Scheme
	for _, option := range schemeOptions {
		option(&scheme)
	}
	theme.Scheme = &scheme
	wdk.InitMaterialThemeInContext(gtx, &theme)
}
