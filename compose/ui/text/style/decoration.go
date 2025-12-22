package style

import (
	"fmt"
	"strings"
)

type TextDecorationMask int

const (
	TextDecorationMaskNone        TextDecorationMask = 0x0
	TextDecorationMaskUnderline   TextDecorationMask = 0x1
	TextDecorationMaskLineThrough TextDecorationMask = 0x2
)

// TextDecoration defines a horizontal line to be drawn on the text.
type TextDecoration struct {
	mask TextDecorationMask
}

var (
	// TextDecorationNone Defines a horizontal line to be drawn on the text.
	TextDecorationNone = TextDecoration{mask: TextDecorationMaskNone}

	// TextDecorationUnderline Draws a horizontal line below the text.
	TextDecorationUnderline = TextDecoration{mask: TextDecorationMaskUnderline}

	// TextDecorationLineThrough Draws a horizontal line over the text.
	TextDecorationLineThrough = TextDecoration{mask: TextDecorationMaskLineThrough}
)

// Combine creates a decoration that includes all the given decorations.
func Combine(decorations []TextDecoration) TextDecoration {
	mask := TextDecorationMaskNone
	for _, decoration := range decorations {
		mask |= decoration.mask
	}
	return TextDecoration{mask: mask}
}

// Plus creates a decoration that includes both of the TextDecorations.
func (t TextDecoration) Plus(decoration TextDecoration) TextDecoration {
	return TextDecoration{mask: t.mask | decoration.mask}
}

// Contains checks whether this TextDecoration contains the given decoration.
func (t TextDecoration) Contains(other TextDecoration) bool {
	return (t.mask | other.mask) == t.mask
}

// String returns a string representation of the TextDecoration.
func (t TextDecoration) String() string {
	if t.mask == 0 {
		return "TextDecoration.None"
	}

	var values []string
	if (t.mask & TextDecorationUnderline.mask) != 0 {
		values = append(values, "Underline")
	}
	if (t.mask & TextDecorationLineThrough.mask) != 0 {
		values = append(values, "LineThrough")
	}

	if len(values) == 1 {
		return "TextDecoration." + values[0]
	}
	return "TextDecoration[" + strings.Join(values, ", ") + "]"
}

// NewTextDecoration constructs a TextDecoration instance from the underlying mask.
// This method ensures the mask is valid.
func NewTextDecoration(mask TextDecorationMask) TextDecoration {
	// Prevent creating an invalid TextDecoration combination.
	// The original Kotlin code checks (mask | 0b11) == 0b11.
	// 0b11 is 3. The valid masks are 0, 1, 2, 3.
	// If mask has bits other than 0 and 1 set, (mask | 3) will be > 3 (or rather have bits set outside last 2).
	// Actually the check `(mask | 0b11) == 0b11` in Kotlin verifies that NO bits outside of 0b11 are set?
	// No, `mask | 0b11` will always have the last two bits set.
	// If mask has other bits set, say 0b100 (4), then 0b100 | 0b011 = 0b111 (7).
	// 7 != 3. So yes, it checks that only the last 2 bits are used.
	if (mask | 0b11) != 0b11 {
		panic(fmt.Sprintf("The given mask=%d is not recognized by TextDecoration.", mask))
	}

	return TextDecoration{mask: mask}
}
