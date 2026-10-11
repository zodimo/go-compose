package font

import (
	"fmt"
	"math"

	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// FontWeight represents the thickness of the glyphs, in a range of [1, 1000].
type FontWeight int

const (
	weightUnspecified = -1

	// fontWeightUnspecified is the sentinel for FontWeight. It is deliberately
	// package-private: the only public way to obtain it is
	// FontWeightCompanion.Unspecified().
	fontWeightUnspecified FontWeight = weightUnspecified
)

// NewFontWeight creates a FontWeight with validation.
// Panics if weight is not in range [1, 1000].
func NewFontWeight(weight int) FontWeight {
	if weight < 1 || weight > 1000 {
		panic(fmt.Sprintf("Font weight can be in range [1, 1000]. Current value: %d", weight))
	}
	return FontWeight(weight)
}

// Weight returns the underlying weight value.
func (w FontWeight) Weight() int {
	return int(w)
}

// Compare compares two FontWeights.
// Returns -1 if w < other, 0 if equal, 1 if w > other.
func (w FontWeight) Compare(other FontWeight) int {
	if !w.IsFontWeight() || !other.IsFontWeight() {
		panic("FontWeight must be specified")
	}
	if w < other {
		return -1
	} else if w > other {
		return 1
	}
	return 0
}

// Equals checks if two FontWeights are equal.
func (w FontWeight) Equals(other FontWeight) bool {
	return w == other
}

// String returns a string representation of the FontWeight.
func (w FontWeight) String() string {
	if !w.IsFontWeight() {
		return "FontWeightUnspecified"
	}
	return fmt.Sprintf("FontWeight(weight=%d)", w)
}

func (w FontWeight) IsFontWeight() bool {
	return isSpecifiedFontWeight(w)
}

// === Sentinel operations ===

func isSpecifiedFontWeight(w FontWeight) bool {
	return w != fontWeightUnspecified
}

func isUnspecifiedFontWeight(w FontWeight) bool {
	return w == fontWeightUnspecified
}

func takeOrElseFontWeight(w, def FontWeight) FontWeight {
	if isSpecifiedFontWeight(w) {
		return w
	}
	return def
}

// mergeFontWeight layers b onto a: a value from b wins only where b is specified.
func mergeFontWeight(a, b FontWeight) FontWeight {
	if isUnspecifiedFontWeight(b) {
		return a
	}
	return b
}

func stringFontWeight(w FontWeight) string {
	if isUnspecifiedFontWeight(w) {
		return "FontWeightUnspecified"
	}
	return fmt.Sprintf("FontWeight(weight=%d)", w)
}

func sameFontWeight(a, b FontWeight) bool {
	return a == b
}

func semanticEqualFontWeight(a, b FontWeight) bool {
	return a == b
}

func equalFontWeight(a, b FontWeight) bool {
	return a == b
}

func copyFontWeight(w FontWeight) FontWeight {
	return w
}

// FontWeightCompanion is the single exported sentinel-family symbol for FontWeight.
var FontWeightCompanion = companions.NewScalarCompanion(
	func() FontWeight { return fontWeightUnspecified },
	stringFontWeight,
	isSpecifiedFontWeight,
	isUnspecifiedFontWeight,
	takeOrElseFontWeight,
	mergeFontWeight,
	sameFontWeight,
	semanticEqualFontWeight,
	equalFontWeight,
	copyFontWeight,
)

// Standard font weight constants
const (
	// FontWeightW100 is the thinnest font weight (Thin)
	FontWeightW100 FontWeight = 100
	// FontWeightW200 is extra light weight
	FontWeightW200 FontWeight = 200
	// FontWeightW300 is light weight
	FontWeightW300 FontWeight = 300
	// FontWeightW400 is normal/regular weight
	FontWeightW400 FontWeight = 400
	// FontWeightW500 is medium weight
	FontWeightW500 FontWeight = 500
	// FontWeightW600 is semi-bold weight
	FontWeightW600 FontWeight = 600
	// FontWeightW700 is bold weight
	FontWeightW700 FontWeight = 700
	// FontWeightW800 is extra-bold weight
	FontWeightW800 FontWeight = 800
	// FontWeightW900 is black (heaviest) weight
	FontWeightW900 FontWeight = 900
)

const (
	// Aliases for standard weights
	FontWeightThin       = FontWeightW100
	FontWeightExtraLight = FontWeightW200
	FontWeightLight      = FontWeightW300
	FontWeightNormal     = FontWeightW400
	FontWeightMedium     = FontWeightW500
	FontWeightSemiBold   = FontWeightW600
	FontWeightBold       = FontWeightW700
	FontWeightExtraBold  = FontWeightW800
	FontWeightBlack      = FontWeightW900
)

// FontWeightValues returns a list of all standard font weights.
func FontWeightValues() []FontWeight {
	return []FontWeight{
		FontWeightW100,
		FontWeightW200,
		FontWeightW300,
		FontWeightW400,
		FontWeightW500,
		FontWeightW600,
		FontWeightW700,
		FontWeightW800,
		FontWeightW900,
	}
}

// LerpFontWeight linearly interpolates between two FontWeights.
// The fraction represents position on the timeline: 0.0 returns start, 1.0 returns stop.
func LerpFontWeight(start, stop FontWeight, fraction float32) FontWeight {
	if !start.IsFontWeight() || !stop.IsFontWeight() {
		panic("FontWeight must be specified")
	}
	weight := lerp.Between32(float32(start), float32(stop), fraction)
	// Coerce to valid range
	intWeight := int(math.Round(float64(weight)))
	if intWeight < 1 {
		intWeight = 1
	} else if intWeight > 1000 {
		intWeight = 1000
	}
	return FontWeight(intWeight)
}
