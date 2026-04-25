package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zodimo/gio-skia/skia"
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/runtime"
	"github.com/zodimo/go-compose/store"
	"github.com/zodimo/go-compose/theme"
	"github.com/zodimo/go-skia-support/skia/base"
	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/impl"
	"github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/unit"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Canvas Demo"))
		w.Option(app.Size(unit.Dp(600), unit.Dp(800)))

		if err := Run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func Run(window *app.Window) error {
	enLocale := system.Locale{Language: "en", Direction: system.LTR}
	var ops op.Ops

	store := store.NewPersistentState()
	store.Subscribe(func() {
		window.Invalidate()
	})

	runtime := runtime.NewRuntime()

	themeManager := theme.GetThemeManager()

	for {
		switch frameEvent := window.Event().(type) {
		case app.DestroyEvent:
			return frameEvent.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, frameEvent)
			gtx.Locale = enLocale
			gtx = themeManager.Material3ThemeInit(gtx)

			composer := compose.NewComposer(api.ComposerWithStore(store))

			skiaCanvas := skia.NewCanvas(gtx.Ops)
			canvasProvidedValue := compose.LocalCanvas.Provides(CanvasWrapper(skiaCanvas))

			uiWithCanvas := compose.CompositionLocalProvider(
				[]api.ProvidedValue{canvasProvidedValue},
				UI(),
			)

			callOp := runtime.Run(gtx, composer, uiWithCanvas)
			callOp.Add(gtx.Ops)

			frameEvent.Frame(gtx.Ops)
		}
	}
}

var _ graphics.Canvas = (*canvasWrapper)(nil)

type canvasWrapper struct {
	canvas skia.Canvas
}

func CanvasWrapper(skiaCanvas skia.Canvas) graphics.Canvas {
	return &canvasWrapper{canvas: skiaCanvas}
}

// Save saves a copy of the current transform and clip on the save stack.
func (c *canvasWrapper) Save() {
	c.canvas.Save()
}

// Restore pops the current save stack, restoring the previous transform and clip.
func (c *canvasWrapper) Restore() {
	c.canvas.Restore()
}

// SaveLayer saves the current state and creates a new group for subsequent operations.
// When restored, the layer is composited into the previous layer using the paint's blend mode.
func (c *canvasWrapper) SaveLayer(bounds geometry.Rect, paint *graphics.Paint) {
	c.canvas.SaveLayer(
		&models.Rect{
			Left:   bounds.Left,
			Top:    bounds.Top,
			Right:  bounds.Right,
			Bottom: bounds.Bottom,
		},
		graphicsPaintToSkiaPaint(paint))
}

// Translate shifts the coordinate space by the given delta.
func (c *canvasWrapper) Translate(dx, dy float32) {
	c.canvas.Translate(dx, dy)
}

// Scale scales the coordinate space by the given factors.
func (c *canvasWrapper) Scale(sx, sy float32) {
	c.canvas.Scale(sx, sy)
}

// Rotate rotates the coordinate space by the given degrees clockwise.
func (c *canvasWrapper) Rotate(degrees float32) {
	c.canvas.Rotate(degrees)
}

// Skew applies an axis-aligned skew transformation.
func (c *canvasWrapper) Skew(sx, sy float32) {
	c.canvas.Skew(sx, sy)
}

// Concat multiplies the current transform by the specified matrix.
func (c *canvasWrapper) Concat(matrix graphics.Matrix) {
	c.canvas.Concat(graphicsMatrixToSkiaMatrix(matrix))
}

// ClipRect reduces the clip region to the intersection of the current clip and the given rectangle.
func (c *canvasWrapper) ClipRect(left, top, right, bottom float32, clipOp graphics.ClipOp) {

	c.canvas.ClipRect(models.Rect{
		Left:   skia.Scalar(left),
		Top:    skia.Scalar(top),
		Right:  skia.Scalar(right),
		Bottom: skia.Scalar(bottom),
	},
		clipOpToSkiaClipOp(clipOp),
		false,
	)
}

