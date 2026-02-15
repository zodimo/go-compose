package layoutnode

import (
	"image"

	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/internal/immap"

	"github.com/zodimo/go-maybe"
)

func NewLayoutNode(id NodeID, key string, slotStore immap.ImmutableMap[any], persistentState PersistentState) LayoutNode {
	return &layoutNode{
		id:           id,
		key:          key,
		children:     []LayoutNode{},
		modifier:     ui.EmptyModifier,
		slots:        slotStore,
		state:        persistentState,
		layoutResult: maybe.None[LayoutResult](),
		idManager:    GetScopedIdentityManager("layout_nodes"),
	}
}

func NewNodeCoordinator(node LayoutNode) NodeCoordinator {

	outNode := &nodeCoordinator{
		LayoutNode:      node,
		elementStore:    EmptyElementStore,
		wrappedChildren: []TreeNode{},
	}
	widget := outNode.GetWidget()

	outNode.layoutCallChain = NewLayoutWidget(widget)
	outNode.pointerCallChain = NewLayoutWidget(widget)
	outNode.WrapChildren()
	return outNode
}

var ZeroGioLayoutWidget = func(gtx LayoutContext) LayoutDimensions {
	return LayoutDimensions{
		Size: image.Pt(0, 0),
	}
}
var ZeroLayoutWidget = NewLayoutWidget(ZeroGioLayoutWidget)

func NewLayoutWidget(innerWidget GioLayoutWidget) LayoutWidget {
	return layoutWidget{
		innerWidget: innerWidget,
	}
}

var EmptyWidgetConstructor = NewLayoutNodeWidgetConstructor(func(node LayoutNode) GioLayoutWidget {
	return ZeroGioLayoutWidget
})

var _ LayoutNodeWidgetConstructor = (*layoutNodeWidgetConstructor)(nil)

type layoutNodeWidgetConstructor struct {
	MakeFunc func(node LayoutNode) GioLayoutWidget
}

func NewLayoutNodeWidgetConstructor(makeFunc func(node LayoutNode) GioLayoutWidget) layoutNodeWidgetConstructor {
	return layoutNodeWidgetConstructor{
		MakeFunc: makeFunc,
	}
}

func NewDummyLayoutNodeWidgetConstructor() LayoutNodeWidgetConstructor {
	return NewLayoutNodeWidgetConstructor(func(node LayoutNode) GioLayoutWidget {
		return func(gtx LayoutContext) LayoutDimensions {
			return LayoutDimensions{}
		}
	})
}

func (c layoutNodeWidgetConstructor) Make(node LayoutNode) GioLayoutWidget {
	return c.MakeFunc(node)
}

func NewLayoutContextReceiver(widget GioLayoutWidget) LayoutContextReceiver {
	return func(gtx LayoutContext) {
		widget(gtx)
	}
}
