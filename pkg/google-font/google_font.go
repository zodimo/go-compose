package googlefont

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/zodimo/go-compose/compose/ui/text/font"
)

// Provider fetches font definitions from Google Fonts.
type Provider struct {
	httpClient *http.Client
}

// NewProvider creates a new Provider.
func NewProvider(httpClient *http.Client) *Provider {
	return &Provider{
		httpClient: httpClient,
	}
}

// FontFamily creates a FontFamily by fetching definitions from Google Fonts.
// Note: This function performs a blocking network call to fetch the CSS.
// In a real application, this might need to be async or cached appropriately,
// but for the purpose of creating a FontFamily definition, we need the URLs.
//
// familyName: e.g. "Roboto", "Open Sans"
// axes: currently supported are "ital" and "wght".
// The implementation constructs a query to fetch all variants (ital/normal, weights 100-900) for simplicity
// or we can allow the user to specify.
//
// For this first pass, we will just take the family name and try to fetch a broad range of styles/weights
// to populate the font family.
func (p *Provider) FontFamily(ctx context.Context, familyName string) (font.FontFamily, error) {
	// Construct URL
	// https://fonts.googleapis.com/css2?family=Roboto:ital,wght@0,100..900;1,100..900&display=swap
	// We ask for full range for now to support dynamic runtime usage.

	// Encode family name
	escapedFamily := url.QueryEscape(familyName)

	// We construct a standard query for 100-900 weights, both italic and normal.
	// 0,100..900 -> normal 100-900
	// 1,100..900 -> italic 100-900
	apiURL := fmt.Sprintf("https://fonts.googleapis.com/css2?family=%s:ital,wght@0,100..900;1,100..900&display=swap", escapedFamily)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Google Fonts requires a user agent usually, or at least handles it better.
	// But standard Go client usually works.

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch google fonts css: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google fonts api returned status: %d", resp.StatusCode)
	}

	cssBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read css response: %w", err)
	}

	parsedFonts := ParseGoogleFontsCss(string(cssBytes))
	if len(parsedFonts) == 0 {
		return nil, fmt.Errorf("no fonts found in css for family: %s", familyName)
	}

	var fonts []font.Font
	for _, pf := range parsedFonts {
		style := font.FontStyleNormal
		if pf.FontStyle == "italic" {
			style = font.FontStyleItalic
		}

		weight := font.FontWeight(pf.FontWeight)
		if weight < font.FontWeight(1) || weight > font.FontWeight(1000) {
			weight = font.FontWeightNormal
		}

		fonts = append(fonts, font.NewUrlFont(pf.Src, weight, style, font.FontLoadingStrategyAsync))
	}

	return font.NewFontListFontFamily(fonts), nil
}
