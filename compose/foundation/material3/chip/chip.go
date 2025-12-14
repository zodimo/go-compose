package chip

import (
	"fmt"
	"go-compose-dev/compose/foundation/layout/box"
	"go-compose-dev/compose/foundation/layout/row"
	"go-compose-dev/compose/foundation/material3/surface"
	"go-compose-dev/compose/foundation/material3/text"
	"go-compose-dev/internal/modifiers/clickable"
	"go-compose-dev/internal/modifiers/padding"

	"gioui.org/widget"
)

// AssistChip represents a Material 3 Assist Chip.
// Assist chips represent smart or automated actions that can span multiple apps,
// such as opening a calendar event from the home screen.
func AssistChip(onClick func(), label string, options ...ChipOption) Composable {
	return Chip(onClick, label, options...)
}

// FilterChip represents a Material 3 Filter Chip.
// Filter chips use tags or descriptive words to filter content.
// They can be a good alternative to toggle buttons or checkboxes.
func FilterChip(onClick func(), label string, options ...ChipOption) Composable {
	// TODO: Add specific FilterChip defaults (selected state styling)
	return Chip(onClick, label, options...)
}

// InputChip represents a Material 3 Input Chip.
// Input chips represent pieces of information entered by a user, such as an event
// attendee or filter term.
func InputChip(onClick func(), label string, options ...ChipOption) Composable {
	// TODO: Add specific InputChip defaults
	return Chip(onClick, label, options...)
}

// SuggestionChip represents a Material 3 Suggestion Chip.
// Suggestion chips help narrow down user intent by presenting dynamically generated suggestions,
// such as possible responses or search filters.
func SuggestionChip(onClick func(), label string, options ...ChipOption) Composable {
	// TODO: Add specific SuggestionChip defaults
	return Chip(onClick, label, options...)
}

// Chip is the internal generic implementation.
func Chip(onClick func(), label string, options ...ChipOption) Composable {
	return func(c Composer) Composer {
		opts := DefaultChipOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		// State for clickable
		key := c.GenerateID()
		path := c.GetPath()
		clickStatePath := fmt.Sprintf("%d/%s/chip_click", key, path)
		clickState := c.State(clickStatePath, func() any { return &widget.Clickable{} })
		gioClickable := clickState.Get().(*widget.Clickable)

		// Layout:
		// Surface (Shape, Border, Color) -> Clickable -> Padding -> Row [Icon, Label, Icon]

		// Padding values (in pixels/dp approximation)
		startPadding := 16
		endPadding := 16
		if opts.LeadingIcon != nil {
			startPadding = 8
		}
		if opts.TrailingIcon != nil {
			endPadding = 8
		}
		verticalPadding := 6

		rowModifier := EmptyModifier.
			Then(padding.Padding(startPadding, verticalPadding, endPadding, verticalPadding))

		// Compose the content
		content := func(c Composer) Composer {
			return row.Row(
				func(c Composer) Composer {
					if opts.LeadingIcon != nil {
						// Add end padding to icon
						iconWrapper := box.Box(
							opts.LeadingIcon,
							box.WithModifier(padding.Padding(0, 0, 8, 0)),
						)
						iconWrapper(c)
					}

					text.Text(label, text.TypestyleLabelLarge)(c)

					if opts.TrailingIcon != nil {
						// Add start padding to icon
						iconWrapper := box.Box(
							opts.TrailingIcon,
							box.WithModifier(padding.Padding(8, 0, 0, 0)),
						)
						iconWrapper(c)
					}
					return c
				},
				row.WithModifier(rowModifier),
				row.WithAlignment(row.Middle),
				row.WithSpacing(row.SpaceSides),
			)(c)
		}

		// Surface options
		surfaceOpts := []surface.SurfaceOption{
			surface.WithShape(opts.Shape),
			surface.WithColor(opts.Color),
			surface.WithBorder(opts.BorderWidth, opts.BorderColor),
			surface.WithShadowElevation(opts.Elevation),
		}

		// Prepare Surface Modifier with Clickable
		// Surface applies its own modifiers (Shadow, Clip, Background, Border, UserModifier)
		// We want Clickable to be part of the UserModifier so it's inside the clip/border but handles events.

		clickableMod := clickable.OnClick(onClick, clickable.WithClickable(gioClickable))

		finalModifier := opts.Modifier.Then(clickableMod)

		surfaceOpts = append(surfaceOpts, surface.WithModifier(finalModifier))

		return surface.Surface(
			content,
			surfaceOpts...,
		)(c)
	}
}
