package box

import (
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/modifiers/box"

	"github.com/zodimo/go-compose/pkg/api"
)

type Composable = api.Composable
type Composer = api.Composer

var MatchParentSizeKey = box.MatchParentSizeKey

// Direction is the alignment of widgets relative to a containing space.
// Defined as a framework-owned type to avoid leaking gioui types.
type Direction uint8

const (
	NW     Direction = iota
	N
	NE
	E
	SE
	S
	SW
	W
	Center
)

// Stack lays out child elements on top of each other, according to an alignment direction.
type Stack struct {
	// Alignment is the direction to align children smaller than the available space.
	Alignment Direction
}

// StackChild represents a child for a Stack layout.
type StackChild struct {
	// Expanded indicates whether the child fills remaining space.
	Expanded bool
	// Widget is the layout function for this child.
	Widget func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions
}

type LayoutContext = layoutnode.LayoutContext
type LayoutDimensions = layoutnode.LayoutDimensions
