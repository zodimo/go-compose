package column

import (
	"go-compose-dev/internal/layoutnode"
	"go-compose-dev/internal/modifiers/weight"

	"gioui.org/layout"
)

func DefaultColumnOptions() ColumnOptions {
	return ColumnOptions{
		Modifier:  EmptyModifier,
		Spacing:   SpaceEnd, // 0
		Alignment: Start,    // 0
	}
}

func Column(content Composable, options ...ColumnOption) Composable {
	opts := DefaultColumnOptions()
	for _, option := range options {
		option(&opts)
	}
	return func(c Composer) Composer {
		c.StartBlock("Column")
		c.Modifier(func(modifier Modifier) Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.WithComposable(content)
		c.SetWidgetConstructor(columnWidgetConstructor(opts))

		return c.EndBlock()
	}
}

func columnWidgetConstructor(options ColumnOptions) layoutnode.LayoutNodeWidgetConstructor {
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
				Axis:      layout.Vertical,
				Spacing:   options.Spacing,
				Alignment: options.Alignment,
			}.Layout(gtx, flexedChildren...)
		}
	})

}
