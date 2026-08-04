package box

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/internal/layoutnode"

	"gioui.org/layout"
)

func Box(content Composable, options ...BoxOption) Composable {
	opts := DefaultBoxOptions()
	for _, option := range options {
		if option == nil {
			continue
		}
		option(&opts)
	}
	return func(c Composer) Composer {
		c.StartBlock("Box")
		c.Modifier(func(modifier ui.Modifier) ui.Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.WithComposable(content)
		c.SetWidgetConstructor(boxWidgetConstructor(opts))

		return c.EndBlock()
	}
}

func boxWidgetConstructor(options BoxOptions) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {

			// Capture backend before entering layout.Stack callbacks
			backend := gtx.DrawBackend()

			// Build framework-owned StackChild values.
			var gioChildren []layout.StackChild
			for _, child := range node.Children() {

				childLayoutNode := child.(layoutnode.NodeCoordinator)

				matchParent := childLayoutNode.Elements().GetElement(MatchParentSizeKey)

				if matchParent.IsSome() {
					gioChildren = append(gioChildren, layout.Expanded(func(gtx layout.Context) layout.Dimensions {
						// MatchParentSize implies matching the size of the container, which is passed in Min constraints
						// by the Stack layout for Expanded children.
						gtx.Constraints.Max = gtx.Constraints.Min
						return layoutnode.ToGioDimensions(childLayoutNode.Layout(layoutnode.NewLayoutContextWithBackend(&gtx, backend)))
					}))
				} else {
					gioChildren = append(gioChildren, layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return layoutnode.ToGioDimensions(childLayoutNode.Layout(layoutnode.NewLayoutContextWithBackend(&gtx, backend)))
					}))
				}
			}

			// Convert framework-owned Direction to gio layout.Direction at the seam.
			return layoutnode.FromGioDimensions(layout.Stack{
				Alignment: layout.Direction(options.Alignment),
			}.Layout(*gtx.ToGio(), gioChildren...))
		}
	})

}
