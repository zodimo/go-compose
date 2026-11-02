//go:build !symbols_extra_fonts

package fonts

import "testing"

func TestSymbolTTF_ExtraFacesAbsentByDefault(t *testing.T) {
	// Without the symbols_extra_fonts tag the optional faces must report
	// unavailable, so the icon package falls back to outlined rather than
	// depending on a binary that is not there.
	for _, style := range []SymbolFace{SymbolFaceRounded, SymbolFaceSharp} {
		if data, ok := SymbolTTF(style); ok || data != nil {
			t.Errorf("style %d reported available without symbols_extra_fonts", style)
		}
	}
}
