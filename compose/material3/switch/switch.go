package mswitch

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/internal/layoutnode"

	"git.sr.ht/~schnwalter/gio-mw/widget/toggle"
)

const Material3SwitchNodeID = "Material3Switch"
const singleSwitchKey = "switch"

// Switch creates a Material3 switch (toggle).
func Switch(
	checked bool,
	onCheckedChange func(bool),
	options ...SwitchOption,
) Composable {
	return func(c Composer) Composer {
		opts := DefaultSwitchOptions()
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

		switchStatePath := fmt.Sprintf("%d/%s/switch", key, path)
		switchValue := c.State(switchStatePath, func() any {
			var values []string
			if checked {
				values = []string{singleSwitchKey}
			}

			// Create a generic toggle widget for strings
			t := toggle.NewToggle(
				[]string{singleSwitchKey},
				values,
				func(selected []string) {
					isChecked := false
					for _, s := range selected {
						if s == singleSwitchKey {
							isChecked = true
							break
						}
					}
					if handlerWrapper.Func != nil {
						handlerWrapper.Func(isChecked)
					}
				},
			)
			return t
		})
		t := switchValue.Get().(*toggle.Toggle[string])

		// Sync external state
		var currentValues []string
		if checked {
			currentValues = []string{singleSwitchKey}
		}
		t.SetValues(currentValues)

		// Note: We don't need to manually update onChange if we used the wrapper pattern correctly above.

		c.StartBlock(Material3SwitchNodeID)
		c.Modifier(func(modifier ui.Modifier) ui.Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.SetWidgetConstructor(switchWidgetConstructor(t))

		return c.EndBlock()
	}
}

func switchWidgetConstructor(t *toggle.Toggle[string]) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			labels := map[string]string{
				singleSwitchKey: "",
			}
			// Calculate the required width for the switch target
			// Width = TrackWidth + 2 * (FocusIndicatorOffset + FocusIndicatorThickness)
			th := toggle.BuildTheme(*gtx.ToGio())
			padding := th.FocusIndicatorOffset + th.FocusIndicatorThickness
			targetWidth := gtx.ToGio().Dp(th.TrackWidth) + 2*gtx.ToGio().Dp(padding)

			// Constrain the width to strictly the target width.
			// This prevents the underlying block.Line from expanding to fill the parent
			// and clips the unused label/spacer area.
			g := *gtx.ToGio()
			g.Constraints.Max.X = targetWidth
			g.Constraints.Min.X = targetWidth

			return layoutnode.FromGioDimensions(t.Layout(g, labels))
		}
	})
}
