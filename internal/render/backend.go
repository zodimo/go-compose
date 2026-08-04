// Package render implements the draw-IR seam: the framework-facing surface of
// the rendering backends. It defines the Backend interface that each rendering
// backend (gio, software, etc.) implements, plus the opaque types that flow
// across the seam (Ops, CallOp, DrawCommand).
//
// The seam is designed so that the framework public API (compose/, modifiers/,
// theme/, runtime/, pkg/) never references gioui.org types. Only this package
// and internal/layoutnode may import engine packages.
//
// Backend interface method → gioui op mapping:
//
//	BeginFrame/EndFrame       → frame lifecycle (op.Record/Stop at frame level)
//	Save/Restore              → state stack (op.StateOp)
//	Translate                 → op.Offset
//	Scale                     → op.Affine (with scale matrix)
//	ClipRect                  → clip.Rect{}.Push/Pop
//	ClipPath                  → clip.Path / clipconvert outline
//	FillRect                  → paint.Fill
//	FillShape                 → paint.FillShape
//	PushOpacity               → paint.PushOpacity
//	DrawTextLayout            → text.Layout (via material.Text or custom)
//	DrawImage                 → widget.Image
//	Record                    → op.Record (begin macro recording)
//	Apply                     → callOp.Add (replay macro)
package render

// Backend is the 2D draw IR interface modeled on Compose's Canvas.
// Each rendering backend (gio, software) implements this interface.
type Backend interface {
	// Lifecycle (per frame)
	BeginFrame(w, h float32, density float32) *Ops
	EndFrame()

	// State stack
	Save() int
	Restore(count int)

	// Transforms (used by scale/, offset/, shadow/ modifiers)
	Translate(dx, dy float32)
	Scale(sx, sy float32)

	// Clipping (used by clip/, shadow/, animation/ modifiers)
	ClipRect(r Rect, radius CornerRadius)
	ClipPath(pathOps []PathOp)

	// Paint (used by background/, border/, shadow/, alpha/ modifiers)
	FillRect(r Rect, c Color)
	FillShape(sh Shape, c Color)
	PushOpacity(alpha float32)

	// Text and images (used by text widgets, icon cache)
	DrawTextLayout(layout TextLayout, at Point, c Color)
	DrawImage(img Image, src, dst Rect)

	// Recording (macro semantics — op.Record/Stop)
	Record() CallOp
	Apply(callOp CallOp)
}

// BackendFactory creates a new Backend instance.
// Registered backends are stored in a global registry and selected by name.
type BackendFactory func() Backend

// --- Framework-owned value types ---
// These are plain structs with no gioui dependencies.

// Point is a 2D coordinate in pixels.
type Point struct{ X, Y float32 }

// Rect is an axis-aligned rectangle defined by two corners.
type Rect struct {
	Min, Max Point
}

// Color is a linear sRGB color with alpha.
type Color struct{ R, G, B, A float32 }

// CornerRadius holds per-corner radii for rounded rectangles.
type CornerRadius struct{ TL, TR, BL, BR float32 }

// ShapeKind classifies the type of shape.
type ShapeKind int

const (
	ShapeRectangle         ShapeKind = iota
	ShapeRoundedRectangle
)

// Shape describes a geometric shape for clipping or filling.
type Shape struct {
	Kind   ShapeKind
	Bounds Rect
	Radius CornerRadius
}

// PathOpKind classifies a single path operation.
type PathOpKind int

const (
	PathMoveTo PathOpKind = iota
	PathLineTo
	PathQuadTo
	PathCubicTo
	PathClose
)

// PathOp is a single operation in a path description.
type PathOp struct {
	Op       PathOpKind
	Point    Point
	Control1 Point // for quad/cubic
	Control2 Point // for cubic only
}

// TextLayout holds text rendering parameters.
// The actual text shaping is engine-specific (design D7).
type TextLayout struct {
	Width  float32
	Height float32
}

// Image is an opaque handle to a raster image.
// The actual image data is backend-specific.
type Image struct {
	Src any
}
