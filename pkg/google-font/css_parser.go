package googlefont

import (
	"regexp"
	"strconv"
	"strings"
)

// ParsedFontFace represents the properties extracted from a @font-face block.
type ParsedFontFace struct {
	FontFamily string
	FontStyle  string
	FontWeight int
	Src        string
}

var (
	fontFaceRegex   = regexp.MustCompile(`@font-face\s*\{([^}]+)\}`)
	fontFamilyRegex = regexp.MustCompile(`font-family:\s*['"]?([^'";]+)['"]?`)
	fontStyleRegex  = regexp.MustCompile(`font-style:\s*([^;]+)`)
	fontWeightRegex = regexp.MustCompile(`font-weight:\s*(\d+)`)
	srcRegex        = regexp.MustCompile(`src:\s*url\(([^)]+)\)`)
)

// ParseGoogleFontsCss parses the CSS returned by Google Fonts and extracts @font-face blocks.
func ParseGoogleFontsCss(css string) []ParsedFontFace {
	var results []ParsedFontFace

	matches := fontFaceRegex.FindAllStringSubmatch(css, -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		block := match[1]

		parsed := ParsedFontFace{
			FontWeight: 400, // Default normal
			FontStyle:  "normal",
		}

		if fm := fontFamilyRegex.FindStringSubmatch(block); len(fm) >= 2 {
			parsed.FontFamily = strings.TrimSpace(fm[1])
		}

		if fs := fontStyleRegex.FindStringSubmatch(block); len(fs) >= 2 {
			parsed.FontStyle = strings.TrimSpace(fs[1])
		}

		if fw := fontWeightRegex.FindStringSubmatch(block); len(fw) >= 2 {
			if w, err := strconv.Atoi(fw[1]); err == nil {
				parsed.FontWeight = w
			}
		}

		if src := srcRegex.FindStringSubmatch(block); len(src) >= 2 {
			parsed.Src = strings.TrimSpace(src[1])
		}

		results = append(results, parsed)
	}

	return results
}