// ClipPath reduces the clip region to the intersection of the current clip and the given path.
func (c *canvasWrapper) ClipPath(path graphics.Path, clipOp graphics.ClipOp) {
	c.canvas.ClipPath(pathToSkPath(path), clipOpToSkiaClipOp(clipOp), false)
}

// DrawLine draws a line between the given points.
func (c *canvasWrapper) DrawLine(p1, p2 geometry.Offset, paint *graphics.Paint) {
	c.canvas.DrawLine(
		models.Point{X: skia.Scalar(p1.X()), Y: skia.Scalar(p1.Y())},
		models.Point{X: skia.Scalar(p2.X()), Y: skia.Scalar(p2.Y())},
		graphicsPaintToSkiaPaint(paint),
	)
}

// DrawRect draws a rectangle.
func (c *canvasWrapper) DrawRect(left, top, right, bottom float32, paint *graphics.Paint) {
	c.canvas.DrawRect(
		models.Rect{
			Left:   skia.Scalar(left),
			Top:    skia.Scalar(top),
			Right:  skia.Scalar(right),
			Bottom: skia.Scalar(bottom),
		},
		graphicsPaintToSkiaPaint(paint),
	)
}

// DrawRoundRect draws a rounded rectangle.
func (c *canvasWrapper) DrawRoundRect(left, top, right, bottom, radiusX, radiusY float32, paint *graphics.Paint) {
	panic("DrawRoundRect not implemented")

	// c.canvas.DrawRRect(
	// 	models.RRect{}

	// 	models.Rect{
	// 		Left:   skia.Scalar(left),
	// 		Top:    skia.Scalar(top),
	// 		Right:  skia.Scalar(right),
	// 		Bottom: skia.Scalar(bottom),
	// 	},
	// 	skia.Scalar(radiusX),
	// 	skia.Scalar(radiusY),
	// 	graphicsPaintToSkiaPaint(paint),
	// )
}

// DrawOval draws an axis-aligned oval that fills the given rectangle.
func (c *canvasWrapper) DrawOval(left, top, right, bottom float32, paint *graphics.Paint) {
	c.canvas.DrawOval(
		models.Rect{
			Left:   skia.Scalar(left),
			Top:    skia.Scalar(top),
			Right:  skia.Scalar(right),
			Bottom: skia.Scalar(bottom),
		},
		graphicsPaintToSkiaPaint(paint),
	)
}

// DrawCircle draws a circle at the given center with the given radius.
func (c *canvasWrapper) DrawCircle(center geometry.Offset, radius float32, paint *graphics.Paint) {
	c.canvas.DrawCircle(
		models.Point{X: skia.Scalar(center.X()), Y: skia.Scalar(center.Y())},
		skia.Scalar(radius),
		graphicsPaintToSkiaPaint(paint),
	)
}

// DrawArc draws an arc scaled to fit inside the given rectangle.
func (c *canvasWrapper) DrawArc(left, top, right, bottom, startAngle, sweepAngle float32, useCenter bool, paint *graphics.Paint) {
	c.canvas.DrawArc(
		models.Rect{
			Left:   skia.Scalar(left),
			Top:    skia.Scalar(top),
			Right:  skia.Scalar(right),
			Bottom: skia.Scalar(bottom),
		},
		skia.Scalar(startAngle),
		skia.Scalar(sweepAngle),
		useCenter,
		graphicsPaintToSkiaPaint(paint),
	)
}

// DrawPath draws the given path.
func (c *canvasWrapper) DrawPath(path graphics.Path, paint *graphics.Paint) {
	c.canvas.DrawPath(
		pathToSkPath(path),
		graphicsPaintToSkiaPaint(paint),
	)
}

