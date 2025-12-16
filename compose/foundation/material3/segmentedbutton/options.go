package segmentedbutton

import (
	"image/color"

	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/pkg/api"

	"gioui.org/unit"
)

// Type aliases for convenience
type Composable = api.Composable
type Composer = api.Composer
type Modifier = modifier.Modifier

// SegmentedButtonRowOptions configures the segmented button row container.
type SegmentedButtonRowOptions struct {
	Modifier modifier.Modifier
	Space    unit.Dp // Overlap/space adjustment between segments
}

type SegmentedButtonRowOption func(*SegmentedButtonRowOptions)

// DefaultSegmentedButtonRowOptions returns default options for the row container.
func DefaultSegmentedButtonRowOptions() SegmentedButtonRowOptions {
	return SegmentedButtonRowOptions{
		Modifier: modifier.EmptyModifier,
		Space:    unit.Dp(0), // Segments typically touch
	}
}

// WithRowModifier sets a custom modifier for the row.
func WithRowModifier(m Modifier) SegmentedButtonRowOption {
	return func(o *SegmentedButtonRowOptions) {
		o.Modifier = o.Modifier.Then(m)
	}
}

// WithSpace sets the space/overlap between segments.
func WithSpace(space unit.Dp) SegmentedButtonRowOption {
	return func(o *SegmentedButtonRowOptions) {
		o.Space = space
	}
}

// SegmentOptions configures an individual segment.
type SegmentOptions struct {
	Modifier               modifier.Modifier
	Icon                   Composable // Optional leading icon
	SelectedIcon           Composable // Icon shown when selected (default: checkmark)
	ShowSelectedIcon       bool       // Whether to show selected icon
	Enabled                bool
	SelectedColor          color.NRGBA // Background color when selected
	UnselectedColor        color.NRGBA // Background color when unselected
	SelectedContentColor   color.NRGBA // Content color when selected
	UnselectedContentColor color.NRGBA // Content color when unselected
	BorderColor            color.NRGBA
	BorderWidth            unit.Dp
}

type SegmentOption func(*SegmentOptions)

// DefaultSegmentOptions returns default options for a segment.
func DefaultSegmentOptions() SegmentOptions {
	return SegmentOptions{
		Modifier:         modifier.EmptyModifier,
		ShowSelectedIcon: true,
		Enabled:          true,
		// Colors will be resolved from theme at render time
		SelectedColor:          color.NRGBA{},                                   // Will use SecondaryContainer
		UnselectedColor:        color.NRGBA{A: 0},                               // Transparent
		SelectedContentColor:   color.NRGBA{},                                   // Will use OnSecondaryContainer
		UnselectedContentColor: color.NRGBA{},                                   // Will use OnSurface
		BorderColor:            color.NRGBA{R: 0x79, G: 0x74, B: 0x7E, A: 0xFF}, // Outline
		BorderWidth:            unit.Dp(1),
	}
}

// WithModifier sets a custom modifier for the segment.
func WithModifier(m Modifier) SegmentOption {
	return func(o *SegmentOptions) {
		o.Modifier = o.Modifier.Then(m)
	}
}

// WithIcon sets the leading icon for the segment.
func WithIcon(icon Composable) SegmentOption {
	return func(o *SegmentOptions) {
		o.Icon = icon
	}
}

// WithSelectedIcon sets the icon shown when the segment is selected.
func WithSelectedIcon(icon Composable) SegmentOption {
	return func(o *SegmentOptions) {
		o.SelectedIcon = icon
	}
}

// WithShowSelectedIcon controls whether to show the selected icon.
func WithShowSelectedIcon(show bool) SegmentOption {
	return func(o *SegmentOptions) {
		o.ShowSelectedIcon = show
	}
}

// WithEnabled controls whether the segment is enabled.
func WithEnabled(enabled bool) SegmentOption {
	return func(o *SegmentOptions) {
		o.Enabled = enabled
	}
}

// WithSelectedColor sets the background color when selected.
func WithSelectedColor(c color.NRGBA) SegmentOption {
	return func(o *SegmentOptions) {
		o.SelectedColor = c
	}
}

// WithUnselectedColor sets the background color when unselected.
func WithUnselectedColor(c color.NRGBA) SegmentOption {
	return func(o *SegmentOptions) {
		o.UnselectedColor = c
	}
}

// WithBorder sets the border width and color.
func WithBorder(width unit.Dp, c color.NRGBA) SegmentOption {
	return func(o *SegmentOptions) {
		o.BorderWidth = width
		o.BorderColor = c
	}
}

// SegmentShape determines the shape of a segment based on its position.
type SegmentShape int

const (
	SegmentShapeStart  SegmentShape = iota // Rounded left corners (TopStart, BottomStart)
	SegmentShapeMiddle                     // No rounded corners
	SegmentShapeEnd                        // Rounded right corners (TopEnd, BottomEnd)
	SegmentShapeOnly                       // Fully rounded (single segment)
)

// GetSegmentShape returns the appropriate RoundedCornerShape for a segment position.
// Uses the standard shape.RoundedCornerShape with per-corner radius support,
// following Jetpack Compose's API pattern.
func GetSegmentShape(radius unit.Dp, position SegmentShape) shape.Shape {
	switch position {
	case SegmentShapeStart:
		// Left segment: TopStart (NW) and BottomStart (SW) rounded
		return shape.RoundedCornerShape{
			TopStart:    radius,
			TopEnd:      0,
			BottomEnd:   0,
			BottomStart: radius,
		}
	case SegmentShapeEnd:
		// Right segment: TopEnd (NE) and BottomEnd (SE) rounded
		return shape.RoundedCornerShape{
			TopStart:    0,
			TopEnd:      radius,
			BottomEnd:   radius,
			BottomStart: 0,
		}
	case SegmentShapeOnly:
		// Single segment: all corners rounded (use uniform Radius)
		return shape.RoundedCornerShape{Radius: radius}
	default: // SegmentShapeMiddle
		// Middle segment: no rounded corners
		return shape.RoundedCornerShape{Radius: 0}
	}
}
