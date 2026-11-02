package icon

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

type IconOptions struct {
	Modifier ui.Modifier
	Color    graphics.Color
	Size     unit.Dp       // Unified size for both IconBytes and SymbolName
	FontSize unit.TextUnit // Deprecated: Use Size instead

	// SymbolStyle selects the Material Symbols typeface for SymbolName icons.
	// It has no effect on IconBytes.
	SymbolStyle SymbolStyle

	// Variable-font axis values for SymbolName icons. Each carries a "has"
	// flag because an explicit 0 (fill 0, grade 0) is a meaningful choice that
	// must stay distinguishable from "not provided". Unprovided axes resolve to
	// the font's default instance, except OpticalSize, which adapts to Size.
	Fill         float32
	hasFill      bool
	Grade        float32
	hasGrade     bool
	OpticalSize  float32
	hasOpticalSz bool
	Weight       float32
	hasWeight    bool

	// Brush paints the icon with a gradient instead of a solid color. When set
	// it takes precedence over Color, which in turn takes precedence over the
	// ambient content color.
	Brush graphics.Brush

	// StrokeWidth outlines the glyph instead of filling it. It is in Dp and does
	// not affect the icon's reported size.
	StrokeWidth unit.Dp
	hasStroke   bool
}

type IconOption func(*IconOptions)

func DefaultIconOptions() IconOptions {
	return IconOptions{
		Modifier:    ui.EmptyModifier,
		Color:       graphics.ColorCompanion.Unspecified(),
		Size:        unit.DpCompanion.Unspecified(),
		FontSize:    unit.TextUnitCompanion.Unspecified(),
		StrokeWidth: unit.DpCompanion.Unspecified(),
	}
}

func WithModifier(m ui.Modifier) IconOption {
	return func(o *IconOptions) {
		o.Modifier = m
	}
}

func WithColor(col graphics.Color) IconOption {
	return func(o *IconOptions) {
		o.Color = col
	}
}

// WithSize sets the icon size in Dp.
// For IconBytes, this constrains the icon dimensions.
// For SymbolName, this sets both the geometry scale and, unless overridden with
// WithOpticalSize, the font's optical size axis.
func WithSize(size unit.Dp) IconOption {
	return func(o *IconOptions) {
		o.Size = size
	}
}

// WithSymbolStyle selects the Material Symbols typeface used for SymbolName
// icons: SymbolOutlined (the default), SymbolRounded or SymbolSharp.
//
// A style whose font asset is not embedded in the build falls back to outlined
// and is reported through SymbolErrorHandler, rather than rendering blank.
func WithSymbolStyle(style SymbolStyle) IconOption {
	return func(o *IconOptions) {
		o.SymbolStyle = style
	}
}

// Deprecated: Use WithSize instead.
// WithSymbolSize sets the font size for symbol icons.
func WithSymbolSize(size unit.TextUnit) IconOption {
	return func(o *IconOptions) {
		o.FontSize = size
	}
}

// WithFill sets the fill axis (0 = outlined, 1 = filled), clamped to [0, 1].
// An explicit 0 is honoured and is distinct from leaving the axis unset.
func WithFill(fill float32) IconOption {
	return func(o *IconOptions) {
		o.Fill = fill
		o.hasFill = true
	}
}

// WithGrade sets the grade axis, which adjusts weight for better legibility
// against busy backgrounds without changing the icon's shape. Values clamp to
// [-50, 200].
func WithGrade(grade float32) IconOption {
	return func(o *IconOptions) {
		o.Grade = grade
		o.hasGrade = true
	}
}

// WithOpticalSize sets the optical size axis explicitly, clamped to [20, 48].
// Without it, optical size adapts to the icon's rendered size.
func WithOpticalSize(opsz float32) IconOption {
	return func(o *IconOptions) {
		o.OpticalSize = opsz
		o.hasOpticalSz = true
	}
}

// WithWeight sets the weight axis, clamped to [100, 700].
//
// Weight changes the glyph design itself. To outline an unchanged icon, use
// WithStroke instead; the two are independent knobs.
func WithWeight(weight float32) IconOption {
	return func(o *IconOptions) {
		o.Weight = weight
		o.hasWeight = true
	}
}

// WithBrush paints the icon with a brush — typically a gradient — instead of a
// solid color.
//
// Precedence is brush over color over the ambient content color. Linear
// gradients of any number of stops are supported; other brush kinds fall back to
// the resolved color so the icon remains visible.
func WithBrush(brush graphics.Brush) IconOption {
	return func(o *IconOptions) {
		o.Brush = brush
	}
}

// WithStroke outlines the icon at width instead of filling it.
//
// Stroke is independent of weight: WithWeight changes the glyph design itself,
// while WithStroke adds an outline of the given width around the glyph as
// designed. A stroke does not change the icon's reported size, so it never
// perturbs layout.
func WithStroke(width unit.Dp) IconOption {
	return func(o *IconOptions) {
		o.StrokeWidth = width
		o.hasStroke = true
	}
}
