package graphics

import (
	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

var canvas Canvas
var graphicsLayer any

// DefaultDensity is a stub density value used as a placeholder within CanvasDrawScope.
// The actual density is provided as a parameter during draw calls.
var DefaultDensity = unit.NewDensity(1.0, 1.0)

// DrawContext provides the dependencies to support a DrawScope drawing environment.
// It provides drawing bounds (size), target canvas, and handles transformations.
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-graphics/src/commonMain/kotlin/androidx/compose/ui/graphics/drawscope/DrawContext.kt
type DrawContext interface {
	// Size returns the current size of the drawing environment.
	Size() geometry.Size

	// SetSize sets the size of the drawing environment.
	SetSize(size geometry.Size)

	// Canvas returns the target canvas to issue drawing commands.
	Canvas() Canvas

	// SetCanvas sets the target canvas.
	SetCanvas(canvas Canvas)

	// Transform returns the controller for issuing transformations.
	Transform() DrawTransform

	// LayoutDirection returns the layout direction of the layout being drawn.
	LayoutDirection() unit.LayoutDirection

	// SetLayoutDirection sets the layout direction.
	SetLayoutDirection(ld unit.LayoutDirection)

	// Density returns the density used for dp/sp conversions.
	Density() unit.DensityScope

	// SetDensity sets the density.
	SetDensity(d unit.DensityScope)

	// GraphicsLayer returns the current graphics layer, if any.
	// May return nil if not drawing into a graphics layer.
	GraphicsLayer() interface{}

	// SetGraphicsLayer sets the current graphics layer.
	SetGraphicsLayer(layer interface{})
}

// drawContextImpl is the default implementation of DrawContext.
type drawContextImpl struct {
	canvas          Canvas
	size            geometry.Size
	transform       DrawTransform
	layoutDirection unit.LayoutDirection
	density         unit.DensityScope
}

// NewDrawContext creates a new DrawContext with the given initial values.
func NewDrawContext(
	canvas Canvas,
	size geometry.Size,
	density unit.DensityScope,
	layoutDirection unit.LayoutDirection,
) DrawContext {
	return &drawContextImpl{
		canvas:          canvas,
		size:            size,
		layoutDirection: layoutDirection,
		density:         density,
	}
}

func (d *drawContextImpl) Size() geometry.Size {
	return d.size
}

func (d *drawContextImpl) SetSize(size geometry.Size) {
	d.size = size
}

func (d *drawContextImpl) Canvas() Canvas {
	return d.canvas
}

func (d *drawContextImpl) SetCanvas(canvas Canvas) {
	d.canvas = canvas
}

func (d *drawContextImpl) Transform() DrawTransform {
	return d.transform
}

func (d *drawContextImpl) setTransform(transform DrawTransform) {
	d.transform = transform
}

func (d *drawContextImpl) LayoutDirection() unit.LayoutDirection {
	return d.layoutDirection
}

func (d *drawContextImpl) SetLayoutDirection(ld unit.LayoutDirection) {
	d.layoutDirection = ld
}

func (d *drawContextImpl) Density() unit.DensityScope {
	return d.density
}

func (d *drawContextImpl) SetDensity(density unit.DensityScope) {
	d.density = density
}

func (d *drawContextImpl) GraphicsLayer() interface{} {
	panic("GraphicsLayer not implemented")
}

func (d *drawContextImpl) SetGraphicsLayer(layer interface{}) {
	panic("SetGraphicsLayer not implemented")
}
