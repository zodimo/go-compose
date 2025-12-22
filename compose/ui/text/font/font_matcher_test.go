package font

import "testing"

// testFont is a simple Font implementation for testing
type testFont struct {
	weight          FontWeight
	style           FontStyle
	loadingStrategy FontLoadingStrategy
}

func (f *testFont) Weight() FontWeight                   { return f.weight }
func (f *testFont) Style() FontStyle                     { return f.style }
func (f *testFont) LoadingStrategy() FontLoadingStrategy { return f.loadingStrategy }

func newTestFont(weight FontWeight, style FontStyle) *testFont {
	return &testFont{
		weight:          weight,
		style:           style,
		loadingStrategy: FontLoadingStrategyBlocking,
	}
}

func TestFontMatcher_ExactMatch(t *testing.T) {
	matcher := NewFontMatcher()

	fonts := []Font{
		newTestFont(FontWeightNormal, FontStyleNormal),
		newTestFont(FontWeightBold, FontStyleNormal),
		newTestFont(FontWeightNormal, FontStyleItalic),
	}

	result := matcher.MatchFont(fonts, FontWeightNormal, FontStyleNormal)
	if len(result) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(result))
	}
	if !result[0].Weight().Equals(FontWeightNormal) || result[0].Style() != FontStyleNormal {
		t.Error("Wrong font matched")
	}
}

func TestFontMatcher_ExactMatch_Bold(t *testing.T) {
	matcher := NewFontMatcher()

	fonts := []Font{
		newTestFont(FontWeightNormal, FontStyleNormal),
		newTestFont(FontWeightBold, FontStyleNormal),
	}

	result := matcher.MatchFont(fonts, FontWeightBold, FontStyleNormal)
	if len(result) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(result))
	}
	if !result[0].Weight().Equals(FontWeightBold) {
		t.Error("Expected Bold weight")
	}
}

func TestFontMatcher_StyleMismatch_UsesStyleFirst(t *testing.T) {
	matcher := NewFontMatcher()

	fonts := []Font{
		newTestFont(FontWeightNormal, FontStyleNormal),
		newTestFont(FontWeightBold, FontStyleItalic),
	}

	// Request Bold Normal - should get Normal Normal (style match)
	// since there's no Bold Normal
	result := matcher.MatchFont(fonts, FontWeightBold, FontStyleNormal)
	if len(result) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(result))
	}
	if result[0].Style() != FontStyleNormal {
		t.Error("Should prefer style match")
	}
}

func TestFontMatcher_WeightBelow400_PreferLower(t *testing.T) {
	matcher := NewFontMatcher()

	// Requesting W200, have W100 and W300
	fonts := []Font{
		newTestFont(FontWeightW100, FontStyleNormal),
		newTestFont(FontWeightW300, FontStyleNormal),
	}

	result := matcher.MatchFont(fonts, FontWeightW200, FontStyleNormal)
	if len(result) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(result))
	}
	// Should prefer W100 (below) over W300 (above)
	if !result[0].Weight().Equals(FontWeightW100) {
		t.Errorf("Expected W100, got %v", result[0].Weight())
	}
}

func TestFontMatcher_WeightAbove500_PreferHigher(t *testing.T) {
	matcher := NewFontMatcher()

	// Requesting W600, have W500 and W700
	fonts := []Font{
		newTestFont(FontWeightW500, FontStyleNormal),
		newTestFont(FontWeightW700, FontStyleNormal),
	}

	result := matcher.MatchFont(fonts, FontWeightW600, FontStyleNormal)
	if len(result) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(result))
	}
	// Should prefer W700 (above) over W500 (below)
	if !result[0].Weight().Equals(FontWeightW700) {
		t.Errorf("Expected W700, got %v", result[0].Weight())
	}
}

func TestFontMatcher_WeightBetween400And500(t *testing.T) {
	matcher := NewFontMatcher()

	// Requesting W450, have W400 and W500
	fonts := []Font{
		newTestFont(FontWeightW400, FontStyleNormal),
		newTestFont(FontWeightW500, FontStyleNormal),
		newTestFont(FontWeightW600, FontStyleNormal),
	}

	result := matcher.MatchFont(fonts, NewFontWeight(450), FontStyleNormal)
	if len(result) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(result))
	}
	// Should prefer W500 (above, within 400-500 range)
	if !result[0].Weight().Equals(FontWeightW500) {
		t.Errorf("Expected W500, got %v", result[0].Weight())
	}
}

func TestFontMatcher_EmptyFonts(t *testing.T) {
	matcher := NewFontMatcher()

	result := matcher.MatchFont([]Font{}, FontWeightNormal, FontStyleNormal)
	if result != nil && len(result) != 0 {
		t.Error("Empty fonts should return empty result")
	}
}

func TestFontMatcher_NoStyleMatch_FallsBackToAll(t *testing.T) {
	matcher := NewFontMatcher()

	// Only have Italic fonts, but requesting Normal
	fonts := []Font{
		newTestFont(FontWeightNormal, FontStyleItalic),
		newTestFont(FontWeightBold, FontStyleItalic),
	}

	result := matcher.MatchFont(fonts, FontWeightNormal, FontStyleNormal)
	if len(result) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(result))
	}
	// Should fallback to matching from all fonts
	if !result[0].Weight().Equals(FontWeightNormal) {
		t.Error("Should match closest weight from all fonts")
	}
}

func TestFontMatcher_MultipleWithSameWeight(t *testing.T) {
	matcher := NewFontMatcher()

	f1 := newTestFont(FontWeightNormal, FontStyleNormal)
	f2 := newTestFont(FontWeightNormal, FontStyleNormal)

	fonts := []Font{f1, f2}

	result := matcher.MatchFont(fonts, FontWeightNormal, FontStyleNormal)
	if len(result) != 2 {
		t.Fatalf("Expected 2 matches, got %d", len(result))
	}
}

func TestFontMatcher_MatchFontFromFamily(t *testing.T) {
	matcher := NewFontMatcher()

	fonts := []Font{
		newTestFont(FontWeightNormal, FontStyleNormal),
		newTestFont(FontWeightBold, FontStyleNormal),
	}
	family := NewFontListFontFamily(fonts)

	result := matcher.MatchFontFromFamily(family, FontWeightBold, FontStyleNormal)
	if len(result) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(result))
	}
	if !result[0].Weight().Equals(FontWeightBold) {
		t.Error("Expected Bold weight")
	}
}
