package layoutnode

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

	IsEmpty() bool
}

type LayoutModifierNode interface {
	LayoutNode
	AttachLayoutModifier(attach func(gtx LayoutContext, widget LayoutWidget) LayoutWidget)
}

type DrawModifierNode interface {
	LayoutNode
	AttachDrawModifier(attach func(gtx LayoutContext, widget LayoutWidget) LayoutWidget)
}

type PointerModifierNode interface {
	LayoutNode
	AttachPointerModifier(attach func(gtx LayoutContext, widget LayoutWidget) LayoutWidget)
}

// Wrapper to allow for inner node expansion
type NodeCoordinator interface {
	LayoutNode
	DrawModifierNode
	LayoutModifierNode
	PointerModifierNode

	LayoutPhase(gtx LayoutContext)
	PointerPhase(gtx LayoutContext)
	DrawPhase(gtx LayoutContext)

	LayoutSelf(gtx LayoutContext) LayoutDimensions

	Elements() ElementStore // ECS style properties

	Expand()
}
