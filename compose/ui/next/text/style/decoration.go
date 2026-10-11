package style

import (
	"fmt"
	"strings"

	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
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

	// textDecorationUnspecified is the sentinel for TextDecoration. It is
	// deliberately package-private: the only public way to obtain it is
	// TextDecorationCompanion.Unspecified().
	textDecorationUnspecified = &TextDecoration{mask: TextDecorationMaskUnspecified}
)

// Combine creates a decoration that includes all the given decorations.
func Combine(decorations []*TextDecoration) *TextDecoration {
	mask := TextDecorationMaskNone
	for _, decoration := range decorations {
		decoration = coalesceTextDecoration(decoration, textDecorationUnspecified)
		if IsSpecifiedTextDecoration(decoration) {
			mask |= decoration.mask
		}
	}
	return &TextDecoration{mask: mask}
}

// Plus creates a decoration that includes both of the TextDecorations.
func (t TextDecoration) Plus(decoration *TextDecoration) *TextDecoration {
	decoration = coalesceTextDecoration(decoration, textDecorationUnspecified)

	return &TextDecoration{mask: t.mask | decoration.mask}
}

// Contains checks whether this TextDecoration contains the given decoration.
func (t TextDecoration) Contains(other *TextDecoration) bool {
	other = coalesceTextDecoration(other, textDecorationUnspecified)

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
	return isSpecifiedTextDecoration(t)
}

// === Sentinel operations ===

func isSpecifiedTextDecoration(t *TextDecoration) bool {
	return t != nil && t != textDecorationUnspecified
}

func isUnspecifiedTextDecoration(t *TextDecoration) bool {
	return !isSpecifiedTextDecoration(t)
}

func takeOrElseTextDecoration(s, def *TextDecoration) *TextDecoration {
	if !isSpecifiedTextDecoration(s) {
		return def
	}
	return s
}

func coalesceTextDecoration(ptr, def *TextDecoration) *TextDecoration {
	if ptr == nil {
		return def
	}
	return ptr
}

func mergeTextDecoration(a, b *TextDecoration) *TextDecoration {
	a = coalesceTextDecoration(a, textDecorationUnspecified)
	b = coalesceTextDecoration(b, textDecorationUnspecified)

	if a == textDecorationUnspecified {
		return b
	}
	if b == textDecorationUnspecified {
		return a
	}

	// Both are custom: allocate new merged style. b wins where specified.
	return &TextDecoration{
		mask: a.mask | b.mask,
	}
}

func stringTextDecoration(t *TextDecoration) string {
	return StringTextDecoration(t)
}

func sameTextDecoration(a, b *TextDecoration) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == textDecorationUnspecified
	}
	if b == nil {
		return a == textDecorationUnspecified
	}
	return a == b
}

func semanticEqualTextDecoration(a, b *TextDecoration) bool {

	a = coalesceTextDecoration(a, textDecorationUnspecified)
	b = coalesceTextDecoration(b, textDecorationUnspecified)

	if a == textDecorationUnspecified && b == textDecorationUnspecified {
		return true
	}

	return a.mask == b.mask
}

func equalTextDecoration(a, b *TextDecoration) bool {
	if !sameTextDecoration(a, b) {
		return semanticEqualTextDecoration(a, b)
	}
	return true
}

func copyTextDecoration(t *TextDecoration, options ...func(*TextDecoration)) *TextDecoration {
	cp := &TextDecoration{}
	if isSpecifiedTextDecoration(t) {
		cp.mask = t.mask
	}
	for _, option := range options {
		option(cp)
	}
	return cp
}

// TextDecorationCompanion is the single exported sentinel-family symbol for TextDecoration.
var TextDecorationCompanion = companions.NewStructCompanion(
	func() *TextDecoration { return textDecorationUnspecified },
	stringTextDecoration,
	isSpecifiedTextDecoration,
	isUnspecifiedTextDecoration,
	takeOrElseTextDecoration,
	mergeTextDecoration,
	sameTextDecoration,
	semanticEqualTextDecoration,
	equalTextDecoration,
	coalesceTextDecoration,
	copyTextDecoration,
)

func StringTextDecoration(t *TextDecoration) string {
	t = coalesceTextDecoration(t, textDecorationUnspecified)

	if !IsSpecifiedTextDecoration(t) {
		return "textDecorationUnspecified"
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
	return takeOrElseTextDecoration(s, def)
}

// Identity (2 ns)
func SameTextDecoration(a, b *TextDecoration) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == textDecorationUnspecified
	}
	if b == nil {
		return a == textDecorationUnspecified
	}
	return a == b
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualTextDecoration(a, b *TextDecoration) bool {

	a = coalesceTextDecoration(a, textDecorationUnspecified)
	b = coalesceTextDecoration(b, textDecorationUnspecified)

	return a.mask == b.mask
}

func EqualTextDecoration(a, b *TextDecoration) bool {
	if !SameTextDecoration(a, b) {
		return SemanticEqualTextDecoration(a, b)
	}
	return true
}

func MergeTextDecoration(a, b *TextDecoration) *TextDecoration {
	a = coalesceTextDecoration(a, textDecorationUnspecified)
	b = coalesceTextDecoration(b, textDecorationUnspecified)

	if a == textDecorationUnspecified {
		return b
	}
	if b == textDecorationUnspecified {
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
