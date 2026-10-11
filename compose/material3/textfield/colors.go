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
	FocusedLabelColor           graphics.Color // Added
	DisabledLabelColor          graphics.Color
	ErrorLabelColor             graphics.Color
	PlaceholderColor            graphics.Color
	DisabledPlaceholderColor    graphics.Color
	SupportingTextColor         graphics.Color
	DisabledSupportingTextColor graphics.Color
	ErrorSupportingTextColor    graphics.Color
	ContainerColor              graphics.Color // Added
	DisabledContainerColor      graphics.Color // Added
}

func DefaultTextFieldColors() TextFieldColors {
	return TextFieldColors{
		TextColor:                   graphics.ColorCompanion.Unspecified(),
		DisabledTextColor:           graphics.ColorCompanion.Unspecified(),
		CursorColor:                 graphics.ColorCompanion.Unspecified(),
		ErrorCursorColor:            graphics.ColorCompanion.Unspecified(),
		SelectionColor:              graphics.ColorCompanion.Unspecified(),
		FocusedIndicatorColor:       graphics.ColorCompanion.Unspecified(),
		UnfocusedIndicatorColor:     graphics.ColorCompanion.Unspecified(),
		DisabledIndicatorColor:      graphics.ColorCompanion.Unspecified(),
		ErrorIndicatorColor:         graphics.ColorCompanion.Unspecified(),
		HoveredIndicatorColor:       graphics.ColorCompanion.Unspecified(),
		LeadingIconColor:            graphics.ColorCompanion.Unspecified(),
		TrailingIconColor:           graphics.ColorCompanion.Unspecified(),
		DisabledLeadingIconColor:    graphics.ColorCompanion.Unspecified(),
		DisabledTrailingIconColor:   graphics.ColorCompanion.Unspecified(),
		LabelColor:                  graphics.ColorCompanion.Unspecified(),
		UnfocusedLabelColor:         graphics.ColorCompanion.Unspecified(),
		FocusedLabelColor:           graphics.ColorCompanion.Unspecified(), // Added
		DisabledLabelColor:          graphics.ColorCompanion.Unspecified(),
		ErrorLabelColor:             graphics.ColorCompanion.Unspecified(),
		PlaceholderColor:            graphics.ColorCompanion.Unspecified(),
		DisabledPlaceholderColor:    graphics.ColorCompanion.Unspecified(),
		SupportingTextColor:         graphics.ColorCompanion.Unspecified(),
		DisabledSupportingTextColor: graphics.ColorCompanion.Unspecified(),
		ErrorSupportingTextColor:    graphics.ColorCompanion.Unspecified(),
		ContainerColor:              graphics.ColorCompanion.Unspecified(), // Added
		DisabledContainerColor:      graphics.ColorCompanion.Unspecified(), // Added
	}
}

func ResolveTextFieldColors(c Composer, colors TextFieldColors) TextFieldColors {

	theme := material3.Theme(c)
	return TextFieldColors{
		TextColor:                   graphics.ColorCompanion.TakeOrElse(colors.TextColor, theme.ColorScheme().OnSurface),
		DisabledTextColor:           graphics.ColorCompanion.TakeOrElse(colors.DisabledTextColor, graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.38)),
		CursorColor:                 graphics.ColorCompanion.TakeOrElse(colors.CursorColor, theme.ColorScheme().Primary),
		ErrorCursorColor:            graphics.ColorCompanion.TakeOrElse(colors.ErrorCursorColor, theme.ColorScheme().Error),
		SelectionColor:              graphics.ColorCompanion.TakeOrElse(colors.SelectionColor, theme.ColorScheme().Primary),
		FocusedIndicatorColor:       graphics.ColorCompanion.TakeOrElse(colors.FocusedIndicatorColor, theme.ColorScheme().Primary),
		UnfocusedIndicatorColor:     graphics.ColorCompanion.TakeOrElse(colors.UnfocusedIndicatorColor, theme.ColorScheme().Outline),
		DisabledIndicatorColor:      graphics.ColorCompanion.TakeOrElse(colors.DisabledIndicatorColor, graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.12)),
		ErrorIndicatorColor:         graphics.ColorCompanion.TakeOrElse(colors.ErrorIndicatorColor, theme.ColorScheme().Error),
		HoveredIndicatorColor:       graphics.ColorCompanion.TakeOrElse(colors.HoveredIndicatorColor, theme.ColorScheme().OnSurface),
		LeadingIconColor:            graphics.ColorCompanion.TakeOrElse(colors.LeadingIconColor, theme.ColorScheme().OnSurfaceVariant),
		TrailingIconColor:           graphics.ColorCompanion.TakeOrElse(colors.TrailingIconColor, theme.ColorScheme().OnSurfaceVariant),
		DisabledLeadingIconColor:    graphics.ColorCompanion.TakeOrElse(colors.DisabledLeadingIconColor, graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.38)),
		DisabledTrailingIconColor:   graphics.ColorCompanion.TakeOrElse(colors.DisabledTrailingIconColor, graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.38)),
		LabelColor:                  graphics.ColorCompanion.TakeOrElse(colors.LabelColor, theme.ColorScheme().OnSurfaceVariant),
		UnfocusedLabelColor:         graphics.ColorCompanion.TakeOrElse(colors.UnfocusedLabelColor, theme.ColorScheme().OnSurfaceVariant),
		FocusedLabelColor:           graphics.ColorCompanion.TakeOrElse(colors.FocusedLabelColor, theme.ColorScheme().Primary), // Added
		DisabledLabelColor:          graphics.ColorCompanion.TakeOrElse(colors.DisabledLabelColor, graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.38)),
		ErrorLabelColor:             graphics.ColorCompanion.TakeOrElse(colors.ErrorLabelColor, theme.ColorScheme().Error),
		PlaceholderColor:            graphics.ColorCompanion.TakeOrElse(colors.PlaceholderColor, theme.ColorScheme().OnSurfaceVariant),
		DisabledPlaceholderColor:    graphics.ColorCompanion.TakeOrElse(colors.DisabledPlaceholderColor, graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.38)),
		SupportingTextColor:         graphics.ColorCompanion.TakeOrElse(colors.SupportingTextColor, theme.ColorScheme().OnSurfaceVariant),
		DisabledSupportingTextColor: graphics.ColorCompanion.TakeOrElse(colors.DisabledSupportingTextColor, graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.38)),
		ErrorSupportingTextColor:    graphics.ColorCompanion.TakeOrElse(colors.ErrorSupportingTextColor, theme.ColorScheme().Error),
		ContainerColor:              graphics.ColorCompanion.TakeOrElse(colors.ContainerColor, theme.ColorScheme().SurfaceVariant),                               // Added - default for filled
		DisabledContainerColor:      graphics.ColorCompanion.TakeOrElse(colors.DisabledContainerColor, graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.12)), // Added - approximate default? Or SurfaceVariant with opacity
	}
}
