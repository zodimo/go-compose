package font

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// PlatformFontLoader implements FontLoader to load fonts from URLs or files.
type PlatformFontLoader struct {
	cacheManager *CacheManager
	httpClient   *http.Client
}

// NewPlatformFontLoader creates a new PlatformFontLoader.
func NewPlatformFontLoader(cacheManager *CacheManager, httpClient *http.Client) *PlatformFontLoader {
	return &PlatformFontLoader{
		cacheManager: cacheManager,
		httpClient:   httpClient,
	}
}

// Load loads the font.
func (l *PlatformFontLoader) Load(ctx context.Context, font Font) (Typeface, error) {
	switch f := font.(type) {
	case *UrlFont:
		return l.loadUrlFont(ctx, f)
	default:
		return nil, fmt.Errorf("unsupported font type: %T", font)
	}
}

func (l *PlatformFontLoader) loadUrlFont(ctx context.Context, font *UrlFont) (Typeface, error) {
	// Check cache first
	if path, exists := l.cacheManager.Get(font.Url()); exists {
		return l.loadTypefaceFromFile(path)
	}

	// Download if not in cache
	req, err := http.NewRequestWithContext(ctx, "GET", font.Url(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download font: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("background font download failed with status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Cache the font
	path, err := l.cacheManager.Put(font.Url(), data)
	if err != nil {
		// If caching fails, we can still try to load from memory or just warn.
		// For now, let's treat the saved file as the source of truth for creating the typeface.
		return nil, fmt.Errorf("failed to cache font: %w", err)
	}

	return l.loadTypefaceFromFile(path)
}

func (l *PlatformFontLoader) loadTypefaceFromFile(path string) (Typeface, error) {
	// TODO: Implement actual parsing of the font file to a Typeface.
	// For now, we return a mock or placeholder since Typeface creation depends on the graphics/text engine (Gio).
	// We might need to extend the Typeface interface or have a specific implementation that holds the bytes/path.
	return &LoadedTypeface{Source: path}, nil
}

// LoadedTypeface is a placeholder implementation of Typeface.
type LoadedTypeface struct {
	Source string
}

func (l *LoadedTypeface) FontFamily() FontFamily {
	return nil // loaded typefaces generally don't have a family name known a priori unless parsed
}
