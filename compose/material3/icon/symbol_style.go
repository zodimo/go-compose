package icon

// SymbolStyle selects which Material Symbols typeface an icon is rendered from.
// The three styles share the same codepoints and variable axes; only the
// letterforms differ.
type SymbolStyle int

const (
	// SymbolOutlined is the default style.
	SymbolOutlined SymbolStyle = iota
	// SymbolRounded has rounded terminals and corners.
	SymbolRounded
	// SymbolSharp has square terminals and corners.
	SymbolSharp
)

// String returns the style's name for diagnostics and error messages.
func (s SymbolStyle) String() string {
	switch s {
	case SymbolOutlined:
		return "outlined"
	case SymbolRounded:
		return "rounded"
	case SymbolSharp:
		return "sharp"
	default:
		return "unknown"
	}
}
