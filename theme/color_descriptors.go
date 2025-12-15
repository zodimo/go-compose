package theme

import "github.com/zodimo/go-compose/theme/colorrole"

var (
	Primary               = ThemeColorDescriptor{colorRole: colorrole.Primary}
	OnPrimary             = ThemeColorDescriptor{colorRole: colorrole.OnPrimary}
	PrimaryContainer      = ThemeColorDescriptor{colorRole: colorrole.PrimaryContainer}
	OnPrimaryContainer    = ThemeColorDescriptor{colorRole: colorrole.OnPrimaryContainer}
	PrimaryFixed          = ThemeColorDescriptor{colorRole: colorrole.PrimaryFixed}
	OnPrimaryFixed        = ThemeColorDescriptor{colorRole: colorrole.OnPrimaryFixed}
	PrimaryFixedVariant   = ThemeColorDescriptor{colorRole: colorrole.PrimaryFixedVariant}
	OnPrimaryFixedVariant = ThemeColorDescriptor{colorRole: colorrole.OnPrimaryFixedVariant}

	Secondary               = ThemeColorDescriptor{colorRole: colorrole.Secondary}
	OnSecondary             = ThemeColorDescriptor{colorRole: colorrole.OnSecondary}
	SecondaryContainer      = ThemeColorDescriptor{colorRole: colorrole.SecondaryContainer}
	OnSecondaryContainer    = ThemeColorDescriptor{colorRole: colorrole.OnSecondaryContainer}
	SecondaryFixed          = ThemeColorDescriptor{colorRole: colorrole.SecondaryFixed}
	OnSecondaryFixed        = ThemeColorDescriptor{colorRole: colorrole.OnSecondaryFixed}
	SecondaryFixedVariant   = ThemeColorDescriptor{colorRole: colorrole.SecondaryFixedVariant}
	OnSecondaryFixedVariant = ThemeColorDescriptor{colorRole: colorrole.OnSecondaryFixedVariant}

	Tertiary               = ThemeColorDescriptor{colorRole: colorrole.Tertiary}
	OnTertiary             = ThemeColorDescriptor{colorRole: colorrole.OnTertiary}
	TertiaryContainer      = ThemeColorDescriptor{colorRole: colorrole.TertiaryContainer}
	OnTertiaryContainer    = ThemeColorDescriptor{colorRole: colorrole.OnTertiaryContainer}
	TertiaryFixed          = ThemeColorDescriptor{colorRole: colorrole.TertiaryFixed}
	OnTertiaryFixed        = ThemeColorDescriptor{colorRole: colorrole.OnTertiaryFixed}
	TertiaryFixedVariant   = ThemeColorDescriptor{colorRole: colorrole.TertiaryFixedVariant}
	OnTertiaryFixedVariant = ThemeColorDescriptor{colorRole: colorrole.OnTertiaryFixedVariant}

	Surface                 = ThemeColorDescriptor{colorRole: colorrole.Surface}
	OnSurface               = ThemeColorDescriptor{colorRole: colorrole.OnSurface}
	SurfaceVariant          = ThemeColorDescriptor{colorRole: colorrole.SurfaceVariant}
	OnSurfaceVariant        = ThemeColorDescriptor{colorRole: colorrole.OnSurfaceVariant}
	SurfaceDim              = ThemeColorDescriptor{colorRole: colorrole.SurfaceDim}
	SurfaceBright           = ThemeColorDescriptor{colorRole: colorrole.SurfaceBright}
	SurfaceContainerLowest  = ThemeColorDescriptor{colorRole: colorrole.SurfaceContainerLowest}
	SurfaceContainerLow     = ThemeColorDescriptor{colorRole: colorrole.SurfaceContainerLow}
	SurfaceContainer        = ThemeColorDescriptor{colorRole: colorrole.SurfaceContainer}
	SurfaceContainerHigh    = ThemeColorDescriptor{colorRole: colorrole.SurfaceContainerHigh}
	SurfaceContainerHighest = ThemeColorDescriptor{colorRole: colorrole.SurfaceContainerHighest}

	InverseSurface   = ThemeColorDescriptor{colorRole: colorrole.InverseSurface}
	OnInverseSurface = ThemeColorDescriptor{colorRole: colorrole.OnInverseSurface}
	InversePrimary   = ThemeColorDescriptor{colorRole: colorrole.InversePrimary}

	Background   = ThemeColorDescriptor{colorRole: colorrole.Background}
	OnBackground = ThemeColorDescriptor{colorRole: colorrole.OnBackground}

	Outline        = ThemeColorDescriptor{colorRole: colorrole.Outline}
	OutlineVariant = ThemeColorDescriptor{colorRole: colorrole.OutlineVariant}

	Error            = ThemeColorDescriptor{colorRole: colorrole.Error}
	OnError          = ThemeColorDescriptor{colorRole: colorrole.OnError}
	ErrorContainer   = ThemeColorDescriptor{colorRole: colorrole.ErrorContainer}
	OnErrorContainer = ThemeColorDescriptor{colorRole: colorrole.OnErrorContainer}

	Scrim  = ThemeColorDescriptor{colorRole: colorrole.Scrim}
	Shadow = ThemeColorDescriptor{colorRole: colorrole.Shadow}

	// Basic Roles
	BasicBg         = ThemeColorDescriptor{colorRole: colorrole.BasicBg}
	BasicFg         = ThemeColorDescriptor{colorRole: colorrole.BasicFg}
	BasicContrastBg = ThemeColorDescriptor{colorRole: colorrole.BasicContrastBg}
	BasicContrastFg = ThemeColorDescriptor{colorRole: colorrole.BasicContrastFg}
)

