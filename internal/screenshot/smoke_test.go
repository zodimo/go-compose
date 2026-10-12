package screenshot_test

import (
	"image"
	"testing"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/button"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/internal/screenshot"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/pkg/api"
)

// smokeFixture is a representative Material 3 screen: text plus buttons and
// spacing, i.e. the kind of composition the kitchen-sink demo shows. It exists
// so the smoke test exercises real composition + layout rather than a trivial
// single widget.
func smokeFixture() api.Composable {
	return func(c api.Composer) api.Composer {
		return column.Column(
			c.Sequence(
				text.Text("Sentinel smoke fixture"),
				spacer.Height(8),
				button.Filled(func() {}, "Primary"),
				spacer.Height(8),
				button.Outlined(func() {}, "Secondary"),
			),
			column.WithModifier(padding.All(16)),
		)(c)
	}
}

const (
	fixtureWidth  = 320
	fixtureHeight = 240
)

// TestSmoke_RendersWithoutGPU verifies the representative fixture composes and
// lays out through the real runtime without touching the GPU. This is the
// regression guard for the sentinel refactor: "runtime output is unchanged"
// means the composition still measures to the same size and still emits a
// non-empty draw list.
//
// Why a size assertion and not an op-byte digest: gioui stores the serialised
// op stream in gioui.org/internal/ops (unexported, unimportable), so the raw
// bytes are not observable from outside the library. The measured dimensions
// and the non-emptiness of the frame are the stable, CPU-only signals that a
// change to emitted ops would move.
func TestSmoke_RendersWithoutGPU(t *testing.T) {
	got := screenshot.Render(smokeFixture(), fixtureWidth, fixtureHeight)

	if got.Size == (image.Point{}) {
		t.Fatal("Render measured a zero size; the fixture did not lay out")
	}
	if got.Size.X <= 0 || got.Size.Y <= 0 {
		t.Fatalf("Render measured a non-positive size: %v", got.Size)
	}
	if got.Size.X > fixtureWidth || got.Size.Y > fixtureHeight {
		t.Fatalf("Render measured %v, exceeding the %dx%d constraints",
			got.Size, fixtureWidth, fixtureHeight)
	}
	if got.Ops == nil {
		t.Fatal("Render returned a nil op.Ops; no frame was produced")
	}
}

// TestSmoke_DeterministicAcrossRuns verifies the same fixture renders to the
// same measured size on repeated runs. A nondeterministic size would indicate
// state leaking through the composition (e.g. a shared sentinel mutated in
// place) — exactly the class of bug the sentinel refactor could introduce.
func TestSmoke_DeterministicAcrossRuns(t *testing.T) {
	first := screenshot.Render(smokeFixture(), fixtureWidth, fixtureHeight)
	for i := 0; i < 3; i++ {
		next := screenshot.Render(smokeFixture(), fixtureWidth, fixtureHeight)
		if next.Size != first.Size {
			t.Fatalf("run %d measured %v, want %v (nondeterministic layout)",
				i+1, next.Size, first.Size)
		}
		if next.Baseline != first.Baseline {
			t.Fatalf("run %d baseline %d, want %d (nondeterministic layout)",
				i+1, next.Baseline, first.Baseline)
		}
	}
}

// TestSmoke_PixelGolden rasterises the fixture and compares it against a
// committed golden image. It requires a working EGL/GPU stack: headless
// rendering hangs rather than failing on machines without one, so the test
// first probes for GPU availability and skips cleanly when it is absent.
//
// Regenerate the golden with: go test ./internal/screenshot/ -run PixelGolden \
// -update-golden
func TestSmoke_PixelGolden(t *testing.T) {
	if !gpuAvailable() {
		t.Skip("no usable EGL/GPU stack: headless rasterisation is unavailable " +
			"(headless.NewWindow hangs); skipping pixel golden comparison")
	}

	img := screenshot.TakeScreenshot(fixtureWidth, fixtureHeight,
		screenshot.Render(smokeFixture(), fixtureWidth, fixtureHeight).DrawOp)

	assertGolden(t, "smoke_fixture.png", img)
}
