package icon

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/zodimo/go-compose/assets/fonts"

	"github.com/go-text/typesetting/font"
)

// symbolAxis tags, in the order the variable font declares them.
const (
	axisFill = "FILL"
	axisGrad = "GRAD"
	axisOpsz = "opsz"
	axisWght = "wght"
)

// symbolFace bundles the parsed variable typeface for one style with the
// metadata the renderer needs to place its glyphs.
type symbolFace struct {
	style SymbolStyle
	font  *font.Font
	// upem is the font's units-per-em; outline coordinates are expressed in
	// these units and must be scaled by sizePx/upem to reach device pixels.
	upem uint16
}

// symbolFaceProvider supplies raw TTF bytes for a style. It is a seam so tests
// can substitute an asset without touching the embedded fonts.
type symbolFaceProvider func(SymbolStyle) ([]byte, bool)

// defaultSymbolFaceProvider reads the variable TTF bytes embedded in
// assets/fonts. Only the outlined face is guaranteed to be present; rounded and
// sharp are optional (see assets/fonts.SymbolTTF).
func defaultSymbolFaceProvider(style SymbolStyle) ([]byte, bool) {
	return fonts.SymbolTTF(fonts.SymbolFace(style))
}

// symbolFonts caches one parsed typeface per style for the lifetime of the
// process. Parsing a ~10MB variable TTF is expensive, so it happens at most
// once per style and the result is shared by every composer.
type symbolFonts struct {
	mu       sync.Mutex
	provider symbolFaceProvider
	cache    map[SymbolStyle]*symbolFace
	errs     map[SymbolStyle]error
}

var (
	globalSymbolFontsMu sync.Mutex
	globalSymbolFonts   *symbolFonts
)

// sharedSymbolFonts returns the process-wide typeface registry.
func sharedSymbolFonts() *symbolFonts {
	globalSymbolFontsMu.Lock()
	defer globalSymbolFontsMu.Unlock()
	if globalSymbolFonts == nil {
		globalSymbolFonts = &symbolFonts{
			provider: defaultSymbolFaceProvider,
			cache:    make(map[SymbolStyle]*symbolFace),
			errs:     make(map[SymbolStyle]error),
		}
	}
	return globalSymbolFonts
}

// face returns the parsed typeface for style, parsing it on first use. The
// returned pointer is stable across calls: a second call never re-parses.
func (s *symbolFonts) face(style SymbolStyle) (*symbolFace, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if f, ok := s.cache[style]; ok {
		return f, nil
	}
	if err, ok := s.errs[style]; ok {
		return nil, err
	}

	data, ok := s.provider(style)
	if !ok || len(data) == 0 {
		err := fmt.Errorf("icon: no Material Symbols font available for style %q", style)
		s.errs[style] = err
		return nil, err
	}

	parsed, err := font.ParseTTF(bytes.NewReader(data))
	if err != nil {
		err = fmt.Errorf("icon: parsing %q Material Symbols font: %w", style, err)
		s.errs[style] = err
		return nil, err
	}

	f := &symbolFace{
		style: style,
		font:  parsed.Font,
		upem:  parsed.Font.Upem(),
	}
	s.cache[style] = f
	return f, nil
}

// available reports whether style can be parsed, without forcing a parse of a
// style that is already known to be missing.
func (s *symbolFonts) available(style SymbolStyle) bool {
	_, err := s.face(style)
	return err == nil
}

// styleErr returns the parsing error remembered for style, or nil. It takes the
// registry lock: face() writes errs while rendering may read it from another
// goroutine, so the map must never be touched outside the mutex.
func (s *symbolFonts) styleErr(style SymbolStyle) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.errs[style]
}

// resolveStyle returns the style to render with and, when it differs from the
// requested style, the requested style so the caller can report the fallback.
// The boolean reports whether a fallback happened. An unavailable requested
// style never wins: it degrades to SymbolOutlined.
func (s *symbolFonts) resolveStyle(requested SymbolStyle) (render SymbolStyle, fellBack bool) {
	if _, err := s.face(requested); err == nil {
		return requested, false
	}
	return SymbolOutlined, true
}
