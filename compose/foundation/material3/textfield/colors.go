package textfield

import "github.com/zodimo/go-compose/theme"

type TextFieldColors struct {
	TextColor                   theme.ColorDescriptor
	DisabledTextColor           theme.ColorDescriptor
	CursorColor                 theme.ColorDescriptor
	ErrorCursorColor            theme.ColorDescriptor
	SelectionColor              theme.ColorDescriptor
	FocusedIndicatorColor       theme.ColorDescriptor
	UnfocusedIndicatorColor     theme.ColorDescriptor
	DisabledIndicatorColor      theme.ColorDescriptor
	ErrorIndicatorColor         theme.ColorDescriptor
	HoveredIndicatorColor       theme.ColorDescriptor
	LeadingIconColor            theme.ColorDescriptor
	TrailingIconColor           theme.ColorDescriptor
	DisabledLeadingIconColor    theme.ColorDescriptor
	DisabledTrailingIconColor   theme.ColorDescriptor
	LabelColor                  theme.ColorDescriptor
	UnfocusedLabelColor         theme.ColorDescriptor
	DisabledLabelColor          theme.ColorDescriptor
	ErrorLabelColor             theme.ColorDescriptor
	PlaceholderColor            theme.ColorDescriptor
	DisabledPlaceholderColor    theme.ColorDescriptor
	SupportingTextColor         theme.ColorDescriptor
	DisabledSupportingTextColor theme.ColorDescriptor
	ErrorSupportingTextColor    theme.ColorDescriptor
}

func DefaultTextFieldColors() TextFieldColors {
	return TextFieldColors{
		TextColor:                   theme.ColorUnspecified,
		DisabledTextColor:           theme.ColorUnspecified,
		CursorColor:                 theme.ColorUnspecified,
		ErrorCursorColor:            theme.ColorUnspecified,
		SelectionColor:              theme.ColorUnspecified,
		FocusedIndicatorColor:       theme.ColorUnspecified, // Active
		UnfocusedIndicatorColor:     theme.ColorUnspecified, // Inactive
		DisabledIndicatorColor:      theme.ColorUnspecified,
		ErrorIndicatorColor:         theme.ColorUnspecified,
		HoveredIndicatorColor:       theme.ColorUnspecified, // Actually M3 says OnSurface for Outline variant hovered? Or just Outline token. Let's use OnSurface for high contrast hover.
		LeadingIconColor:            theme.ColorUnspecified,
		TrailingIconColor:           theme.ColorUnspecified,
		DisabledLeadingIconColor:    theme.ColorUnspecified,
		DisabledTrailingIconColor:   theme.ColorUnspecified,
		LabelColor:                  theme.ColorUnspecified, // Focused Label
		UnfocusedLabelColor:         theme.ColorUnspecified,
		DisabledLabelColor:          theme.ColorUnspecified,
		ErrorLabelColor:             theme.ColorUnspecified,
		PlaceholderColor:            theme.ColorUnspecified,
		DisabledPlaceholderColor:    theme.ColorUnspecified,
		SupportingTextColor:         theme.ColorUnspecified,
		DisabledSupportingTextColor: theme.ColorUnspecified,
		ErrorSupportingTextColor:    theme.ColorUnspecified,
	}
}

func ResolveTextFieldColors(colors TextFieldColors) TextFieldColors {
	selector := theme.ColorHelper.ColorSelector()

	return TextFieldColors{
		TextColor:                   theme.TakeOrElseColor(colors.TextColor, selector.SurfaceRoles.OnSurface),
		DisabledTextColor:           theme.TakeOrElseColor(colors.DisabledTextColor, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		CursorColor:                 theme.TakeOrElseColor(colors.CursorColor, selector.PrimaryRoles.Primary),
		ErrorCursorColor:            theme.TakeOrElseColor(colors.ErrorCursorColor, selector.ErrorRoles.Error),
		SelectionColor:              theme.TakeOrElseColor(colors.SelectionColor, selector.PrimaryRoles.Primary),
		FocusedIndicatorColor:       theme.TakeOrElseColor(colors.FocusedIndicatorColor, selector.PrimaryRoles.Primary),   // Active
		UnfocusedIndicatorColor:     theme.TakeOrElseColor(colors.UnfocusedIndicatorColor, selector.OutlineRoles.Outline), // Inactive
		DisabledIndicatorColor:      theme.TakeOrElseColor(colors.DisabledIndicatorColor, selector.SurfaceRoles.OnSurface.SetOpacity(0.12)),
		ErrorIndicatorColor:         theme.TakeOrElseColor(colors.ErrorIndicatorColor, selector.ErrorRoles.Error),
		HoveredIndicatorColor:       theme.TakeOrElseColor(colors.HoveredIndicatorColor, selector.SurfaceRoles.OnSurface), // Actually M3 says OnSurface for Outline variant hovered? Or just Outline token. Let's use OnSurface for high contrast hover.
		LeadingIconColor:            theme.TakeOrElseColor(colors.LeadingIconColor, selector.SurfaceRoles.OnVariant),
		TrailingIconColor:           theme.TakeOrElseColor(colors.TrailingIconColor, selector.SurfaceRoles.OnVariant),
		DisabledLeadingIconColor:    theme.TakeOrElseColor(colors.DisabledLeadingIconColor, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		DisabledTrailingIconColor:   theme.TakeOrElseColor(colors.DisabledTrailingIconColor, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		LabelColor:                  theme.TakeOrElseColor(colors.LabelColor, selector.PrimaryRoles.Primary), // Focused Label
		UnfocusedLabelColor:         theme.TakeOrElseColor(colors.UnfocusedLabelColor, selector.SurfaceRoles.OnVariant),
		DisabledLabelColor:          theme.TakeOrElseColor(colors.DisabledLabelColor, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		ErrorLabelColor:             theme.TakeOrElseColor(colors.ErrorLabelColor, selector.ErrorRoles.Error),
		PlaceholderColor:            theme.TakeOrElseColor(colors.PlaceholderColor, selector.SurfaceRoles.OnVariant),
		DisabledPlaceholderColor:    theme.TakeOrElseColor(colors.DisabledPlaceholderColor, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		SupportingTextColor:         theme.TakeOrElseColor(colors.SupportingTextColor, selector.SurfaceRoles.OnVariant),
		DisabledSupportingTextColor: theme.TakeOrElseColor(colors.DisabledSupportingTextColor, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		ErrorSupportingTextColor:    theme.TakeOrElseColor(colors.ErrorSupportingTextColor, selector.ErrorRoles.Error),
	}
}
