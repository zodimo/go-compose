// Package fonts provides embedded fonts for the go-compose framework.
// This includes both the Go fonts (for text) and Noto Color Emoji (for emoji).
package fonts

import (
	"embed"
	"sync"

	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
)

//go:embed NotoColorEmoji.ttf
var emojiFont embed.FS

//go:embed MaterialSymbolsOutlined.ttf
var materialSymbolsOutlined embed.FS

var (
	cachedCollection []font.FontFace
	collectionOnce   sync.Once
)

// Collection returns a font collection that includes both Go fonts and Noto Color Emoji.
// This provides full unicode and emoji support in a portable, self-contained way
// without relying on system fonts.
func Collection() []font.FontFace {
	collectionOnce.Do(func() {
		// Start with the Go font collection
		cachedCollection = gofont.Collection()

		// Load and add the emoji font
		emojiData, err := emojiFont.ReadFile("NotoColorEmoji.ttf")
		if err == nil {
			// ParseCollection returns []font.FontFace directly
			emojiFaces, err := opentype.ParseCollection(emojiData)
			if err == nil {
				// Append emoji faces to the collection
				cachedCollection = append(cachedCollection, emojiFaces...)
			}
		}

		// Load and add the Material Symbols font
		symbolsData, err := materialSymbolsOutlined.ReadFile("MaterialSymbolsOutlined.ttf")
		if err == nil {
			symbolsFaces, err := opentype.ParseCollection(symbolsData)
			if err == nil {
				cachedCollection = append(cachedCollection, symbolsFaces...)
			}
		}
	})

	return cachedCollection
}

// GoFontsOnly returns just the Go font collection without emoji.
// Use this if you don't need emoji support and want smaller binary size.
func GoFontsOnly() []font.FontFace {
	return gofont.Collection()
}

// SymbolFace identifies a Material Symbols typeface.
// It mirrors icon.SymbolFace so that this package stays dependency-free.
type SymbolFace int

const (
	// SymbolFaceOutlined is the outlined Material Symbols face.
	SymbolFaceOutlined SymbolFace = iota
	// SymbolFaceRounded is the rounded Material Symbols face.
	SymbolFaceRounded
	// SymbolFaceSharp is the sharp Material Symbols face.
	SymbolFaceSharp
)

// symbolAssets maps a style to its embedded file name.
var symbolAssets = map[SymbolFace]string{
	SymbolFaceOutlined: "MaterialSymbolsOutlined.ttf",
	SymbolFaceRounded:  "MaterialSymbolsRounded.ttf",
	SymbolFaceSharp:    "MaterialSymbolsSharp.ttf",
}

// SymbolTTF returns the raw variable TTF bytes for style, and whether they are
// embedded in this build. Only SymbolFaceOutlined is embedded unconditionally; the
// rounded and sharp faces are present when their assets are included, so a
// binary that renders only outlined icons does not pay for the other two.
//
// The bytes are returned as-is (not parsed) because the caller owns a different
// migration path: go-compose's icon package parses them with
// go-text/typesetting to reach the font's variable axes, which the Gio shaper
// cannot express.
func SymbolTTF(style SymbolFace) ([]byte, bool) {
	name, ok := symbolAssets[style]
	if !ok {
		return nil, false
	}

	// Outlined ships in every build. Rounded and Sharp are opt-in through the
	// symbols_extra_fonts build tag: go:embed embeds every file a pattern matches,
	// so merely shipping the .ttf files would add ~24MB to every binary that
	// imports this package, including ones that never render an icon.
	if style == SymbolFaceOutlined {
		data, err := materialSymbolsOutlined.ReadFile(name)
		if err != nil {
			return nil, false
		}
		return data, true
	}

	return extraSymbolTTF(style)
}

// SymbolCollection returns shaper-ready faces for the requested Material Symbols
// styles, for callers that draw symbol ligatures as text through the Gio shaper
// rather than through the icon package's vector path.
//
// It is deliberately separate from Collection: Collection stays outlined-only so
// a binary that never renders rounded or sharp icons does not grow by the ~24MB
// those two variable fonts cost. Asking for a style whose asset is not embedded
// in the build simply omits it — the returned slice never contains a
// placeholder face — so callers that need a guaranteed face should check the
// length or pass SymbolFaceOutlined.
func SymbolCollection(styles ...SymbolFace) []font.FontFace {
	var out []font.FontFace
	seen := make(map[SymbolFace]bool, len(styles))

	for _, style := range styles {
		if seen[style] {
			continue
		}
		seen[style] = true

		data, ok := SymbolTTF(style)
		if !ok {
			continue
		}
		faces, err := opentype.ParseCollection(data)
		if err != nil {
			continue
		}
		out = append(out, faces...)
	}

	return out
}
