package style

import (
	"testing"
)

func TestLineHeightStyleDefaults(t *testing.T) {
	if DefaultLineHeightStyle.Alignment != LineHeightStyleAlignmentProportional {
		t.Errorf("Default alignment mismatch")
	}
	if DefaultLineHeightStyle.Trim != LineHeightStyleTrimBoth {
		t.Errorf("Default trim mismatch")
	}
	if DefaultLineHeightStyle.Mode != LineHeightStyleModeFixed {
		t.Errorf("Default mode mismatch")
	}
}

func TestLineHeightStyleTrim(t *testing.T) {
	if !LineHeightStyleTrimFirstLineTop.IsTrimFirstLineTop() {
		t.Errorf("FirstLineTop should trim first line top")
	}
	if LineHeightStyleTrimFirstLineTop.IsTrimLastLineBottom() {
		t.Errorf("FirstLineTop should NOT trim last line bottom")
	}

	if !LineHeightStyleTrimLastLineBottom.IsTrimLastLineBottom() {
		t.Errorf("LastLineBottom should trim last line bottom")
	}
	if LineHeightStyleTrimLastLineBottom.IsTrimFirstLineTop() {
		t.Errorf("LastLineBottom should NOT trim first line top")
	}

	if !LineHeightStyleTrimBoth.IsTrimFirstLineTop() {
		t.Errorf("Both should trim first line top")
	}
	if !LineHeightStyleTrimBoth.IsTrimLastLineBottom() {
		t.Errorf("Both should trim last line bottom")
	}

	if LineHeightStyleTrimNone.IsTrimFirstLineTop() {
		t.Errorf("None should NOT trim first line top")
	}
	if LineHeightStyleTrimNone.IsTrimLastLineBottom() {
		t.Errorf("None should NOT trim last line bottom")
	}
}

func TestAlignmentString(t *testing.T) {
	if StringLineHeightStyleAlignment(LineHeightStyleAlignmentTop) != "LineHeightStyle.Alignment.Top" {
		t.Errorf("Top alignment string mismatch")
	}
	// Verify custom alignment
	custom := LineHeightStyleAlignment{TopRatio: 0.25}
	if custom.TopRatio != 0.25 {
		t.Errorf("Custom TopRatio mismatch")
	}
	if StringLineHeightStyleAlignment(&custom) != "LineHeightStyle.Alignment(topRatio = 0.25)" {
		t.Errorf("String() mismatch, got %s", StringLineHeightStyleAlignment(&custom))
	}
}
