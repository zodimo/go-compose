//go:build symbols_extra_fonts

package fonts

import "embed"

// The rounded and sharp Material Symbols variable faces are only embedded when
// this build tag is set. Each is 9-15MB, so they are opt-in rather than paid for
// by every binary that imports this package.
//
// Build with: go build -tags symbols_extra_fonts ./...
//
//go:embed MaterialSymbolsRounded.ttf
//go:embed MaterialSymbolsSharp.ttf
var extraSymbolFonts embed.FS

// extraSymbolTTF returns the embedded bytes for a non-outlined style.
func extraSymbolTTF(style SymbolFace) ([]byte, bool) {
	name, ok := symbolAssets[style]
	if !ok {
		return nil, false
	}
	data, err := extraSymbolFonts.ReadFile(name)
	if err != nil {
		return nil, false
	}
	return data, true
}
