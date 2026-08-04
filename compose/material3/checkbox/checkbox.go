package checkbox

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/internal/layoutnode"

	"git.sr.ht/~schnwalter/gio-mw/widget/checkbox"
)

const Material3CheckboxNodeID = "Material3Checkbox"
const singleCheckboxKey = "checkbox"

// Checkbox creates a Material3 checkbox.
func Checkbox(
	checked bool,
	onCheckedChange func(bool),
	options ...CheckboxOption,
) Composable {
	return func(c Composer) Composer {
		opts := DefaultCheckboxOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		key := c.GenerateID()
		path := c.GetPath()

		// Fix closure capture for dynamic handlers
		handlerWrapperState := c.State(fmt.Sprintf("%d/%s/handler_wrapper", key, path), func() any {
			return &HandlerWrapper{Func: onCheckedChange}
		})
		handlerWrapper := handlerWrapperState.Get().(*HandlerWrapper)
		handlerWrapper.Func = onCheckedChange

		// Persist the generic Checkboxes widget state
		// We use a specific key "checkbox" to manage single boolean state via generic string logic
		checkboxStatePath := fmt.Sprintf("%d/%s/checkbox", key, path)
		checkboxesValue := c.State(checkboxStatePath, func() any {
			// Initial state based on checked param?
			// gio-mw Checkboxes manages its own internal values list.
			// But we want to drive it via props (controlled component).
			// We pass onChange to it.

			// Initial values
			var values []string
			if checked {
				values = []string{singleCheckboxKey}
			}

			// Create a generic checkboxes widget for strings
			cb := checkbox.NewCheckboxes(
				[]string{singleCheckboxKey},
				values,
				func(selected []string) {
					// This callback happens during Layout (Update phase of gio-mw widgets usually).
					// We can't directly effect the parent state here easily without recomposition.
					// But go-compose architecture usually handles callbacks.
					// The callback will be triggered when clicked.

					isChecked := false
					for _, s := range selected {
						if s == singleCheckboxKey {
							isChecked = true
							break
						}
					}
					if handlerWrapper.Func != nil {
						handlerWrapper.Func(isChecked)
					}
				},
			)
			return cb
		})
		cb := checkboxesValue.Get().(*checkbox.Checkboxes[string])

		// Sync external state 'checked' to internal widget state
		// This ensures controlled component behavior
		var currentValues []string
		if checked {
			currentValues = []string{singleCheckboxKey}
		}
		cb.SetValues(currentValues)
		// Also update onChange to capture closest closure?
		// gio-mw stores the callback in struct. We should update it just in case.
		// But `NewCheckboxes` sets it. We can manual set it if exported, but it's private in `checkbox.go` struct?
		// Wait, looking at `checkbox.go` step 196: `onChange func([]T)` is exported field of struct `Checkboxes`.

		c.StartBlock(Material3CheckboxNodeID)
		c.Modifier(func(modifier ui.Modifier) ui.Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.SetWidgetConstructor(checkboxWidgetConstructor(cb))

		return c.EndBlock()
	}
}

func checkboxWidgetConstructor(cb *checkbox.Checkboxes[string]) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			// Update processes events and calls onChange
			cb.Update(*gtx.ToGio())

			// Layout
			// We need to pass labels map
			labels := map[string]string{
				singleCheckboxKey: "", // No label rendered by the checkbox itself, usually wrapped or external in Compose
			}
			// If user wants a label part of the checkbox click area, we might need to expose it.
			// For now, standard Material checkbox is just the box, or box + label?
			// gio-mw Checkboxes renders the set.
			// LayoutWithKind(gtx, LeadingKind, labels)

			return layoutnode.FromGioDimensions(cb.LayoutWithKind(*gtx.ToGio(), checkbox.ButtonKind, labels))
		}
	})
}
