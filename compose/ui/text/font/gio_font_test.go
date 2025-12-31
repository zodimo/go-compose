package font

import (
	"testing"

	giofont "gioui.org/font"
)

func TestToGioWeight(t *testing.T) {
	tests := []struct {
		name     string
		input    FontWeight
		expected giofont.Weight
	}{
		{"Thin/100", FontWeightW100, giofont.Thin},
		{"Normal/400", FontWeightW400, giofont.Normal},
		{"Bold/700", FontWeightW700, giofont.Bold},
		{"Black/900", FontWeightW900, giofont.Black},
		{"Unspecified", FontWeightUnspecified, giofont.Normal},
		{"Custom/250", FontWeight(250), giofont.Weight(-150)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToGioWeight(tt.input)
			if got != tt.expected {
				t.Errorf("ToGioWeight(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestToGioStyle(t *testing.T) {
	tests := []struct {
		name     string
		input    FontStyle
		expected giofont.Style
	}{
		{"Normal", FontStyleNormal, giofont.Regular},
		{"Italic", FontStyleItalic, giofont.Italic},
		{"Unspecified", FontStyleUnspecified, giofont.Regular},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToGioStyle(tt.input)
			if got != tt.expected {
				t.Errorf("ToGioStyle(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestResolveGioTypeface(t *testing.T) {
	tests := []struct {
		name     string
		input    FontFamily
		expected string
	}{
		{"Generic Sans", FontFamilySansSerif, "sans-serif"},
		{"Generic Serif", FontFamilySerif, "serif"},
		{"Generic Monospace", FontFamilyMonospace, "monospace"},
		{"Default", FontFamilyDefault, ""},
		{"Nil", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveGioTypeface(tt.input)
			if got != tt.expected {
				t.Errorf("ResolveGioTypeface(%v) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestToGioFont(t *testing.T) {
	f := FontFamilySansSerif
	w := FontWeightBold
	s := FontStyleItalic

	got := ToGioFont(f, w, s)

	if got.Typeface != "sans-serif" {
		t.Errorf("ToGioFont().Typeface = %q, want %q", got.Typeface, "sans-serif")
	}
	if got.Weight != giofont.Bold {
		t.Errorf("ToGioFont().Weight = %v, want %v", got.Weight, giofont.Bold)
	}
	if got.Style != giofont.Italic {
		t.Errorf("ToGioFont().Style = %v, want %v", got.Style, giofont.Italic)
	}
}
