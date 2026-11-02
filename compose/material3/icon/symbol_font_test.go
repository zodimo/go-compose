package icon

import (
	"errors"
	"sync"
	"testing"
)

// newTestSymbolFonts builds a registry backed by a stub provider so tests never
// depend on which font assets a given build embeds.
func newTestSymbolFonts(available map[SymbolStyle][]byte) *symbolFonts {
	return &symbolFonts{
		provider: func(style SymbolStyle) ([]byte, bool) {
			data, ok := available[style]
			return data, ok
		},
		cache: make(map[SymbolStyle]*symbolFace),
		errs:  make(map[SymbolStyle]error),
	}
}

func TestSymbolFonts_ParsesOncePerStyle(t *testing.T) {
	shared := sharedSymbolFonts()

	first, err := shared.face(SymbolOutlined)
	if err != nil {
		t.Fatalf("face(SymbolOutlined): %v", err)
	}
	second, err := shared.face(SymbolOutlined)
	if err != nil {
		t.Fatalf("face(SymbolOutlined) second call: %v", err)
	}

	if first != second {
		t.Fatalf("registry re-parsed the font: got distinct pointers %p and %p", first, second)
	}
	if first.font == nil {
		t.Fatal("parsed face has a nil font")
	}
}

func TestSymbolFonts_UpemIsStableForBundledFont(t *testing.T) {
	f, err := sharedSymbolFonts().face(SymbolOutlined)
	if err != nil {
		t.Fatalf("face(SymbolOutlined): %v", err)
	}

	// The bundled MaterialSymbolsOutlined.ttf declares 960 units-per-em. This
	// pins the scale factor every outline conversion depends on: if a different
	// font file is embedded, geometry scaling silently changes with it.
	const wantUpem = 960
	if f.upem != wantUpem {
		t.Fatalf("upem = %d, want %d (bundled font changed?)", f.upem, wantUpem)
	}
}

func TestSymbolFonts_MissingAssetIsCachedAsError(t *testing.T) {
	reg := newTestSymbolFonts(map[SymbolStyle][]byte{})

	if reg.available(SymbolRounded) {
		t.Fatal("available(SymbolRounded) = true with an empty provider")
	}
	provided := 0
	reg.provider = func(style SymbolStyle) ([]byte, bool) {
		provided++
		return nil, false
	}

	// Two lookups must not both hit the provider: the error is remembered.
	if _, err := reg.face(SymbolRounded); err == nil {
		t.Fatal("face(SymbolRounded) succeeded with an empty provider")
	}
	if _, err := reg.face(SymbolRounded); err == nil {
		t.Fatal("face(SymbolRounded) succeeded on the second call")
	}
	if provided != 0 {
		t.Fatalf("provider called %d times after the error was cached, want 0", provided)
	}
}

func TestSymbolFonts_ResolveStyleFallsBackToOutlined(t *testing.T) {
	outlinedOnly := map[SymbolStyle][]byte{SymbolOutlined: mustFontBytes(t)}
	reg := newTestSymbolFonts(outlinedOnly)

	got, fellBack := reg.resolveStyle(SymbolRounded)
	if got != SymbolOutlined {
		t.Fatalf("resolveStyle(SymbolRounded) = %v, want SymbolOutlined", got)
	}
	if !fellBack {
		t.Fatal("resolveStyle(SymbolRounded) reported no fallback")
	}

	got, fellBack = reg.resolveStyle(SymbolOutlined)
	if got != SymbolOutlined || fellBack {
		t.Fatalf("resolveStyle(SymbolOutlined) = (%v, %v), want (SymbolOutlined, false)", got, fellBack)
	}
}

// errStyleMissingForTest stands in for a style whose asset is not embedded.
var errStyleMissingForTest = errors.New("test: style asset missing")

// TestSymbolFonts_StyleErrDoesNotExposeTheMap forces a style to fail and asserts
// the error is readable through the lock-taking accessor. Reading errs directly
// (as the fallback path used to) races with face() writing it under s.mu.
func TestSymbolFonts_StyleErrTakesTheLock(t *testing.T) {
	reg := newTestSymbolFonts(map[SymbolStyle][]byte{SymbolOutlined: mustFontBytes(t)})

	if _, err := reg.face(SymbolSharp); err == nil {
		t.Fatal("face(SymbolSharp) succeeded with no embedded asset")
	}
	if err := reg.styleErr(SymbolSharp); err == nil {
		t.Fatal("styleErr reported no error for a style that failed to parse")
	}
	if err := reg.styleErr(SymbolOutlined); err != nil {
		t.Fatalf("styleErr reported %v for a style that parsed", err)
	}
}

// TestSymbolFonts_ConcurrentFallbackRendering hammers the fallback path from
// several goroutines while the registry parses other styles. Run under -race:
// the unsynchronised errs read the renderer used to perform fails here.
func TestSymbolFonts_ConcurrentFallbackRendering(t *testing.T) {
	previous := globalSymbolFonts
	globalSymbolFonts = newTestSymbolFonts(map[SymbolStyle][]byte{SymbolOutlined: mustFontBytes(t)})
	t.Cleanup(func() { globalSymbolFonts = previous })

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				// Missing style: fallback path, which reports the error.
				if _, err := buildSymbolRender(SymbolRounded, SymbolHome, 24, resolveAxes(DefaultIconOptions(), 24)); err != nil {
					t.Errorf("buildSymbolRender: %v", err)
					return
				}
				// Concurrent writer of the same map.
				_, _ = sharedSymbolFonts().face(SymbolSharp)
			}
		}()
	}
	wg.Wait()
}
