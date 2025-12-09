package layoutnode

// The base Node for the Tree
type LayoutNode interface {
	TreeNode

	LayoutNodeChildren() []LayoutNode

	WithChildren(children []LayoutNode) LayoutNode

	Modifier(func(modifier Modifier) Modifier)

	WithSlotsAssoc(k string, v any) LayoutNode

	IsEmpty() bool
}

type WidgetReceiver func(widget LayoutWidget)

type LayoutModifierNode interface {
	LayoutNode
	AttachLayoutModifier(attach func(gtx LayoutContext, widget LayoutWidget) LayoutDimensions)
}

type DrawModifierNode interface {
	LayoutNode
	AttachDrawModifier(attach func(gtx LayoutContext, widget LayoutWidget))
}

type NodeCoordinator interface {
	LayoutNode
	AttachLayoutModifier(attach func(gtx LayoutContext, widget LayoutWidget) LayoutDimensions)
	DrawModifier(attach func(gtx LayoutContext, widget LayoutWidget))
}
