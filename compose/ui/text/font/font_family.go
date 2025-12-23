package font

import "fmt"

// FontFamily is the primary typography interface for Compose applications.
type FontFamily interface {
	fontFamilyMarker()
	Equals(other FontFamily) bool
}

// SystemFontFamily is a base type for FontFamilies installed on the system.
type SystemFontFamily interface {
	FontFamily
	systemFontFamilyMarker()
}

// FileBasedFontFamily is a base type for FontFamilies created from file sources.
type FileBasedFontFamily interface {
	FontFamily
	fileBasedFontFamilyMarker()
}

// GenericFontFamily represents a font family with a generic font family name.
// If the platform cannot find the passed generic font family, it uses the platform default.
type GenericFontFamily struct {
	name           string
	fontFamilyName string
}

func (g *GenericFontFamily) fontFamilyMarker()       {}
func (g *GenericFontFamily) systemFontFamilyMarker() {}
func (g *GenericFontFamily) Equals(other FontFamily) bool {
	if other == nil {
		return false
	}
	otherGeneric, ok := other.(*GenericFontFamily)
	if !ok {
		return false
	}
	return g.name == otherGeneric.name && g.fontFamilyName == otherGeneric.fontFamilyName
}

// Name returns the generic font family name (e.g., "sans-serif", "serif").
func (g *GenericFontFamily) Name() string {
	return g.name
}

// String returns the font family display name.
func (g *GenericFontFamily) String() string {
	return g.fontFamilyName
}

// NewGenericFontFamily creates a new GenericFontFamily.
func NewGenericFontFamily(name, fontFamilyName string) *GenericFontFamily {
	return &GenericFontFamily{name: name, fontFamilyName: fontFamilyName}
}

// DefaultFontFamily is the platform default font family.
type DefaultFontFamily struct{}

func (d *DefaultFontFamily) fontFamilyMarker()       {}
func (d *DefaultFontFamily) systemFontFamilyMarker() {}
func (d *DefaultFontFamily) Equals(other FontFamily) bool {
	if other == nil {
		return false
	}
	_, ok := other.(*DefaultFontFamily)
	return ok
}

// String returns the font family display name.
func (d *DefaultFontFamily) String() string {
	return "FontFamily.Default"
}

// FontListFontFamily defines a font family with a list of Fonts.
type FontListFontFamily struct {
	// Fonts is the fallback list of fonts used for resolving typefaces.
	Fonts []Font
}

func (f *FontListFontFamily) fontFamilyMarker()          {}
func (f *FontListFontFamily) fileBasedFontFamilyMarker() {}
func (f *FontListFontFamily) Equals(other FontFamily) bool {
	if other == nil {
		return false
	}
	otherList, ok := other.(*FontListFontFamily)
	if !ok {
		return false
	}
	return f.equal(otherList)
}

// NewFontListFontFamily creates a FontListFontFamily from a list of fonts.
// Panics if the fonts list is empty.
func NewFontListFontFamily(fonts []Font) *FontListFontFamily {
	if len(fonts) == 0 {
		panic("at least one font should be passed to FontFamily")
	}
	return &FontListFontFamily{Fonts: fonts}
}

// Equals checks if two FontListFontFamilies are equal.
func (f *FontListFontFamily) equal(other *FontListFontFamily) bool {
	if f == other {
		return true
	}
	if other == nil {
		return false
	}
	if len(f.Fonts) != len(other.Fonts) {
		return false
	}
	for i, font := range f.Fonts {
		if font != other.Fonts[i] {
			return false
		}
	}
	return true
}

// String returns a string representation of the FontListFontFamily.
func (f *FontListFontFamily) String() string {
	return fmt.Sprintf("FontListFontFamily(fonts=%v)", f.Fonts)
}

// LoadedFontFamily defines a font family that is already a loaded Typeface.
type LoadedFontFamily struct {
	Typeface Typeface
}

func (l *LoadedFontFamily) fontFamilyMarker() {}

// NewLoadedFontFamily creates a LoadedFontFamily from a typeface.
func NewLoadedFontFamily(typeface Typeface) *LoadedFontFamily {
	return &LoadedFontFamily{Typeface: typeface}
}

// Equals checks if two LoadedFontFamilies are equal.
func (l *LoadedFontFamily) Equals(other FontFamily) bool {

	if other == nil {
		return false
	}
	otherLoadedFontFamily, ok := other.(*LoadedFontFamily)
	if !ok {
		return false
	}
	return l.Typeface == otherLoadedFontFamily.Typeface
}

// String returns a string representation of the LoadedFontFamily.
func (l *LoadedFontFamily) String() string {
	return fmt.Sprintf("LoadedFontFamily(typeface=%v)", l.Typeface)
}

// Standard font family constants
var (
	// FontFamilyDefault is the platform default font.
	FontFamilyDefault FontFamily = &DefaultFontFamily{}

	// FontFamilySansSerif is a font family with low contrast and plain stroke endings.
	FontFamilySansSerif FontFamily = NewGenericFontFamily("sans-serif", "FontFamily.SansSerif")

	// FontFamilySerif is the formal text style for scripts.
	FontFamilySerif FontFamily = NewGenericFontFamily("serif", "FontFamily.Serif")

	// FontFamilyMonospace is a font family where glyphs have the same fixed width.
	FontFamilyMonospace FontFamily = NewGenericFontFamily("monospace", "FontFamily.Monospace")

	// FontFamilyCursive is a cursive, hand-written like font family.
	FontFamilyCursive FontFamily = NewGenericFontFamily("cursive", "FontFamily.Cursive")
)

// FontFamilyFromFonts creates a FontFamily from a list of fonts.
func FontFamilyFromFonts(fonts ...Font) FontFamily {
	return NewFontListFontFamily(fonts)
}

// FontFamilyFromTypeface creates a FontFamily from a loaded typeface.
func FontFamilyFromTypeface(typeface Typeface) FontFamily {
	return NewLoadedFontFamily(typeface)
}
