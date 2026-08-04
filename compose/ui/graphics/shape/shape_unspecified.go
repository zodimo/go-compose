package shape

import (
	"image"
)

var _ Shape = (*ShapeUnspecifiedImpl)(nil)
var ShapeUnspecified Shape = &ShapeUnspecifiedImpl{}

type ShapeUnspecifiedImpl struct {
}

func (s *ShapeUnspecifiedImpl) CreateOutline(size image.Point, metric Metric) Outline {
	return nil
}

func (s *ShapeUnspecifiedImpl) mergeShape(other Shape) Shape {
	return other
}
func (s *ShapeUnspecifiedImpl) sameShape(other Shape) bool {
	return other == ShapeUnspecified
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
