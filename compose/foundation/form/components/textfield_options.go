package components

import (
	"fmt"

	fform "github.com/zodimo/go-compose/compose/foundation/form"
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/pkg/sentinel"

	"github.com/zodimo/go-compose/internal/modifier"
)

type TextFieldComponentOptions struct {
	Modifier  ui.Modifier
	TextStyle *text.TextStyle

	Label      string
	Inline     bool
	HintText   string
	Validators []*fform.ValidationRule[any]
}

type TextFieldComponentOption func(o *TextFieldComponentOptions)

func DefaultTextFieldComponentOptions() TextFieldComponentOptions {
	return TextFieldComponentOptions{
		Modifier:   modifier.EmptyModifier,
		TextStyle:  text.TextStyleUnspecified,
		Label:      sentinel.StringValueUnspecified,
		Inline:     false,
		HintText:   sentinel.StringValueUnspecified,
		Validators: []*fform.ValidationRule[any]{},
	}
}

func TextFieldWithModifier(m ui.Modifier) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		o.Modifier = o.Modifier.Then(m)
	}
}

func TextFieldWithTextStyle(style *text.TextStyle) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		o.TextStyle = style
	}
}

func TextFieldWithTextStyleOption(textStyleOption text.TextStyleOption) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		o.TextStyle = text.CopyTextStyle(o.TextStyle, textStyleOption)
	}
}

func TextFieldWithLabel(label string) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		o.Label = label
	}
}

func TextFieldWithInline(inline bool) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		o.Inline = inline
	}
}

func TextFieldWithHintText(hint string) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		o.HintText = hint
	}
}

func TextFieldWithValidators[T any](validators ...*fform.ValidationRule[T]) TextFieldComponentOption {
	return func(o *TextFieldComponentOptions) {
		anyValidators := make([]*fform.ValidationRule[any], len(validators))

		for i, validator := range validators {
			// Wrap each validator to safely handle `any` and cast it to `T`
			anyValidators[i] = fform.NewValidationRule[any](
				validator.ID(),
				func(v any) error {
					typedVal, ok := v.(T)
					if !ok {
						return fmt.Errorf("invalid value type: expected %T, got %T", *new(T), v)
					}
					return validator.Validate(typedVal)
				},
			)
		}

		o.Validators = anyValidators
	}
}
