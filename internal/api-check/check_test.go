package apicheck

import (
	"testing"
)

func TestIsSeam(t *testing.T) {
	cases := []struct {
		pkg  string
		want bool
	}{
		{"github.com/zodimo/go-compose/internal/layoutnode", true},
		{"github.com/zodimo/go-compose/internal/layoutnode/models", true},
		{"github.com/zodimo/go-compose/internal/render", true},
		{"github.com/zodimo/go-compose/internal/render/gio", true},
		{"github.com/zodimo/go-compose/internal/render/software", true},
		{"github.com/zodimo/go-compose/internal/composer", false},
		{"github.com/zodimo/go-compose/compose", false},
		{"github.com/zodimo/go-compose/cmd/demo", false},
	}
	for _, c := range cases {
		if got := IsSeam(c.pkg); got != c.want {
			t.Errorf("IsSeam(%q) = %v, want %v", c.pkg, got, c.want)
		}
	}
}

func TestIsPublicTree(t *testing.T) {
	cases := []struct {
		pkg  string
		want bool
	}{
		{"github.com/zodimo/go-compose/compose", true},
		{"github.com/zodimo/go-compose/compose/ui/unit", true},
		{"github.com/zodimo/go-compose/modifiers/padding", true},
		{"github.com/zodimo/go-compose/theme", true},
		{"github.com/zodimo/go-compose/runtime", true},
		{"github.com/zodimo/go-compose/pkg/api", true},
		{"github.com/zodimo/go-compose/pkg/x/fileexplorer", true},
		{"github.com/zodimo/go-compose/internal/composer", false},
		{"github.com/zodimo/go-compose/cmd/demo", false},
		{"github.com/zodimo/go-compose/state", false},
	}
	for _, c := range cases {
		if got := isPublicTree(c.pkg); got != c.want {
			t.Errorf("isPublicTree(%q) = %v, want %v", c.pkg, got, c.want)
		}
	}
}

// TestSeamImportWhitelist asserts Rule B behavior in -all mode: seam packages
// that import engine packages are NOT flagged, while a non-seam package that
// imports an engine package IS flagged. In the pre-cleanup (RED) state the
// box package still imports gioui.org/layout directly.
func TestSeamImportWhitelist(t *testing.T) {
	violations, err := Check(Options{All: true})
	if err != nil {
		t.Fatalf("Check(All): %v", err)
	}

	importV := map[string]map[string]bool{}
	for _, v := range violations {
		if v.Kind == KindImport {
			if importV[v.Pkg] == nil {
				importV[v.Pkg] = map[string]bool{}
			}
			importV[v.Pkg][v.GioType] = true
		}
	}

	for _, seam := range []string{
		"github.com/zodimo/go-compose/internal/layoutnode",
		"github.com/zodimo/go-compose/internal/render",
	} {
		if len(importV[seam]) > 0 {
			t.Errorf("seam package %s was flagged for engine imports: %v", seam, importV[seam])
		}
	}

	if !importV["github.com/zodimo/go-compose/compose/foundation/layout/box"]["gioui.org/layout"] {
		t.Errorf("expected box package to be flagged for gioui.org/layout in RED state (import rule) — got %v", importV["github.com/zodimo/go-compose/compose/foundation/layout/box"])
	}
}
