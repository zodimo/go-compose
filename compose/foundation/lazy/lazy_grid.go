package lazy

import (
	"fmt"
	"image"
	"image/color"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/layoutnode"

	"gioui.org/layout"
	"gioui.org/widget"
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
		c.Modifier(func(m ui.Modifier) ui.Modifier {
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
		c.SetWidgetConstructor(lazyGridWidgetConstructor(opts.State, axis, cells, opts.Scrollbar))

		return c.EndBlock()
	}
}

func lazyGridWidgetConstructor(
	state *LazyGridState,
	axis layout.Axis,
	cells GridCells,
	scrollbar bool,
) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			children := node.Children()
			itemCount := len(children)

			if itemCount == 0 {
				return D{}
			}

			var availableSpace int
			var spacing int = 0

			if axis == layout.Vertical {
				availableSpace = gtx.Constraints.Max.X
			} else {
				availableSpace = gtx.Constraints.Max.Y
			}

			cellCount := cells.calculateCrossAxisCellCount(availableSpace, spacing)
			cellSize := cells.calculateCellSize(availableSpace, cellCount, spacing)

			rowCount := (itemCount + cellCount - 1) / cellCount
			if rowCount == 0 {
				return D{}
			}

			state.List.List.Axis = axis

			dims := state.List.List.Layout(gtx, rowCount, func(gtx C, rowIndex int) D {
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
					capturedCellSizeInPx := unit.ToPixelsUnsafe(gtx, unit.Dp(float32(cellSize)))

					flexChildren = append(flexChildren, layout.Rigid(func(gtx C) D {
						// Constrain cell size in the cross-axis direction
						if axis == layout.Vertical {
							gtx.Constraints.Min.X = capturedCellSizeInPx
							gtx.Constraints.Max.X = capturedCellSizeInPx
						} else {
							gtx.Constraints.Min.Y = capturedCellSizeInPx
							gtx.Constraints.Max.Y = capturedCellSizeInPx
						}

						return capturedCoordinator.Layout(gtx)
					}))
				}

				// Add empty spacers for trailing empty cells to maintain alignment
				for i := endIdx - startIdx; i < cellCount; i++ {
					capturedCellSizeInPx := unit.ToPixelsUnsafe(gtx, unit.Dp(float32(cellSize)))
					flexChildren = append(flexChildren, layout.Rigid(func(gtx C) D {
						if axis == layout.Vertical {
							return D{Size: image.Point{X: capturedCellSizeInPx, Y: 0}}
						}
						return D{Size: image.Point{X: 0, Y: capturedCellSizeInPx}}
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

			if scrollbar {
				layoutGridScrollbar(gtx, &state.List, axis, rowCount, dims)
			}

			return dims
		}
	})
}

func layoutGridScrollbar(gtx layout.Context, list *widget.List, axis layout.Axis, rowCount int, dims layout.Dimensions) {
	majorAxisSize := axis.Convert(dims.Size).X
	start, end := fromListPosition(list.Position, rowCount, majorAxisSize)

	if !rangeIsScrollable(start, end) {
		return
	}

	sb := scrollbarStyle{
		scrollbar:  &list.Scrollbar,
		trackPad:   2,
		indicatorW: 6,
		indicatorR: 3,
		color:      color.NRGBA{R: 128, G: 128, B: 128, A: 150},
		hoverColor: color.NRGBA{R: 128, G: 128, B: 128, A: 200},
	}

	anchoring := layout.E
	if axis == layout.Horizontal {
		anchoring = layout.S
	}

	gtx.Constraints.Min = dims.Size
	gtx.Constraints.Max = dims.Size
	anchoring.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		list.Scrollbar.Update(gtx, axis, start, end)

		if list.Scrollbar.IndicatorHovered() {
			sb.color = sb.hoverColor
		}

		return sb.layout(gtx, axis, start, end)
	})

	if delta := list.ScrollDistance(); delta != 0 {
		list.List.ScrollBy(delta * float32(rowCount))
	}
}
