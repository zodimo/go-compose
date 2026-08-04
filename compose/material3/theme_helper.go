package material3

import (
	"github.com/zodimo/go-compose/internal/layoutnode"

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

func UpdateTokenTheme(gtx layoutnode.LayoutContext, schemeOptions []TokenColorSchemeOptions) {
	// Convert framework LayoutContext to engine layout.Context at the seam.
	gioCtx := gtx.ToGio()
	theme := *wdk.GetMaterialTheme(*gioCtx)
	scheme := *theme.Scheme
	for _, option := range schemeOptions {
		option(&scheme)
	}
	theme.Scheme = &scheme
	wdk.InitMaterialThemeInContext(*gioCtx, &theme)
}
