package font

import "testing"

func TestGenericFontFamily_Name(t *testing.T) {
	gff := NewGenericFontFamily("sans-serif", "FontFamily.SansSerif")
	if gff.Name() != "sans-serif" {
		t.Errorf("Expected 'sans-serif', got %s", gff.Name())
	}
}

func TestGenericFontFamily_String(t *testing.T) {
	gff := NewGenericFontFamily("serif", "FontFamily.Serif")
	if gff.String() != "FontFamily.Serif" {
		t.Errorf("Expected 'FontFamily.Serif', got %s", gff.String())
	}
}

func TestDefaultFontFamily_String(t *testing.T) {
	dff := &DefaultFontFamily{}
	if dff.String() != "FontFamily.Default" {
		t.Errorf("Expected 'FontFamily.Default', got %s", dff.String())
	}
}

func TestFontFamilyConstants(t *testing.T) {
	// Test that constants are properly initialized
	if FontFamilyDefault == nil {
		t.Error("FontFamilyDefault should not be nil")
	}
	if FontFamilySansSerif == nil {
		t.Error("FontFamilySansSerif should not be nil")
	}
	if FontFamilySerif == nil {
		t.Error("FontFamilySerif should not be nil")
	}
	if FontFamilyMonospace == nil {
		t.Error("FontFamilyMonospace should not be nil")
	}
	if FontFamilyCursive == nil {
		t.Error("FontFamilyCursive should not be nil")
	}

	// Check GenericFontFamily values
	sansSerif := FontFamilySansSerif.(*GenericFontFamily)
	if sansSerif.Name() != "sans-serif" {
		t.Errorf("SansSerif name should be 'sans-serif', got %s", sansSerif.Name())
	}
}

func TestFontListFontFamily_NewAndFonts(t *testing.T) {
	f := &mockFont{
		weight:          FontWeightNormal,
		style:           FontStyleNormal,
		loadingStrategy: FontLoadingStrategyBlocking,
	}

	family := NewFontListFontFamily([]Font{f})
	if len(family.Fonts) != 1 {
		t.Errorf("Expected 1 font, got %d", len(family.Fonts))
	}
}

func TestFontListFontFamily_Panic_Empty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for empty font list")
		}
	}()
	NewFontListFontFamily([]Font{})
}

func TestFontListFontFamily_Equals(t *testing.T) {
	f1 := &mockFont{weight: FontWeightNormal, style: FontStyleNormal, loadingStrategy: FontLoadingStrategyBlocking}
	f2 := &mockFont{weight: FontWeightBold, style: FontStyleItalic, loadingStrategy: FontLoadingStrategyAsync}

	family1 := NewFontListFontFamily([]Font{f1})
	family2 := NewFontListFontFamily([]Font{f1})
	family3 := NewFontListFontFamily([]Font{f2})
	family4 := NewFontListFontFamily([]Font{f1, f2})

	if !family1.Equals(family1) {
		t.Error("Family should equal itself")
	}

	if !family1.Equals(family2) {
		t.Error("Families with same fonts should be equal")
	}

	if family1.Equals(family3) {
		t.Error("Families with different fonts should not be equal")
	}

	if family1.Equals(family4) {
		t.Error("Families with different number of fonts should not be equal")
	}

	if family1.Equals(nil) {
		t.Error("Family should not equal nil")
	}
}

func TestFontListFontFamily_String(t *testing.T) {
	f := &mockFont{weight: FontWeightNormal, style: FontStyleNormal, loadingStrategy: FontLoadingStrategyBlocking}
	family := NewFontListFontFamily([]Font{f})
	str := family.String()
	if str == "" {
		t.Error("String should not be empty")
	}
}

// mockTypeface implements Typeface for testing
type mockTypeface struct {
	family FontFamily
}

func (m *mockTypeface) FontFamily() FontFamily {
	return m.family
}

func TestLoadedFontFamily(t *testing.T) {
	tf := &mockTypeface{family: FontFamilyDefault}
	family := NewLoadedFontFamily(tf)

	if family.Typeface != tf {
		t.Error("Typeface should match")
	}
}

func TestLoadedFontFamily_Equals(t *testing.T) {
	tf1 := &mockTypeface{family: FontFamilyDefault}
	tf2 := &mockTypeface{family: FontFamilySansSerif}

	family1 := NewLoadedFontFamily(tf1)
	family2 := NewLoadedFontFamily(tf1)
	family3 := NewLoadedFontFamily(tf2)

	if !family1.Equals(family1) {
		t.Error("Family should equal itself")
	}

	if !family1.Equals(family2) {
		t.Error("Families with same typeface should be equal")
	}

	if family1.Equals(family3) {
		t.Error("Families with different typefaces should not be equal")
	}

	if family1.Equals(nil) {
		t.Error("Family should not equal nil")
	}
}

func TestFontFamilyFromFonts(t *testing.T) {
	f := &mockFont{weight: FontWeightNormal, style: FontStyleNormal, loadingStrategy: FontLoadingStrategyBlocking}
	family := FontFamilyFromFonts(f)

	listFamily, ok := family.(*FontListFontFamily)
	if !ok {
		t.Fatal("FontFamilyFromFonts should return FontListFontFamily")
	}

	if len(listFamily.Fonts) != 1 {
		t.Errorf("Expected 1 font, got %d", len(listFamily.Fonts))
	}
}

func TestFontFamilyFromTypeface(t *testing.T) {
	tf := &mockTypeface{family: FontFamilyDefault}
	family := FontFamilyFromTypeface(tf)

	loadedFamily, ok := family.(*LoadedFontFamily)
	if !ok {
		t.Fatal("FontFamilyFromTypeface should return LoadedFontFamily")
	}

	if loadedFamily.Typeface != tf {
		t.Error("Typeface should match")
	}
}

func TestFontFamily_InterfaceImplementation(t *testing.T) {
	// Ensure all types implement FontFamily
	var _ FontFamily = &GenericFontFamily{}
	var _ FontFamily = &DefaultFontFamily{}
	var _ FontFamily = &FontListFontFamily{}
	var _ FontFamily = &LoadedFontFamily{}

	// Ensure system families implement SystemFontFamily
	var _ SystemFontFamily = &GenericFontFamily{}
	var _ SystemFontFamily = &DefaultFontFamily{}

	// Ensure file-based families implement FileBasedFontFamily
	var _ FileBasedFontFamily = &FontListFontFamily{}
}

func TestStringFontFamily(t *testing.T) {
	tests := []struct {
		name     string
		input    FontFamily
		expected string
	}{
		{
			name:     "Nil FontFamily",
			input:    nil,
			expected: "FontFamily.Default",
		},
		{
			name:     "Default FontFamily",
			input:    &DefaultFontFamily{},
			expected: "FontFamily.Default",
		},
		{
			name:     "SansSerif FontFamily",
			input:    FontFamilySansSerif,
			expected: "FontFamily.SansSerif",
		},
		{
			name:     "Custom GenericFontFamily",
			input:    NewGenericFontFamily("foo", "FontFamily.Foo"),
			expected: "FontFamily.Foo",
		},
		{
			name: "List FontFamily",
			input: NewFontListFontFamily([]Font{
				&mockFont{weight: FontWeightNormal, style: FontStyleNormal, loadingStrategy: FontLoadingStrategyBlocking},
			}),
			expected: "FontListFontFamily(fonts=[Font(weight=FontWeight(weight=400), style=Normal)])",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringFontFamily(tt.input)
			if got != tt.expected {
				t.Errorf("StringFontFamily() = %v, want %v", got, tt.expected)
			}
		})
	}
}
