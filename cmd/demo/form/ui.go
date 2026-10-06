package main

import (
	"fmt"
	"strings"

	fform "github.com/zodimo/go-compose/compose/foundation/form"
	"github.com/zodimo/go-compose/compose/foundation/form/components"
	"github.com/zodimo/go-compose/compose/foundation/form/rules"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/button"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/pkg/api"
)

// UI renders the reactive form engine demo: nested groups, a dynamic array,
// cross-field disable (same-as-shipping), touched-gated errors, and a real
// submit/reset/status lifecycle.
func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		formState := fform.RememberFormState(c, buildForm)
		submitted := rememberSubmitted(c)

		return column.Column(
			c.Sequence(
				text.HeadlineMedium("Reactive Form Engine Demo"),
				spacer.Height(8),
				formStatusLine(formState),
				spacer.Height(8),
				fform.Form(
					formState,
					formContent,
					fform.WithModifier(
						padding.All(8).Then(weight.Weight(1)),
					),
				),
				spacer.Height(16),
				formActions(formState, submitted),
				spacer.Height(16),
				submissionSummary(submitted),
			),
			column.WithSpacing(column.SpaceSides),
			column.WithAlignment(column.Middle),
			column.WithModifier(size.FillMax()),
		)(c)
	}
}

// buildForm constructs the remembered form tree. Each leaf control is backed by
// a plain in-memory store: the tree itself is remembered by RememberFormState,
// and the Form composable bumps a version on any tree mutation, so values
// persist across recompositions without a separate MutableValue per field.
func buildForm(c api.Composer) *fform.Group {
	sameAsShipping := fform.NewControl(fform.NewPlainValueStore(false), false)

	billing := fform.NewGroup(map[string]fform.FormNode{
		"street": fform.NewControl(fform.NewPlainValueStore(""), "", rules.StringNotEmpty()),
		"city":   fform.NewControl(fform.NewPlainValueStore(""), "", rules.StringNotEmpty()),
	})

	// Cross-field disable: checking "same as shipping" disables the billing
	// subtree; unchecking restores each child's own setting.
	sameAsShipping.OnValueChange(func(same bool) {
		billing.SetDisabled(same)
	})

	// Password and confirmation demonstrate a cross-field rule: EqualTo reads the
	// sibling's current value at validation time.
	password := fform.NewControl(
		fform.NewPlainValueStore(""),
		"",
		rules.StringNotEmpty(),
		rules.MinLength(8),
	)

	phones := fform.NewArray(phoneRow(""))

	return fform.NewGroup(map[string]fform.FormNode{
		"identity": fform.NewGroup(map[string]fform.FormNode{
			"name": fform.NewControl(fform.NewPlainValueStore(""), "",
				rules.StringNotEmpty(),
				rules.MinLength(3),
				rules.AlphaUnicode(),
			),
			// Age is an optional number: an empty field is distinct from 0, and
			// RequiredOptional rejects only the empty state.
			"age": fform.NewControl[fform.Optional[int]](
				fform.NewPlainValueStore(fform.None[int]()),
				fform.None[int](),
				fform.RequiredOptional[int](),
			),
			"password": password,
			"confirm": fform.NewControl(fform.NewPlainValueStore(""), "",
				rules.EqualTo(func() string { return password.Value() }),
			),
		}),
		// The contact group carries a group-level validator: at least one of
		// email/phone is required. This is a rule about the group, not any single
		// child, so it lives on the group and reports against the group's path.
		"contact": fform.NewGroup(
			map[string]fform.FormNode{
				"email": fform.NewControl(fform.NewPlainValueStore(""), "",
					rules.Email(),
				),
				"phone": fform.NewControl(fform.NewPlainValueStore(""), "",
					rules.Numeric(),
					// Phone is optional (only required if email is blank), so
					// enforce the length only when something was entered.
					rules.MinLengthIfNotEmpty(10),
				),
			},
			fform.WithGroupValidator(rules.GroupAtLeastOneSet[string]("email", "phone")),
		),
		"tier": fform.NewControl(fform.NewPlainValueStore(""), "",
			rules.OneOf("free", "pro", "enterprise"),
		),
		"sameAsShipping": sameAsShipping,
		"billing":        billing,
		"phones":         phones,
	})
}

