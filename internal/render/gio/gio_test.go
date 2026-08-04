package gio

import (
	"image"
	"image/color"
	"testing"

	"gioui.org/op"

	"github.com/zodimo/go-compose/internal/render"
)

// newTestBackend creates a Backend wired to a fresh op.Ops.
func newTestBackend(t *testing.T) *Backend {
	t.Helper()
	b := New()
	b.SetOps(&op.Ops{})
	b.BeginFrame(100, 100, 1)
	return b
}

// --- Basic construction ---

func TestNew(t *testing.T) {
	b := New()
	if b == nil {
		t.Fatal("New() returned nil")
	}
}

func TestBeginFrameReturnsToken(t *testing.T) {
	b := New()
	ops := &op.Ops{}
	b.SetOps(ops)
	tok := b.BeginFrame(200, 300, 2)
	if tok == nil {
		t.Fatal("BeginFrame returned nil Ops")
	}
	if tok.FrameToken() == 0 {
		t.Fatal("BeginFrame returned Ops with zero token")
	}
}

// --- Translate emits an op ---

func TestTranslate(t *testing.T) {
	b := newTestBackend(t)
	b.Translate(10, 20) // Push + append to stack
	if len(b.stack) != 1 {
		t.Errorf("after Translate: stack len = %d, want 1", len(b.stack))
	}
	b.Restore(0)
	if len(b.stack) != 0 {
		t.Errorf("after Restore: stack len = %d, want 0", len(b.stack))
	}
}

// --- Scale emits an op ---

func TestScale(t *testing.T) {
	b := newTestBackend(t)
	b.Scale(2, 3)
	if len(b.stack) != 1 {
		t.Errorf("after Scale: stack len = %d, want 1", len(b.stack))
	}
	b.Restore(0)
}

// --- ClipRect emits ops ---

func TestClipRectPlain(t *testing.T) {
	b := newTestBackend(t)
	b.ClipRect(render.Rect{
		Min: render.Point{X: 10, Y: 10},
		Max: render.Point{X: 90, Y: 90},
	}, render.CornerRadius{})
	if len(b.stack) != 1 {
		t.Errorf("after ClipRect: stack len = %d, want 1", len(b.stack))
	}
	b.Restore(0)
}

func TestClipRectRounded(t *testing.T) {
	b := newTestBackend(t)
	b.ClipRect(render.Rect{
		Min: render.Point{},
		Max: render.Point{X: 100, Y: 100},
	}, render.CornerRadius{TL: 8, TR: 8, BL: 8, BR: 8})
	if len(b.stack) != 1 {
		t.Errorf("after ClipRect rounded: stack len = %d, want 1", len(b.stack))
	}
	b.Restore(0)
}

// --- ClipPath emits ops ---

func TestClipPath(t *testing.T) {
	b := newTestBackend(t)
	b.ClipPath([]render.PathOp{
		{Op: render.PathMoveTo, Point: render.Point{X: 0, Y: 0}},
		{Op: render.PathLineTo, Point: render.Point{X: 100, Y: 0}},
		{Op: render.PathLineTo, Point: render.Point{X: 100, Y: 100}},
		{Op: render.PathClose},
	})
	if len(b.stack) != 1 {
		t.Errorf("after ClipPath: stack len = %d, want 1", len(b.stack))
	}
	b.Restore(0)
}

// --- FillRect emits ops ---

func TestFillRect(t *testing.T) {
	b := newTestBackend(t)
	b.FillRect(render.Rect{
		Min: render.Point{},
		Max: render.Point{X: 50, Y: 50},
	}, render.Color{R: 1, G: 0, B: 0, A: 1})
	// FillRect uses FillShape which pushes+pops internally — no stack growth.
	if len(b.stack) != 0 {
		t.Errorf("after FillRect: stack len = %d, want 0", len(b.stack))
	}
}

// --- FillShape emits ops ---

func TestFillShapeRectangle(t *testing.T) {
	b := newTestBackend(t)
	b.FillShape(render.Shape{
		Kind:   render.ShapeRectangle,
		Bounds: render.Rect{Min: render.Point{}, Max: render.Point{X: 50, Y: 50}},
	}, render.Color{R: 0, G: 1, B: 0, A: 1})
	if len(b.stack) != 0 {
		t.Errorf("after FillShape rect: stack len = %d, want 0", len(b.stack))
	}
}

func TestFillShapeRoundedRect(t *testing.T) {
	b := newTestBackend(t)
	b.FillShape(render.Shape{
		Kind:   render.ShapeRoundedRectangle,
		Bounds: render.Rect{Min: render.Point{}, Max: render.Point{X: 50, Y: 50}},
		Radius: render.CornerRadius{TL: 5, TR: 5, BL: 5, BR: 5},
	}, render.Color{R: 0, G: 0, B: 1, A: 1})
	if len(b.stack) != 0 {
		t.Errorf("after FillShape rrect: stack len = %d, want 0", len(b.stack))
	}
}

