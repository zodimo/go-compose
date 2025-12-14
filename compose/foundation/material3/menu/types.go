package menu

import (
	"github.com/zodimo/go-compose/theme"

	"git.sr.ht/~schnwalter/gio-mw/token"
)

type DropdownMenuColors struct {
	ContainerColor token.MatColor
}

type DropdownMenuItemColors struct {
	TextColor                 token.MatColor
	LeadingIconColor          token.MatColor
	TrailingIconColor         token.MatColor
	DisabledTextColor         token.MatColor
	DisabledLeadingIconColor  token.MatColor
	DisabledTrailingIconColor token.MatColor
}

func DefaultDropdownMenuColors() DropdownMenuColors {
	m3 := theme.GetThemeManager().GetMaterial3Theme()
	return DropdownMenuColors{
		ContainerColor: m3.Scheme.SurfaceContainer, // Elevation Level 2 default
	}
}

func DefaultDropdownMenuItemColors() DropdownMenuItemColors {
	m3 := theme.GetThemeManager().GetMaterial3Theme()
	return DropdownMenuItemColors{
		TextColor:                 m3.Scheme.Surface.OnColor,
		LeadingIconColor:          m3.Scheme.SurfaceVariant.OnColor,
		TrailingIconColor:         m3.Scheme.SurfaceVariant.OnColor,
		DisabledTextColor:         m3.Scheme.Surface.OnColor.SetOpacity(0.38),
		DisabledLeadingIconColor:  m3.Scheme.Surface.OnColor.SetOpacity(0.38),
		DisabledTrailingIconColor: m3.Scheme.Surface.OnColor.SetOpacity(0.38),
	}
}
