package graphics

// DrawStyle defines a way to draw something.
type DrawStyle interface {
}

func EqualDrawStyle(a, b DrawStyle) bool {
	return a == b
}

func TakeOrElseDrawStyle(a, b DrawStyle) DrawStyle {
	if a != nil {
		return a
	}
	return b
}
