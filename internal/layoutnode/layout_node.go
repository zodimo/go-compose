package layoutnode

import "github.com/zodimo/go-maybe"

// The base Node for the Tree

type LayoutNodeWidgetConstructor interface {
	Make(node LayoutNode) GioLayoutWidget
}
type LayoutNode interface {
	TreeNode

	LayoutNodeChildren() []LayoutNode

	WithChildren(children []LayoutNode) LayoutNode

	Modifier(func(modifier Modifier) Modifier)
	UnwrapModifier() Modifier

	WithSlotsAssoc(k string, v Element) LayoutNode // this can be better

	GetWidget() GioLayoutWidget

	SetWidgetConstructor(constructor LayoutNodeWidgetConstructor)
	GetWidgetConstructor() LayoutNodeWidgetConstructor

	IsEmpty() bool

	SetLayoutResult(LayoutResult)
	GetLayoutResult() maybe.Maybe[LayoutResult]

	Layout(gtx LayoutContext) LayoutDimensions
	Draw(gtx LayoutContext) DrawOp
}

type LayoutModifierNode interface {
	LayoutNode
	AttachLayoutModifier(attach func(widget LayoutWidget) LayoutWidget)
}

type DrawModifierNode interface {
	LayoutNode
	AttachDrawModifier(attach func(widget LayoutWidget) LayoutWidget)
}

type PointerModifierNode interface {
	LayoutNode
	AttachPointerModifier(attach func(widget LayoutWidget) LayoutWidget)
}

type ParentDataModifierNode interface {
	LayoutNode
	AttachParentDataModifier(attach func(elements ElementStore) ElementStore)
}

// Wrapper to allow for inner node expansion
type NodeCoordinator interface {
	LayoutNode
	DrawModifierNode
	LayoutModifierNode
	PointerModifierNode
	ParentDataModifierNode

	PointerPhase(gtx LayoutContext)

	Elements() ElementStore

	Expand()
}
