package skiagraphics

import (
	"fmt"
	"image"

	"gioui.org/layout"
	"github.com/zodimo/gio-skia/skia"
	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-skia-support/skia/enums"
	"github.com/zodimo/go-skia-support/skia/impl"
	skiaiface "github.com/zodimo/go-skia-support/skia/interfaces"
	"github.com/zodimo/go-skia-support/skia/models"
)

var _ graphics.Canvas = (*canvasWrapper)(nil)

type canvasWrapper struct {
	canvas skia.Canvas
	// alphaMultiplier is applied to layer paints when saving layers
	alphaMultiplier float32
}

func NewCanvas(gtx layout.Context) graphics.Canvas {
	return &canvasWrapper{canvas: skia.NewCanvas(gtx.Ops), alphaMultiplier: 1.0}
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
	// Apply alpha multiplier to the provided paint if needed
	var effPaint *graphics.Paint
	if paint != nil && c.alphaMultiplier != 1.0 {
		tmp := *paint
		tmp.Alpha = paint.Alpha * c.alphaMultiplier
		effPaint = &tmp
	} else {
		effPaint = paint
	}
	c.canvas.SaveLayer(
		geometryRectToModelsRect(bounds),
		graphicsPaintToSkiaPaint(effPaint))
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
	c.canvas.ClipRect(
		models.Rect{
			Left:   left,
			Top:    top,
			Right:  right,
			Bottom: bottom,
		},
		clipOpToSkiaClipOp(clipOp),
		false,
	)
}

// ClipPath reduces the clip region to the intersection of the current clip and the given path.
func (c *canvasWrapper) ClipPath(path graphics.Path, clipOp graphics.ClipOp) {
	skPath, err := pathToSkPath(path)
	if err != nil {
		return
	}
	c.canvas.ClipPath(skPath, clipOpToSkiaClipOp(clipOp), false)
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
			Left:   left,
			Top:    top,
			Right:  right,
			Bottom: bottom,
		},
		graphicsPaintToSkiaPaint(paint),
	)
}

// DrawRoundRect draws a rounded rectangle.
func (c *canvasWrapper) DrawRoundRect(left, top, right, bottom, radiusX, radiusY float32, paint *graphics.Paint) {
	rrect := models.NewRRect(
		models.Rect{
			Left:   skia.Scalar(left),
			Top:    skia.Scalar(top),
			Right:  skia.Scalar(right),
			Bottom: skia.Scalar(bottom),
		},
		[4]models.Point{
			{X: skia.Scalar(radiusX), Y: skia.Scalar(radiusY)},
			{X: skia.Scalar(radiusX), Y: skia.Scalar(radiusY)},
			{X: skia.Scalar(radiusX), Y: skia.Scalar(radiusY)},
			{X: skia.Scalar(radiusX), Y: skia.Scalar(radiusY)},
		},
	)
	c.canvas.DrawRRect(*rrect, graphicsPaintToSkiaPaint(paint))
}

// DrawOval draws an axis-aligned oval that fills the given rectangle.
func (c *canvasWrapper) DrawOval(left, top, right, bottom float32, paint *graphics.Paint) {
	c.canvas.DrawOval(
		models.Rect{
			Left:   left,
			Top:    top,
			Right:  right,
			Bottom: bottom,
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
			Left:   left,
			Top:    top,
			Right:  right,
			Bottom: bottom,
		},
		startAngle,
		sweepAngle,
		useCenter,
		graphicsPaintToSkiaPaint(paint),
	)
}

