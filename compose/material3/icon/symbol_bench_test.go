package icon

import (
	"image/color"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
)

// BenchmarkSymbolRender_Cold measures resolving and converting an icon with no
// cache, i.e. the cost paid the first time a configuration appears.
func BenchmarkSymbolRender_Cold(b *testing.B) {
	axes := resolveAxes(DefaultIconOptions(), 24)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// A fresh cache per iteration forces a real extraction.
		globalGeometryCache = newGeometryCache(1)
		if _, err := buildSymbolRender(SymbolOutlined, SymbolHome, 24, axes); err != nil {
			b.Fatalf("buildSymbolRender: %v", err)
		}
	}
}

// BenchmarkSymbolRender_Cached measures the steady-state cost, where geometry is
// looked up rather than rebuilt.
func BenchmarkSymbolRender_Cached(b *testing.B) {
	axes := resolveAxes(DefaultIconOptions(), 24)
	globalGeometryCache = newGeometryCache(16)

	// Warm the cache.
	if _, err := buildSymbolRender(SymbolOutlined, SymbolHome, 24, axes); err != nil {
		b.Fatalf("warm: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := buildSymbolRender(SymbolOutlined, SymbolHome, 24, axes); err != nil {
			b.Fatalf("buildSymbolRender: %v", err)
		}
	}
}

// BenchmarkSymbolDraw_CachedOps measures replaying recorded paint ops, which is
// what an already-composed icon costs per frame.
func BenchmarkSymbolDraw_CachedOps(b *testing.B) {
	render, err := buildSymbolRender(SymbolOutlined, SymbolHome, 24, resolveAxes(DefaultIconOptions(), 24))
	if err != nil {
		b.Fatalf("buildSymbolRender: %v", err)
	}

	cache := newOpsCache(16)
	key := paintKey{color: color.NRGBA{A: 255}}
	box := render.size

	var warm op.Ops
	warmGtx := layout.Context{Ops: &warm, Constraints: layout.Exact(box)}
	cache.record(warmGtx, key, func(*op.Ops) layout.Dimensions {
		return render.draw(&warm, fill{color: key.color}, 0)
	})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var ops op.Ops
		gtx := layout.Context{Ops: &ops, Constraints: layout.Exact(box)}
		if _, ok := cache.layout(gtx, key); !ok {
			b.Fatal("expected a cache hit")
		}
	}
}
