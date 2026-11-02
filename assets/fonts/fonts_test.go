package fonts

import "testing"

func TestSymbolTTF_OutlinedAlwaysAvailable(t *testing.T) {
	data, ok := SymbolTTF(SymbolFaceOutlined)
	if !ok {
		t.Fatal("outlined Material Symbols face is not embedded")
	}
	if len(data) == 0 {
		t.Fatal("outlined face is empty")
	}
}

func TestSymbolTTF_UnknownStyleIsUnavailable(t *testing.T) {
	if _, ok := SymbolTTF(SymbolFace(99)); ok {
		t.Fatal("an unknown style reported as available")
	}
}

func TestCollection_DoesNotIncludeExtraSymbolFaces(t *testing.T) {
	// Collection feeds the text shaper and must not grow with the optional icon
	// faces: it returns the go fonts, emoji, and the outlined symbols face only.
	faces := Collection()
	if len(faces) == 0 {
		t.Fatal("Collection returned no faces")
	}

	extra := 0
	for _, f := range faces {
		switch string(f.Font.Typeface) {
		case "Material Symbols Rounded", "Material Symbols Sharp":
			extra++
		}
	}
	if extra != 0 {
		t.Fatalf("Collection includes %d optional symbol faces; it must stay outlined-only", extra)
	}
}
