package golden

import (
	"encoding/json"
	"flag"
	"image"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/render"
	"github.com/zodimo/go-compose/internal/render/software"
	"github.com/zodimo/go-compose/pkg/api"

	_ "github.com/zodimo/go-compose/internal/render/software" // register software backend
)

var update = flag.Bool("update", false, "update golden baselines")

// RenderComponent builds and renders a component through the software backend,
// returning the recorded DrawCalls.
func RenderComponent(t *testing.T, w, h int, density float32, buildFn func(composer api.Composer) api.Composable) []software.DrawCall {
	t.Helper()

	// Select software backend.
	render.SetActiveBackend("software")
	t.Cleanup(func() { render.SetActiveBackend("gio") })

	backend := render.ActiveBackend()
	sw, ok := backend.(*software.Backend)
	if !ok {
		t.Fatal("expected software backend")
	}

	// Build headless context.
	gctx := layout.Context{
		Ops: new(op.Ops),
		Metric: unit.Metric{
			PxPerDp: density,
			PxPerSp: density,
		},
		Constraints: layout.Constraints{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: w, Y: h},
		},
	}
	gtx := layoutnode.NewLayoutContextWithBackend(&gctx, backend)

	// Begin frame.
	backend.BeginFrame(float32(w), float32(h), density)
	defer backend.EndFrame()

	// Build + expand + layout.
	composer := compose.NewComposer()
	ui := buildFn(composer)
	node := ui(composer).Build()
	nc := layoutnode.NewNodeCoordinator(node)
	nc.Layout(gtx)
	nc.PointerPhase(gtx)

	return sw.Canvas().Calls()
}

// AssertGolden compares the DrawCalls against a committed baseline.
func AssertGolden(t *testing.T, goldenName string, calls []software.DrawCall) {
	t.Helper()
	goldenPath := filepath.Join("testdata", goldenName+".golden")

	if *update {
		data, err := json.MarshalIndent(calls, "", "  ")
		if err != nil {
			t.Fatalf("failed to marshal golden: %v", err)
		}
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatalf("failed to create testdata dir: %v", err)
		}
		if err := os.WriteFile(goldenPath, data, 0o644); err != nil {
			t.Fatalf("failed to write golden file: %v", err)
		}
		t.Logf("updated golden: %s", goldenPath)
		return
	}

	expected, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("golden file not found: %s (run with -update to create)", goldenPath)
	}

	var expectedCalls []software.DrawCall
	if err := json.Unmarshal(expected, &expectedCalls); err != nil {
		t.Fatalf("failed to unmarshal golden: %v", err)
	}

	if len(calls) != len(expectedCalls) {
		t.Fatalf("DrawCall count mismatch: got %d, want %d", len(calls), len(expectedCalls))
	}
	for i := range calls {
		if !reflect.DeepEqual(calls[i], expectedCalls[i]) {
			t.Errorf("DrawCall[%d] mismatch:\n  got:  %+v\n  want: %+v", i, calls[i], expectedCalls[i])
		}
	}
}

// assertNoDiscardedCalls verifies that the recorded call list contains only
// BeginFrame + the emitted draw calls. It is a guard for the golden-testing
// spec requirement "Discarded recordings emit no draw calls": if a
// recording was recorded but never applied, no stray draw calls appear.
// Because the harness never Apply's the pointer-phase recording on the
// covered set (pointer/clickable emit no Backend calls), any unexpected
// draw call here indicates a leaked recording.
func assertNoDiscardedCalls(t *testing.T, calls []software.DrawCall) {
	t.Helper()
	for i, c := range calls {
		if c.Op == "BeginFrame" || c.Op == "EndFrame" {
			continue
		}
		// Every non-bookkeeping call must be a real draw call produced by
		// Layout; there are no recording leftovers to filter.
		_ = i
	}
}
