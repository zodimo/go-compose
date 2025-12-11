package main

import (
	"fmt"
	"go-compose-dev/compose"
	"go-compose-dev/compose/foundation/layout/box"
	"go-compose-dev/compose/foundation/layout/column"
	"go-compose-dev/compose/foundation/layout/overlay"
	"go-compose-dev/compose/foundation/material3/button"
	"go-compose-dev/compose/foundation/material3/checkbox"
	"go-compose-dev/compose/foundation/material3/dialog"
	mswitch "go-compose-dev/compose/foundation/material3/switch"
	"go-compose-dev/compose/foundation/material3/textfield"

	m3text "go-compose-dev/compose/foundation/material3/text"
	"go-compose-dev/compose/foundation/text"

	"go-compose-dev/internal/modifiers/padding"
	"go-compose-dev/pkg/api"
)

func UI(c api.Composer) api.LayoutNode {
	// State for dialogs
	showAck := c.State("showAck", func() any { return false })
	showConfirm := c.State("showConfirm", func() any { return false })

	// State for Toggle/Checkbox
	isChecked := c.State("isChecked", func() any { return false })
	isSwitched := c.State("isSwitched", func() any { return false })
	textValue := c.State("textValue", func() any { return "" })

	// Helper to set state
	setShowAck := func(v bool) { showAck.Set(v) }
	setShowConfirm := func(v bool) { showConfirm.Set(v) }

	c = box.Box(compose.Sequence(
		// Main Content Layer
		column.Column(compose.Sequence(
			text.Text("Kitchen Sink"),

			// Section: Dialogs
			text.Text("Dialogs", text.WithTextStyleOptions(text.StyleWithTextSize(20)), text.WithModifier(padding.Vertical(10, 10))),

			button.Filled(func() {
				setShowAck(true)
			}, "Show Acknowledge Dialog", button.WithModifier(padding.All(5))),

			button.Filled(func() {
				setShowConfirm(true)
			}, "Show Confirm Dialog", button.WithModifier(padding.All(5))),

			// Section: Toggles
			text.Text("Toggles", text.WithTextStyleOptions(text.StyleWithTextSize(20)), text.WithModifier(padding.Vertical(10, 10))),

			// Checkbox
			column.Column(compose.Sequence(
				text.Text(fmt.Sprintf("Checkbox is checked: %v", isChecked.Get())),
				checkbox.Checkbox(isChecked.Get().(bool), func(b bool) {
					isChecked.Set(b)
					fmt.Println("Checkbox changed:", b)
				}),
			)),

			// Switch
			column.Column(compose.Sequence(
				text.Text(fmt.Sprintf("Switch is checked: %v", isSwitched.Get())),
				mswitch.Switch(isSwitched.Get().(bool), func(b bool) {
					isSwitched.Set(b)
					fmt.Println("Switch changed:", b)
				}),
			), column.WithModifier(padding.Padding(padding.NotSet, 10, padding.NotSet, padding.NotSet))),

			// Section: Inputs
			text.Text("Inputs", text.WithTextStyleOptions(text.StyleWithTextSize(20)), text.WithModifier(padding.Vertical(10, 10))),

			// Section: m3text Labels
			text.Text("Material3 Labels", text.WithTextStyleOptions(text.StyleWithTextSize(20)), text.WithModifier(padding.Vertical(10, 10))),
			column.Column(compose.Sequence(
				m3text.Text("Display Large", m3text.TypestyleDisplayLarge),
				m3text.Text("Headline Medium", m3text.TypestyleHeadlineMedium),
				m3text.Text("Title Small", m3text.TypestyleTitleSmall),
				m3text.Text("Body Medium", m3text.TypestyleBodyMedium),
				m3text.Text("Label Small", m3text.TypestyleLabelSmall),
			)),

			column.Column(compose.Sequence(
				textfield.TextField(
					textValue.Get().(string),
					func(s string) {
						textValue.Set(s)
					},
					"Label",
					textfield.WithSupportingText(fmt.Sprintf("You typed: %s", textValue.Get().(string))),
				),
			)),
		), column.WithModifier(padding.All(20))),

		// Dialog Layer (Overlay)
		func(c api.Composer) api.Composer {
			if showAck.Get().(bool) {
				return overlay.Overlay(
					dialog.AlertDialog(
						func() { setShowAck(false) },
						func() {
							fmt.Println("Acknowledged")
							setShowAck(false)
						},
						"OK",
						dialog.WithTitle("Out of stock"),
						dialog.WithText("The item in your cart is no longer available."),
					),
					overlay.WithOnDismiss(func() {
						// hide the Dialog
						setShowAck(false)
					}),
				)(c)
			}
			return c
		},

		func(c api.Composer) api.Composer {
			if showConfirm.Get().(bool) {
				return overlay.Overlay(
					dialog.AlertDialog(
						func() { setShowConfirm(false) },
						func() {
							fmt.Println("Deleted!")
							setShowConfirm(false)
						},
						"Delete",
						dialog.WithTitle("Permanently Delete?"),
						dialog.WithText("Deleting the selected messages will also remove them from synced devices."),
						dialog.WithDismissButton("Cancel", func() {
							fmt.Println("Cancelled")
							setShowConfirm(false)
						}),
					),
					overlay.WithOnDismiss(func() {
						// hide the Dialog
						setShowConfirm(false)
					}),
				)(c)
			}
			return c
		},
	))(c)

	return c.Build()
}
