package layout

// Axis is the Horizontal or Vertical direction.
type Axis uint8

const (
	AxisUnspecified Axis = iota
	AxisHorizontal
	AxisVertical
)

func (a Axis) IsSpecified() bool {
	return a != AxisUnspecified
}

// TakeOrElse returns this Dp if Specified, otherwise executes the block.
func (a Axis) TakeOrElse(def Axis) Axis {
	if a.IsSpecified() {
		return a
	}
	return def
}
