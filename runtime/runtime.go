package runtime

import "gioui.org/op"

// Runtime is the interface responsible for executing the layout tree.
// It bridges the gap between the abstract LayoutNode tree produced by composition
// and the actual drawing operations.
type Runtime interface {
	// Run executes the layout logic for the given LayoutNode within the provided LayoutContext.
	// It returns an op.CallOp that contains the drawing operations for the node and its children.
	Run(LayoutContext, LayoutNode) op.CallOp
}
