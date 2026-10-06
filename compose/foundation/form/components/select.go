package components

import (
	"fmt"

	"github.com/zodimo/go-compose/compose"
	fform "github.com/zodimo/go-compose/compose/foundation/form"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/material3/icon"
	"github.com/zodimo/go-compose/compose/material3/menu"
	"github.com/zodimo/go-compose/compose/material3/textfield"
	"github.com/zodimo/go-compose/compose/ui"
	boxm "github.com/zodimo/go-compose/modifiers/box"
	"github.com/zodimo/go-compose/modifiers/clickable"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/pkg/sentinel"
	"github.com/zodimo/go-compose/state"
)

// SelectOption is one selectable entry in a SelectComponent. Value is what the
// bound control receives when the option is chosen; Label is what the user sees.
type SelectOption struct {
	Value string
	Label string
}

// SelectComponent renders an exposed-dropdown select bound to a string control.
//
// It uses a read-only Outlined text field as the anchor (showing the current
// option's label, or the placeholder when unset) wrapped in a clickable box that
// toggles a Material 3 DropdownMenu of the configured options. Choosing an option
// writes the option's Value through the binding, which marks the control touched
// and drives recomposition. Errors are touched-gated like every other field
// component.
//
// The menu's open/closed state is UI-only and lives in a remembered value, never
// in the form tree.
func SelectComponent(binding *fform.FormFieldBinding[string], options []SelectOption, componentOptions ...SelectComponentOption) api.Composable {
	opts := DefaultSelectComponentOptions()
	for _, opt := range componentOptions {
		if opt != nil {
			opt(&opts)
		}
	}

	return func(c api.Composer) api.Composer {
		label := sentinel.TakeOrElseString(opts.Label, "")
		placeholder := sentinel.TakeOrElseString(opts.Placeholder, "")

		message, hasError := fieldError(binding)
		support := supportingText(sentinel.TakeOrElseString(opts.HintText, ""), message, hasError)

		expandedKey := fmt.Sprintf("%d/%s/selectExpanded", c.GenerateID(), c.GetPath())
		expanded := state.MustRemember(c, expandedKey, func() bool { return false })

		display := displayLabel(binding.Value(), options, placeholder)

		menuItems := make([]api.Composable, 0, len(options))
		for _, option := range options {
			opt := option
			menuItems = append(menuItems, menu.DropdownMenuItem(
				opt.Label,
				func() {
					binding.SetValue(opt.Value)
					expanded.Set(false)
				},
			))
		}

		anchor := textfield.Outlined(
			display,
			func(string) {
				// Read-only anchor: text edits are ignored; the value comes from
				// picking a menu item.
			},
			textfield.WithLabel(label),
			textfield.WithSupportingText(support),
			textfield.WithError(hasError),
			textfield.WithEnabled(binding.IsEnabled()),
			textfield.WithReadOnly(true),
			textfield.WithTrailingIcon(icon.Icon(icon.SymbolArrowDropDown)),
			textfield.WithModifier(opts.Modifier),
		)

		// Only an enabled select opens the menu; a disabled one renders the
		// anchor inert with no click modifier.
		var clickMod ui.Modifier = ui.EmptyModifier
		if binding.IsEnabled() {
			clickMod = clickable.OnClick(func() { expanded.Set(!expanded.Get()) })
		}

		return column.Column(
			c.Sequence(
				box.Box(
					c.Sequence(
						anchor,
						box.Box(
							compose.Id(),
							box.WithModifier(
								boxm.MatchParentSize().Then(clickMod),
							),
						),
					),
				),
				menu.DropdownMenu(
					expanded.Get(),
					func() { expanded.Set(false) },
					menuItems,
				),
			),
		)(c)
	}
}

// displayLabel maps a control value back to the label of the matching option,
// falling back to placeholder when the value matches no option (including the
// empty/unset value).
func displayLabel(value string, options []SelectOption, placeholder string) string {
	for _, option := range options {
		if option.Value == value {
			return option.Label
		}
	}
	return placeholder
}
