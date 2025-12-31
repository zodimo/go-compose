package style

import (
	"fmt"
	"strings"
)

type TextDecorationMask int

const (
	TextDecorationMaskUnspecified TextDecorationMask = -1
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
	TextDecorationNone = &TextDecoration{mask: TextDecorationMaskNone}

	// TextDecorationUnderline Draws a horizontal line below the text.
	TextDecorationUnderline = &TextDecoration{mask: TextDecorationMaskUnderline}

	// TextDecorationLineThrough Draws a horizontal line over the text.
	TextDecorationLineThrough = &TextDecoration{mask: TextDecorationMaskLineThrough}

	TextDecorationUnspecified = &TextDecoration{mask: TextDecorationMaskUnspecified}
)

// Combine creates a decoration that includes all the given decorations.
func Combine(decorations []*TextDecoration) *TextDecoration {
	mask := TextDecorationMaskNone
	for _, decoration := range decorations {
		decoration = CoalesceTextDecoration(decoration, TextDecorationUnspecified)
		if IsSpecifiedTextDecoration(decoration) {
			mask |= decoration.mask
		}
	}
	return &TextDecoration{mask: mask}
}

// Plus creates a decoration that includes both of the TextDecorations.
func (t TextDecoration) Plus(decoration *TextDecoration) *TextDecoration {
	decoration = CoalesceTextDecoration(decoration, TextDecorationUnspecified)

	return &TextDecoration{mask: t.mask | decoration.mask}
}

// Contains checks whether this TextDecoration contains the given decoration.
func (t TextDecoration) Contains(other *TextDecoration) bool {
	other = CoalesceTextDecoration(other, TextDecorationUnspecified)

	return (t.mask | other.mask) == t.mask
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

func IsSpecifiedTextDecoration(t *TextDecoration) bool {
	return t != nil && t != TextDecorationUnspecified
}

func StringTextDecoration(t *TextDecoration) string {
	t = CoalesceTextDecoration(t, TextDecorationUnspecified)

	if !IsSpecifiedTextDecoration(t) {
		return "TextDecorationUnspecified"
	}
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

func TakeOrElseTextDecoration(s, def *TextDecoration) *TextDecoration {
	if !IsSpecifiedTextDecoration(s) {
		return def
	}
	return s
}

// Identity (2 ns)
func SameTextDecoration(a, b *TextDecoration) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == TextDecorationUnspecified
	}
	if b == nil {
		return a == TextDecorationUnspecified
	}
	return a == b
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualTextDecoration(a, b *TextDecoration) bool {

	a = CoalesceTextDecoration(a, TextDecorationUnspecified)
	b = CoalesceTextDecoration(b, TextDecorationUnspecified)

	return a.mask == b.mask
}

func EqualTextDecoration(a, b *TextDecoration) bool {
	if !SameTextDecoration(a, b) {
		return SemanticEqualTextDecoration(a, b)
	}
	return true
}

func MergeTextDecoration(a, b *TextDecoration) *TextDecoration {
	a = CoalesceTextDecoration(a, TextDecorationUnspecified)
	b = CoalesceTextDecoration(b, TextDecorationUnspecified)

	if a == TextDecorationUnspecified {
		return b
	}
	if b == TextDecorationUnspecified {
		return a
	}

	// Both are custom: allocate new merged style
	return &TextDecoration{
		mask: a.mask | b.mask,
	}
}

func CoalesceTextDecoration(ptr, def *TextDecoration) *TextDecoration {
	if ptr == nil {
		return def
	}
	return ptr
}
