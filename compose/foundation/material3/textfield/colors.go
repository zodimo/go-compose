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
	selector := theme.ColorHelper.ColorSelector()
	return TextFieldColors{
		TextColor:                   selector.SurfaceRoles.OnSurface,
		DisabledTextColor:           selector.SurfaceRoles.OnSurface.SetOpacity(0.38),
		CursorColor:                 selector.PrimaryRoles.Primary,
		ErrorCursorColor:            selector.ErrorRoles.Error,
		SelectionColor:              selector.PrimaryRoles.Primary,
		FocusedIndicatorColor:       selector.PrimaryRoles.Primary, // Active
		UnfocusedIndicatorColor:     selector.OutlineRoles.Outline, // Inactive
		DisabledIndicatorColor:      selector.SurfaceRoles.OnSurface.SetOpacity(0.12),
		ErrorIndicatorColor:         selector.ErrorRoles.Error,
		HoveredIndicatorColor:       selector.SurfaceRoles.OnSurface, // Actually M3 says OnSurface for Outline variant hovered? Or just Outline token. Let's use OnSurface for high contrast hover.
		LeadingIconColor:            selector.SurfaceRoles.OnVariant,
		TrailingIconColor:           selector.SurfaceRoles.OnVariant,
		DisabledLeadingIconColor:    selector.SurfaceRoles.OnSurface.SetOpacity(0.38),
		DisabledTrailingIconColor:   selector.SurfaceRoles.OnSurface.SetOpacity(0.38),
		LabelColor:                  selector.PrimaryRoles.Primary, // Focused Label
		UnfocusedLabelColor:         selector.SurfaceRoles.OnVariant,
		DisabledLabelColor:          selector.SurfaceRoles.OnSurface.SetOpacity(0.38),
		ErrorLabelColor:             selector.ErrorRoles.Error,
		PlaceholderColor:            selector.SurfaceRoles.OnVariant,
		DisabledPlaceholderColor:    selector.SurfaceRoles.OnSurface.SetOpacity(0.38),
		SupportingTextColor:         selector.SurfaceRoles.OnVariant,
		DisabledSupportingTextColor: selector.SurfaceRoles.OnSurface.SetOpacity(0.38),
		ErrorSupportingTextColor:    selector.ErrorRoles.Error,
	}
}
