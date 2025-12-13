package tooltip

import (
	"fmt"
	"image"

	"go-compose-dev/internal/layoutnode"
	"go-compose-dev/internal/modifier"
	"go-compose-dev/pkg/api"

	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	mwtooltip "git.sr.ht/~schnwalter/gio-mw/widget/tooltip"
)

// Tooltip wraps content and shows a simple text tooltip when hovered.
func Tooltip(
	text string,
	content api.Composable,
	modifiers ...modifier.Modifier,
) api.Composable {
	return func(c api.Composer) api.Composer {
		c.StartBlock("Tooltip")

		for _, m := range modifiers {
			c.Modifier(func(old modifier.Modifier) modifier.Modifier {
				return old.Then(m)
			})
		}

		// State for hover detection
		key := c.GenerateID()
		path := c.GetPath()
		hoverStatePath := fmt.Sprintf("%d/%s/hover", key, path)
		// Pointer to bool to track hover state
		hoveredValue := c.State(hoverStatePath, func() any { v := false; return &v })
		hovered := hoveredValue.Get().(*bool)

		c.WithComposable(content)

		c.SetWidgetConstructor(tooltipWidgetConstructor(text, hovered))

		return c.EndBlock()
	}
}

func tooltipWidgetConstructor(text string, hovered *bool) layoutnode.LayoutNodeWidgetConstructor {
	var gioTooltip *mwtooltip.Tooltip

	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			children := node.Children()
			if len(children) == 0 {
				return layout.Dimensions{}
			}

			contentNode := children[0].(layoutnode.NodeCoordinator)

			// Initialize
			if gioTooltip == nil {
				gioTooltip = mwtooltip.PlainTooltip(text)
			} else if gioTooltip.SupportingText != text {
				gioTooltip.SupportingText = text
			}

			// Measure content
			macroContent := op.Record(gtx.Ops)
			dimsContent := contentNode.Layout(gtx)
			callContent := macroContent.Stop()

			// Hover detection
			checkHover(gtx, dimsContent.Size, hovered)

			// Layout gio-mw tooltip
			return gioTooltip.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				callContent.Add(gtx.Ops)
				return dimsContent
			}, *hovered)
		}
	})
}

func checkHover(gtx layout.Context, size image.Point, tag *bool) {
	// Register input op
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	event.Op(gtx.Ops, tag)

	// Check events
	for {
		ev, ok := gtx.Event(pointer.Filter{Target: tag, Kinds: pointer.Enter | pointer.Leave})
		if !ok {
			break
		}
		if x, ok := ev.(pointer.Event); ok {
			switch x.Kind {
			case pointer.Enter:
				*tag = true
			case pointer.Leave:
				*tag = false
			}
		}
	}
}

// TODO: Known Issue - Tooltips on Clickable elements (like Buttons) do not work correctly
// because the Clickable component 'grabs' the pointer events, preventing the Tooltip
// from receiving Enter/Leave events.
// Future exploration: Investigate using pointer.PassOp to allow events to pass through
// or implement a centralized hover manager.
