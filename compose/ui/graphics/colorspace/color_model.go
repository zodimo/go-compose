package colorspace

import (
	"github.com/zodimo/go-compose/compose/ui/util"
)

// ColorModel is a value class in Kotlin, ported as a struct wrapping a packed int64.
// A color model is required by a [ColorSpace] to describe the way colors can be represented as
// tuples of numbers. A common color model is the [RGB][Rgb] color model which defines a color as
// represented by a tuple of 3 numbers (red, green and blue).
type ColorModel struct {
	// pack both the number of components and an ordinal value to distinguish between different
	// ColorModel types that have the same number of components
	PackedValue int64
}

// ComponentCount returns the number of components for this color model.
// An integer between 1 and 4.
func (cm ColorModel) ComponentCount() int {
	return int(util.UnpackInt1(cm.PackedValue))
}

var (
	// Rgb is a color model with 3 components that refer to the three additive
	// primaries: red, green and blue.
	ColorModelRgb = ColorModel{util.PackInts(3, 0)}

	// Xyz is a color model with 3 components that are used to model human color
	// vision on a basic sensory level.
	ColorModelXyz = ColorModel{util.PackInts(3, 1)}

	// Lab is a color model with 3 components used to describe a color space that is
	// more perceptually uniform than XYZ.
	ColorModelLab = ColorModel{util.PackInts(3, 2)}

	// Cmyk is a color model with 4 components that refer to four inks used in color
	// printing: cyan, magenta, yellow and black (or key). CMYK is a subtractive color model.
	ColorModelCmyk = ColorModel{util.PackInts(4, 3)}
)

func (cm ColorModel) String() string {
	switch cm {
	case ColorModelRgb:
		return "Rgb"
	case ColorModelXyz:
		return "Xyz"
	case ColorModelLab:
		return "Lab"
	case ColorModelCmyk:
		return "Cmyk"
	default:
		return "Unknown"
	}
}
