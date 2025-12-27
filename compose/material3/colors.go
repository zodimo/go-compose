package material3

import "github.com/zodimo/go-compose/theme"

func ContentColorFor(backgroundColor theme.ColorDescriptor) theme.ColorDescriptor {
	switch {
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().PrimaryRoles.Primary):
		return theme.ColorHelper.ColorSelector().PrimaryRoles.OnPrimary
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().SecondaryRoles.Secondary):
		return theme.ColorHelper.ColorSelector().SecondaryRoles.OnSecondary
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().TertiaryRoles.Tertiary):
		return theme.ColorHelper.ColorSelector().TertiaryRoles.OnTertiary
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().BackgroundColorRoles.Background):
		return theme.ColorHelper.ColorSelector().BackgroundColorRoles.OnBackground
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().ErrorRoles.Error):
		return theme.ColorHelper.ColorSelector().ErrorRoles.OnError
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().PrimaryRoles.Container):
		return theme.ColorHelper.ColorSelector().PrimaryRoles.OnContainer
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().SecondaryRoles.Container):
		return theme.ColorHelper.ColorSelector().SecondaryRoles.OnContainer
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().TertiaryRoles.Container):
		return theme.ColorHelper.ColorSelector().TertiaryRoles.OnContainer
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().ErrorRoles.Container):
		return theme.ColorHelper.ColorSelector().ErrorRoles.OnContainer
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().InverseRoles.Surface):
		return theme.ColorHelper.ColorSelector().InverseRoles.OnSurface
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().SurfaceRoles.Surface):
		return theme.ColorHelper.ColorSelector().SurfaceRoles.OnSurface
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().SurfaceRoles.Variant):
		return theme.ColorHelper.ColorSelector().SurfaceRoles.OnVariant
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().SurfaceRoles.Bright):
		return theme.ColorHelper.ColorSelector().SurfaceRoles.OnSurface
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().SurfaceRoles.Container):
		return theme.ColorHelper.ColorSelector().SurfaceRoles.OnSurface
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().SurfaceRoles.ContainerHigh):
		return theme.ColorHelper.ColorSelector().SurfaceRoles.OnSurface
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().SurfaceRoles.ContainerHighest):
		return theme.ColorHelper.ColorSelector().SurfaceRoles.OnSurface
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().SurfaceRoles.ContainerLow):
		return theme.ColorHelper.ColorSelector().SurfaceRoles.OnSurface
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().SurfaceRoles.ContainerLowest):
		return theme.ColorHelper.ColorSelector().SurfaceRoles.OnSurface
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().SurfaceRoles.Dim):
		return theme.ColorHelper.ColorSelector().SurfaceRoles.OnSurface
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().PrimaryRoles.Fixed):
		return theme.ColorHelper.ColorSelector().PrimaryRoles.OnFixed
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().SecondaryRoles.Fixed):
		return theme.ColorHelper.ColorSelector().SecondaryRoles.OnFixed
	case backgroundColor.Compare(theme.ColorHelper.ColorSelector().TertiaryRoles.Fixed):
		return theme.ColorHelper.ColorSelector().TertiaryRoles.OnFixed
	default:
		return theme.ColorHelper.UnspecifiedColor()
	}
}
