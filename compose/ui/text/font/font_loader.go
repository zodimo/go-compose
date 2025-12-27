package font

import "context"

// ResourceLoader allows loading fonts from various sources.
// This interface separates the font definition from the mechanism to load the actual bytes.
type ResourceLoader interface {
	// Load loads the font data.
	// It returns the loaded font/typeface or an error.
	Load(font Font) (interface{}, error)
}

// FontLoader is a legacy name or alias if we want to be specific,
// but ResourceLoader is more generic. For now, let's Stick to ResourceLoader as per Compose pattern often used,
// but since the plan said FontLoader, let's use that for clarity in this package.

// FontLoader loads a Font and returns a Typeface.
type FontLoader interface {
	// Load loads the font.
	// context can be used for cancellation of async loads.
	Load(ctx context.Context, font Font) (Typeface, error)
}
