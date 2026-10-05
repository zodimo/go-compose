package components

import (
	fform "github.com/zodimo/go-compose/compose/foundation/form"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	mswitch "github.com/zodimo/go-compose/compose/material3/switch"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/pkg/sentinel"
)

// SwitchComponent renders a Material 3 switch bound to a bool control, with an
// optional inline label and touched-gated error text beneath it. It mirrors
// CheckboxComponent but for the switch affordance.
func SwitchComponent(binding *fform.FormFieldBinding[bool], options ...SwitchComponentOption) api.Composable {
	opts := DefaultSwitchComponentOptions()
	for _, opt := range options {
		if opt != nil {
			opt(&opts)
		}
	}

	return func(c api.Composer) api.Composer {
		label := sentinel.TakeOrElseString(opts.Label, "")
		hint := sentinel.TakeOrElseString(opts.HintText, "")

		message, hasError := fieldError(binding)
		support := supportingText(hint, message, hasError)

		rowContent := []api.Composable{
			mswitch.Switch(
				binding.Value(),
				func(checked bool) {
					if !binding.IsEnabled() {
						return
					}
					binding.SetValue(checked)
				},
				mswitch.WithModifier(opts.Modifier),
			),
		}
		if label != "" {
			rowContent = append(rowContent, text.BodyMedium(label))
		}

		body := row.Row(
			c.Sequence(rowContent...),
			row.WithAlignment(row.Middle),
		)

		if support == "" {
			return body(c)
		}
		return column.Column(
			c.Sequence(
				body,
				text.BodyMedium(support),
				spacer.Height(opts.ErrorSpacing),
			),
			column.WithModifier(padding.Start(opts.ErrorIndent)),
		)(c)
	}
}
