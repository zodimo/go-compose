package main

import (
	"fmt"
	"strings"

	fform "github.com/zodimo/go-compose/compose/foundation/form"
	"github.com/zodimo/go-compose/compose/foundation/form/components"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/state"

	"github.com/zodimo/go-compose/compose/material3/next/button"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {

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

									nameInputValue := state.MustRemember(c, "name-input", func() string { return "" })

									formFieldState := fform.RememberFormFieldState(c, nameInputValue)

									return components.TextFieldComponent(
										formFieldState,
										components.TextFieldWithLabel("Name"),
										components.TextFieldWithHintText("Dont start with _"),
										components.TextFieldWithValidators(
											fform.NewValidationRule(
												"length",
												func(s string) error {
													if len(s) < 10 {
														return fmt.Errorf("Min length is 10")
													}
													return nil
												},
											),
											fform.NewValidationRule(
												"start with _",
												func(s string) error {
													if strings.HasPrefix(s, "_") {
														return fmt.Errorf("cannot start with _")
													}
													return nil
												},
											),
										),
									)(c)
								})

								fs.Field("submit", button.Outlined(func() {

								}, "submit"))
							},
						)
					},
					fform.WithModifier(padding.All(16)),
				),
			),
			column.WithSpacing(column.SpaceSides),
			column.WithAlignment(column.Middle),
			column.WithModifier(size.FillMax()),
		)(c)
	}
}
