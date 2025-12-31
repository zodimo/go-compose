package font

// Typeface is a class that can be used for changing the font used in text.
type Typeface interface {
	// FontFamily returns the font family used for creating this Typeface.
	// If a platform Typeface was used, it will return nil.
	FontFamily() FontFamily
}
