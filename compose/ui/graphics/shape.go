package graphics

import "fmt"

//https://developer.android.com/reference/kotlin/androidx/compose/ui/graphics/Shape

// zero is the default shape
type Shape int

const (
	ShapeRectangle Shape = iota // default
	// ShapeCircle
)

func (s Shape) String() string {
	switch s {
	case ShapeRectangle:
		return "ShapeRectangle"
	// case ShapeCircle:
	// 	return "ShapeCircle"
	default:
		return fmt.Sprintf("UnknownShape(%d)", s)
	}
}
