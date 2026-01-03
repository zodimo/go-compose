package main

import (
	"gioui.org/unit"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/iconbutton"
	"github.com/zodimo/go-compose/compose/material3/textfield"
	"github.com/zodimo/go-compose/pkg/api"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

func UI(c api.Composer) api.LayoutNode {
	password := c.State("password", func() any { return "" })
	passwordVal := password.Get().(string)

	passwordVisible := c.State("password_visible", func() any { return false })
	isVisible := passwordVisible.Get().(bool)

	// Determine mask based on visibility
	mask := '•'
	if isVisible {
		mask = 0
	}

	// Choose icon based on visibility
	iconData := icons.ActionVisibility
	if isVisible {
		iconData = icons.ActionVisibilityOff
	}

	trailingIcon := iconbutton.Standard(
		func() {
			passwordVisible.Set(!isVisible)
		},
		iconData,
		"Toggle password visibility",
	)

	root := column.Column(
		c.Sequence(
			spacer.Height(int(unit.Dp(20))),

			textfield.SecureTextField(
				passwordVal,
				func(newVal string) {
					password.Set(newVal)
				},
				textfield.WithLabel("Password (Filled)"),
			),

			spacer.Height(int(unit.Dp(20))),

			textfield.SecureTextField(
				passwordVal,
				func(newVal string) {
					password.Set(newVal)
				},
				textfield.WithLabel("Password (Filled with reveal)"),
				textfield.WithMask(mask),
				textfield.WithTrailingIcon(trailingIcon),
			),

			spacer.Height(int(unit.Dp(20))),

			textfield.OutlinedSecureTextField(
				passwordVal,
				func(newVal string) {
					password.Set(newVal)
				},
				textfield.WithLabel("Password (Outlined)"),
			),

			spacer.Height(int(unit.Dp(20))),

			textfield.OutlinedSecureTextField(
				passwordVal,
				func(newVal string) {
					password.Set(newVal)
				},
				textfield.WithLabel("With Reveal"),
				textfield.WithMask(mask),
				textfield.WithTrailingIcon(trailingIcon),
			),
		),
	)

	return root(c).Build()
}
