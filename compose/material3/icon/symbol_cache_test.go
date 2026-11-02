package icon

import (
	"image/color"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
)

// freshGeometryCache installs a cache with the given limit for one test and
// restores the previous one afterwards, so cache assertions cannot leak between
// tests.
func freshGeometryCache(t *testing.T, limit int) *geometryCache {
	t.Helper()
	previous := globalGeometryCache
	globalGeometryCache = newGeometryCache(limit)
	t.Cleanup(func() { globalGeometryCache = previous })
	return globalGeometryCache
}

func freshOpsCache(t *testing.T, limit int) *opsCache {
	t.Helper()
	previous := globalOpsCache
	globalOpsCache = newOpsCache(limit)
	t.Cleanup(func() { globalOpsCache = previous })
	return globalOpsCache
}

func TestGeometryCache_RepeatedLayoutConvertsOnce(t *testing.T) {
	cache := freshGeometryCache(t, 16)
	axes := resolveAxes(DefaultIconOptions(), 24)

	for i := 0; i < 5; i++ {
		if _, err := buildSymbolRender(SymbolOutlined, SymbolHome, 24, axes); err != nil {
			t.Fatalf("buildSymbolRender: %v", err)
		}
	}

	if got, want := cache.conversions(), 1; got != want {
		t.Fatalf("converted %d times over 5 identical layouts, want %d", got, want)
	}
	if got, want := cache.len(), 1; got != want {
		t.Fatalf("cache holds %d entries, want %d", got, want)
	}
}

func TestGeometryCache_DistinctConfigurationsAreSeparate(t *testing.T) {
	cache := freshGeometryCache(t, 16)

	light := resolveAxes(optionsWith(WithWeight(100)), 24)
	bold := resolveAxes(optionsWith(WithWeight(700)), 24)

	r1, err := buildSymbolRender(SymbolOutlined, SymbolHome, 24, light)
	if err != nil {
		t.Fatalf("light: %v", err)
	}
	r2, err := buildSymbolRender(SymbolOutlined, SymbolHome, 24, bold)
	if err != nil {
		t.Fatalf("bold: %v", err)
	}

	if got, want := cache.conversions(), 2; got != want {
		t.Fatalf("converted %d times for two weights, want %d", got, want)
	}
	if r1.geometry.signature == r2.geometry.signature {
		t.Fatal("light and bold share geometry; the cache key ignored the weight axis")
	}
	if got, want := cache.len(), 2; got != want {
		t.Fatalf("cache holds %d entries, want %d", got, want)
	}
}

func TestGeometryCache_ClampedValuesShareAnEntry(t *testing.T) {
	cache := freshGeometryCache(t, 16)

	// Grade 500 clamps to 200, so it must resolve to the same entry as grade 200
	// rather than consuming a second slot.
	clamped := resolveAxes(optionsWith(WithGrade(500)), 24)
	exact := resolveAxes(optionsWith(WithGrade(200)), 24)

	if _, err := buildSymbolRender(SymbolOutlined, SymbolHome, 24, clamped); err != nil {
		t.Fatalf("clamped: %v", err)
	}
	if _, err := buildSymbolRender(SymbolOutlined, SymbolHome, 24, exact); err != nil {
		t.Fatalf("exact: %v", err)
	}

	if got, want := cache.conversions(), 1; got != want {
		t.Fatalf("converted %d times for equivalent clamped values, want %d", got, want)
	}
}

func TestGeometryCache_EvictsLeastRecentlyUsed(t *testing.T) {
	cache := freshGeometryCache(t, 2)

	sizes := []int{20, 24, 28}
	for _, size := range sizes {
		if _, err := buildSymbolRender(SymbolOutlined, SymbolHome, size, resolveAxes(DefaultIconOptions(), size)); err != nil {
			t.Fatalf("size %d: %v", size, err)
		}
	}

	if got, want := cache.len(), 2; got != want {
		t.Fatalf("cache holds %d entries after exceeding the limit, want %d", got, want)
	}

	// The first size was evicted, so re-rendering it must convert again and still
	// produce correct geometry.
	before := cache.conversions()
	got, err := buildSymbolRender(SymbolOutlined, SymbolHome, sizes[0], resolveAxes(DefaultIconOptions(), sizes[0]))
	if err != nil {
		t.Fatalf("re-render evicted entry: %v", err)
	}
	if cache.conversions() != before+1 {
		t.Fatal("evicted entry was served from cache")
	}
	if got.geometry.empty() {
		t.Fatal("re-rendered evicted entry produced empty geometry")
	}
	if got.size.X != sizes[0] {
		t.Fatalf("re-rendered size = %v, want width %d", got.size, sizes[0])
	}
}