// phoneRow builds one phone entry: a group with a required, ≥10-character
// number control. Rows are dynamic, so their stores are plain holders owned by
// the tree rather than remembered composition state.
func phoneRow(initial string) *fform.Group {
	return fform.NewGroup(map[string]fform.FormNode{
		"number": fform.NewControl(fform.NewPlainValueStore(initial), initial,
			rules.StringNotEmpty(),
			rules.Numeric(),
			rules.MinLength(10),
		),
	})
}

// formStatusLine renders the live aggregate status, error count, failing codes,
// and lifecycle flags. The codes line demonstrates the structured error surface:
// it lists the distinct ErrorCodes currently raised across the tree.
func formStatusLine(formState *fform.FormState) api.Composable {
	return func(c api.Composer) api.Composer {
		label := "VALID"
		switch formState.Status() {
		case fform.StatusInvalid:
			label = "INVALID"
		case fform.StatusDisabled:
			label = "DISABLED"
		case fform.StatusPending:
			label = "PENDING"
		}

		summary := fmt.Sprintf(
			"status=%s  errors=%d  touched=%t  dirty=%t",
			label, len(formState.Errors()), formState.IsTouched(), formState.IsDirty(),
		)

		codes := formState.Codes()
		codesLine := "codes=—"
		if len(codes) > 0 {
			parts := make([]string, 0, len(codes))
			for _, code := range codes {
				parts = append(parts, string(code))
			}
			codesLine = "codes=" + strings.Join(parts, ", ")
		}

		return column.Column(
			c.Sequence(
				text.BodyMedium(summary),
				text.BodySmall(codesLine),
			),
		)(c)
	}
}

// formContent binds named tree nodes to their rendered components.
func formContent(fs fform.FormScope) {
	fs.Form("identity", text.TitleLarge("Identity"), func(fs fform.FormScope) {
		fform.ControlFieldOf[string](fs, "name", func(b *fform.FormFieldBinding[string]) api.Composable {
			return components.TextFieldComponent(
				b,
				components.TextFieldWithLabel("Name"),
				components.TextFieldWithHintText("Minimum 3 letters"),
			)
		})
		fform.ControlFieldOf[fform.Optional[int]](fs, "age", func(b *fform.FormFieldBinding[fform.Optional[int]]) api.Composable {
			return components.OptionalNumberFieldComponent(
				b,
				components.NumberFieldWithLabel("Age"),
				components.NumberFieldWithHintText("Leave blank if unknown"),
				components.NumberFieldWithTextStateKey("identity-age-text"),
			)
		})
		fform.ControlFieldOf[string](fs, "password", func(b *fform.FormFieldBinding[string]) api.Composable {
			return components.TextFieldComponent(
				b,
				components.TextFieldWithLabel("Password"),
				components.TextFieldWithHintText("Minimum 8 characters"),
			)
		})
		fform.ControlFieldOf[string](fs, "confirm", func(b *fform.FormFieldBinding[string]) api.Composable {
			return components.TextFieldComponent(
				b,
				components.TextFieldWithLabel("Confirm password"),
				components.TextFieldWithHintText("Must match the password"),
			)
		})
	})

	// The contact group renders its fields, and the group-level validator's
	// failure ("at least one of email/phone") surfaces as an error under the
	// "contact" path.
	fs.Form("contact", text.TitleLarge("Contact"), func(fs fform.FormScope) {
		fform.ControlFieldOf[string](fs, "email", func(b *fform.FormFieldBinding[string]) api.Composable {
			return components.TextFieldComponent(
				b,
				components.TextFieldWithLabel("Email"),
				components.TextFieldWithHintText("name@example.com"),
			)
		})
		fform.ControlFieldOf[string](fs, "phone", func(b *fform.FormFieldBinding[string]) api.Composable {
			return components.TextFieldComponent(
				b,
				components.TextFieldWithLabel("Phone"),
				components.TextFieldWithHintText("At least 10 digits"),
			)
		})
	})

	fform.ControlFieldOf[string](fs, "tier", func(b *fform.FormFieldBinding[string]) api.Composable {
		return components.SelectComponent(
			b,
			[]components.SelectOption{
				{Value: "free", Label: "Free"},
				{Value: "pro", Label: "Pro"},
				{Value: "enterprise", Label: "Enterprise"},
			},
			components.SelectWithLabel("Tier"),
			components.SelectWithPlaceholder("Choose a tier"),
		)
	})

	fform.ControlFieldOf[bool](fs, "sameAsShipping", func(b *fform.FormFieldBinding[bool]) api.Composable {
		return components.CheckboxComponent(
			b,
			components.CheckboxWithLabel("Same as shipping"),
			components.CheckboxWithErrorIndent(40),
		)
	})

	fs.Form("billing", text.TitleLarge("Billing Address"), func(fs fform.FormScope) {
		fform.ControlFieldOf[string](fs, "street", func(b *fform.FormFieldBinding[string]) api.Composable {
			return components.TextFieldComponent(b, components.TextFieldWithLabel("Street"))
		})
		fform.ControlFieldOf[string](fs, "city", func(b *fform.FormFieldBinding[string]) api.Composable {
			return components.TextFieldComponent(b, components.TextFieldWithLabel("City"))
		})
	})

	phoneSection(fs)
}

