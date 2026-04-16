package lazy

import (
	"fmt"
	"image"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/internal/layoutnode"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
)

type C = layout.Context
type D = layout.Dimensions

type lazyChildInfo struct {
	size int
}

// LazyColumn is a vertically scrolling list that only composes and lays out currently visible items.
// Note: Current implementation purely implements Lazy Layout (eager composition).
func LazyColumn(content func(LazyListScope), options ...LazyListOption) compose.Composable {
	return lazyList(layout.Vertical, content, options...)
}

// LazyRow is a horizontally scrolling list that only composes and lays out currently visible items.
// Note: Current implementation purely implements Lazy Layout (eager composition).
func LazyRow(content func(LazyListScope), options ...LazyListOption) compose.Composable {
	return lazyList(layout.Horizontal, content, options...)
}

func lazyList(axis layout.Axis, content func(LazyListScope), options ...LazyListOption) compose.Composable {
	return func(c compose.Composer) compose.Composer {
		opts := DefaultLazyListOptions()
		for _, opt := range options {
			opt(&opts)
		}

		// Ensure state is initialized
		// Note: Ideally state should be passed by user. If not, we create a local one,
		// but since we can't persist it easily without ID generation here (or using Remember internally),
		// we should encourage passing state.
		// For now, if state is nil, we create a new one but it won't be persisted across recompositions
		// unless the user provided a state from RememberLazyListState.
		// To fix this auto-persistence if nil, we would need to look it up in state.
		if opts.State == nil {
			id := c.GenerateID()
			path := c.GetPath()
			key := fmt.Sprintf("%d/%s/lazyListState", id, path)
			// Try to get existing or create new
			opts.State = c.State(key, func() any { return NewLazyListState() }).Get().(*LazyListState)
		}

		c.StartBlock("LazyList")
		c.Modifier(func(m ui.Modifier) ui.Modifier {
			return m.Then(opts.Modifier)
		})

		scope := &lazyListScopeImpl{}
		content(scope)

		// Emit all items as children
		for _, item := range scope.items {
			c.WithComposable(item.Content)
		}

		// Extract sticky indices
		var stickyIndices []int
		for i, item := range scope.items {
			if item.IsSticky {
				stickyIndices = append(stickyIndices, i)
			}
		}

		c.SetWidgetConstructor(lazyListWidgetConstructor(opts.State, axis, stickyIndices, opts.Scrollbar))

		return c.EndBlock()
	}
}

func lazyListWidgetConstructor(state *LazyListState, axis layout.Axis, stickyIndices []int, scrollbar bool) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			state.List.List.Axis = axis

			itemSizes := make(map[int]int)

			dims := state.List.List.Layout(gtx, len(node.Children()), func(gtx C, i int) D {
				if i < 0 || i >= len(node.Children()) {
					return D{}
				}
				child := node.Children()[i].(layoutnode.NodeCoordinator)
				d := child.Layout(gtx)

				size := d.Size.Y
				if axis == layout.Horizontal {
					size = d.Size.X
				}
				itemSizes[i] = size
				return d
			})

			if scrollbar {
				layoutScrollbar(gtx, &state.List, axis, len(node.Children()), dims)
			}

			// Handle Sticky Header
			if len(stickyIndices) > 0 {
				first := state.List.List.Position.First

				// Find the active sticky header (last one <= first)
				stickyIdx := -1
				for _, idx := range stickyIndices {
					if idx <= first {
						stickyIdx = idx
					} else {
						break
					}
				}

				if stickyIdx != -1 {
					// Check for processing next sticky header to push up
					headerOffset := 0

					// Find next sticky header index
					nextStickyIdx := -1
					for _, idx := range stickyIndices {
						if idx > first { // simplified: we just need the first one > first?
							// No, strictly > stickyIdx is enough?
							// Actually it must be > first to be visible below?
							// Logic: next sticky header that is "colliding".
							// It must be visible.
							nextStickyIdx = idx
							break
						}
					}
					// Wait, finding nextStickyIdx from the full list is correct.
					// But we specifically care if it is *visible* and pushing up.
					// Iterate to find next sticky index > stickyIdx (not just > first).
					// Actually, if active is 3 (first=5), next sticky could be 7.
					// If 7 is visible, we calculate its position.

					// Let's refine nextSticky finding:
					nextStickyIdx = -1
					for _, idx := range stickyIndices {
						if idx > stickyIdx {
							nextStickyIdx = idx
							break
						}
					}

					// Refined approach:
					// Always measure active sticky header first.
					// Reset min constraints to allow header to be smaller than the list height
					headerGtx := gtx
					if axis == layout.Vertical {
						headerGtx.Constraints.Min.Y = 0
					} else {
						headerGtx.Constraints.Min.X = 0
					}

					macro := op.Record(gtx.Ops)
					headerNode := node.Children()[stickyIdx].(layoutnode.NodeCoordinator)
					headerDims := headerNode.Layout(headerGtx)
					call := macro.Stop()

					headerSize := headerDims.Size.Y
					if axis == layout.Horizontal {
						headerSize = headerDims.Size.X
					}

					if nextStickyIdx != -1 {
						// Calculate position of nextStickyIdx relative to top
						// Start from First/Offset
						pos := state.List.List.Position.Offset
						found := false

						// We need to traverse from First to nextStickyIdx-1
						// But we only have sizes for *visible* items.
						// If nextStickyIdx is NOT visible, we can't calculate its exact position
						// (and it implies it's far away, so no push needed).

						// So we iterate and check if we have the size.
						current := first
						for current < nextStickyIdx {
							sz, ok := itemSizes[current]
							if !ok {
								// Item not visible/laid out, so nextSticky is far down
								found = false
								break
							}
							pos += sz
							current++
							// If pos becomes large, we can stop? No we need exact collision.
							// But usually viewport is finite.
						}

						if current == nextStickyIdx {
							// effectively found, pos is at the start of nextStickyIdx
							// Now we need the height of the CURRENT sticky header to know if they overlap.
							// We don't know the height of current sticky header unless we measure it!
							// Use a macro to measure it?
							// We will re-layout current sticky header anyway.
							found = true
						}

						if found {
							// Calculate overlap
							// Position of next sticky is `pos`.
							// Bottom of current sticky is `headerOffset + headerSize` (where headerOffset is usually 0).
							// Overlap amount = (headerOffset + headerSize) - pos?
							// Actually we want `headerOffset` such that `headerOffset + headerSize <= pos`.
							// So `headerOffset = min(0, pos - headerSize)`.

							if pos < headerSize {
								headerOffset = pos - headerSize
							}
						}
					}

					// Draw
					// apply offset
					defer clip.Rect{Max: dims.Size}.Push(gtx.Ops).Pop()

					pt := image.Pt(0, headerOffset)
					if axis == layout.Horizontal {
						pt = image.Pt(headerOffset, 0)
					}
					op.Offset(pt).Add(gtx.Ops)
					call.Add(gtx.Ops)

					return dims
				}
			}

			return dims
		}
	})
}
