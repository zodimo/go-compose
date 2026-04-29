package graphics

import (
	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// DefaultBlendMode is the default blend mode used for drawing operations.
var DefaultBlendMode = BlendModeSrcOver

// DrawScope provides a scoped drawing environment with a declarative, stateless API.
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-graphics/src/commonMain/kotlin/androidx/compose/ui/graphics/drawscope/DrawScope.kt
type DrawScope interface {
	unit.DensityScope

	// DrawContext returns the underlying draw context.
	DrawContext() DrawContext

	// LayoutDirection returns the layout direction of the layout being drawn.
	LayoutDirection() unit.LayoutDirection

	// Size returns the current size of the drawing environment.
	Size() geometry.Size

	// Center returns the center of the current drawing environment.
	Center() geometry.Offset

	// DrawLine draws a line between the given points using a color.
	DrawLine(color Color, start, end geometry.Offset, opts ...DrawLineOption)

	// DrawLineWithBrush draws a line between the given points using a brush.
	DrawLineWithBrush(brush Brush, start, end geometry.Offset, opts ...DrawLineOption)

	// DrawRect draws a rectangle with the given color.
	DrawRect(color Color, opts ...DrawRectOption)

	// DrawRectWithBrush draws a rectangle with the given brush.
	DrawRectWithBrush(brush Brush, opts ...DrawRectOption)

	// DrawRoundRect draws a rounded rectangle with the given color.
	DrawRoundRect(color Color, opts ...DrawRoundRectOption)

	// DrawRoundRectWithBrush draws a rounded rectangle with the given brush.
	DrawRoundRectWithBrush(brush Brush, opts ...DrawRoundRectOption)

	// DrawCircle draws a circle with the given color.
	DrawCircle(color Color, opts ...DrawCircleOption)

	// DrawCircleWithBrush draws a circle with the given brush.
	DrawCircleWithBrush(brush Brush, opts ...DrawCircleOption)

	// DrawOval draws an oval with the given color.
	DrawOval(color Color, opts ...DrawOvalOption)

	// DrawOvalWithBrush draws an oval with the given brush.
	DrawOvalWithBrush(brush Brush, opts ...DrawOvalOption)

	// DrawArc draws an arc with the given color.
	DrawArc(color Color, startAngle, sweepAngle float32, useCenter bool, opts ...DrawArcOption)

	// DrawArcWithBrush draws an arc with the given brush.
	DrawArcWithBrush(brush Brush, startAngle, sweepAngle float32, useCenter bool, opts ...DrawArcOption)

	// DrawPath draws a path with the given color.
	DrawPath(path Path, color Color, opts ...DrawPathOption)

	// DrawPathWithBrush draws a path with the given brush.
	DrawPathWithBrush(path Path, brush Brush, opts ...DrawPathOption)

	// DrawPoints draws a sequence of points with the given color.
	DrawPoints(points []geometry.Offset, pointMode PointMode, color Color, opts ...DrawPointsOption)

	// DrawPointsWithBrush draws a sequence of points with the given brush.
	DrawPointsWithBrush(points []geometry.Offset, pointMode PointMode, brush Brush, opts ...DrawPointsOption)

	// DrawImage draws an image at the given offset.
	DrawImage(image ImageBitmap, opts ...DrawImageOption)

	// DrawIntoCanvas provides direct access to the underlying canvas.
	DrawIntoCanvas(block func(Canvas))
}

// ContentDrawScope is a DrawScope that can draw content between other drawing operations.
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-graphics/src/commonMain/kotlin/androidx/compose/ui/graphics/drawscope/ContentDrawScope.kt
type ContentDrawScope interface {
	DrawScope

	// DrawContent causes child drawing operations to run.
	// If not called, the contents of the layout will not be drawn.
	DrawContent()
}
