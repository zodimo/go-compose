package components

import (
	"fmt"
	"strconv"

	fform "github.com/zodimo/go-compose/compose/foundation/form"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/material3/textfield"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/pkg/sentinel"
	"github.com/zodimo/go-compose/state"
)

// OptionalNumberFieldComponent renders an Outlined text field bound to an
// Optional[int] control.
//
// Unlike NumberFieldComponent, the empty string is a first-class value: an empty
// field writes None to the control, so "no input" is distinguishable from a
// legitimate 0. Use it when the field may be left blank (optional numeric input,
// "clear" semantics, RequiredOptional() validation).
//
// The raw text is remembered so partially-typed input ("-", "1e", "abc") can be
// edited freely; the control only ever holds a fully parsed Optional[int]. A
// non-empty, unparsable entry is surfaced as a supporting error once the control
// is touched, exactly like the other field components.
func OptionalNumberFieldComponent(binding *fform.FormFieldBinding[fform.Optional[int]], options ...NumberFieldComponentOption) api.Composable {
	opts := DefaultNumberFieldComponentOptions()
	for _, opt := range options {
		if opt != nil {
			opt(&opts)
		}
	}

	baseKey := opts.TextStateKey
	if baseKey == "" {
		baseKey = "optionalNumberFieldText"
	}
	textKey := baseKey + "/text"
	syncedKey := baseKey + "/synced"

	return func(c api.Composer) api.Composer {
		label := sentinel.TakeOrElseString(opts.Label, "")
		hint := sentinel.TakeOrElseString(opts.HintText, "")

		// The text buffer is the authority for what the user sees while typing.
		// synced holds the last control value the buffer was derived from, so an
		// external change (Reset, programmatic Set) re-seeds the text instead of
		// being shadowed by a stale buffer.
		rawText := state.MustRemember(c, textKey, func() string {
			return optionalText(binding.Value())
		})
		synced := state.MustRemember(c, syncedKey, func() string {
			return optionalKey(binding.Value())
		})

		current := binding.Value()
		if optionalKey(current) != synced.Get() {
			// Control changed underneath us: adopt its value as the displayed
			// text and record the sync point.
			synced.Set(optionalKey(current))
			rawText.Set(optionalText(current))
		}

		message, hasError := fieldError(binding)

		text := rawText.Get()
		if text != "" {
			if _, err := strconv.Atoi(text); err != nil && binding.IsTouched() {
				// Non-empty but unparsable: the user is mid-entry or wrong. Empty
				// is never an error here — that is the whole point of Optional.
				message = "not a valid number"
				hasError = true
			}
		}

		support := supportingText(hint, message, hasError)

		field := textfield.Outlined(
			text,
			func(s string) {
				rawText.Set(s)
				switch {
				case s == "":
					// Empty input is None, not 0.
					binding.SetValue(fform.None[int]())
					synced.Set(optionalKey(binding.Value()))
				default:
					if n, err := strconv.Atoi(s); err == nil {
						binding.SetValue(fform.Some(n))
						synced.Set(optionalKey(binding.Value()))
					} else {
						// Leave the control untouched mid-entry but mark touched
						// so an invalid non-empty value surfaces its error.
						binding.Control().MarkAsTouched()
					}
				}
			},
			textfield.WithLabel(label),
			textfield.WithSupportingText(support),
			textfield.WithError(hasError),
			textfield.WithEnabled(binding.IsEnabled()),
			textfield.WithSingleLine(true),
			textfield.WithModifier(opts.Modifier),
		)

		return column.Column(c.Sequence(field))(c)
	}
}

// optionalText renders an Optional[int] as field text: the number when present,
// otherwise the empty string. This is the mapping that makes an empty field
// possible at all.
func optionalText(v fform.Optional[int]) string {
	if v.IsNone() {
		return ""
	}
	return strconv.Itoa(v.UnwrapOr(0))
}

// optionalKey returns a stable string identity for an Optional[int] used to
// detect external control changes. None is a distinct sentinel from Some(0).
func optionalKey(v fform.Optional[int]) string {
	if v.IsNone() {
		return "\x00none"
	}
	return fmt.Sprintf("some:%d", v.UnwrapOr(0))
}
