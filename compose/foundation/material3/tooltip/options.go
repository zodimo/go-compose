package tooltip

import (
	"go-compose-dev/internal/modifier"
)

// TooltipOptions contains configuration for tooltips.
type TooltipOptions struct {
	Modifier modifier.Modifier
}

// TooltipOption is a functional option for TooltipOptions.
type TooltipOption func(*TooltipOptions)

// DefaultTooltipOptions returns options with an empty modifier.
func DefaultTooltipOptions() TooltipOptions {
	return TooltipOptions{
		Modifier: EmptyModifier,
	}
}

// WithModifier appends a modifier to the tooltip.
func WithModifier(m Modifier) TooltipOption {
	return func(o *TooltipOptions) {
		o.Modifier = o.Modifier.Then(m)
	}
}
