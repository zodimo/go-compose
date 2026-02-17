package snackbar

import (
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/material3/tokens"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
)

// /** Contains the default values used for [Snackbar]. */
// object SnackbarDefaults {
//     /** Default shape of a snackbar. */
//     val shape: Shape
//         @Composable get() = SnackbarTokens.ContainerShape.value

//     /** Default color of a snackbar. */
//     val color: Color
//         @Composable get() = SnackbarTokens.ContainerColor.value

//     /** Default content color of a snackbar. */
//     val contentColor: Color
//         @Composable get() = SnackbarTokens.SupportingTextColor.value

//     /** Default action color of a snackbar. */
//     val actionColor: Color
//         @Composable get() = SnackbarTokens.ActionLabelTextColor.value

//     /** Default action content color of a snackbar. */
//     val actionContentColor: Color
//         @Composable get() = SnackbarTokens.ActionLabelTextColor.value

//     /** Default dismiss action content color of a snackbar. */
//     val dismissActionContentColor: Color
//         @Composable get() = SnackbarTokens.IconColor.value
// }

func SnackbarDefaults(theme material3.ThemeInterface) SnackbarDefaultsData {
	colorScheme := theme.ColorScheme()
	shapes := theme.Shapes()
	return SnackbarDefaultsData{
		Shape:                     shapes.FromToken(tokens.SnackbarTokens.ContainerShape),
		Color:                     colorScheme.FromToken(tokens.SnackbarTokens.ContainerColor),
		ContentColor:              colorScheme.FromToken(tokens.SnackbarTokens.SupportingTextColor),
		ActionColor:               colorScheme.FromToken(tokens.SnackbarTokens.ActionLabelTextColor),
		ActionContentColor:        colorScheme.FromToken(tokens.SnackbarTokens.ActionLabelTextColor),
		DismissActionContentColor: colorScheme.FromToken(tokens.SnackbarTokens.IconColor),
	}
}

type SnackbarDefaultsData struct {
	Shape                     shape.Shape
	Color                     graphics.Color
	ContentColor              graphics.Color
	ActionColor               graphics.Color
	ActionContentColor        graphics.Color
	DismissActionContentColor graphics.Color
}
