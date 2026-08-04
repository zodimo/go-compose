package software

import (
	"reflect"
	"testing"

	"github.com/zodimo/go-compose/internal/render"
)

// TestDeterminism verifies that rendering the same sequence twice produces
// identical DrawCall slices — the core golden-testing guarantee.
func TestDeterminism(t *testing.T) {
	renderSequence := func() []DrawCall {
		b := New().(*Backend)
		b.BeginFrame(100, 100, 1)
		b.FillRect(
			render.Rect{Min: render.Point{X: 10, Y: 10}, Max: render.Point{X: 50, Y: 50}},
			render.Color{R: 1, G: 0, B: 0, A: 1},
		)
		b.Translate(5, 5)
		b.Scale(2, 2)
		b.FillRect(
			render.Rect{Min: render.Point{X: 0, Y: 0}, Max: render.Point{X: 30, Y: 30}},
			render.Color{R: 0, G: 1, B: 0, A: 1},
		)
		b.PushOpacity(0.5)
		b.ClipRect(
			render.Rect{Min: render.Point{X: 0, Y: 0}, Max: render.Point{X: 100, Y: 100}},
			render.CornerRadius{TL: 4, TR: 4, BL: 4, BR: 4},
		)
		b.Save()
		b.Restore(1)
		b.DrawTextLayout(
			render.TextLayout{Width: 100, Height: 20},
			render.Point{X: 10, Y: 10},
			render.Color{R: 0, G: 0, B: 0, A: 1},
		)
		b.EndFrame()
		return b.Canvas().Calls()
	}

	calls1 := renderSequence()
	calls2 := renderSequence()

	if !reflect.DeepEqual(calls1, calls2) {
		t.Errorf("non-deterministic output: got %d calls first time, %d second time", len(calls1), len(calls2))
		if len(calls1) != len(calls2) {
			t.Fatalf("different lengths: %d vs %d", len(calls1), len(calls2))
		}
		for i := range calls1 {
			if !reflect.DeepEqual(calls1[i], calls2[i]) {
				t.Errorf("call %d differs:\n  first:  %+v\n  second: %+v", i, calls1[i], calls2[i])
			}
		}
	}
}

// TestDiscard verifies that a Record() that is never Apply'd emits no draw
// calls — the input-pass discard semantics.
func TestDiscard(t *testing.T) {
	b := New().(*Backend)
	b.BeginFrame(100, 100, 1)

	// Record some calls but never Apply.
	_ = b.Record()
	b.FillRect(
		render.Rect{Min: render.Point{X: 10, Y: 10}, Max: render.Point{X: 50, Y: 50}},
		render.Color{R: 1, G: 0, B: 0, A: 1},
	)
	b.FillRect(
		render.Rect{Min: render.Point{X: 20, Y: 20}, Max: render.Point{X: 60, Y: 60}},
		render.Color{R: 0, G: 1, B: 0, A: 1},
	)
	b.Translate(10, 10)

	b.EndFrame()

	calls := b.Canvas().Calls()
	for _, c := range calls {
		if c.Op == "FillRect" || c.Op == "Translate" {
			t.Errorf("unexpected %q in canvas calls when recording was discarded", c.Op)
		}
	}

	// Should only have BeginFrame and EndFrame.
	expectedOps := []string{"BeginFrame", "EndFrame"}
	if len(calls) != len(expectedOps) {
		t.Fatalf("expected %d calls (BeginFrame+EndFrame), got %d", len(expectedOps), len(calls))
	}
	for i, op := range expectedOps {
		if calls[i].Op != op {
			t.Errorf("call %d: expected %q, got %q", i, op, calls[i].Op)
		}
	}
}

// TestRecordReplay verifies that Record captures calls and Apply replays them
// into the canvas in the correct order.
func TestRecordReplay(t *testing.T) {
	b := New().(*Backend)
	b.BeginFrame(100, 100, 1)

	rec := b.Record()
	b.FillRect(
		render.Rect{Min: render.Point{X: 10, Y: 10}, Max: render.Point{X: 50, Y: 50}},
		render.Color{R: 1, G: 0, B: 0, A: 1},
	)
	b.Translate(5, 5)
	b.FillRect(
		render.Rect{Min: render.Point{X: 20, Y: 20}, Max: render.Point{X: 60, Y: 60}},
		render.Color{R: 0, G: 1, B: 0, A: 1},
	)
	b.Apply(rec)

	b.EndFrame()

	calls := b.Canvas().Calls()

	// Expected: BeginFrame, FillRect, Translate, FillRect, EndFrame
	expectedOps := []string{"BeginFrame", "FillRect", "Translate", "FillRect", "EndFrame"}
	if len(calls) != len(expectedOps) {
		t.Fatalf("expected %d calls, got %d", len(expectedOps), len(calls))
	}
	for i, op := range expectedOps {
		if calls[i].Op != op {
			t.Errorf("call %d: expected Op %q, got %q", i, op, calls[i].Op)
		}
	}

	// Verify the first FillRect has red color.
	if calls[1].Color != (render.Color{R: 1, G: 0, B: 0, A: 1}) {
		t.Errorf("call 1 (FillRect): expected red, got %v", calls[1].Color)
	}

	// Verify Translate.
	if calls[2].Point != (render.Point{X: 5, Y: 5}) {
		t.Errorf("call 2 (Translate): expected (5,5), got %v", calls[2].Point)
	}

	// Verify the second FillRect has green color.
	if calls[3].Color != (render.Color{R: 0, G: 1, B: 0, A: 1}) {
		t.Errorf("call 3 (FillRect): expected green, got %v", calls[3].Color)
	}
}

// TestRasterization verifies that FillRect produces non-empty pixels in the
// rasterized image, enabling pixel-level assertions.
func TestRasterization(t *testing.T) {
	b := New().(*Backend)
	b.BeginFrame(100, 100, 1)
	b.FillRect(
		render.Rect{Min: render.Point{X: 10, Y: 10}, Max: render.Point{X: 50, Y: 50}},
		render.Color{R: 1, G: 0, B: 0, A: 1},
	)
	b.EndFrame()

	img := b.Canvas().Image()
	if img == nil {
		t.Fatal("expected non-nil image from BeginFrame")
	}

	// Check that a pixel inside the filled rect has non-zero alpha.
	_, _, _, a := img.At(25, 25).RGBA()
	if a == 0 {
		t.Error("expected non-transparent pixel at (25, 25) inside filled rect")
	}

	// Check that a pixel outside the filled rect is transparent.
	_, _, _, a = img.At(0, 0).RGBA()
	if a != 0 {
		t.Errorf("expected transparent pixel at (0, 0) outside filled rect, got alpha %d", a)
	}

	// Check another pixel inside the rect.
	_, _, _, a = img.At(49, 49).RGBA()
	if a == 0 {
		t.Error("expected non-transparent pixel at (49, 49) inside filled rect")
	}

	// Check a pixel just outside the rect boundary.
	_, _, _, a = img.At(50, 25).RGBA()
	if a != 0 {
		t.Errorf("expected transparent pixel at (50, 25) outside filled rect, got alpha %d", a)
	}
}
