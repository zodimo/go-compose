package theme

import (
	"image/color"
)

type ResolvableColor interface {
	ResolveColorDescriptor(colorDesc ThemeColorDescriptor) ThemeColor
}

// ResolvedColor is a color that has been resolved from a ThemeColorDescriptor

type ThemeColor interface {
	AsHex() string
	AsNRGBA() color.NRGBA
	SetOpacity(opacity OpacityLevel) ThemeColor
	AsTokenColor() TokenColor
}

type ThemeColorDescriptor struct {
	color     color.NRGBA
	colorRole ColorRole
	isColor   bool
	updates   []func(color TokenColor) TokenColor
}

func (t ThemeColorDescriptor) AppendUpdate(update func(color TokenColor) TokenColor) ThemeColorDescriptor {
	t = t.initUpdates()

	t.updates = append(t.updates, update)
	return t
}

func (t ThemeColorDescriptor) SetOpacity(opacity OpacityLevel) ThemeColorDescriptor {

	update := func(c TokenColor) TokenColor {
		return c.SetOpacity(opacity)
	}
	return t.AppendUpdate(update)
}

func (t ThemeColorDescriptor) initUpdates() ThemeColorDescriptor {
	if t.updates == nil {
		t.updates = make([]func(color TokenColor) TokenColor, 0)
	}
	return t
}
