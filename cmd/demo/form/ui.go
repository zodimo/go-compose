package main

import (
	fform "github.com/zodimo/go-compose/compose/foundation/form"
	"github.com/zodimo/go-compose/compose/foundation/form/components"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/state"

	"github.com/zodimo/go-compose/compose/material3/checkbox"
	"github.com/zodimo/go-compose/compose/material3/next/button"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
)

// UI renders the reactive form engine demo: nested groups, an array node,
// cross-field disable (same-as-shipping), and touched-gated error display.
func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		formState := fform.RememberFormState(c, buildForm)

		return column.Column(
			c.Sequence(
				text.HeadlineMedium("Reactive Form Engine Demo"),
				spacer.Height(16),
				fform.Form(formState, formContent, fform.WithModifier(padding.All(16))),
				button.Outlined(func() {
					// print form data
					//type FormNode interface {
					// formState

				}, "Submit"),
			),
			column.WithSpacing(column.SpaceSides),
			column.WithAlignment(column.Middle),
			column.WithModifier(size.FillMax()),
		)(c)
	}
}

// buildForm constructs the remembered form tree. Each leaf control is backed by a
// remembered typed store so its value survives recomposition.
func buildForm(c api.Composer) *fform.Group {
	nameStore := state.MustRemember(c, "form-demo-name-store", func() string { return "" })
	streetStore := state.MustRemember(c, "form-demo-street-store", func() string { return "" })
	cityStore := state.MustRemember(c, "form-demo-city-store", func() string { return "" })
	phone1Store := state.MustRemember(c, "form-demo-phone1-store", func() string { return "" })
	phone2Store := state.MustRemember(c, "form-demo-phone2-store", func() string { return "" })
	sameAsShippingStore := state.MustRemember(c, "form-demo-same-as-shipping-store", func() bool { return false })

	sameAsShipping := fform.NewControl(fform.NewMutableValueValueStore(sameAsShippingStore), false)

	billing := fform.NewGroup(map[string]fform.FormNode{
		"street": fform.NewControl(fform.NewMutableValueValueStore(streetStore), "", fform.Required("")),
		"city":   fform.NewControl(fform.NewMutableValueValueStore(cityStore), "", fform.Required("")),
	})

	// Cross-field disable: checking "same as shipping" disables the billing subtree.
	sameAsShipping.OnValueChange(func(same bool) {
		billing.SetDisabled(same)
	})

	return fform.NewGroup(map[string]fform.FormNode{
		"identity": fform.NewGroup(map[string]fform.FormNode{
			"name": fform.NewControl(fform.NewMutableValueValueStore(nameStore), "", fform.Required(""), fform.MinLength(3)),
		}),
		"sameAsShipping": sameAsShipping,
		"billing":        billing,
		"phones": fform.NewArray(
			fform.NewGroup(map[string]fform.FormNode{
				"number": fform.NewControl(fform.NewMutableValueValueStore(phone1Store), "", fform.Required("")),
			}),
			fform.NewGroup(map[string]fform.FormNode{
				"number": fform.NewControl(fform.NewMutableValueValueStore(phone2Store), "", fform.Required("")),
			}),
		),
	})
}

// formContent binds named tree nodes to their rendered components.
func formContent(fs fform.FormScope) {
	fs.Form("identity", text.TitleLarge("Identity"), func(fs fform.FormScope) {
		fs.Field("name", func(node fform.FormNode) api.Composable {
			control := node.(*fform.Control[string])
			return components.TextFieldComponent(
				fform.NewFormFieldBinding(control),
				components.TextFieldWithLabel("Name"),
				components.TextFieldWithHintText("Minimum 3 characters"),
			)
		})
	})

	fs.Field("sameAsShipping", func(node fform.FormNode) api.Composable {
		control := node.(*fform.Control[bool])
		return row.Row(
			func(c api.Composer) api.Composer {
				return c.Sequence(
					checkbox.Checkbox(control.Value(), control.Set),
					text.BodyMedium("Same as shipping"),
				)(c)
			},
			row.WithAlignment(row.Middle),
		)
	})

	fs.Form("billing", text.TitleLarge("Billing Address"), func(fs fform.FormScope) {
		fs.Field("street", func(node fform.FormNode) api.Composable {
			control := node.(*fform.Control[string])
			return components.TextFieldComponent(
				fform.NewFormFieldBinding(control),
				components.TextFieldWithLabel("Street"),
			)
		})
		fs.Field("city", func(node fform.FormNode) api.Composable {
			control := node.(*fform.Control[string])
			return components.TextFieldComponent(
				fform.NewFormFieldBinding(control),
				components.TextFieldWithLabel("City"),
			)
		})
	})

	// Array children resolve via full bracketed paths at the root scope: a nested
	// Form scope would join "phones." ahead of the index and break resolution.
	fs.Form("phones", text.TitleLarge("Phone Numbers"), func(fform.FormScope) {})
	fs.Field("phones[0].number", func(node fform.FormNode) api.Composable {
		control := node.(*fform.Control[string])
		return components.TextFieldComponent(
			fform.NewFormFieldBinding(control),
			components.TextFieldWithLabel("Primary phone"),
		)
	})
	fs.Field("phones[1].number", func(node fform.FormNode) api.Composable {
		control := node.(*fform.Control[string])
		return components.TextFieldComponent(
			fform.NewFormFieldBinding(control),
			components.TextFieldWithLabel("Secondary phone"),
		)
	})
}
