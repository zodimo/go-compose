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
	"github.com/zodimo/go-compose/state"
	"github.com/zodimo/go-compose/store"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

// TestScreenshot renders the Tabs Demo UI and saves a screenshot.
func TestScreenshot(t *testing.T) {
	// 1. Setup Context
	var ops op.Ops
	gtx := layout.Context{
		Ops: &ops,
		Metric: unit.Metric{
			PxPerDp: 1.0,
			PxPerSp: 1.0,
		},
		Constraints: layout.Exact(image.Point{X: 600, Y: 400}), // Matching demo window size
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
	rootComposer := UI()(composer)
	layoutNode := rootComposer.Build()

	// 4. Run Layout
	callOp := rt.Run(gtx, layoutNode)

	// 5. Take Screenshot
	img := screenshot.TakeScreenshot(600, 400, callOp)

	// 6. Save to Artifacts
	// Using the provided brain directory
	artifactPath := "/home/jaco/.gemini/antigravity/brain/d80e18d8-0f91-48a2-a106-4ba08b688911/tabs_screenshot.png"
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
