package shape

import (
	"image"

	gioUnit "gioui.org/unit"
)

var _ Shape = (*ShapeUnspecifiedImpl)(nil)

// shapeUnspecified is the sentinel for Shape. It is deliberately package-private:
// the only public way to obtain it is ShapeCompanion.Unspecified().
var shapeUnspecified Shape = &ShapeUnspecifiedImpl{}

type ShapeUnspecifiedImpl struct {
}

func (s *ShapeUnspecifiedImpl) CreateOutline(size image.Point, metric gioUnit.Metric) Outline {
	return nil
}

func (s *ShapeUnspecifiedImpl) mergeShape(other Shape) Shape {
	return other
}
func (s *ShapeUnspecifiedImpl) sameShape(other Shape) bool {
	return other == shapeUnspecified
}
func (s *ShapeUnspecifiedImpl) semanticEqualShape(other Shape) bool {
	if _, ok := other.(*ShapeUnspecifiedImpl); ok {
		return true
	}
	return false
}
func (s *ShapeUnspecifiedImpl) copyShape(options ...ShapeOption) Shape {
	panic("ShapeUnspecifiedImpl.copyShape: cannot copy unspecified shape")
}
func (s *ShapeUnspecifiedImpl) stringShape() string {
	return "Shape{Unspecified}"
}
