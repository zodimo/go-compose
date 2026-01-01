package graphics

import (
	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// CanvasDrawScope is the implementation of DrawScope that issues drawing commands
// into the specified canvas and bounds.
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-graphics/src/commonMain/kotlin/androidx/compose/ui/graphics/drawscope/CanvasDrawScope.kt
type CanvasDrawScope struct {
	drawContext *drawContextImpl
	fillPaint   *Paint
	strokePaint *Paint
}

// NewCanvasDrawScope creates a new CanvasDrawScope.
func NewCanvasDrawScope() *CanvasDrawScope {
	return &CanvasDrawScope{
		drawContext: NewDrawContext().(*drawContextImpl),
		fillPaint:   nil,
		strokePaint: nil,
	}
}

// Draw executes the drawing commands within the provided canvas and bounds.
func (c *CanvasDrawScope) Draw(
	density unit.Density,
	layoutDirection unit.LayoutDirection,
	canvas Canvas,
	size geometry.Size,
	block func(DrawScope),
) {
	c.drawContext.SetDensity(density)
	c.drawContext.SetLayoutDirection(layoutDirection)
	c.drawContext.SetCanvas(canvas)
	c.drawContext.SetSize(size)
	c.drawContext.setTransform(c.asDrawTransform())

	block(c)
}

// DrawWithGraphicsLayer executes drawing with an optional graphics layer.
func (c *CanvasDrawScope) DrawWithGraphicsLayer(
	density unit.Density,
	layoutDirection unit.LayoutDirection,
	canvas Canvas,
	size geometry.Size,
	graphicsLayer interface{},
	block func(DrawScope),
) {
	c.drawContext.SetDensity(density)
	c.drawContext.SetLayoutDirection(layoutDirection)
	c.drawContext.SetCanvas(canvas)
	c.drawContext.SetSize(size)
	c.drawContext.SetGraphicsLayer(graphicsLayer)
	c.drawContext.setTransform(c.asDrawTransform())

	block(c)

	c.drawContext.SetGraphicsLayer(nil)
}

// Density interface implementation
func (c *CanvasDrawScope) Density() float32 {
	return c.drawContext.Density().Density()
}

func (c *CanvasDrawScope) FontScale() float32 {
	return c.drawContext.Density().FontScale()
}

func (c *CanvasDrawScope) DpToPx(dp unit.Dp) float32 {
	return c.drawContext.Density().DpToPx(dp)
}

func (c *CanvasDrawScope) DpRoundToPx(dp unit.Dp) int {
	return c.drawContext.Density().DpRoundToPx(dp)
}

func (c *CanvasDrawScope) TextUnitToPx(tu unit.TextUnit) float32 {
	return c.drawContext.Density().TextUnitToPx(tu)
}

func (c *CanvasDrawScope) TextUnitRoundToPx(tu unit.TextUnit) int {
	return c.drawContext.Density().TextUnitRoundToPx(tu)
}

func (c *CanvasDrawScope) IntToDp(px int) unit.Dp {
	return c.drawContext.Density().IntToDp(px)
}

func (c *CanvasDrawScope) IntToSp(px int) unit.TextUnit {
	return c.drawContext.Density().IntToSp(px)
}

func (c *CanvasDrawScope) FloatToDp(px float32) unit.Dp {
	return c.drawContext.Density().FloatToDp(px)
}

func (c *CanvasDrawScope) FloatToSp(px float32) unit.TextUnit {
	return c.drawContext.Density().FloatToSp(px)
}

func (c *CanvasDrawScope) DpRectToRect(rect unit.DpRect) geometry.Rect {
	return c.drawContext.Density().DpRectToRect(rect)
}

func (c *CanvasDrawScope) DpSizeToSize(size unit.DpSize) geometry.Size {
	return c.drawContext.Density().DpSizeToSize(size)
}

func (c *CanvasDrawScope) SizeToDpSize(size geometry.Size) unit.DpSize {
	return c.drawContext.Density().SizeToDpSize(size)
}

// DrawScope interface implementation
func (c *CanvasDrawScope) DrawContext() DrawContext {
	return c.drawContext
}

func (c *CanvasDrawScope) LayoutDirection() unit.LayoutDirection {
	return c.drawContext.LayoutDirection()
}

func (c *CanvasDrawScope) Size() geometry.Size {
	return c.drawContext.Size()
}

func (c *CanvasDrawScope) Center() geometry.Offset {
	return c.drawContext.Size().Center()
}

// Paint management
func (c *CanvasDrawScope) obtainFillPaint() *Paint {
	if c.fillPaint == nil {
		c.fillPaint = NewPaint()
	}
	return c.fillPaint
}

func (c *CanvasDrawScope) obtainStrokePaint() *Paint {
	if c.strokePaint == nil {
		c.strokePaint = NewPaint()
	}
	return c.strokePaint
}

func (c *CanvasDrawScope) selectPaint(style DrawStyle) *Paint {
	if _, ok := style.(*Stroke); ok {
		return c.obtainStrokePaint()
	}
	return c.obtainFillPaint()
}

func (c *CanvasDrawScope) configurePaintWithColor(
	color Color,
	style DrawStyle,
	alpha float32,
	blendMode BlendMode,
) *Paint {
	paint := c.selectPaint(style)
	paint.Color = modulateColorAlpha(color, alpha)
	paint.BlendMode = blendMode
	paint.Shader = nil

	if stroke, ok := style.(*Stroke); ok {
		paint.StrokeWidth = stroke.Width
	}

	return paint
}

func (c *CanvasDrawScope) configurePaintWithBrush(
	brush Brush,
	style DrawStyle,
	alpha float32,
	blendMode BlendMode,
) *Paint {
	paint := c.selectPaint(style)
	brush.ApplyTo(c.Size(), paint, alpha)
	paint.BlendMode = blendMode

	if stroke, ok := style.(*Stroke); ok {
		paint.StrokeWidth = stroke.Width
	}

	return paint
}

// Drawing methods
func (c *CanvasDrawScope) DrawLine(color Color, start, end geometry.Offset, opts ...DrawLineOption) {
	cfg := defaultDrawLineConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	paint := c.obtainStrokePaint()
	paint.Color = modulateColorAlpha(color, cfg.alpha)
	paint.StrokeWidth = cfg.strokeWidth
	paint.BlendMode = cfg.blendMode

	c.drawContext.Canvas().DrawLine(start, end, paint)
}

func (c *CanvasDrawScope) DrawLineWithBrush(brush Brush, start, end geometry.Offset, opts ...DrawLineOption) {
	cfg := defaultDrawLineConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	paint := c.obtainStrokePaint()
	brush.ApplyTo(c.Size(), paint, cfg.alpha)
	paint.StrokeWidth = cfg.strokeWidth
	paint.BlendMode = cfg.blendMode

	c.drawContext.Canvas().DrawLine(start, end, paint)
}

func (c *CanvasDrawScope) DrawRect(color Color, opts ...DrawRectOption) {
	cfg := defaultDrawRectConfig(c.Size())
	for _, opt := range opts {
		opt(&cfg)
	}

	paint := c.configurePaintWithColor(color, cfg.style, cfg.alpha, cfg.blendMode)
	left := cfg.topLeft.X()
	top := cfg.topLeft.Y()
	right := left + cfg.size.Width()
	bottom := top + cfg.size.Height()
	c.drawContext.Canvas().DrawRect(left, top, right, bottom, paint)
}

func (c *CanvasDrawScope) DrawRectWithBrush(brush Brush, opts ...DrawRectOption) {
	cfg := defaultDrawRectConfig(c.Size())
	for _, opt := range opts {
		opt(&cfg)
	}

	paint := c.configurePaintWithBrush(brush, cfg.style, cfg.alpha, cfg.blendMode)
	left := cfg.topLeft.X()
	top := cfg.topLeft.Y()
	right := left + cfg.size.Width()
	bottom := top + cfg.size.Height()
	c.drawContext.Canvas().DrawRect(left, top, right, bottom, paint)
}

func (c *CanvasDrawScope) DrawRoundRect(color Color, opts ...DrawRoundRectOption) {
	cfg := defaultDrawRoundRectConfig(c.Size())
	for _, opt := range opts {
		opt(&cfg)
	}

	paint := c.configurePaintWithColor(color, cfg.style, cfg.alpha, cfg.blendMode)
	left := cfg.topLeft.X()
	top := cfg.topLeft.Y()
	right := left + cfg.size.Width()
	bottom := top + cfg.size.Height()
	c.drawContext.Canvas().DrawRoundRect(left, top, right, bottom, cfg.cornerRadius.X(), cfg.cornerRadius.Y(), paint)
}

func (c *CanvasDrawScope) DrawRoundRectWithBrush(brush Brush, opts ...DrawRoundRectOption) {
	cfg := defaultDrawRoundRectConfig(c.Size())
	for _, opt := range opts {
		opt(&cfg)
	}

	paint := c.configurePaintWithBrush(brush, cfg.style, cfg.alpha, cfg.blendMode)
	left := cfg.topLeft.X()
	top := cfg.topLeft.Y()
	right := left + cfg.size.Width()
	bottom := top + cfg.size.Height()
	c.drawContext.Canvas().DrawRoundRect(left, top, right, bottom, cfg.cornerRadius.X(), cfg.cornerRadius.Y(), paint)
}

func (c *CanvasDrawScope) DrawCircle(color Color, opts ...DrawCircleOption) {
	cfg := defaultDrawCircleConfig(c.Size())
	for _, opt := range opts {
		opt(&cfg)
	}

	paint := c.configurePaintWithColor(color, cfg.style, cfg.alpha, cfg.blendMode)
	c.drawContext.Canvas().DrawCircle(cfg.center, cfg.radius, paint)
}

func (c *CanvasDrawScope) DrawCircleWithBrush(brush Brush, opts ...DrawCircleOption) {
	cfg := defaultDrawCircleConfig(c.Size())
	for _, opt := range opts {
		opt(&cfg)
	}

	paint := c.configurePaintWithBrush(brush, cfg.style, cfg.alpha, cfg.blendMode)
	c.drawContext.Canvas().DrawCircle(cfg.center, cfg.radius, paint)
}

func (c *CanvasDrawScope) DrawOval(color Color, opts ...DrawOvalOption) {
	cfg := defaultDrawOvalConfig(c.Size())
	for _, opt := range opts {
		opt(&cfg)
	}

	paint := c.configurePaintWithColor(color, cfg.style, cfg.alpha, cfg.blendMode)
	left := cfg.topLeft.X()
	top := cfg.topLeft.Y()
	right := left + cfg.size.Width()
	bottom := top + cfg.size.Height()
	c.drawContext.Canvas().DrawOval(left, top, right, bottom, paint)
}

func (c *CanvasDrawScope) DrawOvalWithBrush(brush Brush, opts ...DrawOvalOption) {
	cfg := defaultDrawOvalConfig(c.Size())
	for _, opt := range opts {
		opt(&cfg)
	}

	paint := c.configurePaintWithBrush(brush, cfg.style, cfg.alpha, cfg.blendMode)
	left := cfg.topLeft.X()
	top := cfg.topLeft.Y()
	right := left + cfg.size.Width()
	bottom := top + cfg.size.Height()
	c.drawContext.Canvas().DrawOval(left, top, right, bottom, paint)
}

func (c *CanvasDrawScope) DrawArc(color Color, startAngle, sweepAngle float32, useCenter bool, opts ...DrawArcOption) {
	cfg := defaultDrawArcConfig(c.Size())
	for _, opt := range opts {
		opt(&cfg)
	}

	paint := c.configurePaintWithColor(color, cfg.style, cfg.alpha, cfg.blendMode)
	left := cfg.topLeft.X()
	top := cfg.topLeft.Y()
	right := left + cfg.size.Width()
	bottom := top + cfg.size.Height()
	c.drawContext.Canvas().DrawArc(left, top, right, bottom, startAngle, sweepAngle, useCenter, paint)
}

func (c *CanvasDrawScope) DrawArcWithBrush(brush Brush, startAngle, sweepAngle float32, useCenter bool, opts ...DrawArcOption) {
	cfg := defaultDrawArcConfig(c.Size())
	for _, opt := range opts {
		opt(&cfg)
	}

	paint := c.configurePaintWithBrush(brush, cfg.style, cfg.alpha, cfg.blendMode)
	left := cfg.topLeft.X()
	top := cfg.topLeft.Y()
	right := left + cfg.size.Width()
	bottom := top + cfg.size.Height()
	c.drawContext.Canvas().DrawArc(left, top, right, bottom, startAngle, sweepAngle, useCenter, paint)
}

func (c *CanvasDrawScope) DrawPath(path Path, color Color, opts ...DrawPathOption) {
	cfg := defaultDrawPathConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	paint := c.configurePaintWithColor(color, cfg.style, cfg.alpha, cfg.blendMode)
	c.drawContext.Canvas().DrawPath(path, paint)
}

func (c *CanvasDrawScope) DrawPathWithBrush(path Path, brush Brush, opts ...DrawPathOption) {
	cfg := defaultDrawPathConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	paint := c.configurePaintWithBrush(brush, cfg.style, cfg.alpha, cfg.blendMode)
	c.drawContext.Canvas().DrawPath(path, paint)
}

func (c *CanvasDrawScope) DrawPoints(points []geometry.Offset, pointMode PointMode, color Color, opts ...DrawPointsOption) {
	cfg := defaultDrawPointsConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	paint := c.obtainStrokePaint()
	paint.Color = modulateColorAlpha(color, cfg.alpha)
	paint.StrokeWidth = cfg.strokeWidth
	paint.BlendMode = cfg.blendMode

	c.drawContext.Canvas().DrawPoints(pointMode, points, paint)
}

func (c *CanvasDrawScope) DrawPointsWithBrush(points []geometry.Offset, pointMode PointMode, brush Brush, opts ...DrawPointsOption) {
	cfg := defaultDrawPointsConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	paint := c.obtainStrokePaint()
	brush.ApplyTo(c.Size(), paint, cfg.alpha)
	paint.StrokeWidth = cfg.strokeWidth
	paint.BlendMode = cfg.blendMode

	c.drawContext.Canvas().DrawPoints(pointMode, points, paint)
}

func (c *CanvasDrawScope) DrawImage(image ImageBitmap, opts ...DrawImageOption) {
	cfg := defaultDrawImageConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	paint := c.obtainFillPaint()
	paint.Alpha = cfg.alpha
	paint.BlendMode = cfg.blendMode

	c.drawContext.Canvas().DrawImage(image, cfg.topLeft, paint)
}

// modulateColorAlpha creates a new color with modulated alpha while keeping RGB.
func modulateColorAlpha(color Color, alpha float32) Color {
	return color.Copy(
		CopyWithAlpha(color.Alpha() * alpha),
	)
}

func (c *CanvasDrawScope) DrawIntoCanvas(block func(Canvas)) {
	block(c.drawContext.Canvas())
}

// asDrawTransform creates a DrawTransform backed by this CanvasDrawScope.
func (c *CanvasDrawScope) asDrawTransform() DrawTransform {
	return &canvasDrawTransform{scope: c}
}

// canvasDrawTransform implements DrawTransform for CanvasDrawScope.
type canvasDrawTransform struct {
	scope *CanvasDrawScope
}

func (t *canvasDrawTransform) Size() geometry.Size {
	return t.scope.drawContext.Size()
}

func (t *canvasDrawTransform) Center() geometry.Offset {
	return t.scope.drawContext.Size().Center()
}

func (t *canvasDrawTransform) Inset(left, top, right, bottom float32) {
	t.scope.drawContext.Canvas().Translate(left, top)
	newWidth := t.Size().Width() - left - right
	newHeight := t.Size().Height() - top - bottom
	t.scope.drawContext.SetSize(geometry.NewSize(newWidth, newHeight))
}

func (t *canvasDrawTransform) ClipRect(left, top, right, bottom float32, clipOp ClipOp) {
	t.scope.drawContext.Canvas().ClipRect(left, top, right, bottom, clipOp)
}

func (t *canvasDrawTransform) ClipPath(path Path, clipOp ClipOp) {
	t.scope.drawContext.Canvas().ClipPath(path, clipOp)
}

func (t *canvasDrawTransform) Translate(left, top float32) {
	t.scope.drawContext.Canvas().Translate(left, top)
}

func (t *canvasDrawTransform) Rotate(degrees float32, pivot geometry.Offset) {
	canvas := t.scope.drawContext.Canvas()
	canvas.Translate(pivot.X(), pivot.Y())
	canvas.Rotate(degrees)
	canvas.Translate(-pivot.X(), -pivot.Y())
}

func (t *canvasDrawTransform) Scale(scaleX, scaleY float32, pivot geometry.Offset) {
	canvas := t.scope.drawContext.Canvas()
	canvas.Translate(pivot.X(), pivot.Y())
	canvas.Scale(scaleX, scaleY)
	canvas.Translate(-pivot.X(), -pivot.Y())
}

func (t *canvasDrawTransform) Transform(matrix Matrix) {
	t.scope.drawContext.Canvas().Concat(matrix)
}

// Compile-time assertion that CanvasDrawScope implements DrawScope
var _ DrawScope = (*CanvasDrawScope)(nil)
