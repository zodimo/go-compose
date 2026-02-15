package main

import (
	"image"
	"testing"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/runtime"
	"github.com/zodimo/go-compose/state"
	"github.com/zodimo/go-compose/store"
	"github.com/zodimo/go-compose/theme"
)

func BenchmarkUI_Composition(b *testing.B) {
	store := store.NewPersistentState(map[string]state.MutableValue{})
	composer := compose.NewComposer(store)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UI()(composer)
	}
}

func BenchmarkUI_Layout(b *testing.B) {
	// Setup
	store := store.NewPersistentState(map[string]state.MutableValue{})
	composer := compose.NewComposer(store)
	rt := runtime.NewRuntime()
	themeManager := theme.GetThemeManager()

	// Prepare context
	var ops op.Ops

	// Create a dummy FrameEvent to initialize context
	size := image.Point{X: 1250, Y: 800}

	// Initialize theme
	// We need a context for theme init
	gtx := layout.Context{
		Ops:         &ops,
		Constraints: layout.Exact(size),
		Metric:      unit.Metric{PxPerDp: 1, PxPerSp: 1},
		Now:         time.Now(),
	}
	// Note: Material3ThemeInit expects a layout.Context and returns one.
	// It basically sets up the theme in the context.
	// We should probably do this once if possible, but it might modify ops?
	// Actually theme init usually just reads from context or sets values in context variable?
	// Let's assume it's cheap or necessary part of layout.

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ops.Reset()
		gtx.Ops = &ops // Ensure ops are reset

		// In the main loop:
		// gtx = themeManager.Material3ThemeInit(gtx)
		gtx = themeManager.Material3ThemeInit(gtx)

		callOp := rt.Run(gtx, composer, UI())
		callOp.Add(gtx.Ops)
	}
}
