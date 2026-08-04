package window

import (
	"image"

	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/pkg/api"

	"gioui.org/gesture"
	"gioui.org/io/event"
	"gioui.org/io/pointer"

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
				e, ok := dismissClick.Update(gtx.ToGio().Source)
					if !ok {
						break
					}
					if e.Kind == gesture.KindClick {
						opts.OnDismissRequest()
					}
				}
			}

			// 1. Record the content layout
			contentMacro := op.Record(gtx.ToGio().Ops)

			// Prepare context for popup content
			g := *gtx.ToGio()
			g.Constraints.Min = image.Point{}

			pGtx := layoutnode.NewLayoutContext(&g)

			// Apply offset if needed
			xPx := pGtx.ToGio().Dp(unit.DpToGioUnitUnsafe(opts.OffsetX))
			yPx := pGtx.ToGio().Dp(unit.DpToGioUnitUnsafe(opts.OffsetY))

			op.Offset(image.Pt(xPx, yPx)).Add(pGtx.ToGio().Ops)

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
				passStack := pointer.PassOp{}.Push(pGtx.ToGio().Ops)
				stack := clip.Rect{Max: contentSize}.Push(pGtx.ToGio().Ops)
				popupClick.Add(pGtx.ToGio().Ops)
				event.Op(pGtx.ToGio().Ops, node)
				stack.Pop()
				passStack.Pop()
			}

			contentCall := contentMacro.Stop()

			// 2. Wrap in Scrim + Content
			finalMacro := op.Record(gtx.ToGio().Ops)

			// Add Scrim (Dismiss Layer)
			if opts.OnDismissRequest != nil {
				// Large rect covering screen
				maxSize := 50000 // Arbitrary large size
				rectSize := image.Pt(maxSize, maxSize)
				// Center the large rect on the anchor
				offset := image.Pt(-maxSize/2, -maxSize/2)

				op.Offset(offset).Add(gtx.ToGio().Ops)

				passStack := pointer.PassOp{}.Push(gtx.ToGio().Ops)
				fullRect := clip.Rect{Max: rectSize}
				clipStack := fullRect.Push(gtx.ToGio().Ops)
				dismissClick.Add(gtx.ToGio().Ops)
				event.Op(gtx.ToGio().Ops, dismissClick)
				clipStack.Pop()
				passStack.Pop()

				// Restore offset for content
				op.Offset(offset.Mul(-1)).Add(gtx.ToGio().Ops)
			}

			// Add Content (on top of scrim)
			contentCall.Add(gtx.ToGio().Ops)

			finalCall := finalMacro.Stop()

			// 3. Defer the execution (draws on top of UI)
			op.Defer(gtx.ToGio().Ops, finalCall)

			return layoutnode.LayoutDimensions{}

		}
	})
}
