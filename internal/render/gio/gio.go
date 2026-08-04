// Package gio implements the render.Backend interface using gioui.org
// rendering primitives. It maps each Backend method to the corresponding
// gio op/clip/paint operations, forming the "gio rendering backend" for
// the go-compose framework.
//
// Usage:
//
//	import _ "github.com/zodimo/go-compose/internal/render/gio"
//
//	b := render.ActiveBackend()
//	gioBackend := b.(*gio.Backend)
//	gioBackend.SetOps(ops)  // set the gio *op.Ops before BeginFrame
package gio

import (
	"image"
	"image/color"
	"sync/atomic"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/zodimo/go-compose/internal/render"
)

func init() {
	render.Register("gio", func() render.Backend {
		return New()
	})
}

// Backend implements render.Backend using gioui.org operations.
//
// The app shell must call SetOps before BeginFrame to provide the
// gioui *op.Ops that the backend writes into.
type Backend struct {
	ops             *op.Ops // set via SetOps; nil outside a frame
	seq            uint64  // monotonic sequence for CallOp IDs
	stack          []stackEntry
	recordingDepth int // >0 when inside a Record block
}

// stackEntry tracks a single state push on the unified state stack.
type stackEntry struct {
	kind      entryKind
	transform op.TransformStack
	clip      clip.Stack
	opacity   paint.OpacityStack
}

type entryKind int

const (
	kindTransform entryKind = iota
	kindClip
	kindOpacity
)

// New creates a new gio Backend ready for use.
func New() *Backend {
	return &Backend{}
}

// SetOps sets the gioui op.Ops that the backend writes into.
// Must be called before BeginFrame.
func (b *Backend) SetOps(ops *op.Ops) {
	b.ops = ops
}

// --- Lifecycle ---

// BeginFrame stores frame parameters and returns a fresh Ops token.
func (b *Backend) BeginFrame(w, h, density float32) *render.Ops {
	// ops must already be set via SetOps.
	return render.NewFrameOps()
}

// EndFrame releases the current frame's op.Ops reference.
func (b *Backend) EndFrame() {
b.ops = nil
	b.stack = b.stack[:0]
	b.recordingDepth = 0
}

// --- State stack ---

// Save snapshots the current stack depth and returns it.
// Restore(toLevel) pops all entries above that level.
//
// Save does NOT push any gio ops — it is a pure bookkeeping snapshot.
func (b *Backend) Save() int {
	return len(b.stack)
}

// Restore pops all state entries back to the given save level.
// Entries are popped in reverse (LIFO) order, correctly undoing
// transforms, clips, and opacity layers.
func (b *Backend) Restore(toLevel int) {
	for len(b.stack) > toLevel {
		top := b.stack[len(b.stack)-1]
		b.stack = b.stack[:len(b.stack)-1]
		switch top.kind {
		case kindTransform:
			top.transform.Pop()
		case kindClip:
			top.clip.Pop()
		case kindOpacity:
			top.opacity.Pop()
		}
	}
}

// --- Transforms ---

// Translate applies a 2D translation. The translation is undone by
// the next Restore call (matching Compose's Canvas semantics).
func (b *Backend) Translate(dx, dy float32) {
	t := op.Affine(f32.AffineId().Offset(f32.Pt(dx, dy))).Push(b.ops)
	if b.recordingDepth == 0 {
b.stack = append(b.stack, stackEntry{kind: kindTransform, transform: t})
}
}

// Scale applies a 2D scale around the origin.
func (b *Backend) Scale(sx, sy float32) {
	t := op.Affine(f32.AffineId().Scale(f32.Point{}, f32.Pt(sx, sy))).Push(b.ops)
	if b.recordingDepth == 0 {
b.stack = append(b.stack, stackEntry{kind: kindTransform, transform: t})
}
}

// --- Clipping ---

// ClipRect clips to a rectangle. If radius is non-zero, a rounded-
// rectangle clip is used (via clip.RRect).
func (b *Backend) ClipRect(r render.Rect, radius render.CornerRadius) {
rect := toImageRect(r)
var cs clip.Stack
if isZeroRadius(radius) {
cs = clip.Rect(rect).Push(b.ops)
} else {
cs = clip.RRect{
Rect: rect,
NW:   int(radius.TL),
NE:   int(radius.TR),
SW:   int(radius.BL),
SE:   int(radius.BR),
}.Push(b.ops)
	}
	if b.recordingDepth == 0 {
b.stack = append(b.stack, stackEntry{kind: kindClip, clip: cs})
}
}

// ClipPath clips to a path described by PathOp entries.
func (b *Backend) ClipPath(pathOps []render.PathOp) {
var p clip.Path
p.Begin(b.ops)
for _, po := range pathOps {
pt := f32.Pt(po.Point.X, po.Point.Y)
switch po.Op {
case render.PathMoveTo:
p.MoveTo(pt)
case render.PathLineTo:
p.LineTo(pt)
case render.PathQuadTo:
p.QuadTo(f32.Pt(po.Control1.X, po.Control1.Y), pt)
case render.PathCubicTo:
p.CubeTo(
f32.Pt(po.Control1.X, po.Control1.Y),
f32.Pt(po.Control2.X, po.Control2.Y),
pt,
)
case render.PathClose:
p.Close()
}
}
pathSpec := p.End()
	cs := clip.Outline{Path: pathSpec}.Op().Push(b.ops)
	if b.recordingDepth == 0 {
b.stack = append(b.stack, stackEntry{kind: kindClip, clip: cs})
}
}

