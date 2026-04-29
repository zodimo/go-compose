package pointer

import (
	"fmt"
	"image"

	"gioui.org/f32"
	"gioui.org/gesture"
	gioPointer "gioui.org/io/pointer"
	"gioui.org/op/clip"
	"github.com/zodimo/go-compose/compose/ui/geometry"
	composePointer "github.com/zodimo/go-compose/compose/ui/input/pointer"
	"github.com/zodimo/go-compose/compose/ui/platform"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/layoutnode"
	node "github.com/zodimo/go-compose/internal/node"
	"github.com/zodimo/go-compose/state"
)

var _ node.ChainNode = (*PointerInputNode)(nil)

type dragState struct {
	dragger      gesture.Drag
	isDragging   bool
	lastPosition f32.Point
}

type PointerInputNode struct {
	node.ChainNode
	key   any
	block func(scope composePointer.PointerInputScope)
}

func NewPointerInputNode(element PointerInputElement) node.ChainNode {
	return &PointerInputNode{
		ChainNode: node.NewChainNode(
			node.NewNodeID(),
			node.NodeKindPointerInput,
			node.PointerInputPhase,
			func(n node.TreeNode) {
				lno := n.(layoutnode.LayoutNode)
				key := lno.GenerateID()
				custKey := fmt.Sprintf("%d_%s", key, element.key)

				dragStatePath := fmt.Sprintf("%s/dragState", custKey)
				dragStateValue := state.MustRemember(
					lno,
					dragStatePath,
					func() *dragState { return &dragState{} },
				)

				no := n.(layoutnode.PointerInputModifierNode)
				no.AttachPointerInputModifier(func(widget layoutnode.LayoutWidget) layoutnode.LayoutWidget {
					return layoutnode.NewLayoutWidget(func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
						dims := widget.Layout(gtx)

						density := unit.DensityFromLayoutContext(gtx)
						viewConfiguration := platform.NewDefaultViewConfiguration()

						handler := composePointer.NewDragGestureDetectorScopeAndHandler()

						pointerInputScope := composePointer.NewPointerInputScope(
							density,
							unit.NewIntSize(dims.Size.X, dims.Size.Y),
							geometry.SizeZero,
							viewConfiguration,
							handler,
						)

						element.block(pointerInputScope)

						ds := dragStateValue.Get()

						area := clip.Rect(image.Rectangle{Max: dims.Size}).Push(gtx.Ops)
						defer area.Pop()

						ds.dragger.Add(gtx.Ops)

						for {
							evt, ok := ds.dragger.Update(gtx.Metric, gtx.Source, gesture.Both)
							if !ok {
								break
							}

							switch evt.Kind {
							case gioPointer.Press:
								ds.isDragging = true
								ds.lastPosition = evt.Position
								handler.DragStart(geometry.NewOffset(evt.Position.X, evt.Position.Y))

							case gioPointer.Drag:
								if ds.isDragging {
									delta := geometry.NewOffset(
										evt.Position.X-ds.lastPosition.X,
										evt.Position.Y-ds.lastPosition.Y,
									)
									ds.lastPosition = evt.Position

									change := composePointer.PointerInputChange{}
									handler.Drag(change, delta)
								}

							case gioPointer.Release:
								if ds.isDragging {
									ds.isDragging = false
									handler.DragEnd()
								}

							case gioPointer.Cancel:
								if ds.isDragging {
									ds.isDragging = false
									handler.DragCancel()
								}
							}
						}

						return dims
					})
				})
			},
		),
	}
}
