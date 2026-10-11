package material3

import (
	"fmt"
	"math"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/material3/tokens"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// Immutable
type ColorSchemeNext struct {
	Primary                 graphics.Color
	OnPrimary               graphics.Color
	PrimaryContainer        graphics.Color
	OnPrimaryContainer      graphics.Color
	InversePrimary          graphics.Color
	Secondary               graphics.Color
	OnSecondary             graphics.Color
	SecondaryContainer      graphics.Color
	OnSecondaryContainer    graphics.Color
	Tertiary                graphics.Color
	OnTertiary              graphics.Color
	TertiaryContainer       graphics.Color
	OnTertiaryContainer     graphics.Color
	Background              graphics.Color
	OnBackground            graphics.Color
	Surface                 graphics.Color
	OnSurface               graphics.Color
	SurfaceVariant          graphics.Color
	OnSurfaceVariant        graphics.Color
	SurfaceTint             graphics.Color
	InverseSurface          graphics.Color
	InverseOnSurface        graphics.Color
	Error                   graphics.Color
	OnError                 graphics.Color
	ErrorContainer          graphics.Color
	OnErrorContainer        graphics.Color
	Outline                 graphics.Color
	OutlineVariant          graphics.Color
	Scrim                   graphics.Color
	SurfaceBright           graphics.Color
	SurfaceContainer        graphics.Color
	SurfaceContainerHigh    graphics.Color
	SurfaceContainerHighest graphics.Color
	SurfaceContainerLow     graphics.Color
	SurfaceContainerLowest  graphics.Color
	SurfaceDim              graphics.Color
	PrimaryFixed            graphics.Color
	PrimaryFixedDim         graphics.Color
	OnPrimaryFixed          graphics.Color
	OnPrimaryFixedVariant   graphics.Color
	SecondaryFixed          graphics.Color
	SecondaryFixedDim       graphics.Color
	OnSecondaryFixed        graphics.Color
	OnSecondaryFixedVariant graphics.Color
	TertiaryFixed           graphics.Color
	TertiaryFixedDim        graphics.Color
	OnTertiaryFixed         graphics.Color
	OnTertiaryFixedVariant  graphics.Color
}

func (c ColorSchemeNext) ContentFor(backgroundColor graphics.Color) graphics.Color {
	switch backgroundColor {
	case c.Primary:
		return c.OnPrimary
	case c.Secondary:
		return c.OnSecondary
	case c.Tertiary:
		return c.OnTertiary
	case c.Background:
		return c.OnBackground
	case c.Error:
		return c.OnError

	case c.PrimaryContainer:
		return c.OnPrimaryContainer
	case c.SecondaryContainer:
		return c.OnSecondaryContainer
	case c.TertiaryContainer:
		return c.OnTertiaryContainer
	case c.ErrorContainer:
		return c.OnErrorContainer

	case c.InverseSurface:
		return c.InverseOnSurface
	case c.Surface:
		return c.OnSurface
	case c.SurfaceVariant:
		return c.OnSurfaceVariant
	case c.SurfaceBright:
		return c.OnSurface
	case c.SurfaceContainer:
		return c.OnSurface
	case c.SurfaceContainerHigh:
		return c.OnSurface
	case c.SurfaceContainerHighest:
		return c.OnSurface
	case c.SurfaceContainerLow:
		return c.OnSurface
	case c.SurfaceContainerLowest:
		return c.OnSurface
	case c.SurfaceDim:
		return c.OnSurface

	case c.PrimaryFixed:
		return c.OnPrimaryFixed
	case c.PrimaryFixedDim:
		return c.OnPrimaryFixed
	case c.SecondaryFixed:
		return c.OnSecondaryFixed
	case c.SecondaryFixedDim:
		return c.OnSecondaryFixed
	case c.TertiaryFixed:
		return c.OnTertiaryFixed
	case c.TertiaryFixedDim:
		return c.OnTertiaryFixed

	default:
		return graphics.ColorCompanion.Unspecified()

	}
}

func (c ColorSchemeNext) SurfaceColorAtElevation(elevation unit.Dp) graphics.Color {
	if elevation == 0 {
		return c.Surface
	}

	alpha := ((4.5 * math.Log(float64(elevation)+1)) + 2) / 100
	return c.SurfaceTint.Copy(graphics.CopyWithAlpha(float32(alpha))).CompositeOver(c.Surface)
}

/**
 * Helper function for component color tokens. Here is an example on how to use component color
 * tokens: ``MaterialTheme.colorScheme.fromToken(ExtendedFabBranded.BrandedContainerColor)``
 */
func (c ColorSchemeNext) FromToken(value tokens.ColorSchemeTokenKey) graphics.Color {
	switch value {
	case tokens.ColorSchemeTokenKeyUnspecified:
		return graphics.ColorCompanion.Unspecified()
	case tokens.ColorSchemeTokenKeyBackground:
		return c.Background
	case tokens.ColorSchemeTokenKeyError:
		return c.Error
	case tokens.ColorSchemeTokenKeyErrorContainer:
		return c.ErrorContainer
	case tokens.ColorSchemeTokenKeyInverseOnSurface:
		return c.InverseOnSurface
	case tokens.ColorSchemeTokenKeyInversePrimary:
		return c.InversePrimary
	case tokens.ColorSchemeTokenKeyInverseSurface:
		return c.InverseSurface
	case tokens.ColorSchemeTokenKeyOnBackground:
		return c.OnBackground
	case tokens.ColorSchemeTokenKeyOnError:
		return c.OnError
	case tokens.ColorSchemeTokenKeyOnErrorContainer:
		return c.OnErrorContainer
	case tokens.ColorSchemeTokenKeyOnPrimary:
		return c.OnPrimary
	case tokens.ColorSchemeTokenKeyOnPrimaryContainer:
		return c.OnPrimaryContainer
	case tokens.ColorSchemeTokenKeyOnPrimaryFixed:
		return c.OnPrimaryFixed
	case tokens.ColorSchemeTokenKeyOnPrimaryFixedVariant:
		return c.OnPrimaryFixedVariant
	case tokens.ColorSchemeTokenKeyOnSecondary:
		return c.OnSecondary
	case tokens.ColorSchemeTokenKeyOnSecondaryContainer:
		return c.OnSecondaryContainer
	case tokens.ColorSchemeTokenKeyOnSecondaryFixed:
		return c.OnSecondaryFixed
	case tokens.ColorSchemeTokenKeyOnSecondaryFixedVariant:
		return c.OnSecondaryFixedVariant
	case tokens.ColorSchemeTokenKeyOnSurface:
		return c.OnSurface
	case tokens.ColorSchemeTokenKeyOnSurfaceVariant:
		return c.OnSurfaceVariant
	case tokens.ColorSchemeTokenKeyOnTertiary:
		return c.OnTertiary
	case tokens.ColorSchemeTokenKeyOnTertiaryContainer:
		return c.OnTertiaryContainer
	case tokens.ColorSchemeTokenKeyOnTertiaryFixed:
		return c.OnTertiaryFixed
	case tokens.ColorSchemeTokenKeyOnTertiaryFixedVariant:
		return c.OnTertiaryFixedVariant
	case tokens.ColorSchemeTokenKeyOutline:
		return c.Outline
	case tokens.ColorSchemeTokenKeyOutlineVariant:
		return c.OutlineVariant
	case tokens.ColorSchemeTokenKeyPrimary:
		return c.Primary
	case tokens.ColorSchemeTokenKeyPrimaryContainer:
		return c.PrimaryContainer
	case tokens.ColorSchemeTokenKeyPrimaryFixed:
		return c.PrimaryFixed
	case tokens.ColorSchemeTokenKeyPrimaryFixedDim:
		return c.PrimaryFixedDim
	case tokens.ColorSchemeTokenKeyScrim:
		return c.Scrim
	case tokens.ColorSchemeTokenKeySecondary:
		return c.Secondary
	case tokens.ColorSchemeTokenKeySecondaryContainer:
		return c.SecondaryContainer
	case tokens.ColorSchemeTokenKeySecondaryFixed:
		return c.SecondaryFixed
	case tokens.ColorSchemeTokenKeySecondaryFixedDim:
		return c.SecondaryFixedDim
	case tokens.ColorSchemeTokenKeySurface:
		return c.Surface
	case tokens.ColorSchemeTokenKeySurfaceBright:
		return c.SurfaceBright
	case tokens.ColorSchemeTokenKeySurfaceContainer:
		return c.SurfaceContainer
	case tokens.ColorSchemeTokenKeySurfaceContainerHigh:
		return c.SurfaceContainerHigh
	case tokens.ColorSchemeTokenKeySurfaceContainerHighest:
		return c.SurfaceContainerHighest
	case tokens.ColorSchemeTokenKeySurfaceContainerLow:
		return c.SurfaceContainerLow
	case tokens.ColorSchemeTokenKeySurfaceContainerLowest:
		return c.SurfaceContainerLowest
	case tokens.ColorSchemeTokenKeySurfaceDim:
		return c.SurfaceDim
	case tokens.ColorSchemeTokenKeySurfaceTint:
		return c.SurfaceTint
	case tokens.ColorSchemeTokenKeySurfaceVariant:
		return c.SurfaceVariant
	case tokens.ColorSchemeTokenKeyTertiary:
		return c.Tertiary
	case tokens.ColorSchemeTokenKeyTertiaryContainer:
		return c.TertiaryContainer
	case tokens.ColorSchemeTokenKeyTertiaryFixed:
		return c.TertiaryFixed
	case tokens.ColorSchemeTokenKeyTertiaryFixedDim:
		return c.TertiaryFixedDim
	default:
		panic(fmt.Sprintf("unknown color scheme token key: %s", value))
	}
}

func LightColorScheme() *ColorSchemeNext {
	return &ColorSchemeNext{
		Primary:                 tokens.ColorLightTokens.Primary,
		OnPrimary:               tokens.ColorLightTokens.OnPrimary,
		PrimaryContainer:        tokens.ColorLightTokens.PrimaryContainer,
		OnPrimaryContainer:      tokens.ColorLightTokens.OnPrimaryContainer,
		InversePrimary:          tokens.ColorLightTokens.InversePrimary,
		Secondary:               tokens.ColorLightTokens.Secondary,
		OnSecondary:             tokens.ColorLightTokens.OnSecondary,
		SecondaryContainer:      tokens.ColorLightTokens.SecondaryContainer,
		OnSecondaryContainer:    tokens.ColorLightTokens.OnSecondaryContainer,
		Tertiary:                tokens.ColorLightTokens.Tertiary,
		OnTertiary:              tokens.ColorLightTokens.OnTertiary,
		TertiaryContainer:       tokens.ColorLightTokens.TertiaryContainer,
		OnTertiaryContainer:     tokens.ColorLightTokens.OnTertiaryContainer,
		Background:              tokens.ColorLightTokens.Background,
		OnBackground:            tokens.ColorLightTokens.OnBackground,
		Surface:                 tokens.ColorLightTokens.Surface,
		OnSurface:               tokens.ColorLightTokens.OnSurface,
		SurfaceVariant:          tokens.ColorLightTokens.SurfaceVariant,
		OnSurfaceVariant:        tokens.ColorLightTokens.OnSurfaceVariant,
		SurfaceTint:             tokens.ColorLightTokens.SurfaceTint,
		InverseSurface:          tokens.ColorLightTokens.InverseSurface,
		InverseOnSurface:        tokens.ColorLightTokens.InverseOnSurface,
		Error:                   tokens.ColorLightTokens.Error,
		OnError:                 tokens.ColorLightTokens.OnError,
		ErrorContainer:          tokens.ColorLightTokens.ErrorContainer,
		OnErrorContainer:        tokens.ColorLightTokens.OnErrorContainer,
		Outline:                 tokens.ColorLightTokens.Outline,
		OutlineVariant:          tokens.ColorLightTokens.OutlineVariant,
		Scrim:                   tokens.ColorLightTokens.Scrim,
		SurfaceBright:           tokens.ColorLightTokens.SurfaceBright,
		SurfaceContainer:        tokens.ColorLightTokens.SurfaceContainer,
		SurfaceContainerHigh:    tokens.ColorLightTokens.SurfaceContainerHigh,
		SurfaceContainerHighest: tokens.ColorLightTokens.SurfaceContainerHighest,
		SurfaceContainerLow:     tokens.ColorLightTokens.SurfaceContainerLow,
		SurfaceContainerLowest:  tokens.ColorLightTokens.SurfaceContainerLowest,
		SurfaceDim:              tokens.ColorLightTokens.SurfaceDim,
		PrimaryFixed:            tokens.ColorLightTokens.PrimaryFixed,
		PrimaryFixedDim:         tokens.ColorLightTokens.PrimaryFixedDim,
		OnPrimaryFixed:          tokens.ColorLightTokens.OnPrimaryFixed,
		OnPrimaryFixedVariant:   tokens.ColorLightTokens.OnPrimaryFixedVariant,
		SecondaryFixed:          tokens.ColorLightTokens.SecondaryFixed,
		SecondaryFixedDim:       tokens.ColorLightTokens.SecondaryFixedDim,
		OnSecondaryFixed:        tokens.ColorLightTokens.OnSecondaryFixed,
		OnSecondaryFixedVariant: tokens.ColorLightTokens.OnSecondaryFixedVariant,
		TertiaryFixed:           tokens.ColorLightTokens.TertiaryFixed,
		TertiaryFixedDim:        tokens.ColorLightTokens.TertiaryFixedDim,
		OnTertiaryFixed:         tokens.ColorLightTokens.OnTertiaryFixed,
		OnTertiaryFixedVariant:  tokens.ColorLightTokens.OnTertiaryFixedVariant,
	}
}

func DarkColorScheme() *ColorSchemeNext {
	return &ColorSchemeNext{
		Primary:                 tokens.ColorDarkTokens.Primary,
		OnPrimary:               tokens.ColorDarkTokens.OnPrimary,
		PrimaryContainer:        tokens.ColorDarkTokens.PrimaryContainer,
		OnPrimaryContainer:      tokens.ColorDarkTokens.OnPrimaryContainer,
		InversePrimary:          tokens.ColorDarkTokens.InversePrimary,
		Secondary:               tokens.ColorDarkTokens.Secondary,
		OnSecondary:             tokens.ColorDarkTokens.OnSecondary,
		SecondaryContainer:      tokens.ColorDarkTokens.SecondaryContainer,
		OnSecondaryContainer:    tokens.ColorDarkTokens.OnSecondaryContainer,
		Tertiary:                tokens.ColorDarkTokens.Tertiary,
		OnTertiary:              tokens.ColorDarkTokens.OnTertiary,
		TertiaryContainer:       tokens.ColorDarkTokens.TertiaryContainer,
		OnTertiaryContainer:     tokens.ColorDarkTokens.OnTertiaryContainer,
		Background:              tokens.ColorDarkTokens.Background,
		OnBackground:            tokens.ColorDarkTokens.OnBackground,
		Surface:                 tokens.ColorDarkTokens.Surface,
		OnSurface:               tokens.ColorDarkTokens.OnSurface,
		SurfaceVariant:          tokens.ColorDarkTokens.SurfaceVariant,
		OnSurfaceVariant:        tokens.ColorDarkTokens.OnSurfaceVariant,
		SurfaceTint:             tokens.ColorDarkTokens.SurfaceTint,
		InverseSurface:          tokens.ColorDarkTokens.InverseSurface,
		InverseOnSurface:        tokens.ColorDarkTokens.InverseOnSurface,
		Error:                   tokens.ColorDarkTokens.Error,
		OnError:                 tokens.ColorDarkTokens.OnError,
		ErrorContainer:          tokens.ColorDarkTokens.ErrorContainer,
		OnErrorContainer:        tokens.ColorDarkTokens.OnErrorContainer,
		Outline:                 tokens.ColorDarkTokens.Outline,
		OutlineVariant:          tokens.ColorDarkTokens.OutlineVariant,
		Scrim:                   tokens.ColorDarkTokens.Scrim,
		SurfaceBright:           tokens.ColorDarkTokens.SurfaceBright,
		SurfaceContainer:        tokens.ColorDarkTokens.SurfaceContainer,
		SurfaceContainerHigh:    tokens.ColorDarkTokens.SurfaceContainerHigh,
		SurfaceContainerHighest: tokens.ColorDarkTokens.SurfaceContainerHighest,
		SurfaceContainerLow:     tokens.ColorDarkTokens.SurfaceContainerLow,
		SurfaceContainerLowest:  tokens.ColorDarkTokens.SurfaceContainerLowest,
		SurfaceDim:              tokens.ColorDarkTokens.SurfaceDim,
		PrimaryFixed:            tokens.ColorDarkTokens.PrimaryFixed,
		PrimaryFixedDim:         tokens.ColorDarkTokens.PrimaryFixedDim,
		OnPrimaryFixed:          tokens.ColorDarkTokens.OnPrimaryFixed,
		OnPrimaryFixedVariant:   tokens.ColorDarkTokens.OnPrimaryFixedVariant,
		SecondaryFixed:          tokens.ColorDarkTokens.SecondaryFixed,
		SecondaryFixedDim:       tokens.ColorDarkTokens.SecondaryFixedDim,
		OnSecondaryFixed:        tokens.ColorDarkTokens.OnSecondaryFixed,
		OnSecondaryFixedVariant: tokens.ColorDarkTokens.OnSecondaryFixedVariant,
		TertiaryFixed:           tokens.ColorDarkTokens.TertiaryFixed,
		TertiaryFixedDim:        tokens.ColorDarkTokens.TertiaryFixedDim,
		OnTertiaryFixed:         tokens.ColorDarkTokens.OnTertiaryFixed,
		OnTertiaryFixedVariant:  tokens.ColorDarkTokens.OnTertiaryFixedVariant,
	}
}

var LocalColorSchemeNext = compose.CompositionLocalOf(func() *ColorSchemeNext {
	return LightColorScheme()
})
