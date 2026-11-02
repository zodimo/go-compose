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
