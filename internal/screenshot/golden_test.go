package screenshot_test

import (
	"bytes"
	"flag"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"gioui.org/gpu/headless"
)

var updateGolden = flag.Bool("update-golden", false,
	"rewrite golden image files instead of comparing against them")

var (
	gpuOnce   sync.Once
	gpuUsable bool
)

// gpuAvailable reports whether headless GPU rasterisation can be used in this
// environment. gioui's headless window opens an EGL context in its init path,
// which blocks (rather than erroring) on machines without a usable display, so
// the probe runs the construction under a timeout and reports false on either
// failure or timeout.
func gpuAvailable() bool {
	gpuOnce.Do(func() {
		done := make(chan bool, 1)
		go func() {
			w, err := headless.NewWindow(1, 1)
			if err != nil {
				done <- false
				return
			}
			w.Release()
			done <- true
		}()
		select {
		case ok := <-done:
			gpuUsable = ok
		case <-time.After(2 * time.Second):
			// NewWindow is stuck opening EGL; treat the GPU as unavailable.
			gpuUsable = false
		}
	})
	return gpuUsable
}

// goldenDir holds the committed golden images.
const goldenDir = "testdata"

// assertGolden compares got against the golden file named name, writing the
// file instead of comparing when -update-golden is set.
func assertGolden(t *testing.T, name string, got image.Image) {
	t.Helper()

	path := filepath.Join(goldenDir, name)

	gotPNG, err := encodePNG(got)
	if err != nil {
		t.Fatalf("encode rendered image: %v", err)
	}

	if *updateGolden {
		if err := os.MkdirAll(goldenDir, 0o755); err != nil {
			t.Fatalf("create golden dir: %v", err)
		}
		if err := os.WriteFile(path, gotPNG, 0o644); err != nil {
			t.Fatalf("write golden %s: %v", path, err)
		}
		t.Logf("updated golden %s (%d bytes)", path, len(gotPNG))
		return
	}

	wantPNG, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("golden %s missing; regenerate with: go test ./internal/screenshot/ "+
				"-run PixelGolden -update-golden", path)
		}
		t.Fatalf("read golden %s: %v", path, err)
	}

	if !bytes.Equal(gotPNG, wantPNG) {
		t.Errorf("rendered frame differs from golden %s (got %d bytes, want %d bytes); "+
			"if the change is intentional, regenerate with: go test ./internal/screenshot/ "+
			"-run PixelGolden -update-golden", path, len(gotPNG), len(wantPNG))
	}
}

func encodePNG(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
