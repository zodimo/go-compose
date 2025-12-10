package padding

import "go-compose-dev/internal/modifier"

type PaddingOptions struct {
	RtlAware bool
}

func DefaultPaddingOptions() PaddingOptions {
	return PaddingOptions{
		RtlAware: false,
	}
}

type PaddingOption func(options *PaddingOptions)

func Padding(start, top, end, bottom int, options ...PaddingOption) Modifier {

	opt := DefaultPaddingOptions()
	for _, option := range options {
		option(&opt)
	}
	return modifier.NewInspectableModifier(
		modifier.NewModifier(
			&paddingElement{
				padding: PaddingData{
					Start:    start,
					Top:      top,
					End:      end,
					Bottom:   bottom,
					RtlAware: opt.RtlAware,
				},
			},
		),
		modifier.NewInspectorInfo(
			"padding",
			map[string]any{
				"Start":   start,
				"Top":     top,
				"End":     end,
				"Bottom":  bottom,
				"options": opt,
			},
		),
	)
}

func PaddingAll(value int, options ...PaddingOption) Modifier {
	return Padding(value, value, value, value, options...)
}
