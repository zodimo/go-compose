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

func TestScreenshot(t *testing.T) {
	// Setup
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

	themeManager := theme.GetThemeManager()
	gtx = themeManager.Material3ThemeInit(gtx)

	store := store.NewPersistentState(map[string]state.MutableValue{})
	rt := runtime.NewRuntime()

	composer := compose.NewComposer(store)
	// Build the UI tree
	rootComposer := UI()(composer)
	layoutNode := rootComposer.Build()

	// Run the runtime to get draw operations
	callOp := rt.Run(gtx, layoutNode)

	// Capture screenshot
	img := screenshot.TakeScreenshot(1024, 768, callOp)

	// Save to artifact directory
	artifactPath := "/home/jaco/.gemini/antigravity/brain/1bffb78c-de3c-4bb7-bca3-655114c36fe2/menu_screenshot.png"
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
