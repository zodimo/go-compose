package button

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/graphics"
)

var ButtonColorsUnspecified = &ButtonColors{
	ContainerColor:         graphics.ColorCompanion.Unspecified(),
	ContentColor:           graphics.ColorCompanion.Unspecified(),
	DisabledContainerColor: graphics.ColorCompanion.Unspecified(),
	DisabledContentColor:   graphics.ColorCompanion.Unspecified(),
}

type ButtonColors struct {
	ContainerColor         graphics.Color
	ContentColor           graphics.Color
	DisabledContainerColor graphics.Color
	DisabledContentColor   graphics.Color
}

func IsSpecifiedButtonColors(c *ButtonColors) bool {
	return c != nil && c != ButtonColorsUnspecified
}

// TakeOrElseButtonColors returns c if specified, otherwise returns defaultColors
func TakeOrElseButtonColors(c, defaultColors *ButtonColors) *ButtonColors {
	if c == nil || c == ButtonColorsUnspecified {
		return defaultColors
	}
	return c
}

// MergeButtonColors merges two ButtonColors, preferring b's specified values over a's
func MergeButtonColors(a, b *ButtonColors) *ButtonColors {
	a = CoalesceButtonColors(a, ButtonColorsUnspecified)
	b = CoalesceButtonColors(b, ButtonColorsUnspecified)

	if a == ButtonColorsUnspecified {
		return b
	}
	if b == ButtonColorsUnspecified {
		return a
	}

	// Both are custom: allocate new merged style
	return &ButtonColors{
		ContainerColor:         graphics.ColorCompanion.TakeOrElse(b.ContainerColor, a.ContainerColor),
		ContentColor:           graphics.ColorCompanion.TakeOrElse(b.ContentColor, a.ContentColor),
		DisabledContainerColor: graphics.ColorCompanion.TakeOrElse(b.DisabledContainerColor, a.DisabledContainerColor),
		DisabledContentColor:   graphics.ColorCompanion.TakeOrElse(b.DisabledContentColor, a.DisabledContentColor),
	}
}

// StringButtonColors returns a string representation of ButtonColors
func StringButtonColors(c *ButtonColors) string {
	if !IsSpecifiedButtonColors(c) {
		return "ButtonColors{Unspecified}"
	}

	return fmt.Sprintf(
		"ButtonColors{ContainerColor: %s, ContentColor: %s, DisabledContainerColor: %s, DisabledContentColor: %s}",
		c.ContainerColor,
		c.ContentColor,
		c.DisabledContainerColor,
		c.DisabledContentColor,
	)
}

// CoalesceButtonColors returns ptr if not nil, otherwise returns def
func CoalesceButtonColors(ptr, def *ButtonColors) *ButtonColors {
	if ptr == nil {
		return def
	}
	return ptr
}

// SameButtonColors returns true if a and b are the same pointer or both unspecified
func SameButtonColors(a, b *ButtonColors) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == ButtonColorsUnspecified
	}
	if b == nil {
		return a == ButtonColorsUnspecified
	}
	return a == b
}

// SemanticEqualButtonColors checks field-by-field equality
func SemanticEqualButtonColors(a, b *ButtonColors) bool {
	a = CoalesceButtonColors(a, ButtonColorsUnspecified)
	b = CoalesceButtonColors(b, ButtonColorsUnspecified)

	return a.ContainerColor == b.ContainerColor &&
		a.ContentColor == b.ContentColor &&
		a.DisabledContainerColor == b.DisabledContainerColor &&
		a.DisabledContentColor == b.DisabledContentColor
}

// EqualButtonColors returns true if a and b are semantically equal
func EqualButtonColors(a, b *ButtonColors) bool {
	if SameButtonColors(a, b) {
		return true
	}
	return SemanticEqualButtonColors(a, b)
}

// ButtonColorsOption is a functional option for CopyButtonColors
type ButtonColorsOption func(*ButtonColors)

// WithButtonContainerColor sets the container color option
func WithButtonContainerColor(c graphics.Color) ButtonColorsOption {
	return func(o *ButtonColors) {
		o.ContainerColor = c
	}
}

// WithButtonContentColor sets the content color option
func WithButtonContentColor(c graphics.Color) ButtonColorsOption {
	return func(o *ButtonColors) {
		o.ContentColor = c
	}
}

// WithButtonDisabledContainerColor sets the disabled container color option
func WithButtonDisabledContainerColor(c graphics.Color) ButtonColorsOption {
	return func(o *ButtonColors) {
		o.DisabledContainerColor = c
	}
}

// WithButtonDisabledContentColor sets the disabled content color option
func WithButtonDisabledContentColor(c graphics.Color) ButtonColorsOption {
	return func(o *ButtonColors) {
		o.DisabledContentColor = c
	}
}

// CopyButtonColors creates a copy with optional modifications
func CopyButtonColors(c *ButtonColors, options ...ButtonColorsOption) *ButtonColors {
	opt := *ButtonColorsUnspecified

	for _, option := range options {
		option(&opt)
	}

	c = CoalesceButtonColors(c, ButtonColorsUnspecified)

	return &ButtonColors{
		ContainerColor:         graphics.ColorCompanion.TakeOrElse(opt.ContainerColor, c.ContainerColor),
		ContentColor:           graphics.ColorCompanion.TakeOrElse(opt.ContentColor, c.ContentColor),
		DisabledContainerColor: graphics.ColorCompanion.TakeOrElse(opt.DisabledContainerColor, c.DisabledContainerColor),
		DisabledContentColor:   graphics.ColorCompanion.TakeOrElse(opt.DisabledContentColor, c.DisabledContentColor),
	}
}
