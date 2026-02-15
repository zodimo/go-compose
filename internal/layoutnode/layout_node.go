package layoutnode

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/state"
	"github.com/zodimo/go-maybe"
)

// The base Node for the Tree

type LayoutNodeWidgetConstructor interface {
	Make(node LayoutNode) GioLayoutWidget
}
type LayoutNode interface {
	TreeNode

	LayoutNodeChildren() []LayoutNode

	WithChildren(children []LayoutNode) LayoutNode

	Modifier(func(modifier ui.Modifier) ui.Modifier)
	UnwrapModifier() ui.Modifier

	WithSlotsAssoc(k string, v any) LayoutNode // this can be better
	FindSlot(k string) maybe.Maybe[any]

	GenerateID() Identifier
	ResetIdentifierKeyCounter()

	state.SupportState

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

type PointerInputModifierNode interface {
	LayoutNode
	AttachPointerInputModifier(attach func(widget LayoutWidget) LayoutWidget)
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
	PointerInputModifierNode
	ParentDataModifierNode

	PointerPhase(gtx LayoutContext)

	Elements() ElementStore

	Expand()
}
