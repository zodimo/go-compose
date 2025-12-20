package textfield

import (
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-maybe"
)

type TextFieldOptions struct {
	Modifier       Modifier
	Enabled        bool
	ReadOnly       bool
	TextStyle      interface{}
	Label          maybe.Maybe[api.Composable]
	Placeholder    maybe.Maybe[api.Composable]
	LeadingIcon    maybe.Maybe[api.Composable]
	TrailingIcon   maybe.Maybe[api.Composable]
	Prefix         maybe.Maybe[api.Composable]
	Suffix         maybe.Maybe[api.Composable]
	SupportingText maybe.Maybe[api.Composable]
	IsError        bool
	SingleLine     bool
	MaxLines       int
	MinLines       int

	// VisualTransformation interface{}
	// KeyboardOptions      interface{}
	// KeyboardActions      interface{}
	// InteractionSource    interface{}

	Shape  interface{}
	Colors interface{}

	OnSubmit func() // Called when Enter is pressed (SingleLine mode)

	// modifier: Modifier = Modifier,
	// enabled: Boolean = true,
	// readOnly: Boolean = false,
	// textStyle: TextStyle = LocalTextStyle.current,
	// label: (@Composable () -> Unit)? = null,
	// placeholder: (@Composable () -> Unit)? = null,
	// leadingIcon: (@Composable () -> Unit)? = null,
	// trailingIcon: (@Composable () -> Unit)? = null,
	// prefix: (@Composable () -> Unit)? = null,
	// suffix: (@Composable () -> Unit)? = null,
	// supportingText: (@Composable () -> Unit)? = null,
	// isError: Boolean = false,
	// visualTransformation: VisualTransformation = VisualTransformation.None,
	// keyboardOptions: KeyboardOptions = KeyboardOptions.Default,
	// keyboardActions: KeyboardActions = KeyboardActions.Default,
	// singleLine: Boolean = false,
	// maxLines: Int = if (singleLine) 1 else Int.MAX_VALUE,
	// minLines: Int = 1,
	// interactionSource: MutableInteractionSource? = null,
	// shape: Shape = TextFieldDefaults.shape,
	// colors: TextFieldColors = TextFieldDefaults.colors()
}

func DefaultTextFieldOptions() TextFieldOptions {
	return TextFieldOptions{
		Modifier:       EmptyModifier,
		Enabled:        true,
		ReadOnly:       false,
		TextStyle:      DefaultTextFieldColors().TextColor,
		Label:          maybe.None[api.Composable](),
		Placeholder:    maybe.None[api.Composable](),
		LeadingIcon:    maybe.None[api.Composable](),
		TrailingIcon:   maybe.None[api.Composable](),
		Prefix:         maybe.None[api.Composable](),
		Suffix:         maybe.None[api.Composable](),
		SupportingText: maybe.None[api.Composable](),
		IsError:        false,
		SingleLine:     true,
		MaxLines:       1,
		MinLines:       1,

		Colors: DefaultTextFieldColors(),
	}
}

type TextFieldOption func(*TextFieldOptions)

func WithModifier(m Modifier) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.Modifier = m
	}
}

func WithEnabled(enabled bool) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.Enabled = enabled
	}
}

func WithReadOnly(readOnly bool) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.ReadOnly = readOnly
	}
}

func WithTextStyle(textStyle interface{}) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.TextStyle = textStyle
	}
}

func WithLabel(label Composable) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.Label = maybe.Some(label)
	}
}

func WithPlaceholder(placeholder Composable) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.Placeholder = maybe.Some(placeholder)
	}
}

func WithLeadingIcon(icon Composable) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.LeadingIcon = maybe.Some(icon)
	}
}

func WithTrailingIcon(icon Composable) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.TrailingIcon = maybe.Some(icon)
	}
}

func WithPrefix(prefix Composable) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.Prefix = maybe.Some(prefix)
	}
}

func WithSuffix(suffix Composable) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.Suffix = maybe.Some(suffix)
	}
}

func WithSupportingText(text Composable) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.SupportingText = maybe.Some(text)
	}
}

func WithError(isError bool) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.IsError = isError
	}
}

func WithSingleLine(singleLine bool) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.SingleLine = singleLine
	}
}

func WithMaxLines(maxLines int) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.MaxLines = maxLines
	}
}

func WithMinLines(minLines int) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.MinLines = minLines
	}
}

func WithShape(shape interface{}) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.Shape = shape
	}
}

func WithColors(colors interface{}) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.Colors = colors
	}
}

func WithOnSubmit(onSubmit func()) TextFieldOption {
	return func(o *TextFieldOptions) {
		o.OnSubmit = onSubmit
	}
}
