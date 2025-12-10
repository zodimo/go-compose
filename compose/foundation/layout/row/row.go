package row

import (
	"go-compose-dev/internal/layoutnode"
	"go-compose-dev/internal/modifiers/weight"

	"gioui.org/layout"
)

func DefaultRowOptions() RowOptions {
	return RowOptions{
		Modifier:  EmptyModifier,
		Spacing:   SpaceEnd, // 0
		Alignment: Start,    // 0
	}
}

func Row(content Composable, options ...RowOption) Composable {
	opts := DefaultRowOptions()
	for _, option := range options {
		option(&opts)
	}
	return func(c Composer) Composer {
		c.StartBlock("Row")
		c.Modifier(func(modifier Modifier) Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.WithComposable(content)
		c.SetWidgetConstructor(rowWidgetConstructor(opts))

		return c.EndBlock()
	}
}

func rowWidgetConstructor(options RowOptions) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			flexedChildren := []layout.FlexChild{}
			for _, child := range node.Children() {

				childLayoutNode := child.(layoutnode.NodeCoordinator)

				elementStore := childLayoutNode.Elements()

				maybeWeightElement := elementStore.GetElement(weight.WeightElementKey)
				if maybeWeightElement.IsNone() {
					flexedChildren = append(flexedChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return childLayoutNode.Layout(gtx)
					}))
				} else {
					weightElement := maybeWeightElement.UnwrapUnsafe().(weight.WeightElement)
					flexedChildren = append(flexedChildren, layout.Flexed(weightElement.WeightData().Weight, func(gtx layout.Context) layout.Dimensions {
						return childLayoutNode.Layout(gtx)
					}))
				}
			}

			return layout.Flex{
				Axis:      layout.Horizontal,
				Spacing:   options.Spacing,
				Alignment: options.Alignment,
			}.Layout(gtx, flexedChildren...)
		}
	})

}
