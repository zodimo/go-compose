package style

import "testing"

func TestLinearity_String(t *testing.T) {
	tests := []struct {
		name     string
		lin      Linearity
		expected string
	}{
		{"Unspecified", LinearityUnspecified, "Linearity.Unspecified"},
		{"Linear", LinearityLinear, "Linearity.Linear"},
		{"FontHinting", LinearityFontHinting, "Linearity.FontHinting"},
		{"None", LinearityNone, "Linearity.None"},
		{"Invalid", Linearity(999), "Invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.lin.String()
			if result != tt.expected {
				t.Errorf("Linearity.String() = %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestTextMotion_Equals(t *testing.T) {
	tests := []struct {
		name     string
		a        *TextMotion
		b        *TextMotion
		expected bool
	}{
		{
			name:     "Static equals Static",
			a:        TextMotionStatic,
			b:        TextMotionStatic,
			expected: true,
		},
		{
			name:     "Animated equals Animated",
			a:        TextMotionAnimated,
			b:        TextMotionAnimated,
			expected: true,
		},
		{
			name:     "Static not equals Animated",
			a:        TextMotionStatic,
			b:        TextMotionAnimated,
			expected: false,
		},
		{
			name: "Custom with same values",
			a: &TextMotion{
				Linearity:               LinearityLinear,
				SubpixelTextPositioning: SubpixelTextPositioningTrue,
			},
			b: &TextMotion{
				Linearity:               LinearityLinear,
				SubpixelTextPositioning: SubpixelTextPositioningTrue,
			},
			expected: true,
		},
		{
			name: "Different linearity",
			a: &TextMotion{
				Linearity:               LinearityLinear,
				SubpixelTextPositioning: SubpixelTextPositioningTrue,
			},
			b: &TextMotion{
				Linearity:               LinearityFontHinting,
				SubpixelTextPositioning: SubpixelTextPositioningTrue,
			},
			expected: false,
		},
		{
			name: "Different subpixel",
			a: &TextMotion{
				Linearity:               LinearityLinear,
				SubpixelTextPositioning: SubpixelTextPositioningTrue,
			},
			b: &TextMotion{
				Linearity:               LinearityLinear,
				SubpixelTextPositioning: SubpixelTextPositioningFalse,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EqualTextMotion(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Equals() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestTextMotion_Copy(t *testing.T) {
	// Copy with no changes
	copied := TextMotionStatic.Copy()
	if !EqualTextMotion(copied, TextMotionStatic) {
		t.Errorf("Copy with no changes should equal original")
	}

	// Copy with linearity change
	newLinearity := LinearityLinear
	copied = TextMotionStatic.Copy(WithLinearity(newLinearity))
	if copied.Linearity != LinearityLinear {
		t.Errorf("Copy should have updated Linearity to LinearityLinear")
	}
	if copied.SubpixelTextPositioning != TextMotionStatic.SubpixelTextPositioning {
		t.Errorf("Copy should have preserved SubpixelTextPositioning as false")
	}

	// Copy with subpixel change
	newSubpixel := true
	copied = TextMotionStatic.Copy(WithSubpixelTextPositioning(newSubpixel))
	if copied.Linearity != LinearityFontHinting {
		t.Errorf("Copy should have preserved Linearity as FontHinting")
	}
	if copied.SubpixelTextPositioning != SubpixelTextPositioningTrue {
		t.Errorf("Copy should have updated SubpixelTextPositioning to true")
	}

	// Copy with both changes
	copied = TextMotionStatic.Copy(WithLinearity(newLinearity), WithSubpixelTextPositioning(newSubpixel))
	if !EqualTextMotion(copied, TextMotionAnimated) {
		t.Errorf("Copy with Linear and subpixel=true should equal TextMotionAnimated")
	}
}

func TestTextMotion_String(t *testing.T) {
	tests := []struct {
		name     string
		motion   *TextMotion
		expected string
	}{

		{"Static", TextMotionStatic, "TextMotion.Static"},
		{"Animated", TextMotionAnimated, "TextMotion.Animated"},
		{
			"Custom",
			&TextMotion{
				Linearity:               LinearityNone,
				SubpixelTextPositioning: SubpixelTextPositioningFalse,
			},
			"TextMotion(Linearity.None, subpixel=False)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringTextMotion(tt.motion)
			if result != tt.expected {
				t.Errorf("String() = %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestTextMotion_Constants(t *testing.T) {
	// Verify Static constant values
	if TextMotionStatic.Linearity != LinearityFontHinting {
		t.Errorf("TextMotionStatic.Linearity should be LinearityFontHinting")
	}
	if TextMotionStatic.SubpixelTextPositioning != SubpixelTextPositioningFalse {
		t.Errorf("TextMotionStatic.SubpixelTextPositioning should be false")
	}

	// Verify Animated constant values
	if TextMotionAnimated.Linearity != LinearityLinear {
		t.Errorf("TextMotionAnimated.Linearity should be LinearityLinear")
	}
	if TextMotionAnimated.SubpixelTextPositioning != SubpixelTextPositioningTrue {
		t.Errorf("TextMotionAnimated.SubpixelTextPositioning should be true")
	}
}

func TestLinearity_Constants(t *testing.T) {
	// Ensure linearity constants have distinct values
	if LinearityLinear == LinearityFontHinting {
		t.Errorf("LinearityLinear should not equal LinearityFontHinting")
	}
	if LinearityLinear == LinearityNone {
		t.Errorf("LinearityLinear should not equal LinearityNone")
	}
	if LinearityFontHinting == LinearityNone {
		t.Errorf("LinearityFontHinting should not equal LinearityNone")
	}

	// Verify expected values match Kotlin implementation
	if LinearityLinear != 1 {
		t.Errorf("LinearityLinear should be 1, got %d", LinearityLinear)
	}
	if LinearityFontHinting != 2 {
		t.Errorf("LinearityFontHinting should be 2, got %d", LinearityFontHinting)
	}
	if LinearityNone != 3 {
		t.Errorf("LinearityNone should be 3, got %d", LinearityNone)
	}
}
