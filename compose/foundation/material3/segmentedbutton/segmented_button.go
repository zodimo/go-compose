package segmentedbutton

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/material3/surface"
	"github.com/zodimo/go-compose/compose/foundation/material3/text"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/modifiers/clickable"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"

	"gioui.org/unit"
	"gioui.org/widget"
)

const SegmentedButtonRowNodeID = "Material3SegmentedButtonRow"
const SegmentedButtonNodeID = "Material3SegmentedButton"

// M3 Spec: Container height is 40dp
const SegmentHeight = unit.Dp(40)

// M3 Spec: Fully rounded corners (half of height)
const SegmentCornerRadius = unit.Dp(20)

// SingleChoiceSegmentedButtonRow creates a segmented button row where only one segment
// can be selected at a time (like radio buttons).
// The content should contain SegmentedButton composables.
func SingleChoiceSegmentedButtonRow(
	content Composable,
	options ...SegmentedButtonRowOption,
) Composable {
	return segmentedButtonRow(content, options...)
}

// MultiChoiceSegmentedButtonRow creates a segmented button row where multiple segments
// can be selected simultaneously (like checkboxes).
// The content should contain SegmentedButton composables.
func MultiChoiceSegmentedButtonRow(
	content Composable,
	options ...SegmentedButtonRowOption,
) Composable {
	// Implementation is the same as single choice - the difference is in how
	// the parent manages the state (single vs multi selection)
	return segmentedButtonRow(content, options...)
}

// segmentedButtonRow is the internal implementation for both row types.
func segmentedButtonRow(
	content Composable,
	options ...SegmentedButtonRowOption,
) Composable {
	return func(c Composer) Composer {
		opts := DefaultSegmentedButtonRowOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		// Row layout for segments
		rowModifier := opts.Modifier.
			Then(size.Height(int(SegmentHeight)))

		return row.Row(
			content,
			row.WithModifier(rowModifier),
			row.WithAlignment(row.Middle),
		)(c)
	}
}

// SegmentedButton creates an individual segment within a segmented button row.
// The `checked` parameter indicates if this segment is currently selected.
// The `onCheckedChange` callback is invoked when the segment is clicked.
// The `shape` parameter determines the corner rounding based on position.
func SegmentedButton(
	checked bool,
	onCheckedChange func(bool),
	label string,
	shape SegmentShape,
	options ...SegmentOption,
) Composable {
	return func(c Composer) Composer {
		opts := DefaultSegmentOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		// State for clickable
		key := c.GenerateID()
		path := c.GetPath()
		clickStatePath := fmt.Sprintf("%d/%s/segment_click", key, path)
		clickState := c.State(clickStatePath, func() any { return &widget.Clickable{} })
		gioClickable := clickState.Get().(*widget.Clickable)

		// Determine colors based on checked state
		bgColor := opts.UnselectedColor
		contentColor := opts.UnselectedContentColor
		if checked {
			bgColor = opts.SelectedColor
			contentColor = opts.SelectedContentColor
		}

		// Get the appropriate shape for this segment
		segmentShape := GetSegmentShape(SegmentCornerRadius, shape)

		// Build onClick handler
		onClick := func() {
			if opts.Enabled && onCheckedChange != nil {
				onCheckedChange(!checked)
			}
		}

		// Content: Row with optional icon(s) and label
		contentComposable := func(c Composer) Composer {
			// Row: [SelectedIcon] [Icon] [Label]
			return row.Row(
				func(c Composer) Composer {
					// Show selected icon (checkmark) when checked
					if checked && opts.ShowSelectedIcon && opts.SelectedIcon != nil {
						box.Box(
							opts.SelectedIcon,
							box.WithModifier(padding.Padding(0, 0, 8, 0)),
						)(c)
					}

					// Show regular icon if present
					if opts.Icon != nil {
						box.Box(
							opts.Icon,
							box.WithModifier(padding.Padding(0, 0, 8, 0)),
						)(c)
					}

					// Apply content color to text
					// Text color is usually handled by providing LocalContentColor
					// We can use surface.WithContentColor if text component supports it,
					// or we can wrap with a component providing content color values.
					// However, surface.Surface handles content color.
					// Let's rely on surface option below unless we need to override.

					// Label
					text.Text(label, text.TypestyleLabelLarge)(c)

					return c
				},
				row.WithAlignment(row.Middle),
			)(c)
		}

		// Padding for content
		contentPadding := padding.Padding(16, 0, 16, 0) // horizontal padding

		// Clickable modifier
		clickableMod := modifier.EmptyModifier
		if opts.Enabled {
			clickableMod = clickable.OnClick(onClick, clickable.WithClickable(gioClickable))
		}

		// Final modifier chain
		finalModifier := opts.Modifier.
			Then(size.Height(int(SegmentHeight))).
			Then(clickableMod)

		// Surface options
		surfaceOpts := []surface.SurfaceOption{
			surface.WithShape(segmentShape),
			surface.WithColor(bgColor),
			surface.WithContentColor(contentColor),
			surface.WithBorder(opts.BorderWidth, opts.BorderColor),
			surface.WithModifier(finalModifier),
		}

		// Wrap content with padding
		paddedContent := func(c Composer) Composer {
			return box.Box(
				contentComposable,
				box.WithModifier(contentPadding.Then(size.FillMaxHeight())),
				box.WithAlignment(box.Center),
			)(c)
		}

		return surface.Surface(
			paddedContent,
			surfaceOpts...,
		)(c)
	}
}
