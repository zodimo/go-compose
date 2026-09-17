package fform

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
)

// type FormItem struct {
// 	Label    api.Composable
// 	HintText api.Composable
// 	Input    api.Composable
// 	Required bool
// 	Inline   bool
// }

type FormItemOptions struct {
	Required  bool
	Inline    bool
	HintText  api.Composable
	Validator []Validator
}

type FormItemOption = func(o *FormItemOptions)

func DefaultFormItemOptions() FormItemOptions {
	return FormItemOptions{
		Required: false,
		Inline:   false,
		HintText: compose.Id(),
	}
}

func FormItem[T any](
	itemState state.MutableValueTyped[T],
	label api.Composable,
	options ...FormItemOption,

) api.Composable {

	opts := DefaultFormItemOptions()
	for _, opt := range options {
		if opt != nil {
			opt(&opts)
		}
	}
	return func(c api.Composer) api.Composer {
		return c.IfLazy(
			opts.Inline,
			func() api.Composable {
				return column.Column(
					c.Sequence(
						row.Row(
							c.Sequence(
								label,
								input,
							),
						),
						opts.HintText,
					),
				)
			},
			func() api.Composable {
				return column.Column(
					c.Sequence(
						label,
						input,
						opts.HintText,
					),
				)
			},
		)(c)
	}

}
