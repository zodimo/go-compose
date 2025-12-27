package font

import "fmt"

// UrlFont represents a font that is loaded from a URL.
type UrlFont struct {
	url             string
	weight          FontWeight
	style           FontStyle
	loadingStrategy FontLoadingStrategy
}

// NewUrlFont creates a new UrlFont.
func NewUrlFont(url string, weight FontWeight, style FontStyle, loadingStrategy FontLoadingStrategy) *UrlFont {
	return &UrlFont{
		url:             url,
		weight:          weight,
		style:           style,
		loadingStrategy: loadingStrategy,
	}
}

// Url returns the URL of the font.
func (u *UrlFont) Url() string {
	return u.url
}

// Weight returns the weight of the font.
func (u *UrlFont) Weight() FontWeight {
	return u.weight
}

// Style returns the style of the font.
func (u *UrlFont) Style() FontStyle {
	return u.style
}

// LoadingStrategy returns the loading strategy of the font.
func (u *UrlFont) LoadingStrategy() FontLoadingStrategy {
	return u.loadingStrategy
}

// String returns a string representation of the UrlFont.
func (u *UrlFont) String() string {
	return fmt.Sprintf("UrlFont(url=%s, weight=%s, style=%s, loadingStrategy=%s)", u.url, u.weight, u.style, u.loadingStrategy)
}
