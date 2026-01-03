package theme

import (
	"fmt"
	"image/color"

	"github.com/zodimo/go-compose/theme/colorrole"
)

type ThemeColorReaderFunc func(theme *Theme) ThemeColor
type ThemeBasicColorReaderFunc func(theme *BasicTheme) ThemeColor

type ThemeColorResolver interface {
	ResolveColorDescriptor(colorDesc ColorDescriptor) ThemeColor
}

type themeColorResolver struct {
	tm ThemeManager
}

func (cr *themeColorResolver) Material3(reader ThemeColorReaderFunc) ThemeColor {
	return reader(cr.tm.getMaterial3Theme())
}

func (cr *themeColorResolver) Material(reader ThemeBasicColorReaderFunc) ThemeColor {
	return reader(cr.tm.materialTheme())
}

func (cr *themeColorResolver) ResolveColorDescriptor(colorDesc ColorDescriptor) ThemeColor {
	//this should panic when colorDesc is a unspecifiedColorDescriptor, its a logical error
	if !IsSpecifiedColor(colorDesc) {
		panic("unspecified color descriptor cannot be resolved")
	}
	colorDescriptor := colorDesc.(colorDescriptor)
	resolvedColor := cr.resolveColorDescriptor(colorDescriptor)
	return cr.applyUpdates(colorDesc.Updates(), resolvedColor.AsTokenColor())
}

func (cr *themeColorResolver) resolveColorDescriptor(colorDesc colorDescriptor) ThemeColor {
	if colorDesc.isColor {
		return ThemeColorFromGraphicsColor(colorDesc.color)
	}

	switch colorDesc.colorRole {
	//Basic
	case colorrole.BasicBg:
		return ThemeColorFromNRGBA(cr.tm.materialTheme().Bg)
	case colorrole.BasicFg:
		return ThemeColorFromNRGBA(cr.tm.materialTheme().Fg)
	case colorrole.BasicContrastBg:
		return ThemeColorFromNRGBA(cr.tm.materialTheme().ContrastBg)
	case colorrole.BasicContrastFg:
		return ThemeColorFromNRGBA(cr.tm.materialTheme().ContrastFg)

	//PrimaryRoles
	case colorrole.Primary:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.Primary.Color)
	case colorrole.OnPrimary:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.Primary.OnColor)
	case colorrole.PrimaryContainer:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.PrimaryContainer.Color)
	case colorrole.OnPrimaryContainer:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.PrimaryContainer.OnColor)
	case colorrole.PrimaryFixed:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.PrimaryFixed.Color)
	case colorrole.OnPrimaryFixed:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.PrimaryFixed.OnColor)
	case colorrole.PrimaryFixedVariant:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.PrimaryFixedVariant.Color)
	case colorrole.OnPrimaryFixedVariant:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.PrimaryFixedVariant.OnColor)

	//SecondaryColorRoles

	case colorrole.Secondary:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.Secondary.Color)
	case colorrole.OnSecondary:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.Secondary.OnColor)
	case colorrole.SecondaryContainer:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.SecondaryContainer.Color)
	case colorrole.OnSecondaryContainer:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.SecondaryContainer.OnColor)
	case colorrole.SecondaryFixed:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.SecondaryFixed.Color)
	case colorrole.OnSecondaryFixed:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.SecondaryFixed.OnColor)
	case colorrole.SecondaryFixedVariant:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.SecondaryFixedVariant.Color)
	case colorrole.OnSecondaryFixedVariant:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.SecondaryFixedVariant.OnColor)

		// TertiaryRoles
	case colorrole.Tertiary:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.Tertiary.Color)
	case colorrole.OnTertiary:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.Tertiary.OnColor)
	case colorrole.TertiaryContainer:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.TertiaryContainer.Color)
	case colorrole.OnTertiaryContainer:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.TertiaryContainer.OnColor)
	case colorrole.TertiaryFixed:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.TertiaryFixed.Color)
	case colorrole.OnTertiaryFixed:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.TertiaryFixed.OnColor)
	case colorrole.TertiaryFixedVariant:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.TertiaryFixedVariant.Color)
	case colorrole.OnTertiaryFixedVariant:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.TertiaryFixedVariant.OnColor)

		// SurfaceRoles
	case colorrole.Surface:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.Surface.Color)
	case colorrole.OnSurface:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.Surface.OnColor)
	case colorrole.SurfaceVariant:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.SurfaceVariant.Color)
	case colorrole.OnSurfaceVariant:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.SurfaceVariant.OnColor)
	case colorrole.SurfaceDim:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.SurfaceDim)
	case colorrole.SurfaceBright:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.SurfaceBright)
	case colorrole.SurfaceContainerLowest:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.SurfaceContainerLowest)
	case colorrole.SurfaceContainerLow:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.SurfaceContainerLow)
	case colorrole.SurfaceContainer:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.SurfaceContainer)
	case colorrole.SurfaceContainerHigh:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.SurfaceContainerHigh)
	case colorrole.SurfaceContainerHighest:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.SurfaceContainerHighest)

		// InverseRoles
	case colorrole.InverseSurface:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.InverseSurface.Color)
	case colorrole.OnInverseSurface:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.InverseSurface.OnColor)
	case colorrole.InversePrimary:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.InversePrimary)
		// BackgroundColorRoles
	case colorrole.Background:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.Background.Color)
	case colorrole.OnBackground:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.Background.OnColor)
		// OutlineRoles
	case colorrole.Outline:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.Outline)
	case colorrole.OutlineVariant:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.OutlineVariant)
		// ErrorRoles
	case colorrole.Error:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.Error.Color)
	case colorrole.OnError:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.Error.OnColor)
	case colorrole.ErrorContainer:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.ErrorContainer.Color)
	case colorrole.OnErrorContainer:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.ErrorContainer.OnColor)
		// ScrimRoles
	case colorrole.Scrim:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.Scrim)
	case colorrole.Shadow:
		return ThemeColorFromTokenColor(cr.tm.getMaterial3Theme().Scheme.Shadow)
	}
	panic(fmt.Sprintf("Unknown color role: %s\n", colorDesc.colorRole))

}
func newThemeColorResolver(tm ThemeManager) ThemeColorResolver {
	return &themeColorResolver{tm: tm}
}

