package colorrole

import "fmt"

// // Bg is the background color atop which content is currently being
// // drawn.
// Bg color.NRGBA

// // Fg is a color suitable for drawing on top of Bg.
// Fg color.NRGBA

// // ContrastBg is a color used to draw attention to active,
// // important, interactive widgets such as buttons.
// ContrastBg color.NRGBA

// // ContrastFg is a color suitable for content drawn on top of
// // ContrastBg.
// ContrastFg color.NRGBA
// type ColorRole int

// const (
// 	Bg ColorRole = iota
// 	Fg
// 	ContrastBg
// 	ContrastFg
// )

// BG = color
// FG = onColor
// CONTRAST_BG = color
// CONTRAST_FG = onColor

// WIP for future
type ColorRole int

const (
	Primary ColorRole = iota
	OnPrimary
	PrimaryContainer
	OnPrimaryContainer
	PrimaryFixed
	OnPrimaryFixed
	PrimaryFixedVariant
	OnPrimaryFixedVariant

	Secondary
	OnSecondary
	SecondaryContainer
	OnSecondaryContainer
	SecondaryFixed
	OnSecondaryFixed
	SecondaryFixedVariant
	OnSecondaryFixedVariant

	Tertiary
	OnTertiary
	TertiaryContainer
	OnTertiaryContainer
	TertiaryFixed
	OnTertiaryFixed
	TertiaryFixedVariant
	OnTertiaryFixedVariant

	Surface
	OnSurface
	SurfaceVariant
	OnSurfaceVariant

	SurfaceDim
	SurfaceBright
	SurfaceContainerLowest
	SurfaceContainerLow
	SurfaceContainer
	SurfaceContainerHigh
	SurfaceContainerHighest

	InverseSurface
	OnInverseSurface
	InversePrimary

	Background
	OnBackground

	Outline
	OutlineVariant

	Error
	OnError
	ErrorContainer
	OnErrorContainer

	Scrim
	Shadow

	BasicBg
	BasicFg
	BasicContrastBg
	BasicContrastFg
)

func (cr ColorRole) String() string {
	switch cr {
	case Primary:
		return "Primary "
	case OnPrimary:
		return "OnPrimary"
	case PrimaryContainer:
		return "PrimaryContainer"
	case OnPrimaryContainer:
		return "OnPrimaryContainer"
	case PrimaryFixed:
		return "PrimaryFixed"
	case OnPrimaryFixed:
		return "OnPrimaryFixed"
	case PrimaryFixedVariant:
		return "PrimaryFixedVariant"
	case OnPrimaryFixedVariant:
		return "OnPrimaryFixedVariant"
	case Secondary:
		return "Secondary"
	case OnSecondary:
		return "OnSecondary"
	case SecondaryContainer:
		return "SecondaryContainer"
	case OnSecondaryContainer:
		return "OnSecondaryContainer"
	case SecondaryFixed:
		return "SecondaryFixed"
	case OnSecondaryFixed:
		return "OnSecondaryFixed"
	case SecondaryFixedVariant:
		return "SecondaryFixedVariant"
	case OnSecondaryFixedVariant:
		return "OnSecondaryFixedVariant"
	case Tertiary:
		return "Tertiary"
	case OnTertiary:
		return "OnTertiary"
	case TertiaryContainer:
		return "TertiaryContainer"
	case OnTertiaryContainer:
		return "OnTertiaryContainer"
	case TertiaryFixed:
		return "TertiaryFixed"
	case OnTertiaryFixed:
		return "OnTertiaryFixed"
	case TertiaryFixedVariant:
		return "TertiaryFixedVariant"
	case OnTertiaryFixedVariant:
		return "OnTertiaryFixedVariant"
	case Surface:
		return "Surface"
	case OnSurface:
		return "OnSurface"
	case SurfaceVariant:
		return "SurfaceVariant"
	case OnSurfaceVariant:
		return "OnSurfaceVariant"
	case SurfaceDim:
		return "SurfaceDim"
	case SurfaceBright:
		return "SurfaceBright"
	case SurfaceContainerLowest:
		return "SurfaceContainerLowest"
	case SurfaceContainerLow:
		return "SurfaceContainerLow"
	case SurfaceContainer:
		return "SurfaceContainer"
	case SurfaceContainerHigh:
		return "SurfaceContainerHigh"
	case SurfaceContainerHighest:
		return "SurfaceContainerHighest"
	case InverseSurface:
		return "InverseSurface"
	case OnInverseSurface:
		return "OnInverseSurface"
	case InversePrimary:
		return "InversePrimary"
	case Background:
		return "Background"
	case OnBackground:
		return "OnBackground"
	case Outline:
		return "Outline"
	case OutlineVariant:
		return "OutlineVariant"
	case Error:
		return "Error"
	case OnError:
		return "OnError"
	case ErrorContainer:
		return "ErrorContainer"
	case OnErrorContainer:
		return "OnErrorContainer"
	case Scrim:
		return "Scrim"
	case Shadow:
		return "Shadow"
	case BasicBg:
		return "BasicBg"
	case BasicFg:
		return "BasicFg"
	case BasicContrastBg:
		return "BasicContrastBg"
	case BasicContrastFg:
		return "BasicContrastFg"

	}
	panic(fmt.Sprintf("unknown color role %v", cr))
}
