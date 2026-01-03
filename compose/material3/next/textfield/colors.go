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
		TextColor:                   graphics.ColorUnspecified,
		DisabledTextColor:           graphics.ColorUnspecified,
		CursorColor:                 graphics.ColorUnspecified,
		ErrorCursorColor:            graphics.ColorUnspecified,
		SelectionColor:              graphics.ColorUnspecified,
		FocusedIndicatorColor:       graphics.ColorUnspecified, // Active
		UnfocusedIndicatorColor:     graphics.ColorUnspecified, // Inactive
		DisabledIndicatorColor:      graphics.ColorUnspecified,
		ErrorIndicatorColor:         graphics.ColorUnspecified,
		HoveredIndicatorColor:       graphics.ColorUnspecified, // Actually M3 says OnSurface for Outline variant hovered? Or just Outline token. Let's use OnSurface for high contrast hover.
		LeadingIconColor:            graphics.ColorUnspecified,
		TrailingIconColor:           graphics.ColorUnspecified,
		DisabledLeadingIconColor:    graphics.ColorUnspecified,
		DisabledTrailingIconColor:   graphics.ColorUnspecified,
		LabelColor:                  graphics.ColorUnspecified, // Focused Label
		UnfocusedLabelColor:         graphics.ColorUnspecified,
		DisabledLabelColor:          graphics.ColorUnspecified,
		ErrorLabelColor:             graphics.ColorUnspecified,
		PlaceholderColor:            graphics.ColorUnspecified,
		DisabledPlaceholderColor:    graphics.ColorUnspecified,
		SupportingTextColor:         graphics.ColorUnspecified,
		DisabledSupportingTextColor: graphics.ColorUnspecified,
		ErrorSupportingTextColor:    graphics.ColorUnspecified,
	}
}

func ResolveTextFieldColors(c Composer, colors TextFieldColors) TextFieldColors {

	theme := material3.Theme(c)
	return TextFieldColors{
		TextColor:                   colors.TextColor.TakeOrElse(theme.ColorScheme().Surface.OnColor),                                              // selector.SurfaceRoles.OnSurface),
		DisabledTextColor:           colors.DisabledTextColor.TakeOrElse(graphics.SetOpacity(theme.ColorScheme().Surface.OnColor, 0.38)),           //, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		CursorColor:                 colors.CursorColor.TakeOrElse(theme.ColorScheme().Primary.Color),                                              //, selector.PrimaryRoles.Primary),
		ErrorCursorColor:            colors.ErrorCursorColor.TakeOrElse(theme.ColorScheme().Error.Color),                                           //, selector.ErrorRoles.Error),
		SelectionColor:              colors.SelectionColor.TakeOrElse(theme.ColorScheme().Primary.Color),                                           //, selector.PrimaryRoles.Primary),
		FocusedIndicatorColor:       colors.FocusedIndicatorColor.TakeOrElse(theme.ColorScheme().Primary.Color),                                    //, selector.PrimaryRoles.Primary),   // Active
		UnfocusedIndicatorColor:     colors.UnfocusedIndicatorColor.TakeOrElse(theme.ColorScheme().Outline),                                        //, selector.OutlineRoles.Outline), // Inactive
		DisabledIndicatorColor:      colors.DisabledIndicatorColor.TakeOrElse(graphics.SetOpacity(theme.ColorScheme().Surface.OnColor, 0.12)),      //, selector.SurfaceRoles.OnSurface.SetOpacity(0.12)),
		ErrorIndicatorColor:         colors.ErrorIndicatorColor.TakeOrElse(theme.ColorScheme().Error.Color),                                        //, selector.ErrorRoles.Error),
		HoveredIndicatorColor:       colors.HoveredIndicatorColor.TakeOrElse(theme.ColorScheme().Surface.OnColor),                                  //, selector.SurfaceRoles.OnSurface), // Actually M3 says OnSurface for Outline variant hovered? Or just Outline token. Let's use OnSurface for high contrast hover.
		LeadingIconColor:            colors.LeadingIconColor.TakeOrElse(theme.ColorScheme().SurfaceVariant.OnColor),                                //, selector.SurfaceRoles.OnVariant),
		TrailingIconColor:           colors.TrailingIconColor.TakeOrElse(theme.ColorScheme().SurfaceVariant.OnColor),                               //, selector.SurfaceRoles.OnVariant),
		DisabledLeadingIconColor:    colors.DisabledLeadingIconColor.TakeOrElse(graphics.SetOpacity(theme.ColorScheme().Surface.OnColor, 0.38)),    //, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		DisabledTrailingIconColor:   colors.DisabledTrailingIconColor.TakeOrElse(graphics.SetOpacity(theme.ColorScheme().Surface.OnColor, 0.38)),   //, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		LabelColor:                  colors.LabelColor.TakeOrElse(theme.ColorScheme().Primary.Color),                                               //, selector.PrimaryRoles.Primary), // Focused Label
		UnfocusedLabelColor:         colors.UnfocusedLabelColor.TakeOrElse(theme.ColorScheme().SurfaceVariant.OnColor),                             //, selector.SurfaceRoles.OnVariant),
		DisabledLabelColor:          colors.DisabledLabelColor.TakeOrElse(graphics.SetOpacity(theme.ColorScheme().Surface.OnColor, 0.38)),          //, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		ErrorLabelColor:             colors.ErrorLabelColor.TakeOrElse(theme.ColorScheme().Error.Color),                                            //, selector.ErrorRoles.Error),
		PlaceholderColor:            colors.PlaceholderColor.TakeOrElse(theme.ColorScheme().SurfaceVariant.OnColor),                                //, selector.SurfaceRoles.OnVariant),
		DisabledPlaceholderColor:    colors.DisabledPlaceholderColor.TakeOrElse(graphics.SetOpacity(theme.ColorScheme().Surface.OnColor, 0.38)),    //, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		SupportingTextColor:         colors.SupportingTextColor.TakeOrElse(theme.ColorScheme().SurfaceVariant.OnColor),                             //, selector.SurfaceRoles.OnVariant),
		DisabledSupportingTextColor: colors.DisabledSupportingTextColor.TakeOrElse(graphics.SetOpacity(theme.ColorScheme().Surface.OnColor, 0.38)), //, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		ErrorSupportingTextColor:    colors.ErrorSupportingTextColor.TakeOrElse(theme.ColorScheme().Error.Color),                                   //, selector.ErrorRoles.Error),
	}
}
