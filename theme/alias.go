package theme

import (
	"image/color"

	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
	"github.com/zodimo/go-compose/theme/colorrole"

	"git.sr.ht/~schnwalter/gio-mw/token"
)

type ColorRole = colorrole.ColorRole

// BasicTheme is the framework-owned representation of a basic Material theme.
// It holds the palette colors used for color role resolution.
type BasicTheme struct {
	Bg         color.NRGBA
	Fg         color.NRGBA
	ContrastBg color.NRGBA
	ContrastFg color.NRGBA
}

// Theme is the framework-owned representation of a Material3 token theme.
// It wraps the color scheme from the gio-mw token package.
type Theme struct {
	Scheme *token.Scheme
}

type TokenColor = token.MatColor

type OpacityLevel = token.OpacityLevel

var colorLerp = lerp.LerpColor
