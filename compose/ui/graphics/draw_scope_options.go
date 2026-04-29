package graphics

import (
	"github.com/zodimo/go-compose/compose/ui/geometry"
)

// DrawLineOption configures optional parameters for DrawLine.
type DrawLineOption func(*drawLineConfig)

type drawLineConfig struct {
	strokeWidth float32
	cap         StrokeCap
	pathEffect  interface{}
	alpha       float32
	colorFilter interface{}
	blendMode   BlendMode
}

func defaultDrawLineConfig() drawLineConfig {
	return drawLineConfig{
		strokeWidth: DefaultStrokeWidth,
		cap:         DefaultStrokeCap,
		pathEffect:  nil,
		alpha:       1.0,
		colorFilter: nil,
		blendMode:   DefaultBlendMode,
	}
}

// WithLineStrokeWidth sets the stroke width for a line.
func WithLineStrokeWidth(width float32) DrawLineOption {
	return func(c *drawLineConfig) { c.strokeWidth = width }
}

// WithLineCap sets the stroke cap for a line.
func WithLineCap(cap StrokeCap) DrawLineOption {
	return func(c *drawLineConfig) { c.cap = cap }
}

// WithLineAlpha sets the alpha for a line.
func WithLineAlpha(alpha float32) DrawLineOption {
	return func(c *drawLineConfig) { c.alpha = alpha }
}

// WithLineBlendMode sets the blend mode for a line.
func WithLineBlendMode(mode BlendMode) DrawLineOption {
	return func(c *drawLineConfig) { c.blendMode = mode }
}

// DrawRectOption configures optional parameters for DrawRect.
type DrawRectOption func(*drawRectConfig)

type drawRectConfig struct {
	topLeft     geometry.Offset
	size        geometry.Size
	alpha       float32
	style       DrawStyle
	colorFilter interface{}
	blendMode   BlendMode
}

func defaultDrawRectConfig(scopeSize geometry.Size) drawRectConfig {
	return drawRectConfig{
		topLeft:     geometry.OffsetZero,
		size:        scopeSize,
		alpha:       1.0,
		style:       Fill,
		colorFilter: nil,
		blendMode:   DefaultBlendMode,
	}
}

// WithRectTopLeft sets the top-left offset for a rectangle.
func WithRectTopLeft(offset geometry.Offset) DrawRectOption {
	return func(c *drawRectConfig) { c.topLeft = offset }
}

// WithRectSize sets the size for a rectangle.
func WithRectSize(size geometry.Size) DrawRectOption {
	return func(c *drawRectConfig) { c.size = size }
}

// WithRectAlpha sets the alpha for a rectangle.
func WithRectAlpha(alpha float32) DrawRectOption {
	return func(c *drawRectConfig) { c.alpha = alpha }
}

// WithRectStyle sets the draw style (Fill or Stroke) for a rectangle.
func WithRectStyle(style DrawStyle) DrawRectOption {
	return func(c *drawRectConfig) { c.style = style }
}

// WithRectBlendMode sets the blend mode for a rectangle.
func WithRectBlendMode(mode BlendMode) DrawRectOption {
	return func(c *drawRectConfig) { c.blendMode = mode }
}

// DrawRoundRectOption configures optional parameters for DrawRoundRect.
type DrawRoundRectOption func(*drawRoundRectConfig)

type drawRoundRectConfig struct {
	topLeft      geometry.Offset
	size         geometry.Size
	cornerRadius geometry.CornerRadius
	alpha        float32
	style        DrawStyle
	colorFilter  interface{}
	blendMode    BlendMode
}

func defaultDrawRoundRectConfig(scopeSize geometry.Size) drawRoundRectConfig {
	return drawRoundRectConfig{
		topLeft:      geometry.OffsetZero,
		size:         scopeSize,
		cornerRadius: geometry.CornerRadiusZero,
		alpha:        1.0,
		style:        Fill,
		colorFilter:  nil,
		blendMode:    DefaultBlendMode,
	}
}

// WithRoundRectTopLeft sets the top-left offset.
func WithRoundRectTopLeft(offset geometry.Offset) DrawRoundRectOption {
	return func(c *drawRoundRectConfig) { c.topLeft = offset }
}

// WithRoundRectSize sets the size.
func WithRoundRectSize(size geometry.Size) DrawRoundRectOption {
	return func(c *drawRoundRectConfig) { c.size = size }
}

