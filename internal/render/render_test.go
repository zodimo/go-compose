package render

import (
	"testing"

	"gioui.org/op"
)

// mockBackend is a minimal Backend implementation for testing.
type mockBackend struct {
	applied  []CallOp
	recorded int
}

func (m *mockBackend) BeginFrame(w, h, density float32) *Ops { return newOps() }
func (m *mockBackend) EndFrame()                             {}
func (m *mockBackend) Save() int                             { return 0 }
func (m *mockBackend) Restore(count int)                     {}
func (m *mockBackend) Translate(dx, dy float32)              {}
func (m *mockBackend) Scale(sx, sy float32)                  {}
func (m *mockBackend) ClipRect(r Rect, radius CornerRadius)  {}
func (m *mockBackend) ClipPath(pathOps []PathOp)             {}
func (m *mockBackend) FillRect(r Rect, c Color)              {}
func (m *mockBackend) FillShape(sh Shape, c Color)           {}
func (m *mockBackend) PushOpacity(alpha float32)             {}
func (m *mockBackend) DrawTextLayout(l TextLayout, at Point, c Color) {}
func (m *mockBackend) DrawImage(img Image, src, dst Rect)    {}
func (m *mockBackend) Record() CallOp {
	m.recorded++
	return CallOp{backendID: "mock", seqID: uint64(m.recorded)}
}
func (m *mockBackend) Apply(callOp CallOp) { m.applied = append(m.applied, callOp) }

// --- Registration and active-backend tests ---

func TestRegisterAndActiveBackend(t *testing.T) {
	resetForTesting()
	defer resetForTesting()

	called := false
	Register("test", func() Backend {
		called = true
		return &mockBackend{}
	})

	SetActiveBackend("test")
	if got := ActiveBackendName(); got != "test" {
		t.Errorf("ActiveBackendName() = %q, want %q", got, "test")
	}

	_ = ActiveBackend()
	if !called {
		t.Error("ActiveBackend() did not call the factory")
	}
}

func TestActiveBackendDefaultIsGio(t *testing.T) {
	resetForTesting()
	defer resetForTesting()

	// Register "gio" so SetActiveBackend doesn't panic
	Register("gio", func() Backend { return &mockBackend{} })

	// Default should be "gio"
	if got := ActiveBackendName(); got != "gio" {
		t.Errorf("default ActiveBackendName() = %q, want %q", got, "gio")
	}
}

func TestSetActiveBackendSwitches(t *testing.T) {
	resetForTesting()
	defer resetForTesting()

	Register("a", func() Backend { return &mockBackend{} })
	Register("b", func() Backend { return &mockBackend{} })

	SetActiveBackend("a")
	if got := ActiveBackendName(); got != "a" {
		t.Errorf("after SetActiveBackend(a): got %q, want %q", got, "a")
	}

	SetActiveBackend("b")
	if got := ActiveBackendName(); got != "b" {
		t.Errorf("after SetActiveBackend(b): got %q, want %q", got, "b")
	}
}

func TestActiveBackendReturnsNewInstance(t *testing.T) {
	resetForTesting()
	defer resetForTesting()

	Register("inst", func() Backend { return &mockBackend{} })
	SetActiveBackend("inst")

	b1 := ActiveBackend()
	b2 := ActiveBackend()
	if b1 == b2 {
		t.Error("ActiveBackend() should return a new instance each call")
	}
}

func TestSetActiveBackendPanicsOnUnknown(t *testing.T) {
	resetForTesting()
	defer resetForTesting()

	defer func() {
		if r := recover(); r == nil {
			t.Error("SetActiveBackend with unknown name did not panic")
		}
	}()
	SetActiveBackend("nonexistent")
}

func TestActiveBackendPanicsWhenNotRegistered(t *testing.T) {
	resetForTesting()
	defer resetForTesting()

	// "gio" is not registered, ActiveBackend should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("ActiveBackend() with unregistered backend did not panic")
		}
	}()
	ActiveBackend()
}

// --- CallOp drop semantics (EC2) ---

func TestCallOpDropSemantics(t *testing.T) {
	m := &mockBackend{}

	// Record a CallOp but never apply it.
	_ = m.Record()

	// Verify nothing was applied — drop = no side effects.
	if len(m.applied) != 0 {
		t.Errorf("mockBackend.applied = %d, want 0 (drop semantics)", len(m.applied))
	}
}

func TestCallOpApplyAfterRecord(t *testing.T) {
	m := &mockBackend{}

	// Record and then apply.
	co := m.Record()
	m.Apply(co)

	if len(m.applied) != 1 {
		t.Errorf("mockBackend.applied = %d, want 1 after Apply", len(m.applied))
	}
	if m.applied[0].seqID != 1 {
		t.Errorf("applied CallOp.seqID = %d, want 1", m.applied[0].seqID)
	}
}

func TestCallOpMultipleDropThenApply(t *testing.T) {
	m := &mockBackend{}

	// Record 3, drop 2, apply 1.
	_ = m.Record()
	_ = m.Record()
	co := m.Record()
	m.Apply(co)

	if len(m.applied) != 1 {
		t.Errorf("mockBackend.applied = %d, want 1", len(m.applied))
	}
	if m.applied[0].seqID != 3 {
		t.Errorf("applied CallOp.seqID = %d, want 3 (third recording)", m.applied[0].seqID)
	}
}

// --- Ops frame token (EC1) ---

func TestFrameTokenUniqueness(t *testing.T) {
	o1 := newOps()
	o2 := newOps()
	if o1.FrameToken() == o2.FrameToken() {
		t.Error("two Ops should have different frame tokens")
	}
	if o1.FrameToken() == 0 {
		t.Error("Ops frame token should not be zero")
	}
}

func TestOpsNilToken(t *testing.T) {
	var o *Ops
	if o.FrameToken() != 0 {
		t.Error("nil Ops.FrameToken() should return 0")
	}
}

// --- DrawCommand compatibility ---

func TestNewDrawCommandReturnsInterface(t *testing.T) {
	// NewDrawCommand should return a DrawCommand interface value
	// that wraps a gioui op.CallOp.
	var cmd DrawCommand = NewDrawCommand(op.CallOp{})
	if cmd == nil {
		t.Fatal("NewDrawCommand returned nil")
	}
}

func TestGioDrawCommandApplyToGioIsCallable(t *testing.T) {
	// Verify ApplyToGio replays into op.Ops.
	cmd := NewDrawCommand(op.CallOp{})
	var ops op.Ops
	ApplyToGio(cmd, &ops)
}

// Ensure gioDrawCommand satisfies the DrawCommand interface at compile time.
var _ DrawCommand = (*gioDrawCommand)(nil)
var _ DrawCommand = (*backendDrawCommand)(nil)
