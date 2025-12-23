package font

import (
	"fmt"

	"github.com/zodimo/go-compose/pkg/floatutils/lerp"
)

const weightUnspecified = -1

var FontWeightUnspecified = FontWeight{weight: weightUnspecified}

// FontWeight represents the thickness of the glyphs, in a range of [1, 1000].
type FontWeight struct {
	weight int
}

// NewFontWeight creates a FontWeight with validation.
// Panics if weight is not in range [1, 1000].
func NewFontWeight(weight int) FontWeight {
	if weight < 1 || weight > 1000 {
		panic(fmt.Sprintf("Font weight can be in range [1, 1000]. Current value: %d", weight))
	}
	return FontWeight{weight: weight}
}

// Weight returns the underlying weight value.
func (w FontWeight) Weight() int {
	return w.weight
}

// Compare compares two FontWeights.
// Returns -1 if w < other, 0 if equal, 1 if w > other.
func (w FontWeight) Compare(other FontWeight) int {
	if !w.IsSpecified() || !other.IsSpecified() {
		panic("FontWeight must be specified")
	}
	if w.weight < other.weight {
		return -1
	} else if w.weight > other.weight {
		return 1
	}
	return 0
}

// Equals checks if two FontWeights are equal.
func (w FontWeight) Equals(other FontWeight) bool {
	return w.weight == other.weight
}

// HashCode returns a hash code for the FontWeight.
func (w FontWeight) HashCode() int {
	return w.weight
}

// String returns a string representation of the FontWeight.
func (w FontWeight) String() string {
	if !w.IsSpecified() {
		return "FontWeightUnspecified"
	}
	return fmt.Sprintf("FontWeight(weight=%d)", w.weight)
}

func (w FontWeight) IsSpecified() bool {
	return w.weight != weightUnspecified
}

func (w FontWeight) TakeOrElse(other FontWeight) FontWeight {
	if !w.IsSpecified() {
		return w
	}
	return other
}

// Standard font weight constants
var (
	// FontWeightW100 is the thinnest font weight (Thin)
	FontWeightW100 = NewFontWeight(100)
	// FontWeightW200 is extra light weight
	FontWeightW200 = NewFontWeight(200)
	// FontWeightW300 is light weight
	FontWeightW300 = NewFontWeight(300)
	// FontWeightW400 is normal/regular weight
	FontWeightW400 = NewFontWeight(400)
	// FontWeightW500 is medium weight
	FontWeightW500 = NewFontWeight(500)
	// FontWeightW600 is semi-bold weight
	FontWeightW600 = NewFontWeight(600)
	// FontWeightW700 is bold weight
	FontWeightW700 = NewFontWeight(700)
	// FontWeightW800 is extra-bold weight
	FontWeightW800 = NewFontWeight(800)
	// FontWeightW900 is black (heaviest) weight
	FontWeightW900 = NewFontWeight(900)

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
	if !start.IsSpecified() || !stop.IsSpecified() {
		panic("FontWeight must be specified")
	}
	weight := lerp.Between32(float32(start.weight), float32(stop.weight), fraction)
	// Coerce to valid range
	if weight < 1 {
		weight = 1
	} else if weight > 1000 {
		weight = 1000
	}
	return FontWeight{weight: int(weight)}
}
