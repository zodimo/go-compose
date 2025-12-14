package theme

import (
	"fmt"

	"github.com/zodimo/go-compose/theme/colorrole"
)

type ThemeColorReaderFunc func(theme *Theme) ThemeColor
type ThemeBasicColorReaderFunc func(theme *BasicTheme) ThemeColor

type ThemeColorResolver interface {
	ResolveColorDescriptor(colorDesc ThemeColorDescriptor) ThemeColor
}

type themeColorResolver struct {
	tm ThemeManager
}

func (cr *themeColorResolver) Material3(reader ThemeColorReaderFunc) ThemeColor {
	return reader(cr.tm.GetMaterial3Theme())
}

func (cr *themeColorResolver) Material(reader ThemeBasicColorReaderFunc) ThemeColor {
	return reader(cr.tm.MaterialTheme())
}

func (cr *themeColorResolver) ResolveColorDescriptor(colorDesc ThemeColorDescriptor) ThemeColor {
	resolvedColor := cr.resolveColorDescriptor(colorDesc)
	return cr.applyUpdates(colorDesc.updates, resolvedColor.AsTokenColor())
}

func (cr *themeColorResolver) resolveColorDescriptor(colorDesc ThemeColorDescriptor) ThemeColor {
	if colorDesc.isColor {
		return ThemeColorFromColor(colorDesc.color)
	}

	switch colorDesc.colorRole {
	//Basic
	case colorrole.BasicBg:
		return ThemeColorFromNRGBA(cr.tm.MaterialTheme().Bg)
	case colorrole.BasicFg:
		return ThemeColorFromNRGBA(cr.tm.MaterialTheme().Fg)
	case colorrole.BasicContrastBg:
		return ThemeColorFromNRGBA(cr.tm.MaterialTheme().ContrastBg)
	case colorrole.BasicContrastFg:
		return ThemeColorFromNRGBA(cr.tm.MaterialTheme().ContrastFg)

	//PrimaryRoles
	case colorrole.Primary:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.Primary.Color)
	case colorrole.OnPrimary:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.Primary.OnColor)
	case colorrole.PrimaryContainer:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.PrimaryContainer.Color)
	case colorrole.OnPrimaryContainer:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.PrimaryContainer.OnColor)
	case colorrole.PrimaryFixed:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.PrimaryFixed.Color)
	case colorrole.OnPrimaryFixed:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.PrimaryFixed.OnColor)
	case colorrole.PrimaryFixedVariant:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.PrimaryFixedVariant.Color)
	case colorrole.OnPrimaryFixedVariant:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.PrimaryFixedVariant.OnColor)

	//SecondaryColorRoles

	case colorrole.Secondary:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.Secondary.Color)
	case colorrole.OnSecondary:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.Secondary.OnColor)
	case colorrole.SecondaryContainer:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.SecondaryContainer.Color)
	case colorrole.OnSecondaryContainer:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.SecondaryContainer.OnColor)
	case colorrole.SecondaryFixed:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.SecondaryFixed.Color)
	case colorrole.OnSecondaryFixed:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.SecondaryFixed.OnColor)
	case colorrole.SecondaryFixedVariant:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.SecondaryFixedVariant.Color)
	case colorrole.OnSecondaryFixedVariant:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.SecondaryFixedVariant.OnColor)

		// TertiaryRoles
	case colorrole.Tertiary:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.Tertiary.Color)
	case colorrole.OnTertiary:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.Tertiary.OnColor)
	case colorrole.TertiaryContainer:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.TertiaryContainer.Color)
	case colorrole.OnTertiaryContainer:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.TertiaryContainer.OnColor)
	case colorrole.TertiaryFixed:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.TertiaryFixed.Color)
	case colorrole.OnTertiaryFixed:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.TertiaryFixed.OnColor)
	case colorrole.TertiaryFixedVariant:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.TertiaryFixedVariant.Color)
	case colorrole.OnTertiaryFixedVariant:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.TertiaryFixedVariant.OnColor)

		// SurfaceRoles
	case colorrole.Surface:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.Surface.Color)
	case colorrole.OnSurface:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.Surface.OnColor)
	case colorrole.SurfaceVariant:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.SurfaceVariant.Color)
	case colorrole.OnSurfaceVariant:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.SurfaceVariant.OnColor)
	case colorrole.SurfaceDim:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.SurfaceDim)
	case colorrole.SurfaceBright:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.SurfaceBright)
	case colorrole.SurfaceContainerLowest:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.SurfaceContainerLowest)
	case colorrole.SurfaceContainerLow:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.SurfaceContainerLow)
	case colorrole.SurfaceContainer:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.SurfaceContainer)
	case colorrole.SurfaceContainerHigh:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.SurfaceContainerHigh)
	case colorrole.SurfaceContainerHighest:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.SurfaceContainerHighest)

		// InverseRoles
	case colorrole.InverseSurface:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.InverseSurface.Color)
	case colorrole.OnInverseSurface:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.InverseSurface.OnColor)
	case colorrole.InversePrimary:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.InversePrimary)
		// BackgroundColorRoles
	case colorrole.Background:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.Background.Color)
	case colorrole.OnBackground:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.Background.OnColor)
		// OutlineRoles
	case colorrole.Outline:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.Outline)
	case colorrole.OutlineVariant:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.OutlineVariant)
		// ErrorRoles
	case colorrole.Error:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.Error.Color)
	case colorrole.OnError:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.Error.OnColor)
	case colorrole.ErrorContainer:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.ErrorContainer.Color)
	case colorrole.OnErrorContainer:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.ErrorContainer.OnColor)
		// ScrimRoles
	case colorrole.Scrim:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.Scrim)
	case colorrole.Shadow:
		return ThemeColorFromTokenColor(cr.tm.GetMaterial3Theme().Scheme.Shadow)
	}
	panic(fmt.Sprintf("Unknown color role: %s\n", colorDesc.colorRole))

}
func newThemeColorResolver(tm ThemeManager) ThemeColorResolver {
	return &themeColorResolver{tm: tm}
}

func (cr *themeColorResolver) applyUpdates(updates []func(color TokenColor) TokenColor, color TokenColor) ThemeColor {
	for _, update := range updates {
		color = update(color)
	}
	return ThemeColorFromTokenColor(color)
}