// --- Paint ---

// FillRect fills a rectangle with a color. Emits clip.Rect + paint.FillShape.
func (b *Backend) FillRect(r render.Rect, c render.Color) {
	paint.FillShape(b.ops, toNRGBA(c), clip.Rect(toImageRect(r)).Op())
}

// FillShape fills a shape with a color. For rounded rectangles, uses
// clip.RRect; for plain rectangles, uses clip.Rect.
func (b *Backend) FillShape(sh render.Shape, c render.Color) {
	rect := toImageRect(sh.Bounds)
	nrgba := toNRGBA(c)
	switch sh.Kind {
	case render.ShapeRoundedRectangle:
		paint.FillShape(b.ops, nrgba, clip.RRect{
			Rect: rect,
			NW:   int(sh.Radius.TL),
			NE:   int(sh.Radius.TR),
			SW:   int(sh.Radius.BL),
			SE:   int(sh.Radius.BR),
		}.Op(b.ops))
	default: // ShapeRectangle and others
		paint.FillShape(b.ops, nrgba, clip.Rect(rect).Op())
	}
}

// PushOpacity pushes an opacity layer. All subsequent drawing is
// composited with the given alpha until the next Restore.
func (b *Backend) PushOpacity(alpha float32) {
	os := paint.PushOpacity(b.ops, clamp01(alpha))
	if b.recordingDepth == 0 {
b.stack = append(b.stack, stackEntry{kind: kindOpacity, opacity: os})
}
}

// --- Text and images (best-effort) ---

// DrawTextLayout draws a text layout as a filled rectangle representing
// the text bounds. The actual text shaping is engine-specific (design D7).
func (b *Backend) DrawTextLayout(layout render.TextLayout, at render.Point, c render.Color) {
	if layout.Width <= 0 || layout.Height <= 0 {
		return
	}
	rect := image.Rectangle{
		Min: image.Point{X: int(at.X), Y: int(at.Y)},
		Max: image.Point{X: int(at.X + layout.Width), Y: int(at.Y + layout.Height)},
	}
	paint.FillShape(b.ops, toNRGBA(c), clip.Rect(rect).Op())
}

// DrawImage draws an image, scaling from src to dst rectangles.
// The image handle must contain an image.Image.
func (b *Backend) DrawImage(img render.Image, src, dst render.Rect) {
	if img.Src == nil {
		return
	}
	im, ok := img.Src.(image.Image)
	if !ok {
		return
	}
	srcW := src.Max.X - src.Min.X
	srcH := src.Max.Y - src.Min.Y
	dstW := dst.Max.X - dst.Min.X
	dstH := dst.Max.Y - dst.Min.Y
	if srcW <= 0 || srcH <= 0 || dstW <= 0 || dstH <= 0 {
		return
	}

	sx := dstW / srcW
	sy := dstH / srcH

	// Translate to dst origin, then scale from src size to dst size.
	stack := op.Affine(f32.AffineId().
		Offset(f32.Pt(dst.Min.X, dst.Min.Y)).
		Scale(f32.Point{}, f32.Pt(sx, sy))).
		Push(b.ops)
	paint.NewImageOp(im).Add(b.ops)
	paint.PaintOp{}.Add(b.ops)
	stack.Pop()
}

// --- Recording (macro semantics) ---

// Record begins recording a macro. Operations drawn after Record are
// captured until Apply is called. A recorded-but-never-applied CallOp
// emits nothing (drop semantics, EC2) — inherent to op.MacroOp.
func (b *Backend) Record() render.CallOp {
	b.recordingDepth++
macro := op.Record(b.ops)
return render.NewCallOp("gio", atomic.AddUint64(&b.seq, 1), macro)
}

// Apply replays a previously recorded CallOp by stopping the macro
// and adding it to the current ops. Only CallOps created by this
// backend (backendID == "gio") are applied.
func (b *Backend) Apply(callOp render.CallOp) {
	if b.recordingDepth > 0 {
		b.recordingDepth--
	}
if render.CallOpBackendID(callOp) != "gio" {
return
}
payload := render.CallOpPayload(callOp)
if payload == nil {
return
}
macro, ok := payload.(op.MacroOp)
if !ok {
return
}
macro.Stop().Add(b.ops)
}

// --- Helpers ---

func toImageRect(r render.Rect) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: int(r.Min.X), Y: int(r.Min.Y)},
		Max: image.Point{X: int(r.Max.X), Y: int(r.Max.Y)},
	}
}

func isZeroRadius(r render.CornerRadius) bool {
	return r.TL == 0 && r.TR == 0 && r.BL == 0 && r.BR == 0
}

func toNRGBA(c render.Color) color.NRGBA {
	return color.NRGBA{
		R: uint8(clamp01(c.R)*255 + 0.5),
		G: uint8(clamp01(c.G)*255 + 0.5),
		B: uint8(clamp01(c.B)*255 + 0.5),
		A: uint8(clamp01(c.A)*255 + 0.5),
	}
}

func clamp01(v float32) float32 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