func (cr *themeColorResolver) applyUpdates(updates []ColorUpdate, tokenColor TokenColor) ThemeColor {
	currentColor := tokenColor.AsNRGBA()

	for _, update := range updates {
		switch update.Action() {
		case SetOpacityColorUpdateAction:
			tokenColor = tokenColor.SetOpacity(GetOpacity(update))
			currentColor = tokenColor.AsNRGBA()
		case LightenColorUpdateAction:
			percentage := GetLighten(update)
			currentColor = lightenColor(currentColor, percentage)
		case DarkenColorUpdateAction:
			percentage := GetDarken(update)
			currentColor = darkenColor(currentColor, percentage)
		case LerpColorUpdateAction:
			params := GetLerp(update)
			resolvedStop := cr.ResolveColorDescriptor(params.Stop)
			currentColorNRGBA := currentColor
			stopColorNRGBA := resolvedStop.AsNRGBA()

			// Use the lerp package to interpolate
			// We need to convert NRGBA to the struct expected by lerp.ColorLerp
			// ColorLerp expects struct{ R, G, B, A float32 }

			startFloat := struct{ R, G, B, A float32 }{
				R: float32(currentColorNRGBA.R),
				G: float32(currentColorNRGBA.G),
				B: float32(currentColorNRGBA.B),
				A: float32(currentColorNRGBA.A),
			}
			stopFloat := struct{ R, G, B, A float32 }{
				R: float32(stopColorNRGBA.R),
				G: float32(stopColorNRGBA.G),
				B: float32(stopColorNRGBA.B),
				A: float32(stopColorNRGBA.A),
			}

			res := colorLerp(startFloat, stopFloat, params.Fraction)

			// Clamp results to 0-255 (though ColorLerp usually handles range if inputs are generic)
			// But here we are dealing with 0-255 mapped to float.
			// Wait, lerp.ColorLerp works on float32. standard NRGBA is uint8.
			// Let's assume standard behavior.

			currentColor = color.NRGBA{
				R: uint8(res.R),
				G: uint8(res.G),
				B: uint8(res.B),
				A: uint8(res.A),
			}
		}
	}

	return ThemeColorFromNRGBA(currentColor)
}
