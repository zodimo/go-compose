package column

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/modifiers/weight"

	"gioui.org/layout"
)

func DefaultColumnOptions() ColumnOptions {
	return ColumnOptions{
		Modifier:  ui.EmptyModifier,
		Spacing:   SpaceEnd, // 0
		Alignment: Start,    // 0
	}
}

func Column(content Composable, options ...ColumnOption) Composable {
	opts := DefaultColumnOptions()
	for _, option := range options {
		if option == nil {
			continue
		}
		option(&opts)
	}
	return func(c Composer) Composer {
		c.StartBlock("Column")
		c.Modifier(func(modifier ui.Modifier) ui.Modifier {
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

			// Capture backend before entering layout.Flex callbacks
			backend := gtx.DrawBackend()

			flexedChildren := []layout.FlexChild{}
			for _, child := range node.Children() {

				childLayoutNode := child.(layoutnode.NodeCoordinator)

				elementStore := childLayoutNode.Elements()

				maybeWeightElement := elementStore.GetElement(weight.WeightElementKey)
				if maybeWeightElement.IsNone() {
				flexedChildren = append(flexedChildren, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// Compose behavior: Cross axis constraints Min is 0
					gtx.Constraints.Min.X = 0
					return layoutnode.ToGioDimensions(childLayoutNode.Layout(layoutnode.NewLayoutContextWithBackend(&gtx, backend)))
				}))
				} else {
					weightElement := maybeWeightElement.UnwrapUnsafe().(weight.WeightElement)
				flexedChildren = append(flexedChildren, layout.Flexed(weightElement.WeightData().Weight, func(gtx layout.Context) layout.Dimensions {
					return layoutnode.ToGioDimensions(childLayoutNode.Layout(layoutnode.NewLayoutContextWithBackend(&gtx, backend)))
				}))
				}
			}

		return layoutnode.FromGioDimensions(layout.Flex{
			Axis:      layout.Vertical,
			Spacing:   layout.Spacing(options.Spacing),
			Alignment: layout.Alignment(options.Alignment),
		}.Layout(*gtx.ToGio(), flexedChildren...))
		}
	})

}
