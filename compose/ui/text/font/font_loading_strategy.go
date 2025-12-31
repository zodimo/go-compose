package font

import "fmt"

// FontLoadingStrategy controls how font loading resolves when displaying text.
type FontLoadingStrategy int

const (
	// FontLoadingStrategyBlocking resolves the font by blocking until loaded.
	// This means the first frame that uses this font will always display using
	// the desired font, and text will never reflow.
	// Should be used for fonts stored on-device.
	FontLoadingStrategyBlocking FontLoadingStrategy = 0

	// FontLoadingStrategyOptionalLocal is best-effort loading from a local resource
	// that MAY be available. Resolution is expected to fail when the resource is
	// not available, which will fallback to the next font in the FontFamily.
	FontLoadingStrategyOptionalLocal FontLoadingStrategy = 1

	// FontLoadingStrategyAsync loads the font in the background without blocking.
	// During loading, the app will use a fallback font. When the font finishes
	// loading, text will reflow with the resolved typeface.
	// Should be used for fonts fetched from a remote source.
	FontLoadingStrategyAsync FontLoadingStrategy = 2
)

// Value returns the underlying integer value.
func (s FontLoadingStrategy) Value() int {
	return int(s)
}

// String returns a string representation of the FontLoadingStrategy.
func (s FontLoadingStrategy) String() string {
	switch s {
	case FontLoadingStrategyBlocking:
		return "Blocking"
	case FontLoadingStrategyOptionalLocal:
		return "OptionalLocal"
	case FontLoadingStrategyAsync:
		return "Async"
	default:
		return fmt.Sprintf("Invalid(value=%d)", s)
	}
}
