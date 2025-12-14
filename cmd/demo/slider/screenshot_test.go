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

func TestSliderScreenshot(t *testing.T) {
	// 1. Setup Context
	var ops op.Ops
	gtx := layout.Context{
		Ops: &ops,
		Metric: unit.Metric{
			PxPerDp: 1.0,
			PxPerSp: 1.0,
		},
		Constraints: layout.Exact(image.Point{X: 400, Y: 800}),
		Now:         time.Now(),
		Locale:      system.Locale{Language: "en", Direction: system.LTR},
	}

	// 2. Init Dependencies
	themeManager := theme.GetThemeManager()
	gtx = themeManager.Material3ThemeInit(gtx)

	store := store.NewPersistentState(map[string]state.MutableValue{})
	composer := compose.NewComposer(store)

	// 3. Compose UI
	rootComposable := UI()
	rootComposable(composer)
	layoutNode := composer.Build()

	// 4. Run Layout
	rt := runtime.NewRuntime()
	callOp := rt.Run(gtx, layoutNode)

	// 5. Take Screenshot
	// Note: TakeScreenshot signature might vary, using the one found in appbar test
	// func TakeScreenshot(width, height int, callOp op.CallOp) image.Image
	img := screenshot.TakeScreenshot(400, 800, callOp)

	// 6. Save to Artifacts
	// Using absolute path as per instructions
	artifactPath := "/home/jaco/.gemini/antigravity/brain/0d7a6678-1671-449e-8b83-91e07a5027b9/slider_screenshot.png"
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
