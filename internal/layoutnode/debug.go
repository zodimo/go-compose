package layoutnode

import (
	"fmt"
	"go-compose-dev/internal/modifier"
	"strings"
)

// DebugLayoutNode returns a formatted string representation of the layout node and its children
// for debugging purposes.
func DebugLayoutNode(node LayoutNode) string {
	if coordinator, ok := node.(*nodeCoordinator); ok {
		node = coordinator.LayoutNode
	}
	layoutNode := node.(*layoutNode)
	var sb strings.Builder
	layoutNode.debugString(&sb, 0)
	return sb.String()
}

// debugString recursively builds the debug string representation
func (ln *layoutNode) debugString(sb *strings.Builder, depth int) {
	// Add indentation based on depth
	indent := strings.Repeat("  ", depth)

	// Node header with basic info
	sb.WriteString(fmt.Sprintf("%s%s {\n", indent, ln.key))
	sb.WriteString(fmt.Sprintf("%s  ID: %s\n", indent, ln.id))

	// Add modifier info if present
	if ln.modifier != nil {
		modifierStr := modifier.DebugModifier(ln.modifier)
		// Indent each line of the modifier string
		modifierLines := strings.Split(modifierStr, "\n")
		sb.WriteString(fmt.Sprintf("%s  Modifiers: %s\n", indent, modifierLines[0]))
		for i := 1; i < len(modifierLines); i++ {
			sb.WriteString(fmt.Sprintf("%s    %s\n", indent, modifierLines[i]))
		}
	} else {
		sb.WriteString(fmt.Sprintf("%s  Modifiers: <none>\n", indent))
	}

	// Add slots info if present
	if len(ln.slots) > 0 {
		sb.WriteString(fmt.Sprintf("%s  Slots: {", indent))
		first := true
		for k, v := range ln.slots {
			if !first {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%s: %v", k, v))
			first = false
		}
		sb.WriteString("}\n")
	}

	// Add children count and preview
	sb.WriteString(fmt.Sprintf("%s  Children: %d", indent, len(ln.children)))
	// if len(ln.children) > 0 {
	// 	sb.WriteString(" [")
	// 	for i, child := range ln.children {
	// 		if i > 0 {
	// 			sb.WriteString(", ")
	// 		}
	// 		if childNode, ok := child.(*layoutNode); ok {
	// 			sb.WriteString(childNode.key)
	// 		} else {
	// 			sb.WriteString(fmt.Sprintf("%T", child))
	// 		}
	// 	}
	// 	sb.WriteString("]")
	// }
	sb.WriteString("\n" + indent + "}")

	// Recursively add children
	if len(ln.children) > 0 {
		// sb.WriteString("\n" + indent + "  Children details:")
		// for _, child := range ln.children {
		// 	sb.WriteString(fmt.Sprintf("%s%s", indent, DebugLayoutNode(child)))
		// }
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
}

// String implements fmt.Stringer for layoutNode
func (ln *layoutNode) String() string {
	return fmt.Sprintf("LayoutNode{id: %s, key: %s, children: %d}", ln.id, ln.key, len(ln.children))
}