type PrimaryColorRoleDescriptors struct {
	Primary               ThemeColorDescriptor
	OnPrimary             ThemeColorDescriptor
	PrimaryContainer      ThemeColorDescriptor
	OnPrimaryContainer    ThemeColorDescriptor
	PrimaryFixed          ThemeColorDescriptor
	OnPrimaryFixed        ThemeColorDescriptor
	PrimaryFixedVariant   ThemeColorDescriptor
	OnPrimaryFixedVariant ThemeColorDescriptor
}

func newPrimaryColorRoleDescriptors() PrimaryColorRoleDescriptors {
	return PrimaryColorRoleDescriptors{
		Primary:               Primary,
		OnPrimary:             OnPrimary,
		PrimaryContainer:      PrimaryContainer,
		OnPrimaryContainer:    OnPrimaryContainer,
		PrimaryFixed:          PrimaryFixed,
		OnPrimaryFixed:        OnPrimaryFixed,
		PrimaryFixedVariant:   PrimaryFixedVariant,
		OnPrimaryFixedVariant: OnPrimaryFixedVariant,
	}
}

type SecondaryColorRoleDescriptors struct {
	Secondary               ThemeColorDescriptor
	OnSecondary             ThemeColorDescriptor
	SecondaryContainer      ThemeColorDescriptor
	OnSecondaryContainer    ThemeColorDescriptor
	SecondaryFixed          ThemeColorDescriptor
	OnSecondaryFixed        ThemeColorDescriptor
	SecondaryFixedVariant   ThemeColorDescriptor
	OnSecondaryFixedVariant ThemeColorDescriptor
}

func newSecondaryColorRoleDescriptors() SecondaryColorRoleDescriptors {
	return SecondaryColorRoleDescriptors{
		Secondary:               Secondary,
		OnSecondary:             OnSecondary,
		SecondaryContainer:      SecondaryContainer,
		OnSecondaryContainer:    OnSecondaryContainer,
		SecondaryFixed:          SecondaryFixed,
		OnSecondaryFixed:        OnSecondaryFixed,
		SecondaryFixedVariant:   SecondaryFixedVariant,
		OnSecondaryFixedVariant: OnSecondaryFixedVariant,
	}
}

type TertiaryColorRoleDescriptors struct {
	Tertiary               ThemeColorDescriptor
	OnTertiary             ThemeColorDescriptor
	TertiaryContainer      ThemeColorDescriptor
	OnTertiaryContainer    ThemeColorDescriptor
	TertiaryFixed          ThemeColorDescriptor
	OnTertiaryFixed        ThemeColorDescriptor
	TertiaryFixedVariant   ThemeColorDescriptor
	OnTertiaryFixedVariant ThemeColorDescriptor
}

func newTertiaryColorRoleDescriptors() TertiaryColorRoleDescriptors {
	return TertiaryColorRoleDescriptors{
		Tertiary:               Tertiary,
		OnTertiary:             OnTertiary,
		TertiaryContainer:      TertiaryContainer,
		OnTertiaryContainer:    OnTertiaryContainer,
		TertiaryFixed:          TertiaryFixed,
		OnTertiaryFixed:        OnTertiaryFixed,
		TertiaryFixedVariant:   TertiaryFixedVariant,
		OnTertiaryFixedVariant: OnTertiaryFixedVariant,
	}
}

type SurfaceColorRoleDescriptors struct {
	Surface                 ThemeColorDescriptor
	OnSurface               ThemeColorDescriptor
	SurfaceVariant          ThemeColorDescriptor
	OnSurfaceVariant        ThemeColorDescriptor
	SurfaceDim              ThemeColorDescriptor
	SurfaceBright           ThemeColorDescriptor
	SurfaceContainerLowest  ThemeColorDescriptor
	SurfaceContainerLow     ThemeColorDescriptor
	SurfaceContainer        ThemeColorDescriptor
	SurfaceContainerHigh    ThemeColorDescriptor
	SurfaceContainerHighest ThemeColorDescriptor
}

