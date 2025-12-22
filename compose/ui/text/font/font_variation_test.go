package font

import "testing"

func TestFontVariationSetting_AxisName(t *testing.T) {
	s := NewFontVariationSetting("test", 1.0)
	if s.AxisName() != "test" {
		t.Errorf("Expected 'test', got %s", s.AxisName())
	}
}

func TestFontVariationSetting_InvalidName(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid axis name length")
		}
	}()
	NewFontVariationSetting("toolong", 1.0)
}

func TestFontVariationItalic(t *testing.T) {
	s := FontVariationItalic(0.5)
	if s.AxisName() != "ital" {
		t.Errorf("Expected 'ital', got %s", s.AxisName())
	}
	if s.ToVariationValue(1.0) != 0.5 {
		t.Errorf("Expected 0.5, got %v", s.ToVariationValue(1.0))
	}
	if s.NeedsDensity() {
		t.Error("Italic should not need density")
	}
}

func TestFontVariationItalic_InvalidRange(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for ital > 1.0")
		}
	}()
	FontVariationItalic(1.5)
}

func TestFontVariationSlant(t *testing.T) {
	s := FontVariationSlant(-12.0)
	if s.AxisName() != "slnt" {
		t.Errorf("Expected 'slnt', got %s", s.AxisName())
	}
	if s.ToVariationValue(1.0) != -12.0 {
		t.Errorf("Expected -12.0, got %v", s.ToVariationValue(1.0))
	}
}

func TestFontVariationSlant_InvalidRange(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for slnt > 90")
		}
	}()
	FontVariationSlant(100.0)
}

func TestFontVariationWidth(t *testing.T) {
	s := FontVariationWidth(100.0)
	if s.AxisName() != "wdth" {
		t.Errorf("Expected 'wdth', got %s", s.AxisName())
	}
	if s.ToVariationValue(1.0) != 100.0 {
		t.Errorf("Expected 100.0, got %v", s.ToVariationValue(1.0))
	}
}

func TestFontVariationWidth_InvalidRange(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for wdth <= 0")
		}
	}()
	FontVariationWidth(0.0)
}

func TestFontVariationWeight(t *testing.T) {
	s := FontVariationWeight(400)
	if s.AxisName() != "wght" {
		t.Errorf("Expected 'wght', got %s", s.AxisName())
	}
	if s.ToVariationValue(1.0) != 400.0 {
		t.Errorf("Expected 400.0, got %v", s.ToVariationValue(1.0))
	}
}

func TestFontVariationWeight_InvalidRange(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for wght > 1000")
		}
	}()
	FontVariationWeight(1001)
}

func TestFontVariationGrade(t *testing.T) {
	s := FontVariationGrade(-50)
	if s.AxisName() != "GRAD" {
		t.Errorf("Expected 'GRAD', got %s", s.AxisName())
	}
	if s.ToVariationValue(1.0) != -50.0 {
		t.Errorf("Expected -50.0, got %v", s.ToVariationValue(1.0))
	}
}

func TestFontVariationGrade_InvalidRange(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for GRAD out of range")
		}
	}()
	FontVariationGrade(1001)
}

func TestFontVariationOpticalSizing(t *testing.T) {
	s := FontVariationOpticalSizing(14.0)
	if s.AxisName() != "opsz" {
		t.Errorf("Expected 'opsz', got %s", s.AxisName())
	}
	if !s.NeedsDensity() {
		t.Error("OpticalSizing should need density")
	}
	// With fontScale 1.5, 14sp should become 21
	if s.ToVariationValue(1.5) != 21.0 {
		t.Errorf("Expected 21.0, got %v", s.ToVariationValue(1.5))
	}
}

func TestNewFontVariationSettings(t *testing.T) {
	settings := NewFontVariationSettings(
		FontVariationWeight(400),
		FontVariationItalic(0.0),
	)

	if len(settings.Settings()) != 2 {
		t.Errorf("Expected 2 settings, got %d", len(settings.Settings()))
	}

	if settings.NeedsDensity() {
		t.Error("These settings should not need density")
	}
}

func TestNewFontVariationSettings_WithDensity(t *testing.T) {
	settings := NewFontVariationSettings(
		FontVariationWeight(400),
		FontVariationOpticalSizing(14.0),
	)

	if !settings.NeedsDensity() {
		t.Error("These settings should need density")
	}
}

func TestNewFontVariationSettings_DuplicatePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for duplicate axis names")
		}
	}()
	NewFontVariationSettings(
		FontVariationWeight(400),
		FontVariationWeight(700),
	)
}

func TestFontVariationSettingsEmpty(t *testing.T) {
	settings := FontVariationSettingsEmpty()
	if !settings.IsEmpty() {
		t.Error("Empty settings should be empty")
	}
	if settings.NeedsDensity() {
		t.Error("Empty settings should not need density")
	}
}

func TestFontVariationSettings_Equals(t *testing.T) {
	s1 := NewFontVariationSettings(FontVariationWeight(400))
	s2 := NewFontVariationSettings(FontVariationWeight(400))
	s3 := NewFontVariationSettings(FontVariationWeight(700))

	if !s1.Equals(s1) {
		t.Error("Settings should equal itself")
	}
	if !s1.Equals(s2) {
		t.Error("Settings with same values should be equal")
	}
	if s1.Equals(s3) {
		t.Error("Settings with different values should not be equal")
	}
	if s1.Equals(nil) {
		t.Error("Settings should not equal nil")
	}
}

func TestFontVariationSettingsFromWeightStyle(t *testing.T) {
	settings := FontVariationSettingsFromWeightStyle(FontWeightBold, FontStyleItalic)

	if len(settings.Settings()) != 2 {
		t.Errorf("Expected 2 settings, got %d", len(settings.Settings()))
	}

	// Find weight and italic settings
	var hasWeight, hasItalic bool
	for _, s := range settings.Settings() {
		if s.AxisName() == "wght" {
			hasWeight = true
			if s.ToVariationValue(1.0) != 700.0 {
				t.Errorf("Weight should be 700, got %v", s.ToVariationValue(1.0))
			}
		}
		if s.AxisName() == "ital" {
			hasItalic = true
			if s.ToVariationValue(1.0) != 1.0 {
				t.Errorf("Italic should be 1.0, got %v", s.ToVariationValue(1.0))
			}
		}
	}

	if !hasWeight {
		t.Error("Should have weight setting")
	}
	if !hasItalic {
		t.Error("Should have italic setting")
	}
}