func TestGeometryCache_RepeatedAccessKeepsEntryAlive(t *testing.T) {
	cache := freshGeometryCache(t, 2)

	// Touch "a" twice so "b" becomes least-recently-used, then insert "c".
	renderAt := func(size int) {
		if _, err := buildSymbolRender(SymbolOutlined, SymbolHome, size, resolveAxes(DefaultIconOptions(), size)); err != nil {
			t.Fatalf("size %d: %v", size, err)
		}
	}
	renderAt(20) // a
	renderAt(24) // b
	renderAt(20) // a again -> b is now LRU
	renderAt(28) // c, evicts b

	// "a" must still be cached.
	before := cache.conversions()
	renderAt(20)
	if cache.conversions() != before {
		t.Fatal("recently used entry was evicted; the cache is not LRU")
	}
}

func TestOpsCache_PaintChangeDoesNotReconvertGeometry(t *testing.T) {
	geomCache := freshGeometryCache(t, 16)
	opsCache := freshOpsCache(t, 16)

	axes := resolveAxes(DefaultIconOptions(), 24)
	render, err := buildSymbolRender(SymbolOutlined, SymbolHome, 24, axes)
	if err != nil {
		t.Fatalf("buildSymbolRender: %v", err)
	}
	geomKey := geometryKey{symbol: SymbolHome, style: SymbolOutlined, axes: axes, sizePx: 24}

	redKey := paintKey{geom: geomKey, color: color.NRGBA{R: 255, A: 255}}
	greenKey := paintKey{geom: geomKey, color: color.NRGBA{G: 255, A: 255}}

	var ops1, ops2 op.Ops
	gtx1 := gtxFor(layout.Context{Ops: &ops1}, render.size)
	gtx2 := gtxFor(layout.Context{Ops: &ops2}, render.size)

	opsCache.record(gtx1, redKey, func(*op.Ops) layout.Dimensions {
		return render.draw(&ops1, fill{color: redKey.color}, 0)
	})
	opsCache.record(gtx2, greenKey, func(*op.Ops) layout.Dimensions {
		return render.draw(&ops2, fill{color: greenKey.color}, 0)
	})

	// Two paints, but the geometry was never rebuilt.
	if got, want := geomCache.conversions(), 1; got != want {
		t.Fatalf("geometry converted %d times after a paint change, want %d", got, want)
	}

	// And the second paint is now a cache hit.
	var ops3 op.Ops
	gtx3 := gtxFor(layout.Context{Ops: &ops3}, render.size)
	if _, ok := opsCache.layout(gtx3, redKey); !ok {
		t.Fatal("previously recorded paint was not a cache hit")
	}
}

func TestOpsCache_ServesRecordedDims(t *testing.T) {
	cache := freshOpsCache(t, 8)

	render, err := buildSymbolRender(SymbolOutlined, SymbolHome, 24, resolveAxes(DefaultIconOptions(), 24))
	if err != nil {
		t.Fatalf("buildSymbolRender: %v", err)
	}
	key := paintKey{color: color.NRGBA{A: 255}}

	var recordOps op.Ops
	gtx := gtxFor(layout.Context{Ops: &recordOps}, render.size)
	want := cache.record(gtx, key, func(*op.Ops) layout.Dimensions {
		return render.draw(&recordOps, fill{color: key.color}, 0)
	})

	var replayOps op.Ops
	replayGtx := gtxFor(layout.Context{Ops: &replayOps}, render.size)
	got, ok := cache.layout(replayGtx, key)
	if !ok {
		t.Fatal("expected a cache hit")
	}
	if got.Size != want.Size || got.Baseline != want.Baseline {
		t.Fatalf("replayed dims %v/%d, want %v/%d", got.Size, got.Baseline, want.Size, want.Baseline)
	}

	hits, _ := cache.stats()
	if hits == 0 {
		t.Fatal("hit counter did not advance")
	}
}

func TestOpsCache_EvictsLeastRecentlyUsed(t *testing.T) {
	cache := freshOpsCache(t, 2)

	render, err := buildSymbolRender(SymbolOutlined, SymbolHome, 24, resolveAxes(DefaultIconOptions(), 24))
	if err != nil {
		t.Fatalf("buildSymbolRender: %v", err)
	}

	record := func(c color.NRGBA) paintKey {
		key := paintKey{color: c}
		var ops op.Ops
		gtx := gtxFor(layout.Context{Ops: &ops}, render.size)
		cache.record(gtx, key, func(*op.Ops) layout.Dimensions {
			return render.draw(&ops, fill{color: c}, 0)
		})
		return key
	}

	first := record(color.NRGBA{R: 1, A: 255})
	record(color.NRGBA{G: 2, A: 255})
	record(color.NRGBA{B: 3, A: 255}) // evicts first

	if got, want := cache.len(), 2; got != want {
		t.Fatalf("ops cache holds %d entries, want %d", got, want)
	}

	var ops op.Ops
	gtx := gtxFor(layout.Context{Ops: &ops}, render.size)
	if _, ok := cache.layout(gtx, first); ok {
		t.Fatal("evicted paint was still served from cache")
	}
}
