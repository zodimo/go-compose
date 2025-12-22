package font

import "fmt"

// FontVariationSetting represents a single point in a font variation axis.
type FontVariationSetting interface {
	// AxisName returns the font variation axis name (e.g., "wght", "ital").
	AxisName() string

	// ToVariationValue converts the setting to a final value for use as a font variation setting.
	// The fontScale parameter is used for settings that need density (like optical sizing).
	ToVariationValue(fontScale float32) float32

	// NeedsDensity returns true if this setting requires density to resolve.
	NeedsDensity() bool
}

// FontVariationSettings is a collection of settings to apply to a single font.
// Settings must be unique on AxisName.
type FontVariationSettings struct {
	settings     []FontVariationSetting
	needsDensity bool
}

// Settings returns all settings in this collection.
func (s *FontVariationSettings) Settings() []FontVariationSetting {
	return s.settings
}

// NeedsDensity returns true if any setting requires density to resolve.
func (s *FontVariationSettings) NeedsDensity() bool {
	return s.needsDensity
}

// IsEmpty returns true if there are no settings.
func (s *FontVariationSettings) IsEmpty() bool {
	return len(s.settings) == 0
}

// Equals checks if two FontVariationSettings are equal.
func (s *FontVariationSettings) Equals(other *FontVariationSettings) bool {
	if s == other {
		return true
	}
	if other == nil {
		return false
	}
	if len(s.settings) != len(other.settings) {
		return false
	}
	for i, setting := range s.settings {
		otherSetting := other.settings[i]
		if setting.AxisName() != otherSetting.AxisName() {
			return false
		}
		// Compare values at fontScale 1.0
		if setting.ToVariationValue(1.0) != otherSetting.ToVariationValue(1.0) {
			return false
		}
	}
	return true
}

// NewFontVariationSettings creates a FontVariationSettings from a list of settings.
// Panics if duplicate axis names are provided.
func NewFontVariationSettings(settings ...FontVariationSetting) *FontVariationSettings {
	// Check for duplicates
	seen := make(map[string]bool)
	for _, s := range settings {
		if seen[s.AxisName()] {
			panic(fmt.Sprintf("'%s' must be unique", s.AxisName()))
		}
		seen[s.AxisName()] = true
	}

	needsDensity := false
	for _, s := range settings {
		if s.NeedsDensity() {
			needsDensity = true
			break
		}
	}

	return &FontVariationSettings{
		settings:     settings,
		needsDensity: needsDensity,
	}
}

// FontVariationSettingsEmpty returns an empty FontVariationSettings.
func FontVariationSettingsEmpty() *FontVariationSettings {
	return &FontVariationSettings{
		settings:     nil,
		needsDensity: false,
	}
}

// FontVariationSettingsFromWeightStyle creates settings that configure FontWeight and FontStyle.
func FontVariationSettingsFromWeightStyle(weight FontWeight, style FontStyle, additional ...FontVariationSetting) *FontVariationSettings {
	settings := []FontVariationSetting{
		FontVariationWeight(weight.Weight()),
		FontVariationItalic(float32(style.Value())),
	}
	settings = append(settings, additional...)
	return NewFontVariationSettings(settings...)
}

// settingFloat is a simple float-valued setting.
type settingFloat struct {
	axisName string
	value    float32
}

func (s *settingFloat) AxisName() string                           { return s.axisName }
func (s *settingFloat) ToVariationValue(fontScale float32) float32 { return s.value }
func (s *settingFloat) NeedsDensity() bool                         { return false }
func (s *settingFloat) String() string {
	return fmt.Sprintf("FontVariation.Setting(axisName='%s', value=%v)", s.axisName, s.value)
}

// settingInt is an int-valued setting stored as float.
type settingInt struct {
	axisName string
	value    int
}

func (s *settingInt) AxisName() string                           { return s.axisName }
func (s *settingInt) ToVariationValue(fontScale float32) float32 { return float32(s.value) }
func (s *settingInt) NeedsDensity() bool                         { return false }
func (s *settingInt) String() string {
	return fmt.Sprintf("FontVariation.Setting(axisName='%s', value=%d)", s.axisName, s.value)
}

// settingTextUnit is a text-unit valued setting that needs fontScale.
type settingTextUnit struct {
	axisName string
	value    float32 // in sp
}

func (s *settingTextUnit) AxisName() string { return s.axisName }
func (s *settingTextUnit) ToVariationValue(fontScale float32) float32 {
	return s.value * fontScale
}
func (s *settingTextUnit) NeedsDensity() bool { return true }
func (s *settingTextUnit) String() string {
	return fmt.Sprintf("FontVariation.Setting(axisName='%s', value=%vsp)", s.axisName, s.value)
}

// NewFontVariationSetting creates a font variation setting for any axis.
// The name must be exactly 4 characters.
func NewFontVariationSetting(name string, value float32) FontVariationSetting {
	if len(name) != 4 {
		panic(fmt.Sprintf("Name must be exactly four characters. Actual: '%s'", name))
	}
	return &settingFloat{axisName: name, value: value}
}

// FontVariationItalic creates an italic setting ('ital').
// Value should be in range [0.0, 1.0] where 0.0 is upright and 1.0 is italic.
func FontVariationItalic(value float32) FontVariationSetting {
	if value < 0.0 || value > 1.0 {
		panic(fmt.Sprintf("'ital' must be in 0.0..1.0. Actual: %v", value))
	}
	return &settingFloat{axisName: "ital", value: value}
}

// FontVariationSlant creates a slant setting ('slnt').
// Value is an angle in range [-90, 90] where 0 is upright.
func FontVariationSlant(value float32) FontVariationSetting {
	if value < -90.0 || value > 90.0 {
		panic(fmt.Sprintf("'slnt' must be in -90..90. Actual: %v", value))
	}
	return &settingFloat{axisName: "slnt", value: value}
}

// FontVariationWidth creates a width setting ('wdth').
// Value must be > 0.0.
func FontVariationWidth(value float32) FontVariationSetting {
	if value <= 0.0 {
		panic(fmt.Sprintf("'wdth' must be strictly > 0.0. Actual: %v", value))
	}
	return &settingFloat{axisName: "wdth", value: value}
}

// FontVariationWeight creates a weight setting ('wght').
// Value must be in range [1, 1000].
func FontVariationWeight(value int) FontVariationSetting {
	if value < 1 || value > 1000 {
		panic(fmt.Sprintf("'wght' value must be in [1, 1000]. Actual: %d", value))
	}
	return &settingInt{axisName: "wght", value: value}
}

// FontVariationGrade creates a grade setting ('GRAD').
// Value must be in range [-1000, 1000].
func FontVariationGrade(value int) FontVariationSetting {
	if value < -1000 || value > 1000 {
		panic(fmt.Sprintf("'GRAD' must be in -1000..1000. Actual: %d", value))
	}
	return &settingInt{axisName: "GRAD", value: value}
}

// FontVariationOpticalSizing creates an optical sizing setting ('opsz').
// Value is the font size in sp (scaled pixels).
func FontVariationOpticalSizing(valueSp float32) FontVariationSetting {
	return &settingTextUnit{axisName: "opsz", value: valueSp}
}
