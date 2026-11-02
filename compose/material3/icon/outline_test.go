package icon

import (
	"testing"

	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/font/opentype"
)

// homeOutline returns the default-instance outline for the "home" symbol.
func homeOutline(t *testing.T) (font.GlyphOutline, uint16) {
	t.Helper()

	face, err := sharedSymbolFonts().face(SymbolOutlined)
	if err != nil {
		t.Fatalf("load outlined face: %v", err)
	}

	glyph, err := resolveGlyph(font.NewFace(face.font), SymbolOutlined, SymbolHome, 24)
	if err != nil {
		t.Fatalf("resolve home: %v", err)
	}
	outline, ok := glyph.outline()
	if !ok {
		t.Fatal("home glyph has no outline")
	}
	return outline, face.upem
}

func TestResolveGlyph_ValidNamesResolveToOneGlyph(t *testing.T) {
	face, err := sharedSymbolFonts().face(SymbolOutlined)
	if err != nil {
		t.Fatalf("load outlined face: %v", err)
	}

	for _, name := range []SymbolName{SymbolHome, SymbolSearch, SymbolSettings, SymbolStyleIcon} {
		glyph, err := resolveGlyph(font.NewFace(face.font), SymbolOutlined, name, 24)
		if err != nil {
			t.Errorf("resolve %q: %v", name, err)
			continue
		}
		if _, ok := glyph.outline(); !ok {
			t.Errorf("%q resolved to gid %d with no outline", name, glyph.glyphID)
		}
	}
}

func TestResolveGlyph_UnknownNameIsReported(t *testing.T) {
	face, err := sharedSymbolFonts().face(SymbolOutlined)
	if err != nil {
		t.Fatalf("load outlined face: %v", err)
	}

	// An unknown name shapes to one glyph per letter. Rendering those would
	// silently draw text, so resolution must fail instead.
	if _, err := resolveGlyph(font.NewFace(face.font), SymbolOutlined, "not_a_real_icon", 24); err == nil {
		t.Fatal("resolve of an unknown name succeeded; it must be reported as unresolved")
	}
}

func TestResolveGlyph_EmptyNameIsReported(t *testing.T) {
	face, err := sharedSymbolFonts().face(SymbolOutlined)
	if err != nil {
		t.Fatalf("load outlined face: %v", err)
	}
	if _, err := resolveGlyph(font.NewFace(face.font), SymbolOutlined, "", 24); err == nil {
		t.Fatal("resolve of an empty name succeeded")
	}
}

func TestOutlineToPath_HomeSegmentCount(t *testing.T) {
	outline, upem := homeOutline(t)

	// The bundled MaterialSymbolsOutlined.ttf renders "home" as 3 contours
	// totalling 30 segments. Pinning the count catches a font swap or an outline
	// extraction regression, either of which would silently change every icon.
	const wantSegments = 30
	if len(outline.Segments) != wantSegments {
		t.Fatalf("home outline has %d segments, want %d (bundled font changed?)",
			len(outline.Segments), wantSegments)
	}

	got := outlineToPath(outline.Segments, upem, 24)
	if got.empty() {
		t.Fatal("outlineToPath produced an empty path")
	}
	if got.commands != wantSegments {
		t.Fatalf("emitted %d path commands, want %d (one per segment)", got.commands, wantSegments)
	}
}

func TestOutlineToPath_EmptyInput(t *testing.T) {
	if !outlineToPath(nil, 960, 24).empty() {
		t.Fatal("outlineToPath(nil) produced a non-empty path")
	}
}

func TestOutlineToPath_ZeroUpemIsSafe(t *testing.T) {
	outline, _ := homeOutline(t)
	if !outlineToPath(outline.Segments, 0, 24).empty() {
		t.Fatal("outlineToPath with upem 0 produced a path; it must bail out instead")
	}
}

func TestOutlineToPath_InkBoundsScaleWithSize(t *testing.T) {
	outline, upem := homeOutline(t)

	// In a 960 upem the home ink spans x 160..800, y 120..840, i.e.
	// 0.667em wide and 0.75em tall, sitting entirely above the baseline.
	// At 24px that is 16px wide and 18px tall.
	small := outlineToPath(outline.Segments, upem, 24).ink
	if got, want := small.Dx(), 16; got != want {
		t.Errorf("ink width at 24px = %d, want %d", got, want)
	}
	if got, want := small.Dy(), 18; got != want {
		t.Errorf("ink height at 24px = %d, want %d", got, want)
	}
	if small.Max.Y > 0 {
		t.Errorf("ink extends below the baseline (Max.Y=%d); the Y flip is wrong", small.Max.Y)
	}

	// The ink box must scale linearly with the requested size.
	large := outlineToPath(outline.Segments, upem, 48).ink
	if got, want := large.Dx(), small.Dx()*2; got != want {
		t.Errorf("ink width at 48px = %d, want %d", got, want)
	}
	if got, want := large.Dy(), small.Dy()*2; got != want {
		t.Errorf("ink height at 48px = %d, want %d", got, want)
	}
}

func TestOutlineBounds_HomeMatchesFontUnits(t *testing.T) {
	outline, upem := homeOutline(t)
	minX, minY, maxX, maxY, ok := outlineBounds(outline.Segments)
	if !ok {
		t.Fatal("outlineBounds reported no bounds")
	}

	// Recorded from the bundled font, in font units.
	wantMinX, wantMaxX := float32(160), float32(800)
	wantMinY, wantMaxY := float32(120), float32(840)
	if minX != wantMinX || maxX != wantMaxX || minY != wantMinY || maxY != wantMaxY {
		t.Fatalf("home ink bounds = (%v,%v)-(%v,%v), want (%v,%v)-(%v,%v) in upem %d",
			minX, minY, maxX, maxY, wantMinX, wantMinY, wantMaxX, wantMaxY, upem)
	}
}

func TestOutlineBounds_EmptyInput(t *testing.T) {
	if _, _, _, _, ok := outlineBounds(nil); ok {
		t.Fatal("outlineBounds(nil) claimed to have bounds")
	}
}

// everySegmentOpOutline builds a synthetic outline exercising each segment op,
// including curves, so conversion coverage does not depend on which glyphs a
// given font happens to draw with lines only.
func everySegmentOpOutline() []font.Segment {
	return []font.Segment{
		{Op: opentype.SegmentOpMoveTo, Args: [3]font.SegmentPoint{{X: 0, Y: 0}}},
		{Op: opentype.SegmentOpLineTo, Args: [3]font.SegmentPoint{{X: 10, Y: 0}}},
		{Op: opentype.SegmentOpQuadTo, Args: [3]font.SegmentPoint{{X: 15, Y: 5}, {X: 10, Y: 10}}},
		{Op: opentype.SegmentOpCubeTo, Args: [3]font.SegmentPoint{{X: 7, Y: 12}, {X: 3, Y: 12}, {X: 0, Y: 10}}},
		{Op: opentype.SegmentOpMoveTo, Args: [3]font.SegmentPoint{{X: 2, Y: 2}}},
		{Op: opentype.SegmentOpLineTo, Args: [3]font.SegmentPoint{{X: 8, Y: 2}}},
	}
}

func TestOutlineToPath_AllSegmentOps(t *testing.T) {
	in := everySegmentOpOutline()
	got := outlineToPath(in, 1000, 48)
	if got.empty() {
		t.Fatal("outlineToPath dropped an outline containing every segment op")
	}
	if got.commands != len(in) {
		t.Fatalf("emitted %d commands for %d segments", got.commands, len(in))
	}
}
