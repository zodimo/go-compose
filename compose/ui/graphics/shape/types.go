package shape

import (
	"image"

	"gioui.org/op"
	"gioui.org/op/clip"
	gioUnit "gioui.org/unit"
)

// https://developer.android.com/reference/kotlin/androidx/compose/ui/graphics/Shape

type Shape interface {
	CreateOutline(size image.Point, metric gioUnit.Metric) Outline
}

type Outline interface {
	Push(ops *op.Ops) clip.Stack
	Op(ops *op.Ops) clip.Op
	Path(ops *op.Ops) clip.PathSpec
}
