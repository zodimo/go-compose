package theme

type PrimaryColorRoles interface {
	Primary() ColorRole
	OnPrimary() ColorRole
	PrimaryContainer() ColorRole
	OnPrimaryContainer() ColorRole
	PrimaryFixed() ColorRole
	OnPrimaryFixed() ColorRole
	PrimaryFixedVariant() ColorRole
	OnPrimaryFixedVariant() ColorRole
}

type SecondaryColorRoles interface {
	Secondary() ColorRole
	OnSecondary() ColorRole
	SecondaryContainer() ColorRole
	OnSecondaryContainer() ColorRole
	SecondaryFixed() ColorRole
	OnSecondaryFixed() ColorRole
	SecondaryFixedVariant() ColorRole
	OnSecondaryFixedVariant() ColorRole
}

type TertiaryColorRoles interface {
	Tertiary() ColorRole
	OnTertiary() ColorRole
	TertiaryContainer() ColorRole
	OnTertiaryContainer() ColorRole
	TertiaryFixed() ColorRole
	OnTertiaryFixed() ColorRole
	TertiaryFixedVariant() ColorRole
	OnTertiaryFixedVariant() ColorRole
}
type SurfaceColorRoles interface {
	Surface() ColorRole
	OnSurface() ColorRole
	SurfaceVariant() ColorRole
	OnSurfaceVariant() ColorRole
	SurfaceDim() ColorRole
	SurfaceBright() ColorRole
	SurfaceContainerLowest() ColorRole
	SurfaceContainerLow() ColorRole
	SurfaceContainer() ColorRole
	SurfaceContainerHigh() ColorRole
	SurfaceContainerHighest() ColorRole
}

type InverseColorRoles interface {
	InverseSurface() ColorRole
	OnInverseSurface() ColorRole
	InversePrimary() ColorRole
}

type BackgroundColorRoles interface {
	Background() ColorRole
	OnBackground() ColorRole
}

type OutlineColorRoles interface {
	Outline() ColorRole
	OutlineVariant() ColorRole
}

type ErrorColorRoles interface {
	Error() ColorRole
	OnError() ColorRole
	ErrorContainer() ColorRole
	OnErrorContainer() ColorRole
}
type ScrimColorRoles interface {
	Scrim() ColorRole
	Shadow() ColorRole
}

// Roles for Original Gioui Matrial Theme Palette
type BasicColorRoles interface {
	BasicBg() ColorRole
	BasicFg() ColorRole
	BasicContrastBg() ColorRole
	BasicContrastFg() ColorRole
}

type M3ColorRoles interface {
	PrimaryRoles() PrimaryColorRoles
	SecondaryRoles() SecondaryColorRoles
	TertiaryRoles() TertiaryColorRoles
	SurfaceRoles() SurfaceColorRoles
	InverseRoles() InverseColorRoles
	BackgroundColorRoles() BackgroundColorRoles
	OutlineRoles() OutlineColorRoles
	ErrorRoles() ErrorColorRoles
	ScrimRoles() ScrimColorRoles
}

type ColorRoles interface {
	M3ColorRoles
	BasicColorRoles
}
