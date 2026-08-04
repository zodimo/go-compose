package row

// Spacing determines the spacing mode for a Flex layout.
// Defined as a framework-owned type to avoid leaking gioui types.
type Spacing uint8

const (
	// SpaceEnd leaves space at the end.
	SpaceEnd Spacing = iota
	// SpaceStart leaves space at the start.
	SpaceStart
	// SpaceSides shares space between the start and end.
	SpaceSides
	// SpaceAround distributes space evenly between children,
	// with half as much space at the start and end.
	SpaceAround
	// SpaceBetween distributes space evenly between children,
	// leaving no space at the start and end.
	SpaceBetween
	// SpaceEvenly distributes space evenly between children and
	// at the start and end.
	SpaceEvenly
)

// Alignment is the mutual alignment of a list of widgets.
// Defined as a framework-owned type to avoid leaking gioui types.
type Alignment uint8

const (
	Start Alignment = iota
	End
	Middle
	Baseline
)