// DrawPath draws the given path.
func (c *canvasWrapper) DrawPath(path graphics.Path, paint *graphics.Paint) {
	skPath, err := pathToSkPath(path)
	if err != nil {
		return
	}
	c.canvas.DrawPath(
		skPath,
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
			Left:   float32(srcOffset.X),
			Top:    float32(srcOffset.Y),
			Right:  float32(srcOffset.X + srcSize.Width),
			Bottom: float32(srcOffset.Y + srcSize.Height),
		},
		models.Rect{
			Left:   float32(dstOffset.X),
			Top:    float32(dstOffset.Y),
			Right:  float32(dstOffset.X + dstSize.Width),
			Bottom: float32(dstOffset.Y + dstSize.Height),
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

func graphicsPaintToSkiaPaint(paint *graphics.Paint) skia.SkPaint {
	skPaint := impl.NewPaint()

	if paint.Color.Alpha() > 0 {
		skPaint.SetColor(colorToSkColor(paint.Color))
	}

	if paint.Alpha < 1.0 {
		skPaint.SetAlpha(uint8(paint.Alpha * 255))
	}

	if paint.Style != nil {
		if paint.Style == graphics.Fill {
			skPaint.SetStyle(enums.PaintStyleFill)
		} else if stroke, isStroke := paint.Style.(*graphics.Stroke); isStroke {
			skPaint.SetStyle(enums.PaintStyleStroke)
			skPaint.SetStrokeWidth(skia.Scalar(stroke.Width))
			skPaint.SetStrokeCap(graphicsStrokeCapToSkiaPaintCap(stroke.Cap))
			skPaint.SetStrokeJoin(graphicsStrokeJoinToSkiaPaintJoin(stroke.Join))
			skPaint.SetStrokeMiter(stroke.Miter)
		}
	}

	// Apply shader if any (conversion currently a stub)
	if paint != nil && paint.Shader != nil {
		if shader := graphicsShaderToSkiaShader(paint.Shader); shader != nil {
			skPaint.SetShader(shader)
		}
	}

	// Apply color filter and mask filter if available
	// Note: Declaring the fields on Paint allows future mapping into Skia filters.
	// The actual Skia API calls are intentionally omitted here to keep
	// compatibility with the current binding surface.

	skPaint.SetBlendMode(graphicsBlendModeToSkiaBlendMode(paint.BlendMode))

	return skPaint
}

// graphicsShaderToSkiaShader converts a graphics.Shader into a Skia shader when possible.
// This is a minimal stub to allow the rest of the pipeline to operate pending
// full gradient shader support in the Skia binding.
func graphicsShaderToSkiaShader(_s graphics.Shader) skiaiface.Shader {
	// Gradient and composite shader conversions are intentionally left as a stub
	// for now. When the underlying Skia bindings are extended, implement the
	// conversion logic here depending on the concrete Shader type.
	return nil
}

func graphicsBlendModeToSkiaBlendMode(mode graphics.BlendMode) enums.BlendMode {
	switch mode {
	case graphics.BlendModeClear:
		return enums.BlendModeClear
	case graphics.BlendModeSrc:
		return enums.BlendModeSrc
	case graphics.BlendModeDst:
		return enums.BlendModeDst
	case graphics.BlendModeSrcOver:
		return enums.BlendModeSrcOver
	case graphics.BlendModeDstOver:
		return enums.BlendModeDstOver
	case graphics.BlendModeSrcIn:
		return enums.BlendModeSrcIn
	case graphics.BlendModeDstIn:
		return enums.BlendModeDstIn
	case graphics.BlendModeSrcOut:
		return enums.BlendModeSrcOut
	case graphics.BlendModeDstOut:
		return enums.BlendModeDstOut
	case graphics.BlendModeSrcAtop:
		return enums.BlendModeSrcATop
	case graphics.BlendModeDstAtop:
		return enums.BlendModeDstATop
	case graphics.BlendModeXor:
		return enums.BlendModeXor
	case graphics.BlendModePlus:
		return enums.BlendModePlus
	case graphics.BlendModeModulate:
		return enums.BlendModeModulate
	case graphics.BlendModeScreen:
		return enums.BlendModeScreen
	case graphics.BlendModeOverlay:
		return enums.BlendModeOverlay
	case graphics.BlendModeDarken:
		return enums.BlendModeDarken
	case graphics.BlendModeLighten:
		return enums.BlendModeLighten
	case graphics.BlendModeColorDodge:
		return enums.BlendModeColorDodge
	case graphics.BlendModeColorBurn:
		return enums.BlendModeColorBurn
	case graphics.BlendModeHardlight:
		return enums.BlendModeHardLight
	case graphics.BlendModeSoftlight:
		return enums.BlendModeSoftLight
	case graphics.BlendModeDifference:
		return enums.BlendModeDifference
	case graphics.BlendModeExclusion:
		return enums.BlendModeExclusion
	case graphics.BlendModeMultiply:
		return enums.BlendModeMultiply
	case graphics.BlendModeHue:
		return enums.BlendModeHue
	case graphics.BlendModeSaturation:
		return enums.BlendModeSaturation
	case graphics.BlendModeColor:
		return enums.BlendModeColor
	case graphics.BlendModeLuminosity:
		return enums.BlendModeLuminosity
	default:
		return enums.BlendModeSrcOver
	}
}

func graphicsMatrixToSkiaMatrix(matrix graphics.Matrix) skia.SkMatrix {
	return impl.NewMatrixAll(
		matrix[0], matrix[4], matrix[12],
		matrix[1], matrix[5], matrix[13],
		matrix[3], matrix[7], matrix[15],
	)
}

func clipOpToSkiaClipOp(clipOp graphics.ClipOp) skia.ClipOp {
	switch clipOp {
	case graphics.ClipOpIntersect:
		return enums.ClipOpIntersect
	case graphics.ClipOpDifference:
		return enums.ClipOpDifference
	default:
		panic(fmt.Sprintf("Unknown ClipOp: %s", clipOp))
	}
}

func geometryRectToModelsRect(r geometry.Rect) *models.Rect {
	return &models.Rect{
		Left:   r.Left,
		Top:    r.Top,
		Right:  r.Right,
		Bottom: r.Bottom,
	}
}

func pathToSkPath(path graphics.Path) (skia.SkPath, error) {
	if pw, ok := path.(*pathWrapper); ok {
		return pw.path, nil
	}
	return skia.NewPath(), fmt.Errorf("path is not a Skia path: %T", path)
}

func imageToSkiaImageBitmap(img graphics.ImageBitmap) skia.SkImage {
	if stdImg, ok := img.(image.Image); ok {
		return imageToSkiaImageFromImage(stdImg)
	}
	if converter, ok := img.(interface{ ToSkiaImage() skia.SkImage }); ok {
		return converter.ToSkiaImage()
	}
	panic("imageToSkiaImageBitmap: unsupported ImageBitmap type")
}

func imageToSkiaImageFromImage(img image.Image) skia.SkImage {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width <= 0 || height <= 0 {
		return nil
	}

	var rgbaImage *image.RGBA
	switch m := img.(type) {
	case *image.RGBA:
		rgbaImage = m
	default:
		rgbaImage = image.NewRGBA(bounds)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				rgbaImage.Set(x, y, m.At(x, y))
			}
		}
	}

	info := models.NewImageInfo(width, height, enums.ColorTypeRGBA8888, enums.AlphaTypePremul)
	return impl.MakeRasterData(info, rgbaImage.Pix, rgbaImage.Stride)
}

func pointModeToSkiaPointMode(pointMode graphics.PointMode) skia.PointMode {
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
	return skia.Color4fFromNRGBA(nrgba)
}

func graphicsStrokeCapToSkiaPaintCap(cap graphics.StrokeCap) enums.PaintCap {
	switch cap {
	case graphics.StrokeCapButt:
		return enums.PaintCapButt
	case graphics.StrokeCapRound:
		return enums.PaintCapRound
	case graphics.StrokeCapSquare:
		return enums.PaintCapSquare
	default:
		return enums.PaintCapButt
	}
}

func graphicsStrokeJoinToSkiaPaintJoin(join graphics.StrokeJoin) enums.PaintJoin {
	switch join {
	case graphics.StrokeJoinMiter:
		return enums.PaintJoinMiter
	case graphics.StrokeJoinRound:
		return enums.PaintJoinRound
	case graphics.StrokeJoinBevel:
		return enums.PaintJoinBevel
	default:
		return enums.PaintJoinMiter
	}
}
