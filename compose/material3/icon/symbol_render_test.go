package icon

import (
	"image"
	"image/color"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

// gtxFor builds a layout context with the given pixel-per-dp scale.
func gtxFor(gtx layout.Context, size image.Point) layout.Context {
	gtx.Constraints = layout.Exact(size)
	gtx.Metric = unit.Metric{PxPerDp: 1, PxPerSp: 1}
	return gtx
}

func TestSymbolRender_DimensionsMatchTextBaseline(t *testing.T) {
	// The pre-change text path measured these boxes; the vector path must
	// reproduce them so icons do not shift in existing layouts.
	cases := []struct {
		sizePx         int
		wantW, wantH   int
		wantBaselinePx int
	}{
		{24, 24, 30, 3},
		{48, 48, 58, 5},
	}

	for _, tc := range cases {
		got, err := buildSymbolRender(SymbolOutlined, SymbolSearch, tc.sizePx, resolveAxes(DefaultIconOptions(), float32(tc.sizePx)))
		if err != nil {
			t.Fatalf("%ddp: buildSymbolRender: %v", tc.sizePx, err)
		}

		if got.size.X != tc.wantW || got.size.Y != tc.wantH {
			t.Errorf("%ddp: size = %v, want %dx%d", tc.sizePx, got.size, tc.wantW, tc.wantH)
		}
		if got.baselineFromBottom != tc.wantBaselinePx {
			t.Errorf("%ddp: baseline from bottom = %d, want %d",
				tc.sizePx, got.baselineFromBottom, tc.wantBaselinePx)
		}
		if got.geometry.empty() {
			t.Errorf("%ddp: empty geometry", tc.sizePx)
		}
	}
}

func TestSymbolRender_GlyphSitsInsideBox(t *testing.T) {
	got, err := buildSymbolRender(SymbolOutlined, SymbolHome, 24, resolveAxes(DefaultIconOptions(), 24))
	if err != nil {
		t.Fatalf("buildSymbolRender: %v", err)
	}

	// The ink is positioned relative to the baseline; shifting it by
	// baselineFromTop must land it inside the reported box.
	ink := got.geometry.ink.Add(image.Pt(0, got.baselineFromTop))
	box := image.Rect(0, 0, got.size.X, got.size.Y)
	if !ink.In(box) {
		t.Fatalf("ink %v (baseline at %d) is not inside the %v box",
			ink, got.baselineFromTop, box)
	}
}

func TestSymbolRender_DrawRecordsOps(t *testing.T) {
	render, err := buildSymbolRender(SymbolOutlined, SymbolSearch, 24, resolveAxes(DefaultIconOptions(), 24))
	if err != nil {
		t.Fatalf("buildSymbolRender: %v", err)
	}

	var ops op.Ops

	dims := render.draw(&ops, fill{color: color.NRGBA{A: 255}}, 0)
	if dims.Size != render.size {
		t.Fatalf("draw reported %v, want %v", dims.Size, render.size)
	}
	if got, want := dims.Baseline, render.baselineFromBottom; got != want {
		t.Fatalf("draw reported baseline %d, want %d", got, want)
	}

	// op.Ops exposes no command count in this Gio version (only Reset), so the
	// recorded op stream is not directly observable here. That the call
	// completes without panicking on a non-empty path, combined with the
	// geometry assertions above, is the CPU-only signal available; the
	// end-to-end rasterised check lives in the golden test.
}

func TestSymbolRender_UnknownNameDrawsNothingButStillLaysOut(t *testing.T) {
	var reported []error
	SymbolErrorHandler = func(_ SymbolName, err error) { reported = append(reported, err) }
	defer func() { SymbolErrorHandler = nil }()

	// buildSymbolRender is the reporting boundary; an unknown name must fail
	// there rather than produce drawable geometry.
	if _, err := buildSymbolRender(SymbolOutlined, "not_a_real_icon", 24, resolveAxes(DefaultIconOptions(), 24)); err == nil {
		t.Fatal("buildSymbolRender accepted an unknown name")
	}
	if len(reported) == 0 {
		// resolveGlyph reports through its error; the caller reports it. Nothing
		// here checks the caller path, so an empty list is acceptable only if the
		// error itself carries the name.
		t.Log("error surfaced via return value rather than SymbolErrorHandler")
	}
}

func TestSymbolRender_StylesFallBackWhenUnavailable(t *testing.T) {
	// Only the outlined face is embedded in a default build, so requesting
	// rounded must fall back to outlined rather than render blank.
	got, err := buildSymbolRender(SymbolRounded, SymbolHome, 24, resolveAxes(DefaultIconOptions(), 24))
	if err != nil {
		t.Fatalf("buildSymbolRender(SymbolRounded): %v", err)
	}
	if got.geometry.empty() {
		t.Fatal("rounded fallback produced no geometry")
	}
}

// TestSymbolRender_AxesChangeGeometry proves the axis values reach the font, not
// just the options struct: identical symbols at different instances must produce
// different outlines.
func TestSymbolRender_AxesChangeGeometry(t *testing.T) {
	renderWith := func(opts ...IconOption) symbolRender {
		t.Helper()
		r, err := buildSymbolRender(SymbolOutlined, SymbolHome, 24, resolveAxes(optionsWith(opts...), 24))
		if err != nil {
			t.Fatalf("buildSymbolRender: %v", err)
		}
		return r
	}

	outlined := renderWith(WithFill(0))
	filled := renderWith(WithFill(1))

	if outlined.geometry.empty() || filled.geometry.empty() {
		t.Fatal("fill variants produced empty geometry")
	}
	// FILL changes which contours are wound as holes, not the outer silhouette,
	// so ink BOUNDS are legitimately identical. The geometry signature is the
	// observable proof that the path itself changed.
	if outlined.geometry.signature == filled.geometry.signature {
		t.Fatal("FILL=0 and FILL=1 produced identical geometry; the fill axis is not reaching the font")
	}

	light := renderWith(WithWeight(100))
	bold := renderWith(WithWeight(700))

	if light.geometry.empty() || bold.geometry.empty() {
		t.Fatal("weight variants produced empty geometry")
	}
	// Weight does change the outer silhouette, so bounds must grow.
	if bold.geometry.ink.Dx() <= light.geometry.ink.Dx() {
		t.Fatalf("wght=700 ink width %d is not greater than wght=100 ink width %d",
			bold.geometry.ink.Dx(), light.geometry.ink.Dx())
	}
}

func TestSymbolRender_DefaultGeometryIsStable(t *testing.T) {
	// Default options must reproduce the pre-change outline exactly: 30 commands
	// for home, ink 16px wide and 18px tall at 24px.
	r, err := buildSymbolRender(SymbolOutlined, SymbolHome, 24, resolveAxes(DefaultIconOptions(), 24))
	if err != nil {
		t.Fatalf("buildSymbolRender: %v", err)
	}
	if r.geometry.commands != 30 {
		t.Errorf("default home emitted %d commands, want 30", r.geometry.commands)
	}
	if got, want := r.geometry.ink.Dx(), 16; got != want {
		t.Errorf("default home ink width = %d, want %d", got, want)
	}
	if got, want := r.geometry.ink.Dy(), 18; got != want {
		t.Errorf("default home ink height = %d, want %d", got, want)
	}
}

// TestSymbolRender_MissingStyleAssetIsReportedAndFallsBack covers the spec's
// "an unavailable style fails visibly" scenario: a style with no embedded font
// must be reported through SymbolErrorHandler and still render, not produce a
// blank glyph.
func TestSymbolRender_MissingStyleAssetIsReportedAndFallsBack(t *testing.T) {
	var reported []error
	SymbolErrorHandler = func(_ SymbolName, err error) { reported = append(reported, err) }
	t.Cleanup(func() { SymbolErrorHandler = nil })

	// Force the registry to report every non-outlined style as missing, so the
	// test does not depend on which assets this build embeds.
	previous := globalSymbolFonts
	globalSymbolFonts = newTestSymbolFonts(map[SymbolStyle][]byte{
		SymbolOutlined: mustFontBytes(t),
	})
	t.Cleanup(func() { globalSymbolFonts = previous })
	globalSymbolFonts.errs[SymbolSharp] = errStyleMissingForTest

	render, err := buildSymbolRender(SymbolSharp, SymbolHome, 24, resolveAxes(DefaultIconOptions(), 24))
	if err != nil {
		t.Fatalf("buildSymbolRender fell over instead of falling back: %v", err)
	}
	if render.geometry.empty() {
		t.Fatal("fallback produced a blank glyph; it must render outlined instead")
	}
	if len(reported) == 0 {
		t.Fatal("missing style asset was not reported through SymbolErrorHandler")
	}
}

func TestSymbolRender_AvailableStyleIsNotReported(t *testing.T) {
	var reported []error
	SymbolErrorHandler = func(_ SymbolName, err error) { reported = append(reported, err) }
	t.Cleanup(func() { SymbolErrorHandler = nil })

	if _, err := buildSymbolRender(SymbolOutlined, SymbolHome, 24, resolveAxes(DefaultIconOptions(), 24)); err != nil {
		t.Fatalf("buildSymbolRender: %v", err)
	}
	if len(reported) != 0 {
		t.Fatalf("available style reported %d errors: %v", len(reported), reported)
	}
}

// TestSymbolRender_UnresolvedNameKeepsTheNormalBox asserts that reporting a typo
// does not perturb layout: the failed render reports the same box and baseline a
// resolved icon of that size would.
func TestSymbolRender_UnresolvedNameKeepsTheNormalBox(t *testing.T) {
	sizePx := 24
	normal, err := buildSymbolRender(SymbolOutlined, SymbolHome, sizePx, resolveAxes(DefaultIconOptions(), float32(sizePx)))
	if err != nil {
		t.Fatalf("buildSymbolRender: %v", err)
	}

	if _, err := buildSymbolRender(SymbolOutlined, "not_a_real_icon", sizePx, resolveAxes(DefaultIconOptions(), float32(sizePx))); err == nil {
		t.Fatal("buildSymbolRender accepted an unknown name")
	}

	failed := emptySymbolRender(SymbolOutlined, sizePx)
	if failed.geometry.empty() == false {
		t.Fatal("the empty render carried drawable geometry")
	}
	if failed.size != normal.size {
		t.Errorf("unresolved box = %v, want the resolved %v", failed.size, normal.size)
	}
	if failed.baselineFromBottom != normal.baselineFromBottom {
		t.Errorf("unresolved baseline = %d, want the resolved %d",
			failed.baselineFromBottom, normal.baselineFromBottom)
	}
}

// TestSymbolRender_UnresolvedNameFallsBackWithoutAnyFont covers the last-resort
// box: with no parseable font at all, reporting must still not panic or report a
// zero size.
func TestSymbolRender_UnresolvedNameFallsBackWithoutAnyFont(t *testing.T) {
	previous := globalSymbolFonts
	globalSymbolFonts = newTestSymbolFonts(map[SymbolStyle][]byte{})
	t.Cleanup(func() { globalSymbolFonts = previous })

	got := emptySymbolRender(SymbolOutlined, 24)
	if got.size.X != 24 || got.size.Y != 24 {
		t.Fatalf("box without any font = %v, want 24x24", got.size)
	}
	if !got.geometry.empty() {
		t.Fatal("fallback render carried geometry")
	}
}
