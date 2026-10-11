package textfield

import (
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui/graphics"
)

type TextFieldColors struct {
	TextColor                   graphics.Color
	DisabledTextColor           graphics.Color
	CursorColor                 graphics.Color
	ErrorCursorColor            graphics.Color
	SelectionColor              graphics.Color
	FocusedIndicatorColor       graphics.Color
	UnfocusedIndicatorColor     graphics.Color
	DisabledIndicatorColor      graphics.Color
	ErrorIndicatorColor         graphics.Color
	HoveredIndicatorColor       graphics.Color
	LeadingIconColor            graphics.Color
	TrailingIconColor           graphics.Color
	DisabledLeadingIconColor    graphics.Color
	DisabledTrailingIconColor   graphics.Color
	LabelColor                  graphics.Color
	UnfocusedLabelColor         graphics.Color
	DisabledLabelColor          graphics.Color
	ErrorLabelColor             graphics.Color
	PlaceholderColor            graphics.Color
	DisabledPlaceholderColor    graphics.Color
	SupportingTextColor         graphics.Color
	DisabledSupportingTextColor graphics.Color
	ErrorSupportingTextColor    graphics.Color
}

func DefaultTextFieldColors() TextFieldColors {
	return TextFieldColors{
		TextColor:                   graphics.ColorCompanion.Unspecified(),
		DisabledTextColor:           graphics.ColorCompanion.Unspecified(),
		CursorColor:                 graphics.ColorCompanion.Unspecified(),
		ErrorCursorColor:            graphics.ColorCompanion.Unspecified(),
		SelectionColor:              graphics.ColorCompanion.Unspecified(),
		FocusedIndicatorColor:       graphics.ColorCompanion.Unspecified(), // Active
		UnfocusedIndicatorColor:     graphics.ColorCompanion.Unspecified(), // Inactive
		DisabledIndicatorColor:      graphics.ColorCompanion.Unspecified(),
		ErrorIndicatorColor:         graphics.ColorCompanion.Unspecified(),
		HoveredIndicatorColor:       graphics.ColorCompanion.Unspecified(), // Actually M3 says OnSurface for Outline variant hovered? Or just Outline token. Let's use OnSurface for high contrast hover.
		LeadingIconColor:            graphics.ColorCompanion.Unspecified(),
		TrailingIconColor:           graphics.ColorCompanion.Unspecified(),
		DisabledLeadingIconColor:    graphics.ColorCompanion.Unspecified(),
		DisabledTrailingIconColor:   graphics.ColorCompanion.Unspecified(),
		LabelColor:                  graphics.ColorCompanion.Unspecified(), // Focused Label
		UnfocusedLabelColor:         graphics.ColorCompanion.Unspecified(),
		DisabledLabelColor:          graphics.ColorCompanion.Unspecified(),
		ErrorLabelColor:             graphics.ColorCompanion.Unspecified(),
		PlaceholderColor:            graphics.ColorCompanion.Unspecified(),
		DisabledPlaceholderColor:    graphics.ColorCompanion.Unspecified(),
		SupportingTextColor:         graphics.ColorCompanion.Unspecified(),
		DisabledSupportingTextColor: graphics.ColorCompanion.Unspecified(),
		ErrorSupportingTextColor:    graphics.ColorCompanion.Unspecified(),
	}
}

