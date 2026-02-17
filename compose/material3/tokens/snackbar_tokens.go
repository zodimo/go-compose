package tokens

import "github.com/zodimo/go-compose/compose/ui/unit"

// package androidx.compose.material3.tokens

// import androidx.compose.ui.unit.dp

// internal object SnackbarTokens {
//     val ActionFocusLabelTextColor = ColorSchemeKeyTokens.InversePrimary
//     val ActionHoverLabelTextColor = ColorSchemeKeyTokens.InversePrimary
//     val ActionLabelTextColor = ColorSchemeKeyTokens.InversePrimary
//     val ActionLabelTextFont = TypographyKeyTokens.LabelLarge
//     val ActionPressedLabelTextColor = ColorSchemeKeyTokens.InversePrimary
//     val ContainerColor = ColorSchemeKeyTokens.InverseSurface
//     val ContainerElevation = ElevationTokens.Level3
//     val ContainerShape = ShapeKeyTokens.CornerExtraSmall
//     val IconColor = ColorSchemeKeyTokens.InverseOnSurface
//     val FocusIconColor = ColorSchemeKeyTokens.InverseOnSurface
//     val HoverIconColor = ColorSchemeKeyTokens.InverseOnSurface
//     val PressedIconColor = ColorSchemeKeyTokens.InverseOnSurface
//     val IconSize = 24.0.dp
//     val SupportingTextColor = ColorSchemeKeyTokens.InverseOnSurface
//     val SupportingTextFont = TypographyKeyTokens.BodyMedium
//     val SingleLineContainerHeight = 48.0.dp
//     val TwoLinesContainerHeight = 68.0.dp
// }

var SnackbarTokens = SnackbarTokensData{
	ActionFocusLabelTextColor:   ColorSchemeKeyTokens.InversePrimary,
	ActionHoverLabelTextColor:   ColorSchemeKeyTokens.InversePrimary,
	ActionLabelTextColor:        ColorSchemeKeyTokens.InversePrimary,
	ActionLabelTextFont:         TypographyKeyTokens.LabelLarge,
	ActionPressedLabelTextColor: ColorSchemeKeyTokens.InversePrimary,
	ContainerColor:              ColorSchemeKeyTokens.InverseSurface,
	ContainerElevation:          ElevationTokens.Level3,
	ContainerShape:              ShapeKeyTokens.CornerExtraSmall,
	IconColor:                   ColorSchemeKeyTokens.InverseOnSurface,
	FocusIconColor:              ColorSchemeKeyTokens.InverseOnSurface,
	HoverIconColor:              ColorSchemeKeyTokens.InverseOnSurface,
	PressedIconColor:            ColorSchemeKeyTokens.InverseOnSurface,
	IconSize:                    unit.Dp(24),
	SupportingTextColor:         ColorSchemeKeyTokens.InverseOnSurface,
	SupportingTextFont:          TypographyKeyTokens.BodyMedium,
	SingleLineContainerHeight:   unit.Dp(48),
	TwoLinesContainerHeight:     unit.Dp(68),
}

type SnackbarTokensData struct {
	ActionFocusLabelTextColor   ColorSchemeTokenKey
	ActionHoverLabelTextColor   ColorSchemeTokenKey
	ActionLabelTextColor        ColorSchemeTokenKey
	ActionLabelTextFont         TypographyTokenKey
	ActionPressedLabelTextColor ColorSchemeTokenKey
	ContainerColor              ColorSchemeTokenKey
	ContainerElevation          unit.Dp
	ContainerShape              ShapeTokenKey
	IconColor                   ColorSchemeTokenKey
	FocusIconColor              ColorSchemeTokenKey
	HoverIconColor              ColorSchemeTokenKey
	PressedIconColor            ColorSchemeTokenKey
	IconSize                    unit.Dp
	SupportingTextColor         ColorSchemeTokenKey
	SupportingTextFont          TypographyTokenKey
	SingleLineContainerHeight   unit.Dp
	TwoLinesContainerHeight     unit.Dp
}