// WithRoundRectCornerRadius sets the corner radius.
func WithRoundRectCornerRadius(radius geometry.CornerRadius) DrawRoundRectOption {
	return func(c *drawRoundRectConfig) { c.cornerRadius = radius }
}

// WithRoundRectAlpha sets the alpha.
func WithRoundRectAlpha(alpha float32) DrawRoundRectOption {
	return func(c *drawRoundRectConfig) { c.alpha = alpha }
}

// WithRoundRectStyle sets the draw style.
func WithRoundRectStyle(style DrawStyle) DrawRoundRectOption {
	return func(c *drawRoundRectConfig) { c.style = style }
}

// WithRoundRectBlendMode sets the blend mode.
func WithRoundRectBlendMode(mode BlendMode) DrawRoundRectOption {
	return func(c *drawRoundRectConfig) { c.blendMode = mode }
}

// DrawCircleOption configures optional parameters for DrawCircle.
type DrawCircleOption func(*drawCircleConfig)

type drawCircleConfig struct {
	radius      float32
	center      geometry.Offset
	alpha       float32
	style       DrawStyle
	colorFilter interface{}
	blendMode   BlendMode
}

func defaultDrawCircleConfig(scopeSize geometry.Size) drawCircleConfig {
	center := scopeSize.Center()
	minDim := scopeSize.MinDimension()
	return drawCircleConfig{
		radius:      minDim / 2,
		center:      center,
		alpha:       1.0,
		style:       Fill,
		colorFilter: nil,
		blendMode:   DefaultBlendMode,
	}
}

// WithCircleRadius sets the radius.
func WithCircleRadius(radius float32) DrawCircleOption {
	return func(c *drawCircleConfig) { c.radius = radius }
}

// WithCircleCenter sets the center.
func WithCircleCenter(center geometry.Offset) DrawCircleOption {
	return func(c *drawCircleConfig) { c.center = center }
}

// WithCircleAlpha sets the alpha.
func WithCircleAlpha(alpha float32) DrawCircleOption {
	return func(c *drawCircleConfig) { c.alpha = alpha }
}

// WithCircleStyle sets the draw style.
func WithCircleStyle(style DrawStyle) DrawCircleOption {
	return func(c *drawCircleConfig) { c.style = style }
}

// WithCircleBlendMode sets the blend mode.
func WithCircleBlendMode(mode BlendMode) DrawCircleOption {
	return func(c *drawCircleConfig) { c.blendMode = mode }
}

// DrawOvalOption configures optional parameters for DrawOval.
type DrawOvalOption func(*drawOvalConfig)

type drawOvalConfig struct {
	topLeft     geometry.Offset
	size        geometry.Size
	alpha       float32
	style       DrawStyle
	colorFilter interface{}
	blendMode   BlendMode
}

func defaultDrawOvalConfig(scopeSize geometry.Size) drawOvalConfig {
	return drawOvalConfig{
		topLeft:     geometry.OffsetZero,
		size:        scopeSize,
		alpha:       1.0,
		style:       Fill,
		colorFilter: nil,
		blendMode:   DefaultBlendMode,
	}
}

// WithOvalTopLeft sets the top-left offset.
func WithOvalTopLeft(offset geometry.Offset) DrawOvalOption {
	return func(c *drawOvalConfig) { c.topLeft = offset }
}

// WithOvalSize sets the size.
func WithOvalSize(size geometry.Size) DrawOvalOption {
	return func(c *drawOvalConfig) { c.size = size }
}

// WithOvalAlpha sets the alpha.
func WithOvalAlpha(alpha float32) DrawOvalOption {
	return func(c *drawOvalConfig) { c.alpha = alpha }
}

// WithOvalStyle sets the draw style.
func WithOvalStyle(style DrawStyle) DrawOvalOption {
	return func(c *drawOvalConfig) { c.style = style }
}

// WithOvalBlendMode sets the blend mode.
func WithOvalBlendMode(mode BlendMode) DrawOvalOption {
	return func(c *drawOvalConfig) { c.blendMode = mode }
}

// DrawArcOption configures optional parameters for DrawArc.
type DrawArcOption func(*drawArcConfig)

type drawArcConfig struct {
	topLeft     geometry.Offset
	size        geometry.Size
	alpha       float32
	style       DrawStyle
	colorFilter interface{}
	blendMode   BlendMode
}

func defaultDrawArcConfig(scopeSize geometry.Size) drawArcConfig {
	return drawArcConfig{
		topLeft:     geometry.OffsetZero,
		size:        scopeSize,
		alpha:       1.0,
		style:       Fill,
		colorFilter: nil,
		blendMode:   DefaultBlendMode,
	}
}

