package size

import (
	"go-compose-dev/internal/modifier"
)

type SizeOptions struct {
	Required bool
}

func DefaultSizeOptions() SizeOptions {
	return SizeOptions{
		Required: false,
	}
}

func SizeRequired() SizeOption {
	return func(options *SizeOptions) {
		options.Required = true
	}
}

type SizeOption func(options *SizeOptions)

func Size(width, height int, options ...SizeOption) Modifier {

	opt := DefaultSizeOptions()
	for _, option := range options {
		option(&opt)
	}

	return modifier.NewInspectableModifier(
		modifier.NewModifier(
			&SizeElement{
				size: SizeData{
					Width:    width,
					Height:   height,
					Required: opt.Required,
				},
			},
		),
		modifier.NewInspectorInfo(
			"size",
			map[string]any{
				"width":   width,
				"height":  height,
				"options": opt,
			},
		),
	)
}

func FillMax() Modifier {

	return modifier.NewInspectableModifier(
		modifier.NewModifier(
			&SizeElement{
				size: SizeData{
					Width:   NotSet,
					Height:  NotSet,
					FillMax: true,
				},
			},
		),
		modifier.NewInspectorInfo(
			"size",
			map[string]any{
				"FillMax": true,
			},
		),
	)
}

func FillMaxWidth() Modifier {

	return modifier.NewInspectableModifier(
		modifier.NewModifier(
			&SizeElement{
				size: SizeData{
					Width:        NotSet,
					Height:       NotSet,
					FillMaxWidth: true,
				},
			},
		),
		modifier.NewInspectorInfo(
			"size",
			map[string]any{
				"FillMaxWidth": true,
			},
		),
	)
}

func FillMaxHeight() Modifier {

	return modifier.NewInspectableModifier(
		modifier.NewModifier(
			&SizeElement{
				size: SizeData{
					Width:         NotSet,
					Height:        NotSet,
					FillMaxHeight: true,
				},
			},
		),
		modifier.NewInspectorInfo(
			"size",
			map[string]any{
				"FillMaxHeight": true,
			},
		),
	)
}

func WrapContent() Modifier {

	return modifier.NewInspectableModifier(
		modifier.NewModifier(
			&SizeElement{
				size: SizeData{
					Width:       NotSet,
					Height:      NotSet,
					WrapContent: true,
				},
			},
		),
		modifier.NewInspectorInfo(
			"size",
			map[string]any{
				"WrapContent": true,
			},
		),
	)
}

func Width(width int, options ...SizeOption) Modifier {

	opt := DefaultSizeOptions()
	for _, option := range options {
		option(&opt)
	}

	return modifier.NewInspectableModifier(
		modifier.NewModifier(
			&SizeElement{
				size: SizeData{
					Width:    width,
					Height:   NotSet,
					Required: opt.Required,
				},
			},
		),
		modifier.NewInspectorInfo(
			"size",
			map[string]any{
				"width":   width,
				"options": opt,
			},
		),
	)
}

func Height(height int, options ...SizeOption) Modifier {

	opt := DefaultSizeOptions()
	for _, option := range options {
		option(&opt)
	}

	return modifier.NewInspectableModifier(
		modifier.NewModifier(
			&SizeElement{
				size: SizeData{
					Width:    NotSet,
					Height:   height,
					Required: opt.Required,
				},
			},
		),
		modifier.NewInspectorInfo(
			"size",
			map[string]any{
				"height":  height,
				"options": opt,
			},
		),
	)
}
