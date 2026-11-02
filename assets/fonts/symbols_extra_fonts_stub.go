//go:build !symbols_extra_fonts

package fonts

// extraSymbolTTF reports that no extra style assets are embedded. Without the
// symbols_extra_fonts build tag, only the outlined face is available and the
// icon package falls back to it.
func extraSymbolTTF(SymbolFace) ([]byte, bool) {
	return nil, false
}
