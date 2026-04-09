package divider

import (
	fDivider "github.com/zodimo/go-compose/compose/foundation/divider"
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui/layout"
	"github.com/zodimo/go-compose/pkg/api"
)

func Divider(options ...DividerOption) api.Composable {
	return func(c api.Composer) api.Composer {

		opts := DefaultDividerOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		theme := material3.Theme(c)
		opts.Color = opts.Color.TakeOrElse(theme.ColorScheme().OutlineVariant)
		options = append(options, WithColor(opts.Color))
		if opts.Axis == layout.AxisUnspecified {
			opts.Axis = layout.AxisHorizontal
			options = append(options, WithAxis(opts.Axis))
		}

		return fDivider.Divider(options...)(c)
	}
}

func HorizontalDivider(options ...DividerOption) api.Composable {
	return func(c api.Composer) api.Composer {

		opts := DefaultDividerOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		opts.Axis = layout.AxisHorizontal

		options = append(options, WithAxis(opts.Axis))

		return Divider(options...)(c)
	}
}

func VerticalDivider(options ...DividerOption) api.Composable {
	return func(c api.Composer) api.Composer {

		opts := DefaultDividerOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		opts.Axis = layout.AxisVertical

		options = append(options, WithAxis(opts.Axis))

		return Divider(options...)(c)
	}
}
