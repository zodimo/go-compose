package material3

import (
	"git.sr.ht/~schnwalter/gio-mw/defaults/schemes"
	"github.com/zodimo/go-compose/compose"
)

type ThemeInterface interface {
	ColorScheme() *ColorScheme
	Typography() *Typography
	// Shapes() *Shapes
	// MotionScheme() *MotionScheme
}

func Theme(c compose.Composer) ThemeInterface {
	return themeImpl{
		composer: c,
	}
}

type themeImpl struct {
	composer compose.Composer
}

func (t themeImpl) ColorScheme() *ColorScheme {
	return LocalColorScheme.Current(t.composer)
}

func (t themeImpl) Typography() *Typography {
	return LocalTypography.Current(t.composer)
}

func (t themeImpl) Shapes() *Shapes {
	return LocalShapes.Current(t.composer)
}

func (t themeImpl) MotionScheme() *MotionScheme {
	return LocalMotionScheme.Current(t.composer)
}

var LocalColorScheme = compose.CompositionLocalOf(func() *ColorScheme {
	return ColorSchemeFromTokens(schemes.SchemeBaselineLight())
})

var LocalTypography = compose.CompositionLocalOf(func() *Typography {
	return nil
})

var LocalShapes = compose.CompositionLocalOf(func() *Shapes {
	return nil
})

var LocalMotionScheme = compose.CompositionLocalOf(func() *MotionScheme {
	return &DefaultMotionScheme
})