func ResolveTextFieldColors(c Composer, colors TextFieldColors) TextFieldColors {

	theme := material3.Theme(c)
	return TextFieldColors{
		TextColor:                   graphics.ColorCompanion.TakeOrElse(colors.TextColor, theme.ColorScheme().OnSurface),                                              // selector.SurfaceRoles.OnSurface),
		DisabledTextColor:           graphics.ColorCompanion.TakeOrElse(colors.DisabledTextColor, graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.38)),           //, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		CursorColor:                 graphics.ColorCompanion.TakeOrElse(colors.CursorColor, theme.ColorScheme().Primary),                                              //, selector.PrimaryRoles.Primary),
		ErrorCursorColor:            graphics.ColorCompanion.TakeOrElse(colors.ErrorCursorColor, theme.ColorScheme().Error),                                           //, selector.ErrorRoles.Error),
		SelectionColor:              graphics.ColorCompanion.TakeOrElse(colors.SelectionColor, theme.ColorScheme().Primary),                                           //, selector.PrimaryRoles.Primary),
		FocusedIndicatorColor:       graphics.ColorCompanion.TakeOrElse(colors.FocusedIndicatorColor, theme.ColorScheme().Primary),                                    //, selector.PrimaryRoles.Primary),   // Active
		UnfocusedIndicatorColor:     graphics.ColorCompanion.TakeOrElse(colors.UnfocusedIndicatorColor, theme.ColorScheme().Outline),                                  //, selector.OutlineRoles.Outline), // Inactive
		DisabledIndicatorColor:      graphics.ColorCompanion.TakeOrElse(colors.DisabledIndicatorColor, graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.12)),      //, selector.SurfaceRoles.OnSurface.SetOpacity(0.12)),
		ErrorIndicatorColor:         graphics.ColorCompanion.TakeOrElse(colors.ErrorIndicatorColor, theme.ColorScheme().Error),                                        //, selector.ErrorRoles.Error),
		HoveredIndicatorColor:       graphics.ColorCompanion.TakeOrElse(colors.HoveredIndicatorColor, theme.ColorScheme().OnSurface),                                  //, selector.SurfaceRoles.OnSurface), // Actually M3 says OnSurface for Outline variant hovered? Or just Outline token. Let's use OnSurface for high contrast hover.
		LeadingIconColor:            graphics.ColorCompanion.TakeOrElse(colors.LeadingIconColor, theme.ColorScheme().OnSurfaceVariant),                                //, selector.SurfaceRoles.OnVariant),
		TrailingIconColor:           graphics.ColorCompanion.TakeOrElse(colors.TrailingIconColor, theme.ColorScheme().OnSurfaceVariant),                               //, selector.SurfaceRoles.OnVariant),
		DisabledLeadingIconColor:    graphics.ColorCompanion.TakeOrElse(colors.DisabledLeadingIconColor, graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.38)),    //, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		DisabledTrailingIconColor:   graphics.ColorCompanion.TakeOrElse(colors.DisabledTrailingIconColor, graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.38)),   //, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		LabelColor:                  graphics.ColorCompanion.TakeOrElse(colors.LabelColor, theme.ColorScheme().Primary),                                               //, selector.PrimaryRoles.Primary), // Focused Label
		UnfocusedLabelColor:         graphics.ColorCompanion.TakeOrElse(colors.UnfocusedLabelColor, theme.ColorScheme().OnSurfaceVariant),                             //, selector.SurfaceRoles.OnVariant),
		DisabledLabelColor:          graphics.ColorCompanion.TakeOrElse(colors.DisabledLabelColor, graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.38)),          //, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		ErrorLabelColor:             graphics.ColorCompanion.TakeOrElse(colors.ErrorLabelColor, theme.ColorScheme().Error),                                            //, selector.ErrorRoles.Error),
		PlaceholderColor:            graphics.ColorCompanion.TakeOrElse(colors.PlaceholderColor, theme.ColorScheme().OnSurfaceVariant),                                //, selector.SurfaceRoles.OnVariant),
		DisabledPlaceholderColor:    graphics.ColorCompanion.TakeOrElse(colors.DisabledPlaceholderColor, graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.38)),    //, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		SupportingTextColor:         graphics.ColorCompanion.TakeOrElse(colors.SupportingTextColor, theme.ColorScheme().OnSurfaceVariant),                             //, selector.SurfaceRoles.OnVariant),
		DisabledSupportingTextColor: graphics.ColorCompanion.TakeOrElse(colors.DisabledSupportingTextColor, graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.38)), //, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		ErrorSupportingTextColor:    graphics.ColorCompanion.TakeOrElse(colors.ErrorSupportingTextColor, theme.ColorScheme().Error),                                   //, selector.ErrorRoles.Error),
	}
}
