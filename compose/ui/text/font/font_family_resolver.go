package font

import "github.com/zodimo/go-compose/state"

// FontFamilyResolver is the main interface for resolving FontFamily into a platform-specific typeface.
type FontFamilyResolver interface {
	// Preload resolves and caches all fonts reachable in a FontFamily.
	Preload(fontFamily FontFamily)

	// Resolve resolves a typeface using any appropriate logic for the FontFamily.
	// Returns a state.Value that contains the platform-specific Typeface.
	Resolve(
		fontFamily FontFamily,
		fontWeight FontWeight,
		fontStyle FontStyle,
		fontSynthesis FontSynthesis,
	) state.Value
}
