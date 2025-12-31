package font

import (
	"testing"
)

func TestMaximumAsyncTimeoutMillis(t *testing.T) {
	if MaximumAsyncTimeoutMillis != 15000 {
		t.Errorf("Expected 15000, got %d", MaximumAsyncTimeoutMillis)
	}
}

// mockFont implements the Font interface for testing
type mockFont struct {
	weight          FontWeight
	style           FontStyle
	loadingStrategy FontLoadingStrategy
}

func (m *mockFont) Weight() FontWeight {
	return m.weight
}

func (m *mockFont) Style() FontStyle {
	return m.style
}

func (m *mockFont) LoadingStrategy() FontLoadingStrategy {
	return m.loadingStrategy
}

func TestToFontFamily(t *testing.T) {
	f := &mockFont{
		weight:          FontWeightNormal,
		style:           FontStyleNormal,
		loadingStrategy: FontLoadingStrategyBlocking,
	}

	family := ToFontFamily(f)
	if family == nil {
		t.Fatal("ToFontFamily returned nil")
	}

	listFamily, ok := family.(*FontListFontFamily)
	if !ok {
		t.Fatal("ToFontFamily should return FontListFontFamily")
	}

	if len(listFamily.Fonts) != 1 {
		t.Errorf("Expected 1 font, got %d", len(listFamily.Fonts))
	}

	if listFamily.Fonts[0] != f {
		t.Error("Font in family should be the original font")
	}
}

func TestFont_Interface(t *testing.T) {
	f := &mockFont{
		weight:          FontWeightBold,
		style:           FontStyleItalic,
		loadingStrategy: FontLoadingStrategyAsync,
	}

	if !f.Weight().Equals(FontWeightBold) {
		t.Error("Weight should be Bold")
	}
	if f.Style() != FontStyleItalic {
		t.Error("Style should be Italic")
	}
	if f.LoadingStrategy() != FontLoadingStrategyAsync {
		t.Error("LoadingStrategy should be Async")
	}
}
