package box

import (
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/modifiers/box"

	"github.com/zodimo/go-compose/pkg/api"

	"gioui.org/layout"
)

type Composable = api.Composable
type Composer = api.Composer

var MatchParentSizeKey = box.MatchParentSizeKey

// Direction is the alignment of widgets relative to a containing
// space.
type Direction = layout.Direction

const (
	NW     Direction = layout.NW
	N      Direction = layout.N
	NE     Direction = layout.NE
	E      Direction = layout.E
	SE     Direction = layout.SE
	S      Direction = layout.S
	SW     Direction = layout.SW
	W      Direction = layout.W
	Center Direction = layout.Center
)

type Stack = layout.Stack
type StackChild = layout.StackChild
type LayoutContext = layoutnode.LayoutContext
type LayoutDimensions = layoutnode.LayoutDimensions
