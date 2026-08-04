// Package software implements a deterministic recording backend for golden tests.
// It records every draw operation as an ordered DrawCall slice, enabling
// pixel-level and call-sequence assertions. Output is fully deterministic:
// no map iteration, no wall-clock time, no pointer addresses in payloads.
package software

import (
	"image"
	"image/color"

	"github.com/zodimo/go-compose/internal/render"
)

// DrawCall represents a single recorded draw operation.
// Each field corresponds to a Backend method parameter.
// Zero-valued fields are unused for a given Op.
type DrawCall struct {
	Op           string
	Rect         render.Rect
	Color        render.Color
	Point        render.Point
	ScaleX       float32
	ScaleY       float32
	Alpha        float32
	Shape        render.Shape
	TextLayout   render.TextLayout
	Image        render.Image
	Src          render.Rect
	Dst          render.Rect
	CornerRadius render.CornerRadius
	PathOps      []render.PathOp
	Count        int
}

// recording holds a sequence of DrawCalls captured during a Record session.
type recording struct {
	calls      []DrawCall
	frameToken uint64
}

// Canvas provides access to the recorded draw calls and optional rasterization.
type Canvas struct {
	calls []DrawCall
	img   *image.RGBA
}

// Calls returns a copy of the recorded draw calls.
func (c *Canvas) Calls() []DrawCall {
	if len(c.calls) == 0 {
		return nil
	}
	dst := make([]DrawCall, len(c.calls))
	copy(dst, c.calls)
	return dst
}

// Image returns the rasterized image, or nil if BeginFrame was not called.
func (c *Canvas) Image() *image.RGBA {
	return c.img
}

// Backend implements render.Backend as a deterministic recording backend
// for golden tests. It records every draw operation as an ordered DrawCall
// slice. Output is fully deterministic: no map iteration, no wall-clock time,
// no pointer addresses in payloads.
type Backend struct {
	canvas     Canvas
	recording  *recording // non-nil when in recording mode
	recSeqID   uint64     // monotonically increasing sequence ID for recordings
	frameToken uint64     // current frame token
	saveStack  []int      // each entry is the call count at Save time
}

// New returns a new software Backend that implements render.Backend.
func New() render.Backend {
	return &Backend{}
}

func init() {
	render.Register("software", func() render.Backend { return New() })
}

// Canvas returns the backend's Canvas for inspection and assertion.
func (b *Backend) Canvas() *Canvas {
	return &b.canvas
}

// append adds a DrawCall to the current target (recording buffer or main canvas).
func (b *Backend) append(dc DrawCall) {
	if b.recording != nil {
		b.recording.calls = append(b.recording.calls, dc)
	} else {
		b.canvas.calls = append(b.canvas.calls, dc)
	}
}

// BeginFrame starts a new frame, allocating the rasterization target.
func (b *Backend) BeginFrame(w, h float32, density float32) *render.Ops {
	// Reset state for the new frame.
	b.recording = nil
	b.canvas.calls = b.canvas.calls[:0]
	b.saveStack = b.saveStack[:0]
	b.frameToken++

	// Allocate rasterization target.
	iw, ih := int(w), int(h)
	if iw > 0 && ih > 0 {
		b.canvas.img = image.NewRGBA(image.Rect(0, 0, iw, ih))
	} else {
		b.canvas.img = nil
	}

	b.append(DrawCall{
		Op:   "BeginFrame",
		Rect: render.Rect{Min: render.Point{}, Max: render.Point{X: w, Y: h}},
	})

	return render.NewFrameOps()
}

// EndFrame finishes the current frame.
func (b *Backend) EndFrame() {
	b.recording = nil
	b.append(DrawCall{Op: "EndFrame"})
}

// Save pushes the current state and returns a token for Restore.
func (b *Backend) Save() int {
	b.saveStack = append(b.saveStack, len(b.canvas.calls))
	b.append(DrawCall{Op: "Save"})
	return len(b.saveStack)
}

// Restore pops count states from the stack.
func (b *Backend) Restore(count int) {
	b.append(DrawCall{Op: "Restore", Count: count})
	if count > 0 && count <= len(b.saveStack) {
		b.saveStack = b.saveStack[:len(b.saveStack)-count]
	}
}

