package material3

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
)

type ColorSchemeOption func(*ColorScheme)

func WithPrimary(colSet *ColorSet) ColorSchemeOption {
	if !IsSpecifiedColorSet(colSet) {
		panic("primary cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.Primary = colSet
	}
}

func WithPrimaryContainer(colSet *ColorSet) ColorSchemeOption {
	if !IsSpecifiedColorSet(colSet) {
		panic("primaryContainer cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.PrimaryContainer = colSet
	}
}

func WithPrimaryFixed(colSet *ColorSet) ColorSchemeOption {
	if !IsSpecifiedColorSet(colSet) {
		panic("primaryFixed cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.PrimaryFixed = colSet
	}
}

func WithPrimaryFixedVariant(colSet *ColorSet) ColorSchemeOption {
	if !IsSpecifiedColorSet(colSet) {
		panic("primaryFixedVariant cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.PrimaryFixedVariant = colSet
	}
}

func WithSecondary(colSet *ColorSet) ColorSchemeOption {
	if !IsSpecifiedColorSet(colSet) {
		panic("secondary cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.Secondary = colSet
	}
}

func WithSecondaryContainer(colSet *ColorSet) ColorSchemeOption {
	if !IsSpecifiedColorSet(colSet) {
		panic("secondaryContainer cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.SecondaryContainer = colSet
	}
}

func WithSecondaryFixed(colSet *ColorSet) ColorSchemeOption {
	if !IsSpecifiedColorSet(colSet) {
		panic("secondaryFixed cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.SecondaryFixed = colSet
	}
}

func WithSecondaryFixedVariant(colSet *ColorSet) ColorSchemeOption {
	if !IsSpecifiedColorSet(colSet) {
		panic("secondaryFixedVariant cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.SecondaryFixedVariant = colSet
	}
}

func WithTertiary(colSet *ColorSet) ColorSchemeOption {
	if !IsSpecifiedColorSet(colSet) {
		panic("tertiary cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.Tertiary = colSet
	}
}

func WithTertiaryContainer(colSet *ColorSet) ColorSchemeOption {
	if !IsSpecifiedColorSet(colSet) {
		panic("tertiaryContainer cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.TertiaryContainer = colSet
	}
}

func WithTertiaryFixed(colSet *ColorSet) ColorSchemeOption {
	if !IsSpecifiedColorSet(colSet) {
		panic("tertiaryFixed cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.TertiaryFixed = colSet
	}
}

func WithTertiaryFixedVariant(colSet *ColorSet) ColorSchemeOption {
	if !IsSpecifiedColorSet(colSet) {
		panic("tertiaryFixedVariant cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.TertiaryFixedVariant = colSet
	}
}

func WithSurface(colSet *ColorSet) ColorSchemeOption {
	if !IsSpecifiedColorSet(colSet) {
		panic("surface cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.Surface = colSet
	}
}

func WithSurfaceVariant(colSet *ColorSet) ColorSchemeOption {
	if !IsSpecifiedColorSet(colSet) {
		panic("surfaceVariant cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.SurfaceVariant = colSet
	}
}

func WithSurfaceDim(col Color) ColorSchemeOption {
	if !graphics.ColorCompanion.IsSpecified(col) {
		panic("surfaceDim cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.SurfaceDim = col
	}
}

func WithSurfaceBright(col Color) ColorSchemeOption {
	if !graphics.ColorCompanion.IsSpecified(col) {
		panic("surfaceBright cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.SurfaceBright = col
	}
}

func WithSurfaceContainerLowest(col Color) ColorSchemeOption {
	if !graphics.ColorCompanion.IsSpecified(col) {
		panic("surfaceContainerLowest cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.SurfaceContainerLowest = col
	}
}

func WithSurfaceContainerLow(col Color) ColorSchemeOption {
	if !graphics.ColorCompanion.IsSpecified(col) {
		panic("surfaceContainerLow cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.SurfaceContainerLow = col
	}
}

func WithSurfaceContainer(col Color) ColorSchemeOption {
	if !graphics.ColorCompanion.IsSpecified(col) {
		panic("surfaceContainer cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.SurfaceContainer = col
	}
}

func WithSurfaceContainerHigh(col Color) ColorSchemeOption {
	if !graphics.ColorCompanion.IsSpecified(col) {
		panic("surfaceContainerHigh cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.SurfaceContainerHigh = col
	}
}

func WithSurfaceContainerHighest(col Color) ColorSchemeOption {
	if !graphics.ColorCompanion.IsSpecified(col) {
		panic("surfaceContainerHighest cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.SurfaceContainerHighest = col
	}
}

// colset
func WithInverseSurface(colSet *ColorSet) ColorSchemeOption {
	if !IsSpecifiedColorSet(colSet) {
		panic("inverseSurface cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.InverseSurface = colSet
	}
}

func WithInversePrimary(col Color) ColorSchemeOption {
	if !graphics.ColorCompanion.IsSpecified(col) {
		panic("inversePrimary cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.InversePrimary = col
	}
}

func WithBackground(colSet *ColorSet) ColorSchemeOption {
	if !IsSpecifiedColorSet(colSet) {
		panic("background cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.Background = colSet
	}
}

func WithOutline(col Color) ColorSchemeOption {
	if !graphics.ColorCompanion.IsSpecified(col) {
		panic("outline cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.Outline = col
	}
}

func WithOutlineVariant(col Color) ColorSchemeOption {
	if !graphics.ColorCompanion.IsSpecified(col) {
		panic("outlineVariant cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.OutlineVariant = col
	}
}

func WithError(colSet *ColorSet) ColorSchemeOption {
	if !IsSpecifiedColorSet(colSet) {
		panic("error cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.Error = colSet
	}
}

func WithErrorContainer(colSet *ColorSet) ColorSchemeOption {
	if !IsSpecifiedColorSet(colSet) {
		panic("errorContainer cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.ErrorContainer = colSet
	}
}

func WithScrim(col Color) ColorSchemeOption {
	if !graphics.ColorCompanion.IsSpecified(col) {
		panic("scrim cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.Scrim = col
	}
}

func WithShadow(col Color) ColorSchemeOption {
	if !graphics.ColorCompanion.IsSpecified(col) {
		panic("shadow cannot be unspecified")
	}
	return func(c *ColorScheme) {
		c.Shadow = col
	}
}
