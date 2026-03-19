package window

import (
	"image"

	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/pkg/api"

	"gioui.org/gesture"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
)

// PopupAlignment determines how the popup is aligned relative to its anchor point.
// For now, we assume the anchor point is the top-left of where this Popup is called.
type PopupAlignment int

const (
	AlignTopLeft PopupAlignment = iota
)

// Popup shows content overlaid on top of other content.
// It uses op.Defer to effectively "break out" of the z-order (though it respects parent clipping).
// The content is laid out with loose constraints (0 to Max).
func Popup(
	content api.Composable,
	options ...PopupOption,
) api.Composable {
	return func(c api.Composer) api.Composer {
		opts := DefaultPopupOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		c.StartBlock("Popup")
		c.WithComposable(content)
		c.SetWidgetConstructor(popupWidgetConstructor(opts))
		return c.EndBlock()
	}
}

func popupWidgetConstructor(opts PopupOptions) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			popupClick := node.State("popupClick", func() any { return &gesture.Click{} }).Get().(*gesture.Click)
			dismissClick := node.State("dismissClick", func() any { return &gesture.Click{} }).Get().(*gesture.Click)

			// Handle dismiss events
			if opts.OnDismissRequest != nil {
				for {
					e, ok := dismissClick.Update(gtx.Source)
					if !ok {
						break
					}
					if e.Kind == gesture.KindClick {
						opts.OnDismissRequest()
					}
				}
			}

			// 1. Record the content layout
			contentMacro := op.Record(gtx.Ops)

			// Prepare context for popup content
			pGtx := gtx
			pGtx.Constraints.Min = image.Point{}

			// Apply offset if needed
			xPx := pGtx.Dp(unit.DpToGioUnitUnsafe(opts.OffsetX))
			yPx := pGtx.Dp(unit.DpToGioUnitUnsafe(opts.OffsetY))

			op.Offset(image.Pt(xPx, yPx)).Add(pGtx.Ops)

			// Layout children and track content size
			var contentSize image.Point
			for _, child := range node.Children() {
				childLayoutNode := child.(layoutnode.NodeCoordinator)
				dims := childLayoutNode.Layout(pGtx)
				if dims.Size.X > contentSize.X {
					contentSize.X = dims.Size.X
				}
				if dims.Size.Y > contentSize.Y {
					contentSize.Y = dims.Size.Y
				}
			}

			// Add blocking handler at the popup content area
			if contentSize.X > 0 && contentSize.Y > 0 {
				passStack := pointer.PassOp{}.Push(pGtx.Ops)
				stack := clip.Rect{Max: contentSize}.Push(pGtx.Ops)
				popupClick.Add(pGtx.Ops)
				event.Op(pGtx.Ops, node)
				stack.Pop()
				passStack.Pop()
			}

			contentCall := contentMacro.Stop()

			// 2. Wrap in Scrim + Content
			finalMacro := op.Record(gtx.Ops)

			// Add Scrim (Dismiss Layer)
			if opts.OnDismissRequest != nil {
				// Large rect covering screen
				maxSize := 50000 // Arbitrary large size
				rectSize := image.Pt(maxSize, maxSize)
				// Center the large rect on the anchor
				offset := image.Pt(-maxSize/2, -maxSize/2)

				op.Offset(offset).Add(gtx.Ops)

				passStack := pointer.PassOp{}.Push(gtx.Ops)
				fullRect := clip.Rect{Max: rectSize}
				clipStack := fullRect.Push(gtx.Ops)
				dismissClick.Add(gtx.Ops)
				event.Op(gtx.Ops, dismissClick)
				clipStack.Pop()
				passStack.Pop()

				// Restore offset for content
				op.Offset(offset.Mul(-1)).Add(gtx.Ops)
			}

			// Add Content (on top of scrim)
			contentCall.Add(gtx.Ops)

			finalCall := finalMacro.Stop()

			// 3. Defer the execution (draws on top of UI)
			op.Defer(gtx.Ops, finalCall)

			return layout.Dimensions{}
		}
	})
}
