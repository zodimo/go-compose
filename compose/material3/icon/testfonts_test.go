package icon

import (
	"testing"

	"github.com/zodimo/go-compose/assets/fonts"
)

// mustFontBytes returns the embedded outlined TTF for tests that need a real
// parseable font.
func mustFontBytes(t *testing.T) []byte {
	t.Helper()
	data, ok := fonts.SymbolTTF(fonts.SymbolFaceOutlined)
	if !ok || len(data) == 0 {
		t.Fatal("no embedded outlined Material Symbols font")
	}
	return data
}
