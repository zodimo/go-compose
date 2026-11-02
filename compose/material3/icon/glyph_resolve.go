package icon

import (
	"fmt"
	"sync"

	"golang.org/x/image/math/fixed"

	"github.com/go-text/typesetting/di"
	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/language"
	"github.com/go-text/typesetting/shaping"
)

// glyphLookup resolves a symbol name to its glyph id and remembers the result.
//
// Names are resolved against the font's DEFAULT instance, never against the
// requested axis values. That is not an optimization: Material Symbols declines
// to substitute some ligatures at extreme axis values. At exactly FILL=1 the
// bundled font stops mapping "home", "search" and "settings" to their icon
// glyphs and shapes them as several letter glyphs instead, while FILL=0 and
// FILL=0.5 resolve normally. Resolving at the default instance and then applying
// axes only to outline extraction keeps an icon's identity independent of how it
// is styled.
type glyphLookup struct {
	mu      sync.Mutex
	entries map[string]font.GID
}

var (
	globalGlyphLookupMu sync.Mutex
	globalGlyphLookup   *glyphLookup
)

func sharedGlyphLookup() *glyphLookup {
	globalGlyphLookupMu.Lock()
	defer globalGlyphLookupMu.Unlock()
	if globalGlyphLookup == nil {
		globalGlyphLookup = &glyphLookup{entries: make(map[string]font.GID)}
	}
	return globalGlyphLookup
}

// resolvedGlyph is a glyph id together with the face its outline is extracted
// from. The face carries the requested variation settings.
type resolvedGlyph struct {
	face    *font.Face
	glyphID font.GID
}

// outline returns the glyph's vector outline under the face's variation
// settings, or false when the glyph has no outline data.
func (r resolvedGlyph) outline() (font.GlyphOutline, bool) {
	switch data := r.face.GlyphData(r.glyphID).(type) {
	case font.GlyphOutline:
		return data, len(data.Segments) > 0
	case font.GlyphSVG:
		// SVG-only glyphs carry a fallback outline; use it.
		return data.Outline, len(data.Outline.Segments) > 0
	default:
		return font.GlyphOutline{}, false
	}
}

// resolveGlyph finds the glyph for name and pairs it with face, which supplies
// the variation instance its outline is extracted at.
func resolveGlyph(face *font.Face, style SymbolStyle, name SymbolName, sizePx int) (resolvedGlyph, error) {
	gid, err := sharedGlyphLookup().glyphFor(style, name, sizePx)
	if err != nil {
		return resolvedGlyph{}, err
	}
	return resolvedGlyph{face: face, glyphID: gid}, nil
}

// glyphFor returns the default-instance glyph id for name, shaping it on first
// use and caching the result.
func (l *glyphLookup) glyphFor(style SymbolStyle, name SymbolName, sizePx int) (font.GID, error) {
	key := style.String() + "\x00" + string(name)

	l.mu.Lock()
	if gid, ok := l.entries[key]; ok {
		l.mu.Unlock()
		return gid, nil
	}
	l.mu.Unlock()

	sf, err := sharedSymbolFonts().face(style)
	if err != nil {
		return 0, err
	}

	gid, err := shapeSymbolName(font.NewFace(sf.font), name, sizePx)
	if err != nil {
		return 0, err
	}

	l.mu.Lock()
	l.entries[key] = gid
	l.mu.Unlock()

	return gid, nil
}

// shapeSymbolName shapes a symbol name into the single glyph it represents.
//
// Symbol names resolve through the font's ligature table, not a codepoint
// lookup: SymbolName holds the icon name ("home"), and the font maps that name
// to a glyph. A name that is not a Material icon shapes to one glyph per letter,
// so anything other than exactly one glyph is reported as unresolved rather than
// being drawn as text.
func shapeSymbolName(face *font.Face, name SymbolName, sizePx int) (font.GID, error) {
	runes := []rune(string(name))
	if len(runes) == 0 {
		return 0, fmt.Errorf("icon: empty symbol name")
	}

	input := shaping.Input{
		Text:      runes,
		RunStart:  0,
		RunEnd:    len(runes),
		Direction: di.DirectionLTR,
		Face:      face,
		// Size affects only reported advances and bounds, not glyph selection,
		// but it is set so shape metrics stay meaningful.
		Size:     fixed.I(sizePx),
		Script:   language.Latin,
		Language: language.NewLanguage("en"),
	}

	out := (&shaping.HarfbuzzShaper{}).Shape(input)
	if len(out.Glyphs) != 1 {
		return 0, fmt.Errorf(
			"icon: %q is not a Material Symbol (resolved to %d glyphs, want 1)",
			name, len(out.Glyphs))
	}

	return out.Glyphs[0].GlyphID, nil
}
