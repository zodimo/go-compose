package clickable

import (
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/modifier"
	node "github.com/zodimo/go-compose/internal/node"
	internalclickable "github.com/zodimo/go-compose/internal/clickable"
	"github.com/zodimo/go-compose/modifiers/helpers"
)

type Element = modifier.Element
type InspectableModifier = modifier.InspectableModifier
type ModifierInspectorInfo = modifier.InspectorInfo

type Node = node.Node
type TreeNode = node.TreeNode
type ChainNode = node.ChainNode

type DrawModifierNode = layoutnode.DrawModifierNode
type LayoutModifierNode = layoutnode.LayoutModifierNode
type PointerModifierNode = layoutnode.PointerInputModifierNode

type LayoutContext = layoutnode.LayoutContext
type LayoutWidget = layoutnode.LayoutWidget

var ToNRGBA = helpers.ToNRGBA

// GioClickable wraps widget.Clickable behind a framework-owned type.
// Methods convert framework-owned LayoutContext to the engine layout.Context
// at the seam boundary.
type GioClickable struct {
	impl *internalclickable.Impl
}

// NewGioClickable creates a new GioClickable backed by widget.Clickable.
func NewGioClickable() *GioClickable {
	return &GioClickable{impl: internalclickable.New()}
}

// Clicked reports whether the clickable was clicked in this frame.
func (gc *GioClickable) Clicked(gtx layoutnode.LayoutContext) bool {
	return gc.impl.Clicked(*gtx.ToGio())
}

// Hovered reports whether the pointer is hovering over the clickable.
func (gc *GioClickable) Hovered() bool {
	return gc.impl.Hovered()
}

// Pressed reports whether the pointer is pressed down on the clickable.
func (gc *GioClickable) Pressed() bool {
	return gc.impl.Pressed()
}
