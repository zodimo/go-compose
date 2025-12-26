package graphics

// DrawStyle defines a way to draw something.
type DrawStyle interface {
}

func EqualDrawStyle(a, b DrawStyle) bool {
	panic("EqualDrawStyle not implemented")
}

func TakeOrElseDrawStyle(a, b DrawStyle) DrawStyle {
	panic("TakeOrElseDrawStyle not implemented")
}
