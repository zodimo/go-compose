package layoutnode

import (
	"fmt"
	"strings"
)

// DebugLayoutNode returns a formatted string representation of the layout node and its children
// for debugging purposes.
func DebugLayoutNode(node LayoutNode) string {
	layoutNode := node.(*layoutNode)
	var sb strings.Builder
	layoutNode.debugString(&sb, 0)
	return sb.String()
}

// debugString recursively builds the debug string representation
func (ln *layoutNode) debugString(sb *strings.Builder, depth int) {
	// Add indentation based on depth
	indent := strings.Repeat("  ", depth)

	// Node header
	sb.WriteString(fmt.Sprintf("%sLayoutNode{id: %s, key: %s", indent, ln.id, ln.key))

	// Add modifier info if present
	if ln.modifier != nil {
		sb.WriteString(fmt.Sprintf(", modifier: %T", ln.modifier))
	}

	// Add children count
	sb.WriteString(fmt.Sprintf(", children: %d", len(ln.children)))

	// Close the node
	sb.WriteString("}")

	// Add slots info if present
	if len(ln.slots) > 0 {
		sb.WriteString(" [")
		first := true
		for k, v := range ln.slots {
			if !first {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%s: %v", k, v))
			first = false
		}
		sb.WriteString("]")
	}
	sb.WriteString("\n")

	// Recursively add children
	for _, child := range ln.children {
		if childNode, ok := child.(*layoutNode); ok {
			childNode.debugString(sb, depth+1)
		} else {
			sb.WriteString(fmt.Sprintf("%s  %T{...}\n", indent, child))
		}
	}
}

// String implements fmt.Stringer for layoutNode
func (ln *layoutNode) String() string {
	return fmt.Sprintf("LayoutNode{id: %s, key: %s, children: %d}", ln.id, ln.key, len(ln.children))
}
