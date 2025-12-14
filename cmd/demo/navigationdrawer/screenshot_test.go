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

// TestScreenshot renders the NavigationDrawer UI and saves a screenshot.
func TestScreenshot(t *testing.T) {
	// 1. Setup Context
	var ops op.Ops
	gtx := layout.Context{
		Ops: &ops,
		Metric: unit.Metric{
			PxPerDp: 1.0,
			PxPerSp: 1.0,
		},
		Constraints: layout.Exact(image.Point{X: 1024, Y: 768}),
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
	img := screenshot.TakeScreenshot(1024, 768, callOp)

	// 6. Save to Artifacts
	// Using absolute path for artifacts directory
	artifactPath := "/home/jaco/.gemini/antigravity/brain/13deee75-9b8b-4e0b-bc18-325327911b78/navigation_drawer_screenshot.png"
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