// WithArcTopLeft sets the top-left offset.
func WithArcTopLeft(offset geometry.Offset) DrawArcOption {
	return func(c *drawArcConfig) { c.topLeft = offset }
}

// WithArcSize sets the size.
func WithArcSize(size geometry.Size) DrawArcOption {
	return func(c *drawArcConfig) { c.size = size }
}

// WithArcAlpha sets the alpha.
func WithArcAlpha(alpha float32) DrawArcOption {
	return func(c *drawArcConfig) { c.alpha = alpha }
}

// WithArcStyle sets the draw style.
func WithArcStyle(style DrawStyle) DrawArcOption {
	return func(c *drawArcConfig) { c.style = style }
}

// WithArcBlendMode sets the blend mode.
func WithArcBlendMode(mode BlendMode) DrawArcOption {
	return func(c *drawArcConfig) { c.blendMode = mode }
}

// DrawPathOption configures optional parameters for DrawPath.
type DrawPathOption func(*drawPathConfig)

type drawPathConfig struct {
	alpha       float32
	style       DrawStyle
	colorFilter interface{}
	blendMode   BlendMode
}

func defaultDrawPathConfig() drawPathConfig {
	return drawPathConfig{
		alpha:       1.0,
		style:       NewStroke(1),
		colorFilter: nil,
		blendMode:   DefaultBlendMode,
	}
}

// WithPathAlpha sets the alpha.
func WithPathAlpha(alpha float32) DrawPathOption {
	return func(c *drawPathConfig) { c.alpha = alpha }
}

// WithPathStyle sets the draw style.
func WithPathStyle(style DrawStyle) DrawPathOption {
	return func(c *drawPathConfig) { c.style = style }
}

// WithPathBlendMode sets the blend mode.
func WithPathBlendMode(mode BlendMode) DrawPathOption {
	return func(c *drawPathConfig) { c.blendMode = mode }
}

// DrawPointsOption configures optional parameters for DrawPoints.
type DrawPointsOption func(*drawPointsConfig)

type drawPointsConfig struct {
	strokeWidth float32
	cap         StrokeCap
	pathEffect  interface{}
	alpha       float32
	colorFilter interface{}
	blendMode   BlendMode
}

func defaultDrawPointsConfig() drawPointsConfig {
	return drawPointsConfig{
		strokeWidth: DefaultStrokeWidth,
		cap:         StrokeCapButt,
		pathEffect:  nil,
		alpha:       1.0,
		colorFilter: nil,
		blendMode:   DefaultBlendMode,
	}
}

// WithPointsStrokeWidth sets the stroke width.
func WithPointsStrokeWidth(width float32) DrawPointsOption {
	return func(c *drawPointsConfig) { c.strokeWidth = width }
}

// WithPointsCap sets the stroke cap.
func WithPointsCap(cap StrokeCap) DrawPointsOption {
	return func(c *drawPointsConfig) { c.cap = cap }
}

// WithPointsAlpha sets the alpha.
func WithPointsAlpha(alpha float32) DrawPointsOption {
	return func(c *drawPointsConfig) { c.alpha = alpha }
}

// WithPointsBlendMode sets the blend mode.
func WithPointsBlendMode(mode BlendMode) DrawPointsOption {
	return func(c *drawPointsConfig) { c.blendMode = mode }
}

// DrawImageOption configures optional parameters for DrawImage.
type DrawImageOption func(*drawImageConfig)

type drawImageConfig struct {
	topLeft     geometry.Offset
	alpha       float32
	style       DrawStyle
	colorFilter interface{}
	blendMode   BlendMode
}

func defaultDrawImageConfig() drawImageConfig {
	return drawImageConfig{
		topLeft:     geometry.OffsetZero,
		alpha:       1.0,
		style:       Fill,
		colorFilter: nil,
		blendMode:   DefaultBlendMode,
	}
}

// WithImageTopLeft sets the top-left offset.
func WithImageTopLeft(offset geometry.Offset) DrawImageOption {
	return func(c *drawImageConfig) { c.topLeft = offset }
}

// WithImageAlpha sets the alpha.
func WithImageAlpha(alpha float32) DrawImageOption {
	return func(c *drawImageConfig) { c.alpha = alpha }
}

// WithImageBlendMode sets the blend mode.
func WithImageBlendMode(mode BlendMode) DrawImageOption {
	return func(c *drawImageConfig) { c.blendMode = mode }
}
