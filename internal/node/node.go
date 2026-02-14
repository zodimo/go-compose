package node

// NodeKind identifies the type of modifier node.
// These are bit flags that can be combined.
type NodeKind uint16

const (
	NodeKindAny NodeKind = 1 << iota
	NodeKindLayout
	NodeKindDraw
	NodeKindSemantics
	NodeKindPointerInput
	NodeKindLocals
	NodeKindParentData
	NodeKindLayoutAware
	NodeKindGlobalPositionAware
	NodeKindIntermediateMeasure
	NodeKindFocusTarget
	NodeKindFocusProperties
	NodeKindFocusEvent
	NodeKindKeyInput
	NodeKindRotaryInput
	NodeKindCompositionLocalConsumer
)

func (nk NodeKind) String() string {
	switch nk {
	case NodeKindAny:
		return "NodeKindAny"
	case NodeKindLayout:
		return "NodeKindLayout"
	case NodeKindDraw:
		return "NodeKindDraw"
	case NodeKindSemantics:
		return "NodeKindSemantics"
	case NodeKindPointerInput:
		return "NodeKindPointerInput"
	case NodeKindLocals:
		return "NodeKindLocals"
	case NodeKindParentData:
		return "NodeKindParentData"
	case NodeKindLayoutAware:
		return "NodeKindLayoutAware"
	case NodeKindGlobalPositionAware:
		return "NodeKindGlobalPositionAware"
	case NodeKindIntermediateMeasure:
		return "NodeKindIntermediateMeasure"
	case NodeKindFocusTarget:
		return "NodeKindFocusTarget"
	case NodeKindFocusProperties:
		return "NodeKindFocusProperties"
	case NodeKindFocusEvent:
		return "NodeKindFocusEvent"
	case NodeKindKeyInput:
		return "NodeKindKeyInput"
	case NodeKindRotaryInput:
		return "NodeKindRotaryInput"
	case NodeKindCompositionLocalConsumer:
		return "NodeKindCompositionLocalConsumer"
	default:
		return "Unknown NodeKind"
	}
}

// Bitset of phases a node cares about.
type Phases uint8

const (
	LayoutPhase Phases = 1 << iota
	DrawPhase
	PointerInputPhase
	// SemanticsPhase, ParentDataPhase, etc.
)

type Node interface {
	GetID() NodeID
}

type TreeNode interface {
	Node

	Children() []TreeNode
}

type ChainNode interface {
	Node

	//NodeKind identifies the type of modifier node.
	Kind() NodeKind

	// Bitset of phases a node cares about.
	Phases() Phases

	// Attach is called when the node is attached to a node tree
	Attach(node TreeNode)

	// Detach is called when the node is detached from a node tree
	Detach()

	// Next returns the next node in the chain
	Next() ChainNode

	// SetNext sets the next node in the chain
	SetNext(node ChainNode)
}
