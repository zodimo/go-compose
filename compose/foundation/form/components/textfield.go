package components

import (
	"errors"

	fform "github.com/zodimo/go-compose/compose/foundation/form"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	ftext "github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/material3/textfield"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/pkg/sentinel"
)

func TextFieldComponent(fieldState *fform.FormFieldState[string], options ...TextFieldComponentOption) api.Composable {
	opts := DefaultTextFieldComponentOptions()
	for _, opt := range options {
		if opt != nil {
			opt(&opts)
		}
	}

	validate := func(v string) error {
		var errs []error
		for _, validator := range opts.Validators {
			if err := validator.Validate(v); err != nil {
				errs = append(errs, err)
			}
		}

		if len(errs) > 0 {
			return errors.Join(errs...)
		}

		return nil
	}

	return func(c api.Composer) api.Composer {
		return c.IfLazy(
			opts.Inline,
			func() api.Composable {
				return column.Column(
					c.Sequence(
						row.Row(
							c.Sequence(
								textfield.Outlined(
									fieldState.Value(),
									func(s string) {
										fieldState.Touch()
										fieldState.SetValue(s)
										fieldState.SetError(validate(s))
									},
									textfield.WithLabel(sentinel.TakeOrElseString(opts.Label, "")),
									textfield.WithSupportingText(sentinel.TakeOrElseString(opts.HintText, "")),
									textfield.WithError(fieldState.HasError()),
									textfield.WithModifier(opts.Modifier),
								),
							),
						),
						c.WhenLazy(fieldState.HasError(), func() api.Composable {
							return text.BodySmall(fieldState.Error().Error(), ftext.WithColor(graphics.ColorRed))
						}),
					),
				)
			},
			func() api.Composable {
				return column.Column(
					c.Sequence(
						textfield.Outlined(
							fieldState.Value(),
							func(s string) {
								fieldState.Touch()
								fieldState.SetValue(s)
								fieldState.SetError(validate(s))
							},
							textfield.WithLabel(sentinel.TakeOrElseString(opts.Label, "")),
							textfield.WithSupportingText(sentinel.TakeOrElseString(opts.HintText, "")),
							textfield.WithError(fieldState.HasError()),
							textfield.WithModifier(opts.Modifier),
						),

						c.WhenLazy(fieldState.HasError(), func() api.Composable {
							return text.BodySmall(fieldState.Error().Error(), ftext.WithColor(graphics.ColorRed))
						}),
					),
				)
			},
		)(c)
	}

}