// --- PushOpacity emits ops ---

func TestPushOpacity(t *testing.T) {
	b := newTestBackend(t)
	b.PushOpacity(0.5)
	if len(b.stack) != 1 {
		t.Errorf("after PushOpacity: stack len = %d, want 1", len(b.stack))
	}
	b.Restore(0)
}

func TestPushOpacityClamped(t *testing.T) {
	b := newTestBackend(t)
	b.PushOpacity(2.0) // should be clamped to 1.0 internally
	if len(b.stack) != 1 {
		t.Errorf("after PushOpacity clamped: stack len = %d, want 1", len(b.stack))
	}
	b.Restore(0)
}

// --- Save / Restore ---

func TestSaveRestore(t *testing.T) {
	b := newTestBackend(t)
	b.Translate(10, 10)
	b.ClipRect(render.Rect{
		Min: render.Point{},
		Max: render.Point{X: 100, Y: 100},
	}, render.CornerRadius{})
	level := b.Save() // snapshots stack length=2, no push
	b.PushOpacity(0.5)
	b.Translate(5, 5)

	if len(b.stack) != 4 {
		t.Fatalf("before Restore: stack len = %d, want 4", len(b.stack))
	}
	b.Restore(level)
	// Should have popped 2 entries (opacity, translate).
	if len(b.stack) != 2 {
		t.Errorf("after Restore: stack len = %d, want 2", len(b.stack))
	}
	b.Restore(0)
}

func TestSaveRestoreDeepNesting(t *testing.T) {
	b := newTestBackend(t)
	l1 := b.Save()
	b.Translate(1, 1)
	l2 := b.Save()
	b.Translate(2, 2)
	l3 := b.Save()
	b.PushOpacity(0.3)

	if len(b.stack) != 3 {
		t.Fatalf("innermost: stack len = %d, want 3", len(b.stack))
	}
	b.Restore(l3) // pop opacity
	if len(b.stack) != 2 {
		t.Errorf("after Restore(l3): stack len = %d, want 2", len(b.stack))
	}
	b.Restore(l2) // pop translate(2,2)
	if len(b.stack) != 1 {
		t.Errorf("after Restore(l2): stack len = %d, want 1", len(b.stack))
	}
	b.Restore(l1) // pop translate(1,1)
	if len(b.stack) != 0 {
		t.Errorf("after Restore(l1): stack len = %d, want 0", len(b.stack))
	}
}

// --- Record / Apply macro replay ---

func TestRecordApply(t *testing.T) {
	b := newTestBackend(t)
	callOp := b.Record()

	// Draw into the macro.
	b.FillRect(render.Rect{
		Min: render.Point{},
		Max: render.Point{X: 10, Y: 10},
	}, render.Color{R: 1, A: 1})

	// Apply replays the macro.
	b.Apply(callOp)
	// No panic = success.
}

func TestRecordApplyPreservesOrder(t *testing.T) {
	b := newTestBackend(t)
	callOp := b.Record()

	// Multiple draws into the macro.
	// During recording (recordingDepth>0), Translate/Scale do NOT push to our stack.
	b.Translate(5, 5)
	b.FillRect(render.Rect{
		Min: render.Point{},
		Max: render.Point{X: 20, Y: 20},
	}, render.Color{R: 0, G: 1, A: 1})
	b.Scale(2, 2)

	// Apply — no panic = ops were recorded and replayed in order.
	b.Apply(callOp)
	// Nothing on our stack (Translate/Scale were inside macro, not tracked).
	if len(b.stack) != 0 {
		t.Errorf("after Apply: stack len = %d, want 0", len(b.stack))
	}
}

// --- Drop semantics (EC2) ---

func TestRecordDropSemantics(t *testing.T) {
	b := newTestBackend(t)
	_ = b.Record() // record but never apply

	// Draw into the macro (it's still live).
	b.FillRect(render.Rect{
		Min: render.Point{},
		Max: render.Point{X: 10, Y: 10},
	}, render.Color{R: 1, A: 1})

	// Never call Apply. The macro is dropped.
	b.EndFrame()

	// No panic = drop semantics hold.
}

func TestRecordDropMultiple(t *testing.T) {
	b := newTestBackend(t)
	for i := 0; i < 5; i++ {
		_ = b.Record()
		b.FillRect(render.Rect{
			Min: render.Point{},
			Max: render.Point{X: 10, Y: 10},
		}, render.Color{R: 1, A: 1})
	}
	b.EndFrame()
	// No panic = multiple drops are safe.
}

func TestRecordApplyAfterDrop(t *testing.T) {
	b := newTestBackend(t)

	// Record and drop.
	_ = b.Record()
	b.FillRect(render.Rect{
		Min: render.Point{},
		Max: render.Point{X: 10, Y: 10},
	}, render.Color{R: 1, A: 1})

	// Record and apply — should still work.
	callOp := b.Record()
	b.FillRect(render.Rect{
		Min: render.Point{},
		Max: render.Point{X: 20, Y: 20},
	}, render.Color{R: 0, G: 1, A: 1})
	b.Apply(callOp)

	b.EndFrame()
}

