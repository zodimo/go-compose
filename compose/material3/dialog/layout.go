package dialog

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/overlay"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/stopclickthru"
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/material3/surface"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
)

// DialogContent creates the internal layout structure for a Material3 dialog.
// It arranges the icon, title, text, and buttons according to Material Design 3 specifications.
//
// Layout structure:
//   - Icon (optional, centered)
//   - Title (optional, centered if icon present, start-aligned otherwise)
//   - Text (optional, start-aligned)
//   - Buttons (end-aligned in a row)
func DialogContent(
	icon api.Composable,
	title api.Composable,
	content api.Composable,
	buttons api.Composable,
) api.Composable {
	return func(c api.Composer) api.Composer {
		theme := material3.Theme(c)
		colorScheme := theme.ColorScheme()

		// Build content items
		var contentItems []api.Composable

		// Icon slot (optional, centered)
		if icon != nil {
			iconItem := box.Box(
				compose.CompositionLocalProvider(
					[]api.ProvidedValue{material3.LocalContentColor.Provides(colorScheme.Secondary)},
					icon,
				),
				box.WithAlignment(box.N), // Center horizontally at top
				box.WithModifier(padding.Padding(0, 0, 0, int(DialogPadding.IconBottom))),
			)
			contentItems = append(contentItems, iconItem)
		}

		// Title slot (optional)
		if title != nil {
			// Center title if icon is present, otherwise start-align
			alignment := box.NW
			if icon != nil {
				alignment = box.N
			}
			titleItem := box.Box(
				compose.CompositionLocalProvider(
					[]api.ProvidedValue{material3.LocalContentColor.Provides(colorScheme.OnSurface)},
					title,
				),
				box.WithAlignment(alignment),
				box.WithModifier(padding.Padding(0, 0, 0, int(DialogPadding.TitleBottom))),
			)
			contentItems = append(contentItems, titleItem)
		}

		if content != nil {
			contentItem := box.Box(
				compose.CompositionLocalProvider(
					[]api.ProvidedValue{material3.LocalContentColor.Provides(colorScheme.OnSurfaceVariant)},
					content,
				),
				box.WithAlignment(box.NW),
				box.WithModifier(padding.Padding(0, 0, 0, int(DialogPadding.TextBottom))),
			)
			contentItems = append(contentItems, contentItem)
		}

		// Buttons slot (end-aligned in row)
		if buttons != nil {
			buttonsItem := box.Box(
				buttons,
				box.WithAlignment(box.NE), // Align to end
				box.WithModifier(size.FillMaxWidth()),
			)
			contentItems = append(contentItems, buttonsItem)
		}

		// Create the content column and wrap in surface
		return column.Column(
			compose.Sequence(contentItems...),
		)(c)
	}
}

// DialogButtonRow creates a row layout for dialog buttons with proper spacing.
// Buttons are arranged end-aligned with 8dp spacing between them.
func DialogButtonRow(buttons ...api.Composable) api.Composable {
	if len(buttons) == 0 {
		return nil
	}

	// Interleave buttons with spacers
	var rowItems []api.Composable
	for i, button := range buttons {
		if button == nil {
			continue
		}
		if i > 0 && len(rowItems) > 0 {
			rowItems = append(rowItems, spacer.Width(int(DialogPadding.ButtonSpacing)))
		}
		rowItems = append(rowItems, button)
	}

	return row.Row(
		compose.Sequence(rowItems...),
		row.WithSpacing(row.SpaceStart), // Push to end by leaving space at start
	)
}

// DialogSurface wraps content in a dialog-styled surface with proper
// Material3 styling (rounded corners, elevation, background color).
// Used by BasicAlertDialog for custom dialog content.
func DialogSurface(content api.Composable, onDismiss func()) api.Composable {
	return func(c api.Composer) api.Composer {
		theme := material3.Theme(c)
		colorScheme := theme.ColorScheme()

		return overlay.Overlay(
			surface.Surface(
				stopclickthru.StopClickThru(
					//content with padding
					box.Box(
						content,
						box.WithModifier(
							padding.All(int(DialogPadding.All)).
								Then(size.WidthIn(int(DialogDefaults.MinWidth), int(DialogDefaults.MaxWidth))),
						),
					),
				),
				surface.WithShape(DialogDefaults.Shape),
				surface.WithColor(colorScheme.SurfaceContainerHigh),
				surface.WithShadowElevation(DialogDefaults.ShadowElevation),
			),
			overlay.WithOnClick(onDismiss),
		)(c)
	}
}
