package button

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation"
	foundationLayout "github.com/zodimo/go-compose/compose/foundation/layout"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/material3/surface"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	text "github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/modifiers/clickable"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
)

func OutlinedButton(onClick func(), content Composable, options ...ButtonOption) Composable {
	return func(c Composer) Composer {
		opts := DefaultButtonOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		enabled := opts.Enabled
		// Default to enabled if not specified? Or is default explicitly handled in options?
		// In Go zero value for bool is false. We might want default true.
		// Let's assume the caller sets it or we should have a "WithEnabled(true)" default.
		// However, struct zero value is false.
		// We can check if option was applied or default to true logic.
		// For now let's assume if it's false it might be intentional, but usually buttons are enabled by default.
		// A better way is to set Enabled = true in separate DefaultButtonOptions constructor or handle it here.
		// Since we don't have "IsSet" flags, we might want to change DefaultButtonOptions to set Enabled = true.
		// But let's check `DefaultButtonOptions` in `button_defaults.go` which currently only sets Modifier.
		// I will modify `DefaultButtonOptions` to set Enabled = true later or assume true if we want.
		// Actually, let's just default to true in `DefaultButtonOptions` if I edit it, but for now I'll hardcode true if not sure?
		// No, let's fix `DefaultButtonOptions` in `button_defaults.go` separately.
		// For now, I will use `opts.Enabled`.

		// Resolve Defaults if nil/zero
		colors := opts.Colors
		if !IsSpecifiedButtonColors(colors) {
			colors = ButtonDefaults.OutlinedButtonColors(c)
		}

		// elevation := opts.Elevation
		// if (elevation == ButtonElevation{}) {
		// 	elevation = ButtonDefaults.OutlinedButtonElevation(c)
		// }

		border := opts.Border
		if !foundation.IsSpecifiedBorderStroke(border) {
			border = ButtonDefaults.OutlinedButtonBorder(c, enabled)
		}

		// Resolve shapes with precedence: opts.Shapes > opts.Shape > defaults
		// This supports both the simple Shape API and expressive Shapes API
		var buttonShapes *ButtonShapes
		if IsSpecifiedButtonShapes(opts.Shapes) {
			buttonShapes = opts.Shapes
		} else if shape.IsSpecifiedShape(opts.Shape) {
			// Single shape provided - use it for both states (no morphing)
			buttonShapes = &ButtonShapes{
				Shape:        opts.Shape,
				PressedShape: opts.Shape,
			}
		} else {
			// Use defaults with shape morphing
			buttonShapes = ButtonDefaults.OutlinedButtonShapes(c)
		}

		contentPadding := opts.ContentPadding
		if (contentPadding == foundationLayout.PaddingValues{}) {
			contentPadding = ButtonDefaults.ContentPadding()
		}

		// State for interaction
		// We need a persistent clickable state to track pressed/hovered state
		key := c.GenerateID()
		path := c.GetPath()
		clickableStatePath := fmt.Sprintf("%d/%s/buttonClickable", key, path)
		clickableState := c.State(clickableStatePath, func() any { return &clickable.GioClickable{} })
		clickState := clickableState.Get().(*clickable.GioClickable)

		// Determine colors and shape based on state
		var containerColor, contentColor graphics.Color
		var activeShape shape.Shape

		if !enabled {
			containerColor = colors.DisabledContainerColor
			contentColor = colors.DisabledContentColor
			// Use default shape when disabled (no pressed shape)
			activeShape = buttonShapes.Shape
		} else {
			containerColor = colors.ContainerColor
			contentColor = colors.ContentColor

			// Select shape based on pressed state
			if clickState.Pressed() {
				// Use pressed shape when button is being pressed
				activeShape = shape.TakeOrElseShape(buttonShapes.PressedShape, buttonShapes.Shape)
			} else {
				// Use default shape when not pressed
				activeShape = buttonShapes.Shape
			}
		}

		return surface.Surface(
			// Apply Clickable Modifier to Surface
			// We need a way to pass the state to the modifier.
			// `clickable.Clickable` function?
			// `modifiers/clickable/clickable.go` defines `ClickableData`.
			// We need a `WithResult` or similar if we want to extract the state?
			// Or `modifiers/clickable` should accept a state object.

			// Wait, looking at `modifiers/clickable/constructor.go` (I didn't read it but inferred from `node.go`).
			// `clickble.Clickable()`?
			// I need to check how to create a clickable modifier with external state.
			// In `node.go`: `element.clickableData.Clickable = clickable`.
			// So `ClickableElement` has `Clickable *GioClickable`.
			// `clickable.Clickable()` factory probably takes options or struct.

			// Let's assume `clickable.ClickableWithState(state)` or similar exists or I need to add it.
			// Currently `clickable.Clickable(onClick)`.

			// For now, let's just use `clickable.Clickable(onClick)` and assume we can't read state properly yet
			// without further refactoring of clickable modifier to expose state.
			// This is a blocker for "perfect" M3 interaction states (visual feedback),
			// but we can implement the structure.

			// Actually, let's look at `node.go` again.
			// `state` is created internally if nil.
			// Unless `ClickableElement` can be constructed with one.

			// Temporary: Use internal Clickable for logic, but we miss visual state updates for now
			// unless we refactor clickable.

			// Using `clickable.Clickable(onClick)`:
			func(c Composer) Composer {
				return row.Row(
					// Content Wrapper for Text Style
					func(c Composer) Composer {
						// Apply LocalTextStyle
						// We need `CompositionLocalProvider`.
						// Which is `c.Provider(...)`.

						textStyle := material3.Theme(c).Typography().LabelLarge
						// Copy with content color (requires copy helper or new style)
						// textStyle.Copy doesn't exist? We use TextStyleFromOptions(WithColor)
						// merging with existing style? Or just overriding color.
						// TextStyle is immutable/pointer.
						// We can use text.CopyTextStyle(textStyle, text.WithColor(contentColor))
						textStyle = text.CopyTextStyle(textStyle, text.WithColor(contentColor))

						return c.StartProviders([]api.ProvidedValue{
							material3.LocalTextStyle.Provides(textStyle),
							material3.LocalContentColor.Provides(contentColor),
						}).
							WithComposable(content).
							EndProviders()
					},

					row.WithModifier(
						size.Height(40). // hardcoded min height 40? unit.Dp(40) -> 40 int assumption
									Then(padding.Padding(
								int(contentPadding.Start),
								int(contentPadding.Top),
								int(contentPadding.End),
								int(contentPadding.Bottom),
							)),
					),
					row.WithAlignment(row.Middle),    // Center Vertically
					row.WithSpacing(row.SpaceAround), // SpaceAround? No, "Around" doesn't exist? Check Step 683: SpaceAround
				)(c)
			},

			surface.WithModifier(
				opts.Modifier.
					Then(clickable.OnClick(onClick, clickable.WithClickable(clickState))),
			),
			surface.WithShape(activeShape),
			surface.WithColor(containerColor),
			surface.WithBorder(
				border.Width,
				graphics.AsSolidColor(border.Brush).Value,
			),
		)(c)
	}
}