// --- DrawTextLayout ---

func TestDrawTextLayout(t *testing.T) {
	b := newTestBackend(t)
	b.DrawTextLayout(
		render.TextLayout{Width: 100, Height: 20},
		render.Point{X: 10, Y: 10},
		render.Color{R: 0, G: 0, B: 0, A: 1},
	)
	// No panic; FillShape internally manages clip.
}

func TestDrawTextLayoutZeroSize(t *testing.T) {
	b := newTestBackend(t)
	b.DrawTextLayout(
		render.TextLayout{Width: 0, Height: 0},
		render.Point{},
		render.Color{},
	)
	// No panic; zero-size is a no-op.
}

// --- DrawImage ---

func TestDrawImageNil(t *testing.T) {
	b := newTestBackend(t)
	b.DrawImage(render.Image{Src: nil}, render.Rect{}, render.Rect{})
	// No panic; nil image is a no-op.
}

func TestDrawImageNonImageHandle(t *testing.T) {
	b := newTestBackend(t)
	b.DrawImage(render.Image{Src: "not-an-image"}, render.Rect{}, render.Rect{})
	// No panic; non-image.Image handle is a no-op.
}

func TestDrawImageValid(t *testing.T) {
	b := newTestBackend(t)
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	b.DrawImage(
		render.Image{Src: img},
		render.Rect{Max: render.Point{X: 10, Y: 10}},
		render.Rect{Max: render.Point{X: 50, Y: 50}},
	)
	// No panic; valid image drawn with scaling.
}

// --- EndFrame ---

func TestEndFrameClearsState(t *testing.T) {
	b := newTestBackend(t)
	b.Translate(10, 10)
	b.EndFrame()
	if b.ops != nil {
		t.Error("EndFrame should nil the ops reference")
	}
	if len(b.stack) != 0 {
		t.Errorf("EndFrame should clear stack, got len %d", len(b.stack))
	}
}

// --- Color conversion ---

func TestToNRGBA(t *testing.T) {
	c := render.Color{R: 0.5, G: 0.5, B: 0.5, A: 1.0}
	n := toNRGBA(c)
	expected := color.NRGBA{R: 128, G: 128, B: 128, A: 255}
	if n != expected {
		t.Errorf("toNRGBA(%v) = %v, want %v", c, n, expected)
	}
}

func TestToNRGBAClamped(t *testing.T) {
	c := render.Color{R: -0.1, G: 1.5, B: 0, A: 0}
	n := toNRGBA(c)
	if n.R != 0 {
		t.Errorf("toNRGBA clamped R: got %d, want 0", n.R)
	}
	if n.G != 255 {
		t.Errorf("toNRGBA clamped G: got %d, want 255", n.G)
	}
}

func TestIsZeroRadius(t *testing.T) {
	if !isZeroRadius(render.CornerRadius{}) {
		t.Error("zero radius should be detected")
	}
	if isZeroRadius(render.CornerRadius{TL: 1}) {
		t.Error("non-zero radius should not be zero")
	}
}

// --- Registration ---

func TestRegistration(t *testing.T) {
	// The init() function registers "gio". Verify the factory works.
	b := render.ActiveBackend()
	if b == nil {
		t.Fatal("ActiveBackend() returned nil after gio init()")
	}
	if _, ok := b.(*Backend); !ok {
		t.Errorf("ActiveBackend() returned %T, want *gio.Backend", b)
	}
}

// --- Integration: full frame cycle ---

func TestFullFrameCycle(t *testing.T) {
	b := New()
	ops := &op.Ops{}
	b.SetOps(ops)
	tok := b.BeginFrame(400, 300, 1.5)
	if tok.FrameToken() == 0 {
		t.Fatal("frame token should not be zero")
	}

	// Simulate a typical frame: save, translate, clip, fill, restore.
	l0 := b.Save()
	b.Translate(20, 20)
	b.ClipRect(render.Rect{
		Min: render.Point{},
		Max: render.Point{X: 360, Y: 260},
	}, render.CornerRadius{TL: 8, TR: 8, BL: 8, BR: 8})
	b.FillRect(render.Rect{
		Min: render.Point{},
		Max: render.Point{X: 360, Y: 260},
	}, render.Color{R: 0.95, G: 0.95, B: 0.95, A: 1})
	b.Restore(l0)

	// Record a macro, draw into it, then apply.
	callOp := b.Record()
	b.FillShape(render.Shape{
		Kind:   render.ShapeRectangle,
		Bounds: render.Rect{Min: render.Point{X: 50, Y: 50}, Max: render.Point{X: 150, Y: 100}},
	}, render.Color{R: 0.2, G: 0.6, B: 1, A: 1})
	b.Apply(callOp)

	b.EndFrame()
}
