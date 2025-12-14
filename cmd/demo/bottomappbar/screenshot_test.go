package main

import (
	"image"
	"image/png"
	"os"
	"testing"
	"time"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/runtime"
	"github.com/zodimo/go-compose/internal/screenshot"
	"github.com/zodimo/go-compose/internal/state"
	"github.com/zodimo/go-compose/internal/store"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

// TestScreenshot renders the Bottom Bar UI and saves a screenshot.
func TestScreenshot(t *testing.T) {
	// 1. Setup Context
	var ops op.Ops
	gtx := layout.Context{
		Ops: &ops,
		Metric: unit.Metric{
			PxPerDp: 1.0,
			PxPerSp: 1.0,
		},
		Constraints: layout.Exact(image.Point{X: 400, Y: 800}), // Matching demo window size
		Now:         time.Now(),
		Locale:      system.Locale{Language: "en", Direction: system.LTR},
	}

	// 2. Init Dependencies
	themeManager := theme.GetThemeManager()
	gtx = themeManager.Material3ThemeInit(gtx)

	store := store.NewPersistentState(map[string]state.MutableValue{})
	rt := runtime.NewRuntime()

	// 3. Compose UI
	composer := compose.NewComposer(store)
	// UI() returns api.Composable which is func(c api.Composer) api.Composer
	rootComposer := UI()(composer)
	layoutNode := rootComposer.Build()

	// 4. Run Layout
	callOp := rt.Run(gtx, layoutNode)

	// 5. Take Screenshot
	img := screenshot.TakeScreenshot(400, 800, callOp)

	// 6. Save to Artifacts
	artifactPath := "/home/jaco/.gemini/antigravity/brain/ab4d5f06-8966-4915-a391-9420cf49154d/bottomappbar_screenshot.png"
	f, err := os.Create(artifactPath)
	if err != nil {
		t.Fatalf("Failed to create screenshot file: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatalf("Failed to encode screenshot: %v", err)
	}

	t.Logf("Screenshot saved to %s", artifactPath)
}
