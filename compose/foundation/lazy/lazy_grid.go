package lazy

import (
	"fmt"
	"image"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/modifier"

	"gioui.org/layout"
)

// LazyVerticalGrid is a vertically scrolling grid that lays out items in columns.
// The columns parameter determines how items are arranged horizontally.
func LazyVerticalGrid(
	columns GridCells,
	content func(LazyGridScope),
	options ...LazyGridOption,
) compose.Composable {
	return lazyGrid(layout.Vertical, columns, content, options...)
}

// LazyHorizontalGrid is a horizontally scrolling grid that lays out items in rows.
// The rows parameter determines how items are arranged vertically.
func LazyHorizontalGrid(
	rows GridCells,
	content func(LazyGridScope),
	options ...LazyGridOption,
) compose.Composable {
	return lazyGrid(layout.Horizontal, rows, content, options...)
}

func lazyGrid(axis layout.Axis, cells GridCells, content func(LazyGridScope), options ...LazyGridOption) compose.Composable {
	return func(c compose.Composer) compose.Composer {
		opts := DefaultLazyGridOptions()
		for _, opt := range options {
			opt(&opts)
		}

		// Ensure state is initialized
		if opts.State == nil {
			id := c.GenerateID()
			path := c.GetPath()
			key := fmt.Sprintf("%d/%s/lazyGridState", id, path)
			opts.State = c.State(key, func() any { return NewLazyGridState() }).Get().(*LazyGridState)
		}

		c.StartBlock("LazyGrid")
		c.Modifier(func(m modifier.Modifier) modifier.Modifier {
			return m.Then(opts.Modifier)
		})

		// Collect all items
		scope := &lazyGridScopeImpl{}
		content(scope)

		// Emit all items as children (same as LazyList)
		for _, item := range scope.items {
			c.WithComposable(item.Content)
		}

		// Store cells and axis for widget constructor
		c.SetWidgetConstructor(lazyGridWidgetConstructor(opts.State, axis, cells))

		return c.EndBlock()
	}
}

func lazyGridWidgetConstructor(
	state *LazyGridState,
	axis layout.Axis,
	cells GridCells,
) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			children := node.Children()
			itemCount := len(children)

			if itemCount == 0 {
				return D{}
			}

			// Calculate the cross-axis cell count based on available space
			var availableSpace int
			var spacing int = 0 // Could be made configurable

			if axis == layout.Vertical {
				availableSpace = gtx.Constraints.Max.X
			} else {
				availableSpace = gtx.Constraints.Max.Y
			}

			cellCount := cells.calculateCrossAxisCellCount(availableSpace, spacing)
			cellSize := cells.calculateCellSize(availableSpace, cellCount, spacing)

			// Group items into rows (for vertical) or columns (for horizontal)
			rowCount := (itemCount + cellCount - 1) / cellCount
			if rowCount == 0 {
				return D{}
			}

			// Set up list axis
			state.List.List.Axis = axis

			return state.List.List.Layout(gtx, rowCount, func(gtx C, rowIndex int) D {
				// Calculate range of items for this row
				startIdx := rowIndex * cellCount
				endIdx := startIdx + cellCount
				if endIdx > itemCount {
					endIdx = itemCount
				}

				// Build row of cells using layout.Flex
				flexChildren := make([]layout.FlexChild, 0, cellCount)

				for i := startIdx; i < endIdx; i++ {
					child := children[i]
					childCoordinator := child.(layoutnode.NodeCoordinator)

					// Capture for closure
					capturedCoordinator := childCoordinator
					capturedCellSize := cellSize

					flexChildren = append(flexChildren, layout.Rigid(func(gtx C) D {
						// Constrain cell size in the cross-axis direction
						if axis == layout.Vertical {
							gtx.Constraints.Min.X = capturedCellSize
							gtx.Constraints.Max.X = capturedCellSize
						} else {
							gtx.Constraints.Min.Y = capturedCellSize
							gtx.Constraints.Max.Y = capturedCellSize
						}

						return capturedCoordinator.Layout(gtx)
					}))
				}

				// Add empty spacers for trailing empty cells to maintain alignment
				for i := endIdx - startIdx; i < cellCount; i++ {
					capturedCellSize := cellSize
					flexChildren = append(flexChildren, layout.Rigid(func(gtx C) D {
						if axis == layout.Vertical {
							return D{Size: image.Point{X: capturedCellSize, Y: 0}}
						}
						return D{Size: image.Point{X: 0, Y: capturedCellSize}}
					}))
				}

				// Layout the row (cross-axis to the scroll direction)
				var rowAxis layout.Axis
				if axis == layout.Vertical {
					rowAxis = layout.Horizontal
				} else {
					rowAxis = layout.Vertical
				}

				return layout.Flex{Axis: rowAxis}.Layout(gtx, flexChildren...)
			})
		}
	})
}
