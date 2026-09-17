package main

import (
	"fmt"

	fform "github.com/zodimo/go-compose/compose/foundation/form"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/state"

	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/material3/textfield"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		// alphaState := state.MustState(c, "alpha", func() float32 {
		// 	return 1.0
		// })

		// formSubmitCallback := func( // Callback function to handle form submission
		// 	values map[string]string,
		// ) {
		// 	for key, value := range values {
		// 		println(key + ": " + value)
		// 	}
		// }

		// cancelCallback := func() {}

		formState := fform.RememberFormState(c)

		return column.Column(
			c.Sequence(
				text.HeadlineMedium("Form Component Demo"),
				spacer.Height(16),
				fform.Form(
					formState,
					func(fs fform.FormScope) {
						fs.Form(
							"example_form",
							text.TitleLarge("This is a simple form."),
							func(fs fform.FormScope) {
								fs.Field("name", func(c api.Composer) api.Composer {
									// label
									//input
									// hint
									//error
									//onchange

									nameInputValue, err := state.Remember(c, "name-input", func() string { return "" })
									if err != nil {
										panic(fmt.Errorf("Could not rememeber form item name: %w", err))
									}
									formItemState := fform.RememberFormItemState(c)

									return fform.FormItem(
										formItemState,
										text.BodyLarge("name"),
										textfield.Filled(nameInputValue.Get(), func(s string) {
											fmt.Printf("Name changeed: %s\n", s)
											nameInputValue.Set(s)

										}),
									)(c)
								})
							},
						)
					},
				),
			),
			column.WithSpacing(column.SpaceSides),
			column.WithAlignment(column.Middle),
			column.WithModifier(size.FillMax()),
		)(c)
	}
}

// forms.NewTextField(
// 	"name", // Id
// 	[]forms.DisplayCondition{&forms.AlwaysDisplay{}}, // Display conditions (always display)
// 	[]forms.Validator{&forms.NotEmptyValidator{}},    // Validators (not empty)
// 	"John Doe", // Placeholder
// 	"Name: ",   // Prompt
// 	"",         // Default value
// ),