// phoneSection renders the dynamic phone-number array: a header with an "Add"
// action, then one row per phone with a remove action. Add/remove mutate the
// Array node directly, which fires the tree notification and recomposes.
func phoneSection(fs fform.FormScope) {
	phones, ok := fs.Array("phones")
	if !ok {
		return
	}

	fs.Field("phones", func(fform.FormNode) api.Composable {
		return row.Row(
			func(c api.Composer) api.Composer {
				return c.Sequence(
					text.TitleLarge("Phone Numbers"),
					spacer.Width(8),
					button.Text(func() {
						phones.Add(phoneRow(""))
					}, "Add"),
				)(c)
			},
			row.WithAlignment(row.Middle),
			row.WithModifier(size.FillMaxWidth()),
		)
	})

	fs.FormArray("phones", func(index int, rowScope fform.FormScope) api.Composable {
		field := fform.ControlViewOf[string](rowScope, "number", func(b *fform.FormFieldBinding[string]) api.Composable {
			return components.TextFieldComponent(
				b,
				components.TextFieldWithLabel(fmt.Sprintf("Phone %d", index+1)),
			)
		})

		return row.Row(
			func(c api.Composer) api.Composer {
				return c.Sequence(
					flexWeight(field, 1),
					button.Text(func() {
						phones.Remove(index)
					}, "Remove"),
				)(c)
			},
			row.WithAlignment(row.Middle),
			row.WithModifier(size.FillMaxWidth()),
		)
	})
}

// flexWeight wraps content so it flexes within a row.
func flexWeight(content api.Composable, w int) api.Composable {
	return column.Column(
		content,
		column.WithModifier(weight.Weight(w)),
	)
}

// formActions renders Submit and Reset buttons wired to the form lifecycle.
func formActions(formState *fform.FormState, submitted *submittedState) api.Composable {
	return func(c api.Composer) api.Composer {
		return row.Row(
			c.Sequence(
				button.Filled(func() {
					// Submit validates, marks every control touched so errors
					// render, and only proceeds when valid.
					if !formState.Submit(func(value map[string]any) {
						submitted.Set(value)
					}) {
						submitted.Clear()
					}
				}, "Submit"),
				spacer.Width(8),
				button.Outlined(func() {
					formState.Reset()
					submitted.Clear()
				}, "Reset"),
			),
			row.WithAlignment(row.Middle),
		)(c)
	}
}

// submissionSummary renders the last submitted value, or a hint when nothing has
// been submitted successfully.
func submissionSummary(submitted *submittedState) api.Composable {
	return func(c api.Composer) api.Composer {
		value, ok := submitted.Get()
		if !ok {
			return text.BodyMedium("Submit when valid to see the serialized value.")(c)
		}
		return column.Column(
			c.Sequence(
				text.TitleMedium("Submitted value"),
				text.BodyMedium(fmt.Sprintf("%v", value)),
			),
			column.WithModifier(padding.All(8)),
		)(c)
	}
}
