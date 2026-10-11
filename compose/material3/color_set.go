package material3

import (
	"git.sr.ht/~schnwalter/gio-mw/token"
	"github.com/zodimo/go-compose/compose/ui/graphics"
)

var ColorSetUnspecified = &ColorSet{
	Color:   ColorUnspecified,
	OnColor: ColorUnspecified,
}

type ColorSet struct {
	Color   Color
	OnColor Color
}

func (c *ColorSet) ToTokens() token.MatColorSet {
	return token.MatColorSet{
		Color:   ColorToToken(c.Color),
		OnColor: ColorToToken(c.OnColor),
	}
}

func CoalesceColorSet(ptr, def *ColorSet) *ColorSet {
	if ptr == nil {
		return def
	}
	return ptr
}

func IsSpecifiedColorSet(s *ColorSet) bool {
	return s != nil && s != ColorSetUnspecified
}

func TakeOrElseColorSet(s, def *ColorSet) *ColorSet {
	if !IsSpecifiedColorSet(s) {
		return def
	}
	return s
}

// Identity (2 ns)
func SameColorSet(a, b *ColorSet) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == ColorSetUnspecified
	}
	if b == nil {
		return a == ColorSetUnspecified
	}
	return a == b
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualColorSet(a, b *ColorSet) bool {

	a = CoalesceColorSet(a, ColorSetUnspecified)
	b = CoalesceColorSet(b, ColorSetUnspecified)

	return a.Color == b.Color &&
		a.OnColor == b.OnColor
}

func EqualColorSet(a, b *ColorSet) bool {
	if !SameColorSet(a, b) {
		return SemanticEqualColorSet(a, b)
	}
	return true
}

func MergeColorSet(a, b *ColorSet) *ColorSet {
	if b == nil {
		return a
	}
	if a == nil {
		return b
	}

	a = CoalesceColorSet(a, ColorSetUnspecified)
	b = CoalesceColorSet(b, ColorSetUnspecified)

	return &ColorSet{
		Color:   graphics.ColorCompanion.TakeOrElse(b.Color, a.Color),
		OnColor: graphics.ColorCompanion.TakeOrElse(b.OnColor, a.OnColor),
	}
}

func ColorSetFromTokens(token token.MatColorSet) *ColorSet {
	return &ColorSet{
		Color:   graphics.FromNRGBA(token.Color.AsNRGBA()),
		OnColor: graphics.FromNRGBA(token.OnColor.AsNRGBA()),
	}
}