// DrawImage draws an image at the given offset.
func (c *canvasWrapper) DrawImage(
	image graphics.ImageBitmap,
	topLeftOffset geometry.Offset,
	paint *graphics.Paint,
) {
	c.canvas.DrawImage(
		imageToSkiaImageBitmap(image),
		skia.Scalar(topLeftOffset.X()),
		skia.Scalar(topLeftOffset.Y()),
		graphicsPaintToSkiaPaint(paint),
	)
}

// DrawImageRect draws a portion of an image into a destination rectangle.
func (c *canvasWrapper) DrawImageRect(
	image graphics.ImageBitmap,
	srcOffset graphics.IntOffset,
	srcSize graphics.IntSize,
	dstOffset graphics.IntOffset,
	dstSize graphics.IntSize,
	paint *graphics.Paint,
) {
	c.canvas.DrawImageRect(
		imageToSkiaImageBitmap(image),
		&models.Rect{
			Left:   skia.Scalar(srcOffset.X),
			Top:    skia.Scalar(srcOffset.Y),
			Right:  skia.Scalar(srcOffset.X + srcSize.Width),
			Bottom: skia.Scalar(srcOffset.Y + srcSize.Height),
		},
		models.Rect{
			Left:   skia.Scalar(dstOffset.X),
			Top:    skia.Scalar(dstOffset.Y),
			Right:  skia.Scalar(dstOffset.X + dstSize.Width),
			Bottom: skia.Scalar(dstOffset.Y + dstSize.Height),
		},
		graphicsPaintToSkiaPaint(paint),
	)
}

// DrawPoints draws a sequence of points according to the given PointMode.
func (c *canvasWrapper) DrawPoints(
	pointMode graphics.PointMode,
	points []geometry.Offset,
	paint *graphics.Paint,
) {
	skiaPoints := make([]models.Point, len(points))
	for i, p := range points {
		skiaPoints[i] = models.Point{X: skia.Scalar(p.X()), Y: skia.Scalar(p.Y())}
	}
	c.canvas.DrawPoints(pointModeToSkiaPointMode(pointMode), skiaPoints, graphicsPaintToSkiaPaint(paint))
}

// EnableZ enables Z-ordering for 3D-like effects (Android-specific, may be no-op on other platforms).
func (c *canvasWrapper) EnableZ() {}

// DisableZ disables Z-ordering.
func (c *canvasWrapper) DisableZ() {}

func graphicsPaintToSkiaPaint(paint *graphics.Paint) interfaces.SkPaint {
	fmt.Printf("paint: %v\n", paint)

	skPaint := impl.NewPaint()
	if paint.Color.Alpha() > 0 {
		skPaint.SetColor(colorToSkColor(paint.Color))
	}

	return skPaint
}

func graphicsMatrixToSkiaMatrix(matrix graphics.Matrix) interfaces.SkMatrix {
	panic("graphicsMatrixToSkiaMatrix Not implemented")
}

func clipOpToSkiaClipOp(clipOp graphics.ClipOp) enums.ClipOp {
	panic("clipOpToSkiaClipOp Not implemented")
}

func pathToSkPath(path graphics.Path) interfaces.SkPath {
	panic("pathToSkPath Not implemented")
}

func imageToSkiaImageBitmap(image graphics.ImageBitmap) interfaces.SkImage {
	panic("imageToSkiaImageBitmap Not implemented")
}

func pointModeToSkiaPointMode(pointMode graphics.PointMode) enums.PointMode {
	switch pointMode {
	case graphics.PointModePoints:
		return enums.PointModePoints
	case graphics.PointModeLines:
		return enums.PointModeLines
	case graphics.PointModePolygon:
		return enums.PointModePolygon
	}
	return enums.PointModePoints
}

func colorToSkColor(color graphics.Color) models.Color4f {
	nrgba := color.ToNRGBA()
	return models.Color4f{
		R: base.Scalar(float32(nrgba.R) / 255.0),
		G: base.Scalar(float32(nrgba.G) / 255.0),
		B: base.Scalar(float32(nrgba.B) / 255.0),
		A: base.Scalar(float32(nrgba.A) / 255.0),
	}
}
