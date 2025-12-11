package size

import (
	node "go-compose-dev/internal/Node"
	"go-compose-dev/internal/layoutnode"
	"image"

	"gioui.org/layout"
	"gioui.org/op"
)

var _ ChainNode = (*SizeNode)(nil)

// NodeKind should also implement the interface of the LayoutNode for that phase

type SizeNode struct {
	ChainNode
	size SizeData
}

// NewSizeNode creates a new size node
func NewSizeNode(sizeData SizeData) ChainNode {
	return SizeNode{
		ChainNode: node.NewChainNode(
			node.NewNodeID(),
			node.NodeKindLayout,
			node.LayoutPhase,
			//OnAttach
			func(n TreeNode) {
				// how should the tree now be updated when attached
				// tree nde is the layout tree

				no := n.(layoutnode.LayoutModifierNode)
				// we can now work with the layoutNode
				no.AttachLayoutModifier(func(widget layoutnode.LayoutWidget) layoutnode.LayoutWidget {
					return layoutnode.NewLayoutWidget(
						func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
							// 1. Calculate constraints to pass to child.
							childConstraints := ApplySizeDataToConstraints(gtx.Constraints, sizeData)

							// 2. Measure child.
							macro := op.Record(gtx.Ops)
							// Create a context with modified constraints for the child
							childGtx := gtx
							childGtx.Constraints = childConstraints
							childDims := widget.Layout(childGtx)
							call := macro.Stop()

							// 3. Determine my size.
							mySize := image.Point{
								X: childDims.Size.X,
								Y: childDims.Size.Y,
							}

							// Handle Width overrides
							if sizeData.Width != NotSet {
								// Fixed width overrides child measurement
								mySize.X = sizeData.Width
							} else if sizeData.FillMaxWidth || sizeData.FillMax {
								// Fill behavior uses max constraints
								mySize.X = gtx.Constraints.Max.X
							} else {
								// Default/Wrap behavior: respect incoming constraints
								// If we are wrapping, we wanted min=0 for child, but our size
								// must still respect our parent's min constraints.
								mySize.X = Clamp(mySize.X, gtx.Constraints.Min.X, gtx.Constraints.Max.X)
							}

							// Handle Height overrides
							if sizeData.Height != NotSet {
								mySize.Y = sizeData.Height
							} else if sizeData.FillMaxHeight || sizeData.FillMax {
								mySize.Y = gtx.Constraints.Max.Y
							} else {
								mySize.Y = Clamp(mySize.Y, gtx.Constraints.Min.Y, gtx.Constraints.Max.Y)
							}

							// 4. Align
							if sizeData.Alignment != nil {
								// Calculate offset
								offset := sizeData.Alignment.Align(childDims.Size, mySize, layoutnode.LayoutDirectionLTR)

								// Apply offset
								// We put the offset operation before replaying the child recording
								defer op.Offset(offset).Push(gtx.Ops).Pop()
							}

							// Add the child operations
							call.Add(gtx.Ops)

							return layout.Dimensions{
								Size: mySize,
								// We should probably merge baselines here if needed, but keeping simple for now
								Baseline: childDims.Baseline,
							}
						},
					)
				})

			},
		),
		size: sizeData,
	}
}

func GetSizeConstraintsAndSizeData(constraints layout.Constraints, sizeData SizeData) image.Point {
	// This function seems to be legacy or used for strict size calculation.
	// The logic is now embedded in NewSizeNode.
	// We keep it for backward compatibility if used elsewhere,
	// but purely based on constraints (ignoring child).
	size := image.Point{
		X: constraints.Min.X,
		Y: constraints.Min.Y,
	}

	if sizeData.Width != NotSet {
		if sizeData.Required {
			size.X = sizeData.Width
		} else {
			size.X = Clamp(sizeData.Width, constraints.Min.X, constraints.Max.X)
		}
	}
	if sizeData.Height != NotSet {
		if sizeData.Required {
			size.Y = sizeData.Height
		} else {
			size.Y = Clamp(sizeData.Height, constraints.Min.Y, constraints.Max.Y)
		}
	}

	if sizeData.FillMaxWidth {
		size.X = constraints.Max.X
	}
	if sizeData.FillMaxHeight {
		size.Y = constraints.Max.Y
	}

	if sizeData.FillMax {
		size.X = constraints.Max.X
		size.Y = constraints.Max.Y
	}
	return size

}

func ApplySizeDataToConstraints(constraints layout.Constraints, sizeData SizeData) layout.Constraints {

	// Start with incoming constraints
	c := constraints

	// Apply Fixed Width/Height Logic to Min/Max
	if sizeData.Width != NotSet {
		if sizeData.Required {
			c.Min.X = sizeData.Width
			c.Max.X = sizeData.Width
		} else {
			c.Min.X = Clamp(sizeData.Width, c.Min.X, c.Max.X)
			c.Max.X = Clamp(sizeData.Width, c.Min.X, c.Max.X)
		}
	}
	if sizeData.Height != NotSet {
		if sizeData.Required {
			c.Min.Y = sizeData.Height
			c.Max.Y = sizeData.Height
		} else {
			c.Min.Y = Clamp(sizeData.Height, c.Min.Y, c.Max.Y)
			c.Max.Y = Clamp(sizeData.Height, c.Min.Y, c.Max.Y)
		}
	}

	// Fill Logic
	if sizeData.FillMaxWidth {
		c.Min.X = c.Max.X
	}
	if sizeData.FillMaxHeight {
		c.Min.Y = c.Max.Y
	}

	if sizeData.FillMax {
		c.Min.X = c.Max.X
		c.Min.Y = c.Max.Y
	}

	// Wrap Logic overrides Min constraints to allow shrinking
	if sizeData.WrapWidth {
		c.Min.X = 0
	}
	if sizeData.WrapHeight {
		c.Min.Y = 0
	}

	// Unbounded Logic overrides Max constraints
	if sizeData.Unbounded {
		// Use a large constant for infinity
		const Inf = 1e6
		if sizeData.WrapWidth {
			c.Max.X = Inf
		}
		if sizeData.WrapHeight {
			c.Max.Y = Inf
		}
		// Or if generalized unbounded:
		// c.Max.X = Inf
		// c.Max.Y = Inf
	}

	return c
}
