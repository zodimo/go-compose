package main

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/overlay"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/material3/button"
	"github.com/zodimo/go-compose/compose/foundation/material3/checkbox"
	"github.com/zodimo/go-compose/compose/foundation/material3/dialog"
	"github.com/zodimo/go-compose/compose/foundation/material3/radiobutton"
	mswitch "github.com/zodimo/go-compose/compose/foundation/material3/switch"
	"github.com/zodimo/go-compose/compose/foundation/material3/textfield"

	m3text "github.com/zodimo/go-compose/compose/foundation/material3/text"
	"github.com/zodimo/go-compose/compose/foundation/text"

	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"

	"github.com/zodimo/go-compose/compose/foundation/material3/progress"
)

func UI(c api.Composer) api.LayoutNode {
	// State for dialogs
	showAck := c.State("showAck", func() any { return false })
	showConfirm := c.State("showConfirm", func() any { return false })

	// State for Toggle/Checkbox
	isChecked := c.State("isChecked", func() any { return false })
	isSwitched := c.State("isSwitched", func() any { return false })
	textValue := c.State("textValue", func() any { return "" })
	radioOption := c.State("radioOption", func() any { return 0 })
	progressVal := c.State("progressVal", func() any { return float32(0.5) })

	// Helper to set state
	setShowAck := func(v bool) { showAck.Set(v) }
	setShowConfirm := func(v bool) { showConfirm.Set(v) }

	c = box.Box(c.Sequence(
		// Main Content Layer
		column.Column(c.Sequence(
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
			row.Row(c.Sequence(
				checkbox.Checkbox(isChecked.Get().(bool), func(b bool) {
					isChecked.Set(b)
					fmt.Println("Checkbox changed:", b)
				}),
				text.Text(fmt.Sprintf("Checkbox is checked: %v", isChecked.Get())),
			), row.WithAlignment(row.Middle)),

			// Switch
			column.Column(c.Sequence(
				text.Text(fmt.Sprintf("Switch is checked: %v", isSwitched.Get())),
				mswitch.Switch(isSwitched.Get().(bool), func(b bool) {
					isSwitched.Set(b)
					fmt.Println("Switch changed:", b)
				}),
			), column.WithModifier(padding.Padding(padding.NotSet, 10, padding.NotSet, padding.NotSet))),

			// Radio Buttons
			column.Column(c.Sequence(
				text.Text(fmt.Sprintf("Radio selection: %d", radioOption.Get())),
				// Row 1
				row.Row(c.Sequence(
					radiobutton.RadioButton(radioOption.Get().(int) == 0, func() { radioOption.Set(0) }),
					text.Text("Option 0"),
				), row.WithAlignment(row.Middle)),
				// Row 2
				row.Row(c.Sequence(
					radiobutton.RadioButton(radioOption.Get().(int) == 1, func() { radioOption.Set(1) }),
					text.Text("Option 1"),
				), row.WithAlignment(row.Middle)),
				// Row 3 (Disabled)
				row.Row(c.Sequence(
					radiobutton.RadioButton(radioOption.Get().(int) == 2, func() { radioOption.Set(2) }, radiobutton.WithEnabled(false)),
					text.Text("Option 2 (Disabled)"),
				), row.WithAlignment(row.Middle)),
			), column.WithModifier(padding.Vertical(10, 10))),

			// Section: Progress Indicators
			text.Text("Progress Indicators", text.WithTextStyleOptions(text.StyleWithTextSize(20)), text.WithModifier(padding.Vertical(10, 10))),
			column.Column(c.Sequence(
				row.Row(c.Sequence(
					progress.CircularProgressIndicator(progressVal.Get().(float32)),
					progress.LinearProgressIndicator(progressVal.Get().(float32), progress.WithModifier(size.Width(200)), progress.WithModifier(padding.Horizontal(20, padding.NotSet))),
				), row.WithAlignment(row.Middle)),
				row.Row(c.Sequence(
					button.Text(func() {
						p := progressVal.Get().(float32) + 0.1
						if p > 1 {
							p = 1
						}
						progressVal.Set(p)
					}, "Add 0.1"),
					button.Text(func() {
						progressVal.Set(float32(0))
					}, "Reset"),
				)),
			), column.WithModifier(padding.Vertical(10, 10))),

			// Section: Inputs

			text.Text("Inputs", text.WithTextStyleOptions(text.StyleWithTextSize(20)), text.WithModifier(padding.Vertical(10, 10))),

			// Section: m3text Labels
			text.Text("Material3 Labels", text.WithTextStyleOptions(text.StyleWithTextSize(20)), text.WithModifier(padding.Vertical(10, 10))),
			column.Column(c.Sequence(
				m3text.Text("Display Large", m3text.TypestyleDisplayLarge),
				m3text.Text("Headline Medium", m3text.TypestyleHeadlineMedium),
				m3text.Text("Title Small", m3text.TypestyleTitleSmall),
				m3text.Text("Body Medium", m3text.TypestyleBodyMedium),
				m3text.Text("Label Small", m3text.TypestyleLabelSmall),
			)),

			column.Column(c.Sequence(
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
		c.When(showAck.Get().(bool), overlay.Overlay(
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
		)),

		c.When(showConfirm.Get().(bool), overlay.Overlay(
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
		)),
	))(c)

	return c.Build()
}
