# Screenshot of Recorded Ops

Example from gio tests```go
import 	(
    "gioui.org/gpu/headless"
)

width := lineWidth
height := 100
cap := image.NewRGBA(image.Rect(0, 0, width, height))
w, _ := headless.NewWindow(width, height)
defer w.Release()

ops := new(op.Ops)
gtx := layout.Context{
    Constraints: layout.Constraints{Max: image.Pt(width, height)},
    Ops:         ops,
}

w.Frame(ops)
w.Screenshot(cap)
b := new(bytes.Buffer)
_ = png.Encode(b, cap)
screenshotName := tc.name + ".png"
_ = os.WriteFile(screenshotName, b.Bytes(), 0o644)
t.Logf("wrote %q", screenshotName)
```

# Create a Recording
```go

import 	(
    "gioui.org/op"
)
//gtx layout.Context
macro := op.Record(gtx.Ops)
// call node.Layout(gtx) or similar
callOp := macro.Stop()
```

# Smoke tests (render.go, smoke_test.go)

`Render` runs a `Composable` through the real composition + layout path and
returns the measured size, baseline, and draw ops — **without touching the
GPU**. It mirrors `runtime.Run`, with the one addition of keeping the
`layout.Dimensions` that `runtime.Run` discards.

This exists because GPU headless rasterisation is unavailable in some
environments: `headless.NewWindow` opens an EGL context in its init path and
*blocks* (rather than erroring) where no display is usable. So verification of
rendered output comes in two layers:

- **CPU-only (always runs):** `TestSmoke_RendersWithoutGPU` and
  `TestSmoke_DeterministicAcrossRuns` assert the fixture measures to a positive,
  stable size across repeated runs. This is the regression guard for changes
  that must not alter emitted UI (e.g. the sentinel API refactor).
- **Pixel golden (skips without a GPU):** `TestSmoke_PixelGolden` rasterises via
  `TakeScreenshot` and compares against `testdata/smoke_fixture.png`. It first
  probes for a usable EGL/GPU stack and skips cleanly when there is none.
  Regenerate the golden on a GPU-capable machine with:

  ```
  go test ./internal/screenshot/ -run PixelGolden -update-golden
  ```

Note: gioui stores the serialised op stream in `gioui.org/internal/ops`
(unexported, unimportable), so the raw op bytes are not observable from outside
the library — which is why the CPU-only layer asserts measured dimensions rather
than hashing the op stream.