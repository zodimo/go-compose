package font

import (
	"github.com/zodimo/go-compose/state"
)

// DefaultFontFamilyResolver is a basic FontFamilyResolver implementation
// that provides system-default font handling. It does not do any actual
// font resolution but provides a no-op implementation suitable for basic use cases.
type DefaultFontFamilyResolver struct{}

// NewDefaultFontFamilyResolver creates a new DefaultFontFamilyResolver.
func NewDefaultFontFamilyResolver() *DefaultFontFamilyResolver {
	return &DefaultFontFamilyResolver{}
}

// Preload is a no-op in the default resolver.
func (r *DefaultFontFamilyResolver) Preload(fontFamily FontFamily) {
	// No-op: default resolver doesn't cache fonts
}

// Resolve returns a state.Value containing the resolved typeface.
// The default resolver returns a static value representing the default system font.
func (r *DefaultFontFamilyResolver) Resolve(
	fontFamily FontFamily,
	fontWeight FontWeight,
	fontStyle FontStyle,
	fontSynthesis FontSynthesis,
) state.Value {
	// Return a static value - the actual font resolution happens at render time
	// through Gio's font.Font which is already handled by the widget layer
	return &immutableValue{value: ResolvedTypeface{
		FontFamily: fontFamily,
		FontWeight: fontWeight,
		FontStyle:  fontStyle,
	}}
}

// ResolvedTypeface represents a resolved font configuration.
type ResolvedTypeface struct {
	FontFamily FontFamily
	FontWeight FontWeight
	FontStyle  FontStyle
}

// immutableValue is a simple state.Value implementation.
type immutableValue struct {
	value any
}

func (v *immutableValue) Get() any {
	return v.value
}
