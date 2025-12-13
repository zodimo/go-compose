package badge

import (
	"image/color"

	"go-compose-dev/internal/modifier"
	"go-compose-dev/pkg/api"
)

type BadgeOptions struct {
	Content        api.Composable
	ContainerColor color.NRGBA
	ContentColor   color.NRGBA
	Modifier       modifier.Modifier
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
func WithText(text string) BadgeOption {
	return func(o *BadgeOptions) {
		o.Content = Text(text, TypestyleLabelSmall)
	}
}

func WithContainerColor(c color.NRGBA) BadgeOption {
	return func(o *BadgeOptions) {
		o.ContainerColor = c
	}
}

func WithContentColor(c color.NRGBA) BadgeOption {
	return func(o *BadgeOptions) {
		o.ContentColor = c
	}
}

func WithModifier(modifier Modifier) BadgeOption {
	return func(o *BadgeOptions) {
		o.Modifier = o.Modifier.Then(modifier)
	}
}

func DefaultBadgeOptions() BadgeOptions {
	return BadgeOptions{
		Modifier: EmptyModifier,
	}
}
