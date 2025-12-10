package size

import (
	node "go-compose-dev/internal/Node"
	"go-compose-dev/internal/layoutnode"
	"image"

	"gioui.org/layout"
)

var _ ChainNode = (*SizeNode)(nil)

// NodeKind should also implement the interface of the LayoutNode for that phase

type SizeNode struct {
	ChainNode
	size SizeData
}

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
							size := GetSizeConstraintsAndSizeData(gtx.Constraints, sizeData)
							// if size.
							childConstraints := gtx.Constraints
							childConstraints = ApplySizeDataToConstraints(childConstraints, sizeData)
							gtx.Constraints = childConstraints
							widget.Layout(gtx)
							return layout.Dimensions{
								Size: size,
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

	// fmt.Printf("ApplySizeOptionsToConstraints: constraints before: %v\n", constraints)
	// fmt.Printf("ApplySizeOptionsToConstraints: opts: %s\n", opts)

	if sizeData.Width != NotSet {
		if sizeData.Required {
			constraints.Min.X = sizeData.Width
			constraints.Max.X = sizeData.Width
		} else {
			constraints.Min.X = Clamp(sizeData.Width, constraints.Min.X, constraints.Max.X)
			constraints.Max.X = Clamp(sizeData.Width, constraints.Min.X, constraints.Max.X)
		}

	}
	if sizeData.Height != NotSet {
		if sizeData.Required {
			constraints.Min.Y = sizeData.Height
			constraints.Max.Y = sizeData.Height
		} else {
			constraints.Min.Y = Clamp(sizeData.Height, constraints.Min.Y, constraints.Max.Y)
			constraints.Max.Y = Clamp(sizeData.Height, constraints.Min.Y, constraints.Max.Y)
		}

	}

	if sizeData.FillMaxWidth {
		constraints.Min.X = constraints.Max.X
	}
	if sizeData.FillMaxHeight {
		constraints.Min.Y = constraints.Max.Y
	}

	if sizeData.FillMax {
		constraints.Min.X = constraints.Max.X
		constraints.Min.Y = constraints.Max.Y
	}
	return constraints
}