// Translate records a translation.
func (b *Backend) Translate(dx, dy float32) {
	b.append(DrawCall{
		Op:    "Translate",
		Point: render.Point{X: dx, Y: dy},
	})
}

// Scale records a scale.
func (b *Backend) Scale(sx, sy float32) {
	b.append(DrawCall{
		Op:     "Scale",
		ScaleX: sx,
		ScaleY: sy,
	})
}

// ClipRect records a rectangular clip.
func (b *Backend) ClipRect(r render.Rect, radius render.CornerRadius) {
	b.append(DrawCall{
		Op:           "ClipRect",
		Rect:         r,
		CornerRadius: radius,
	})
}

// ClipPath records a path clip.
func (b *Backend) ClipPath(pathOps []render.PathOp) {
	b.append(DrawCall{
		Op:      "ClipPath",
		PathOps: pathOps,
	})
}

// FillRect fills a rectangle with a color. Also rasterizes into the image.
func (b *Backend) FillRect(r render.Rect, c render.Color) {
	b.append(DrawCall{
		Op:    "FillRect",
		Rect:  r,
		Color: c,
	})
	b.rasterizeRect(r, c)
}

// FillShape fills a shape with a color. Also rasterizes the shape bounds.
func (b *Backend) FillShape(sh render.Shape, c render.Color) {
	b.append(DrawCall{
		Op:    "FillShape",
		Shape: sh,
		Color: c,
	})
	b.rasterizeRect(sh.Bounds, c)
}

// PushOpacity records an opacity change.
func (b *Backend) PushOpacity(alpha float32) {
	b.append(DrawCall{
		Op:    "PushOpacity",
		Alpha: alpha,
	})
}

// DrawTextLayout records a text draw operation.
func (b *Backend) DrawTextLayout(layout render.TextLayout, at render.Point, c render.Color) {
	b.append(DrawCall{
		Op:         "DrawTextLayout",
		TextLayout: layout,
		Point:      at,
		Color:      c,
	})
}

// DrawImage records an image draw operation.
func (b *Backend) DrawImage(img render.Image, src, dst render.Rect) {
	b.append(DrawCall{
		Op:    "DrawImage",
		Image: img,
		Src:   src,
		Dst:   dst,
	})
}

// Record starts a recording session. Draw calls following Record are captured
// into a buffer that is returned as a CallOp. A Record that is never Apply'd
// emits nothing (discard semantics — the input-pass discards the recording).
func (b *Backend) Record() render.CallOp {
	b.recSeqID++
	rec := &recording{frameToken: b.frameToken}
	b.recording = rec
	return render.NewCallOp("software", b.recSeqID, rec)
}

// Apply replays a recording's draw calls into the canvas.
func (b *Backend) Apply(callOp render.CallOp) {
	payload := render.CallOpPayload(callOp)
	if payload == nil {
		return
	}
	rec, ok := payload.(*recording)
	if !ok {
		return
	}
	b.canvas.calls = append(b.canvas.calls, rec.calls...)
	b.recording = nil
}

// rasterizeRect performs a simple pixel fill of a rectangle in the image.
func (b *Backend) rasterizeRect(r render.Rect, c render.Color) {
	img := b.canvas.img
	if img == nil {
		return
	}

	bounds := img.Bounds()
	x0 := int(r.Min.X)
	y0 := int(r.Min.Y)
	x1 := int(r.Max.X)
	y1 := int(r.Max.Y)

	// Clamp to image bounds.
	if x0 < bounds.Min.X {
		x0 = bounds.Min.X
	}
	if y0 < bounds.Min.Y {
		y0 = bounds.Min.Y
	}
	if x1 > bounds.Max.X {
		x1 = bounds.Max.X
	}
	if y1 > bounds.Max.Y {
		y1 = bounds.Max.Y
	}
	if x0 >= x1 || y0 >= y1 {
		return
	}

	// Convert linear sRGB [0,1] to NRGBA [0,255] and set pixel.
	nc := color.NRGBA{R: clampByte(c.R), G: clampByte(c.G), B: clampByte(c.B), A: clampByte(c.A)}

	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			img.Set(x, y, nc)
		}
	}
}

// clampByte converts a [0,1] float to a [0,255] uint8.
func clampByte(f float32) uint8 {
	if f <= 0 {
		return 0
	}
	if f >= 1 {
		return 255
	}
	return uint8(f*255 + 0.5)
}
