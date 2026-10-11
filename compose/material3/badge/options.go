package badge

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/pkg/api"
)

type BadgeOptions struct {
	Content        api.Composable
	ContainerColor graphics.Color
	ContentColor   graphics.Color
	Modifier       ui.Modifier
}

type BadgeOption func(*BadgeOptions)

// WithContent sets custom composable content for the badge.
// For simple text badges, use WithText instead.
// The badge will apply the appropriate content color automatically.
func WithContent(content api.Composable) BadgeOption {
	return func(o *BadgeOptions) {
		o.Content = content
	}
}

// WithText is a convenience function to create a badge with text content.
// The text will use TypestyleLabelSmall and the badge's content color.
func WithText(label string) BadgeOption {
	return func(o *BadgeOptions) {
		o.Content = badgeText(label)
	}
}

func WithContainerColor(c graphics.Color) BadgeOption {
	return func(o *BadgeOptions) {
		o.ContainerColor = c
	}
}

func WithContentColor(c graphics.Color) BadgeOption {
	return func(o *BadgeOptions) {
		o.ContentColor = c
	}
}

func WithModifier(m ui.Modifier) BadgeOption {
	return func(o *BadgeOptions) {
		o.Modifier = m
	}
}

func DefaultBadgeOptions() BadgeOptions {
	return BadgeOptions{
		Modifier: ui.EmptyModifier,

		ContainerColor: graphics.ColorCompanion.Unspecified(),
		ContentColor:   graphics.ColorCompanion.Unspecified(),

		// ContainerColor: theme.ColorHelper.ColorSelector().ErrorRoles.Error,
		// ContentColor:   theme.ColorHelper.ColorSelector().ErrorRoles.OnError,
	}
}
