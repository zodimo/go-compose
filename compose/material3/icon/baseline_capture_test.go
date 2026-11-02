package icon_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	m3icon "github.com/zodimo/go-compose/compose/material3/icon"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/screenshot"
	"github.com/zodimo/go-compose/pkg/api"
)

// baselineCapture is the measurement recorded for one icon configuration.
type baselineCapture struct {
	Symbol     string `json:"symbol"`
	SizeDp     int    `json:"sizeDp"`
	WidthPx    int    `json:"widthPx"`
	HeightPx   int    `json:"heightPx"`
	BaselinePx int    `json:"baselinePx"`
}

// defaultBaselineSymbols is the representative set rendered with no options.
// It deliberately mixes shapes: a solid glyph (home), a glyph with interior
// holes (search), and a glyph with thin strokes (settings), so a regression in
// contour closing or hole winding shows up.
var defaultBaselineSymbols = []struct {
	name string
	sym  m3icon.SymbolName
}{
	{"home", m3icon.SymbolHome},
	{"search", m3icon.SymbolSearch},
	{"settings", m3icon.SymbolSettings},
}

var defaultBaselineSizes = []unit.Dp{24, 48}

// fixtureFor stacks one icon per row at a fixed size.
func fixtureFor(sym m3icon.SymbolName, size unit.Dp) api.Composable {
	return func(c api.Composer) api.Composer {
		return column.Column(
			c.Sequence(
				m3icon.Icon(sym, m3icon.WithSize(size)),
			),
			column.WithAlignment(column.Start),
		)(c)
	}
}

// TestCaptureDefaultRenderBaseline writes the measurement of the default render
// path to testdata. It is a tool, not an assertion: run it against the
// pre-change code to produce the file that TestDefaultRenderMatchesBaseline
// compares against afterwards.
func TestCaptureDefaultRenderBaseline(t *testing.T) {
	if os.Getenv("CAPTURE_ICON_BASELINE") == "" {
		t.Skip("set CAPTURE_ICON_BASELINE=1 to record the default-render baseline")
	}

	out := make([]baselineCapture, 0, len(defaultBaselineSymbols)*len(defaultBaselineSizes))
	for _, entry := range defaultBaselineSymbols {
		for _, size := range defaultBaselineSizes {
			got := screenshot.Render(fixtureFor(entry.sym, size), 128, 128)
			out = append(out, baselineCapture{
				Symbol:     entry.name,
				SizeDp:     int(size),
				WidthPx:    got.Size.X,
				HeightPx:   got.Size.Y,
				BaselinePx: got.Baseline,
			})
			t.Logf("%s @ %ddp -> %v baseline=%d", entry.name, int(size), got.Size, got.Baseline)
		}
	}

	path := filepath.Join("testdata", "default_render_baseline.json")
	if err := os.MkdirAll("testdata", 0o755); err != nil {
		t.Fatalf("create testdata: %v", err)
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		t.Fatalf("marshal baseline: %v", err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
	t.Logf("wrote %s (%d entries)", path, len(out))
}

// TestDefaultRenderMatchesBaseline is the regression guard for the switch from
// glyph-text rendering to vector rendering.
//
// What it asserts, and why not more: the committed baseline was captured from
// the text path, which Gio rasterises through its own glyph cache and hinting.
// The vector path fills the same outline through clip.Outline, so the two cannot
// be pixel-identical; sub-pixel edge coverage differs by design. What MUST stay
// identical is everything layout depends on: the box the icon reports and the
// baseline it reports, at every tested size. A one-pixel drift here would shift
// every icon in every existing layout, so it is treated as a failure; a
// rasterised coverage difference is not.
func TestDefaultRenderMatchesBaseline(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "default_render_baseline.json"))
	if err != nil {
		t.Fatalf("read baseline: %v (regenerate with CAPTURE_ICON_BASELINE=1 "+
			"against a known-good build)", err)
	}

	var want []baselineCapture
	if err := json.Unmarshal(data, &want); err != nil {
		t.Fatalf("parse baseline: %v", err)
	}
	if len(want) == 0 {
		t.Fatal("baseline is empty")
	}

	type key struct {
		symbol string
		sizeDp int
	}
	byKey := make(map[key]baselineCapture, len(want))
	for _, entry := range want {
		byKey[key{entry.Symbol, entry.SizeDp}] = entry
	}

	for _, entry := range defaultBaselineSymbols {
		for _, size := range defaultBaselineSizes {
			wantEntry, ok := byKey[key{entry.name, int(size)}]
			if !ok {
				t.Fatalf("baseline has no entry for %s @ %ddp", entry.name, int(size))
			}

			got := screenshot.Render(fixtureFor(entry.sym, size), 128, 128)

			if got.Size.X != wantEntry.WidthPx || got.Size.Y != wantEntry.HeightPx {
				t.Errorf("%s @ %ddp: box = %dx%d, baseline was %dx%d",
					entry.name, int(size),
					got.Size.X, got.Size.Y, wantEntry.WidthPx, wantEntry.HeightPx)
			}
			if got.Baseline != wantEntry.BaselinePx {
				t.Errorf("%s @ %ddp: baseline = %d, baseline was %d",
					entry.name, int(size), got.Baseline, wantEntry.BaselinePx)
			}
		}
	}
}
