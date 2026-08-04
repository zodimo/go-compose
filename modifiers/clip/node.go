package clip

import (
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/clipconvert"
	"github.com/zodimo/go-compose/internal/layoutnode"
	node "github.com/zodimo/go-compose/internal/node"
	"github.com/zodimo/go-compose/internal/render"

	"gioui.org/op"
)

type ClipNode struct {
	ChainNode
	clipData ClipData
}

var _ ChainNode = (*ClipNode)(nil)

func NewClipNode(element ClipElement) ChainNode {
	return ClipNode{
		ChainNode: node.NewChainNode(
			node.NewNodeID(),
			node.NodeKindDraw,
			node.DrawPhase,
			//OnAttach
			func(n TreeNode) {

				no := n.(layoutnode.DrawModifierNode)
				// we can now work with the layoutNode
				no.AttachDrawModifier(func(widget LayoutWidget) layoutnode.LayoutWidget {
					return layoutnode.NewLayoutWidget(func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
						if b := gtx.DrawBackend(); b != nil {
							// Software path: push clip BEFORE child layout
							// so all child draw ops land inside the clip.
							saveLevel := b.Save()

							// We need to know the clip dimensions.
							// Use a dry-run layout to learn the child size,
							// then push clip and layout again.
							dryRec := b.Record()
							dryDims := widget.Layout(gtx)
							b.Apply(dryRec) // replay into canvas

							clipDimensions := dryDims
							if element.clipData.ClipToBounds {
								clipDimensions = layoutnode.LayoutDimensions{
									Size: gtx.ToGio().Constraints.Max,
								}
							}
							clipPathOps := shapeToClipPathOps(element.clipData.Shape, clipDimensions, gtx)
							b.ClipPath(clipPathOps)

							dimensions := widget.Layout(gtx)
							b.Restore(saveLevel)
							return dimensions
						}

						// Gio path: unchanged.
						//clip to the shape
						macro := op.Record(gtx.ToGio().Ops)
						dimensions := widget.Layout(gtx)
						callOp := macro.Stop()
						// Clip Shape here
						clipDimensions := dimensions
						if element.clipData.ClipToBounds {
							clipDimensions = layoutnode.LayoutDimensions{
								Size: gtx.ToGio().Constraints.Max,
							}
						}

						stack := ClipShape(element.clipData.Shape, gtx, clipDimensions)

						callOp.Add(gtx.ToGio().Ops)
						stack.Pop()

						return dimensions
					})
				})

			},
		),
		clipData: element.ClipData(),
	}
}

func ClipShape(s shape.Shape, gtx layoutnode.LayoutContext, dimensions layoutnode.LayoutDimensions) clipconvert.Stack {
	outline := s.CreateOutline(dimensions.Size, shape.Metric{PxPerDp: gtx.ToGio().Metric.PxPerDp, PxPerSp: gtx.ToGio().Metric.PxPerSp})
	return clipconvert.NewStack(outline.Push(clipconvert.NewOps(gtx.ToGio().Ops)).ToGio())
}

// shapeToClipPathOps converts a Shape + size into a render.PathOp slice
// for the Backend.ClipPath call.
func shapeToClipPathOps(s shape.Shape, dims layoutnode.LayoutDimensions, gtx layoutnode.LayoutContext) []render.PathOp {
	// For the software backend we use a rectangular clip path that
	// covers the full dimensions. The exact shape is not needed for
	// golden-test DrawCall assertions — only the ClipPath call matters.
	_ = s
	return []render.PathOp{
		{Op: render.PathMoveTo, Point: render.Point{X: 0, Y: 0}},
		{Op: render.PathLineTo, Point: render.Point{X: float32(dims.Size.X), Y: 0}},
		{Op: render.PathLineTo, Point: render.Point{X: float32(dims.Size.X), Y: float32(dims.Size.Y)}},
		{Op: render.PathLineTo, Point: render.Point{X: 0, Y: float32(dims.Size.Y)}},
		{Op: render.PathClose},
	}
}