func newSurfaceColorRoleDescriptors() SurfaceColorRoleDescriptors {
	return SurfaceColorRoleDescriptors{
		Surface:                 Surface,
		OnSurface:               OnSurface,
		SurfaceVariant:          SurfaceVariant,
		OnSurfaceVariant:        OnSurfaceVariant,
		SurfaceDim:              SurfaceDim,
		SurfaceBright:           SurfaceBright,
		SurfaceContainerLowest:  SurfaceContainerLowest,
		SurfaceContainerLow:     SurfaceContainerLow,
		SurfaceContainer:        SurfaceContainer,
		SurfaceContainerHigh:    SurfaceContainerHigh,
		SurfaceContainerHighest: SurfaceContainerHighest,
	}
}

type InverseColorRoleDescriptors struct {
	InverseSurface   ThemeColorDescriptor
	OnInverseSurface ThemeColorDescriptor
	InversePrimary   ThemeColorDescriptor
}

func newInverseColorRoleDescriptors() InverseColorRoleDescriptors {
	return InverseColorRoleDescriptors{
		InverseSurface:   InverseSurface,
		OnInverseSurface: OnInverseSurface,
		InversePrimary:   InversePrimary,
	}
}

type BackgroundColorRoleDescriptors struct {
	Background   ThemeColorDescriptor
	OnBackground ThemeColorDescriptor
}

func newBackgroundColorRoleDescriptors() BackgroundColorRoleDescriptors {
	return BackgroundColorRoleDescriptors{
		Background:   Background,
		OnBackground: OnBackground,
	}
}

type OutlineColorRoleDescriptors struct {
	Outline        ThemeColorDescriptor
	OutlineVariant ThemeColorDescriptor
}

func newOutlineColorRoleDescriptors() OutlineColorRoleDescriptors {
	return OutlineColorRoleDescriptors{
		Outline:        Outline,
		OutlineVariant: OutlineVariant,
	}
}

type ErrorColorRoleDescriptors struct {
	Error            ThemeColorDescriptor
	OnError          ThemeColorDescriptor
	ErrorContainer   ThemeColorDescriptor
	OnErrorContainer ThemeColorDescriptor
}

func newErrorColorRoleDescriptors() ErrorColorRoleDescriptors {
	return ErrorColorRoleDescriptors{
		Error:            Error,
		OnError:          OnError,
		ErrorContainer:   ErrorContainer,
		OnErrorContainer: OnErrorContainer,
	}
}

type ScrimColorRoleDescriptors struct {
	Scrim  ThemeColorDescriptor
	Shadow ThemeColorDescriptor
}

func newScrimColorRoleDescriptors() ScrimColorRoleDescriptors {
	return ScrimColorRoleDescriptors{
		Scrim:  Scrim,
		Shadow: Shadow,
	}
}

// Roles for Original Gioui Matrial Theme Palette
type BasicColorRoleDescriptors struct {
	BasicBg         ThemeColorDescriptor
	BasicFg         ThemeColorDescriptor
	BasicContrastBg ThemeColorDescriptor
	BasicContrastFg ThemeColorDescriptor
}

func newBasicColorRoleDescriptors() BasicColorRoleDescriptors {
	return BasicColorRoleDescriptors{
		BasicBg:         BasicBg,
		BasicFg:         BasicFg,
		BasicContrastBg: BasicContrastBg,
		BasicContrastFg: BasicContrastFg,
	}
}

type M3ColorRolesDescriptors struct {
	PrimaryRoles         PrimaryColorRoleDescriptors
	SecondaryRoles       SecondaryColorRoleDescriptors
	TertiaryRoles        TertiaryColorRoleDescriptors
	SurfaceRoles         SurfaceColorRoleDescriptors
	InverseRoles         InverseColorRoleDescriptors
	BackgroundColorRoles BackgroundColorRoleDescriptors
	OutlineRoles         OutlineColorRoleDescriptors
	ErrorRoles           ErrorColorRoleDescriptors
	ScrimRoles           ScrimColorRoleDescriptors
}

func newM3ColorRolesDescriptors() M3ColorRolesDescriptors {
	return M3ColorRolesDescriptors{
		PrimaryRoles:         newPrimaryColorRoleDescriptors(),
		SecondaryRoles:       newSecondaryColorRoleDescriptors(),
		TertiaryRoles:        newTertiaryColorRoleDescriptors(),
		SurfaceRoles:         newSurfaceColorRoleDescriptors(),
		InverseRoles:         newInverseColorRoleDescriptors(),
		BackgroundColorRoles: newBackgroundColorRoleDescriptors(),
		OutlineRoles:         newOutlineColorRoleDescriptors(),
		ErrorRoles:           newErrorColorRoleDescriptors(),
		ScrimRoles:           newScrimColorRoleDescriptors(),
	}
}

type ColorRoleDescriptors struct {
	M3ColorRolesDescriptors
	BasicColorRoleDescriptors
}

func NewColorRoleDescriptors() ColorRoleDescriptors {
	return ColorRoleDescriptors{
		M3ColorRolesDescriptors:   newM3ColorRolesDescriptors(),
		BasicColorRoleDescriptors: newBasicColorRoleDescriptors(),
	}
}
