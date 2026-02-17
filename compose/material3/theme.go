package material3

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/ui/graphics"
)

type ThemeInterface interface {
	ColorScheme() *ColorSchemeNext
	Typography() *Typography
	Shapes() *Shapes
	MotionScheme() *MotionScheme
}

func Theme(c compose.Composer) ThemeInterface {
	return themeImpl{
		composer: c,
	}
}

type themeImpl struct {
	composer compose.Composer
}

func (t themeImpl) ColorScheme() *ColorSchemeNext {
	return LocalColorSchemeNext.Current(t.composer)
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

func (t themeImpl) ContentColorFor(backgroundColor graphics.Color) graphics.Color {
	return t.ColorScheme().ContentFor(backgroundColor)
}

var LocalShapes = compose.CompositionLocalOf(func() *Shapes {
	return DefaultShapes
})

var LocalMotionScheme = compose.CompositionLocalOf(func() *MotionScheme {
	return DefaultMotionScheme
})
