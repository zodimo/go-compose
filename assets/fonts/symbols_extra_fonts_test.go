//go:build symbols_extra_fonts

package fonts

import "testing"

func TestSymbolTTF_ExtraFacesEmbedded(t *testing.T) {
	for _, style := range []SymbolFace{SymbolFaceRounded, SymbolFaceSharp} {
		data, ok := SymbolTTF(style)
		if !ok {
			t.Errorf("style %d is not embedded under symbols_extra_fonts", style)
			continue
		}
		if len(data) == 0 {
			t.Errorf("style %d embedded zero bytes", style)
		}
	}
}

func TestSymbolCollection_IncludesRequestedExtraStyles(t *testing.T) {
	faces := SymbolCollection(SymbolFaceRounded, SymbolFaceSharp)
	if len(faces) == 0 {
		t.Fatal("SymbolCollection returned no faces for rounded and sharp")
	}

	got := map[string]bool{}
	for _, f := range faces {
		got[string(f.Font.Typeface)] = true
	}
	for _, want := range []string{"Material Symbols Rounded", "Material Symbols Sharp"} {
		if !got[want] {
			t.Errorf("SymbolCollection is missing %q; got %v", want, got)
		}
	}
}

func TestSymbolCollection_OutlinedIsOptIn(t *testing.T) {
	for _, f := range SymbolCollection(SymbolFaceRounded) {
		if string(f.Font.Typeface) == "Material Symbols Outlined" {
			t.Fatal("outlined face appeared without being requested")
		}
	}
}
