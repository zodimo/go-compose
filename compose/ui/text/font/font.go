package font

import "fmt"

// MaximumAsyncTimeoutMillis is the global timeout for fetching an async font.
// After this timeout, a font load may no longer trigger text reflow.
const MaximumAsyncTimeoutMillis = 15_000

// Font is the interface of the font resource.
type Font interface {
	// Weight returns the weight of the font.
	// The system uses this to match a font to a font request.
	Weight() FontWeight

	// Style returns the style of the font, normal or italic.
	// The system uses this to match a font to a font request.
	Style() FontStyle

	// LoadingStrategy returns the loading strategy for this font.
	LoadingStrategy() FontLoadingStrategy
}

// ToFontFamily creates a FontFamily from a single Font.
func ToFontFamily(f Font) FontFamily {
	return NewFontListFontFamily([]Font{f})
}

func StringFont(f Font) string {
	return fmt.Sprintf("Font(weight=%s, style=%s)", f.Weight(), f.Style())
}
