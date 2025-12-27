package graphics

import (
	"fmt"
	"image/color"

	"github.com/zodimo/go-compose/compose/ui/graphics/colorspace"
	"github.com/zodimo/go-compose/compose/ui/util"
)

// Color is a 64-bit value representing a color.
// The value is packed as follows:
// - For sRGB colors:
//   - Bits 0-31: Unused (or Alpha/Red/Green/Blue for 32-bit int representation, but here it's 64-bit)
//   - Actually, Kotlin docs say:
//     "For sRGB colors, the value is packed as:
//   - Bits 32-63: ARGB (8 bits per component)"
//     "For other color spaces:
//   - Bits 0-15: Red (Float16)
//   - Bits 16-31: Green (Float16)
//   - Bits 32-47: Blue (Float16)
//   - Bits 48-57: Alpha (10 bits)
//   - Bits 58-63: ColorSpace ID (6 bits)"
type Color uint64

const (
	// ColorUnspecified represents an uninitialized or unspecified color.
	ColorUnspecified Color = 0x10

	// Internal constants for packing
	colorSpaceMask  = 0x3F
	colorSpaceShift = 58
	alphaMask       = 0x3FF
	alphaShift      = 48
	minId           = colorspace.MinId
	maxId           = colorspace.MaxId
	srgbId          = colorspace.ColorSpaceSrgb
)

// Standard Colors
const (
	ColorBlack Color = 0xFF00000000000000 // sRGB Black? No, wait.
	// Kotlin: Color(0xFF000000) for Black.
	// In Kotlin, Color(int) creates an sRGB color.
	// sRGB packing in Kotlin: (value.toULong() and 0xFFFFFFFFUL) shl 32
	// So 0xFF000000 -> 0xFF00000000000000.
	// Wait, 0xFF000000 is ARGB. Alpha=255.
	// 0xFF000000 << 32 = 0xFF00000000000000?
	// 0xFF (Alpha) << 24 | 0x00 (R) << 16 ...
	// Yes.
	ColorDarkGray    Color = 0xFF44444400000000
	ColorGray        Color = 0xFF88888800000000
	ColorLightGray   Color = 0xFFCCCCCC00000000
	ColorWhite       Color = 0xFFFFFFFF00000000
	ColorRed         Color = 0xFFFF000000000000
	ColorGreen       Color = 0xFF00FF0000000000
	ColorBlue        Color = 0xFF0000FF00000000
	ColorYellow      Color = 0xFFFFFF0000000000
	ColorCyan        Color = 0xFF00FFFF00000000
	ColorMagenta     Color = 0xFFFF00FF00000000
	ColorTransparent Color = 0x0000000000000000
)

// NewColorSrgb creates a new sRGB Color from 8-bit components.
func NewColorSrgb(r, g, b, a int) Color {
	// Pack ARGB to high 32 bits
	val := (uint32(a&0xFF) << 24) | (uint32(r&0xFF) << 16) | (uint32(g&0xFF) << 8) | (uint32(b & 0xFF))
	return Color(uint64(val) << 32)
}

// NewColor creates a new Color from float components and a ColorSpace.
func NewColor(r, g, b, a float32, space colorspace.ColorSpace) Color {
	if space.Id() == colorspace.ColorSpaceSrgb {
		// Convert float to 8-bit sRGB
		// Assume floats are 0..1? Yes.
		// Pack sRGB
		ir := int(r*255.0 + 0.5)
		ig := int(g*255.0 + 0.5)
		ib := int(b*255.0 + 0.5)
		ia := int(a*255.0 + 0.5)
		return NewColorSrgb(
			util.FastCoerceInInt(ir, 0, 255),
			util.FastCoerceInInt(ig, 0, 255),
			util.FastCoerceInInt(ib, 0, 255),
			util.FastCoerceInInt(ia, 0, 255),
		)
	}

	// Pack other color space
	// 3x Float16 + 10-bit Alpha + 6-bit ID
	r16 := uint64(util.FloatToHalf(r))
	g16 := uint64(util.FloatToHalf(g))
	b16 := uint64(util.FloatToHalf(b))

	// Alpha 10-bit: (a * 1023 + 0.5) clamped
	a10 := int(a*1023.0 + 0.5)
	a10 = util.FastCoerceInInt(a10, 0, 1023)

	sid := space.Id()
	if sid == colorspace.MinId { // Unspecified?
		return ColorUnspecified // Or handle as error?
	}
	id := uint64(sid)

	// Kotlin bit layout for non-sRGB:
	// Bits 48-63: Red (Half-float)
	// Bits 32-47: Green (Half-float)
	// Bits 16-31: Blue (Half-float)
	// Bits 6-15: Alpha (10-bit)
	// Bits 0-5: ColorSpace ID (6-bit)

	val := (r16 << 48) | (g16 << 32) | (b16 << 16) | (uint64(a10) << 6) | id
	return Color(val)
}

// Implement basic getters
func (c Color) Alpha() float32 {
	if c&0x3F == 0 { // Check if sRGB (ID shift? No. 63 is max ID. 0 is sRGB)
		// Wait, for sRGB, the lower 32 bits MUST be 0?
		// "For sRGB colors, the value is packed as: Bits 32-63: ARGB"
		// If bits 0-31 are 0, it MIGHT be sRGB? No.
		// There's a check.
		// Kotlin `isSrgb`: `(value and 0x3FUL) == 0UL` (Checking low 6 bits ID matches 0?)
		// But sRGB packed color has 0 in lower bits?
		// Yes, `(val.toLong() and 0x3F) == 0` check assumes user doesn't put garbage in lower bits
		// when creating sRGB color via `NewColorSrgb`.
		// `NewColorSrgb` shifts by 32, so lower 32 bits are 0.
		// ColorSpace 0 is sRGB.
		// So checking lowest 6 bits being 0 implies sRGB ONLY IF we respect the "unused" bits are 0.
		// Let's assume standard creation.
	}

	// ID check logic from Kotlin
	if (c & 0x3F) == 0 {
		return float32((c>>56)&0xFF) / 255.0
	}
	// Non-sRGB: Alpha is at bits 6-15
	return float32((c>>6)&0x3FF) / 1023.0
}

func (c Color) Red() float32 {
	if (c & 0x3F) == 0 {
		return float32((c>>48)&0xFF) / 255.0
	}
	// Non-sRGB: Red is at bits 48-63
	return util.HalfToFloat(uint16((c >> 48) & 0xFFFF))
}

func (c Color) Green() float32 {
	if (c & 0x3F) == 0 {
		return float32((c>>40)&0xFF) / 255.0
	}
	// Non-sRGB: Green is at bits 32-47
	return util.HalfToFloat(uint16((c >> 32) & 0xFFFF))
}

func (c Color) Blue() float32 {
	if (c & 0x3F) == 0 {
		return float32((c>>32)&0xFF) / 255.0
	}
	// Non-sRGB: Blue is at bits 16-31
	return util.HalfToFloat(uint16((c >> 16) & 0xFFFF))
}

func (c Color) ColorSpaceId() int {
	// ColorSpace ID is at bits 0-5
	return int(c & 0x3F)
}

func (c Color) IsSpecified() bool {
	return c != ColorUnspecified
}

func (c Color) IsUnspecified() bool {
	return c == ColorUnspecified
}

// TakeOrElse returns this color if specified, otherwise returns the block result.
// Direct value version for Go zero-allocation.
func (c Color) TakeOrElse(block Color) Color {
	if c.IsSpecified() {
		return c
	}
	return block
}

// Copy creates a new color with modified components.
func (c Color) Copy(alpha, red, green, blue float32) Color {
	id := c.ColorSpaceId()
	space := colorspace.Get(id)
	return NewColor(red, green, blue, alpha, space)
}

// ColorSpace returns the ColorSpace object for this color.
func (c Color) ColorSpace() colorspace.ColorSpace {
	return colorspace.Get(c.ColorSpaceId())
}

// Convert transforms this color to another color space.
func (c Color) Convert(destColorSpace colorspace.ColorSpace) Color {
	srcSpace := c.ColorSpace()
	if srcSpace.Id() == destColorSpace.Id() {
		return c
	}

	connector := colorspace.NewConnector(srcSpace, destColorSpace, colorspace.RenderIntentPerceptual)
	v := []float32{c.Red(), c.Green(), c.Blue()}
	result := connector.Transform(v)
	return NewColor(result[0], result[1], result[2], c.Alpha(), destColorSpace)
}

// GetComponents returns the color components as [red, green, blue, alpha].
func (c Color) GetComponents() [4]float32 {
	return [4]float32{c.Red(), c.Green(), c.Blue(), c.Alpha()}
}

// Component accessors for destructuring
func (c Color) Component1() float32               { return c.Red() }
func (c Color) Component2() float32               { return c.Green() }
func (c Color) Component3() float32               { return c.Blue() }
func (c Color) Component4() float32               { return c.Alpha() }
func (c Color) Component5() colorspace.ColorSpace { return c.ColorSpace() }

// ToArgb converts to 32-bit ARGB.
// For non-sRGB colors, converts to sRGB first.
func (c Color) ToArgb() uint32 {
	if (c & 0x3F) == 0 {
		return uint32(c >> 32)
	}
	// Convert to sRGB first
	srgbColor := c.Convert(colorspace.Srgb)
	return uint32(srgbColor >> 32)
}

func (c Color) String() string {
	cs := c.ColorSpace()
	return fmt.Sprintf("Color(%f, %f, %f, %f, %s)", c.Red(), c.Green(), c.Blue(), c.Alpha(), cs.Name())
}

// CompositeOver composites this color over the background color.
// The result is in the background's color space.
func (c Color) CompositeOver(background Color) Color {
	// Convert foreground to background's color space
	fg := c.Convert(background.ColorSpace())

	bgA := background.Alpha()
	fgA := fg.Alpha()
	a := fgA + (bgA * (1.0 - fgA))

	if a == 0.0 {
		return ColorTransparent
	}

	r := compositeComponent(fg.Red(), background.Red(), fgA, bgA, a)
	g := compositeComponent(fg.Green(), background.Green(), fgA, bgA, a)
	b := compositeComponent(fg.Blue(), background.Blue(), fgA, bgA, a)

	return UncheckedColor(r, g, b, a, background.ColorSpace())
}

// compositeComponent performs the Porter-Duff 'source over' calculation for a single component.
func compositeComponent(fgC, bgC, fgA, bgA, a float32) float32 {
	if a == 0 {
		return 0
	}
	return ((fgC * fgA) + ((bgC * bgA) * (1.0 - fgA))) / a
}

// Luminance returns the relative luminance of the color.
// Based on WCAG 2.0 formula.
func (c Color) Luminance() float32 {
	cs := c.ColorSpace()

	// Require RGB color model
	if cs.Model() != colorspace.ColorModelRgb {
		panic(fmt.Sprintf("The specified color must be encoded in an RGB color space. The supplied color space is %v", cs.Model()))
	}

	// Get EOTF (Electro-Optical Transfer Function) from RGB color space
	rgb, ok := cs.(*colorspace.Rgb)
	if !ok {
		// Fallback for sRGB
		rgb = colorspace.Srgb
	}

	r := rgb.Eotf(float64(c.Red()))
	g := rgb.Eotf(float64(c.Green()))
	b := rgb.Eotf(float64(c.Blue()))

	lum := float32((0.2126 * r) + (0.7152 * g) + (0.0722 * b))
	return util.FastCoerceIn(lum, 0.0, 1.0)
}

// NewColorLong creates a new sRGB Color from a 32-bit ARGB long.
// Useful for specifying colors with alpha > 0x80 without sign issues.
func NewColorLong(argb int64) Color {
	return Color((uint64(argb) & 0xFFFFFFFF) << 32)
}

// UncheckedColor creates a color without validation.
// Used for performance-critical code like lerp where values are known to be valid.
func UncheckedColor(r, g, b, a float32, space colorspace.ColorSpace) Color {
	if space.IsSrgb() {
		argb := (int(a*255.0+0.5) << 24) |
			(int(r*255.0+0.5) << 16) |
			(int(g*255.0+0.5) << 8) |
			int(b*255.0+0.5)
		return Color(uint64(argb) << 32)
	}

	r16 := uint64(util.FloatToHalf(r))
	g16 := uint64(util.FloatToHalf(g))
	b16 := uint64(util.FloatToHalf(b))

	a10 := int(maxf(0.0, minf(a, 1.0))*1023.0 + 0.5)
	id := uint64(space.Id())

	val := (r16 << 48) | (g16 << 32) | (b16 << 16) | (uint64(a10) << 6) | id
	return Color(val)
}

// Lerp linearly interpolates between two colors.
// Interpolation is done in Oklab color space for perceptually uniform results.
// The result is converted to the stop color's color space.
func Lerp(start, stop Color, fraction float32) Color {
	oklab := colorspace.OklabInstance
	startColor := start.Convert(oklab)
	endColor := stop.Convert(oklab)

	startAlpha := startColor.Alpha()
	startL := startColor.Red()   // L in Oklab
	startA := startColor.Green() // a in Oklab
	startB := startColor.Blue()  // b in Oklab

	endAlpha := endColor.Alpha()
	endL := endColor.Red()
	endA := endColor.Green()
	endB := endColor.Blue()

	// Clamp fraction to avoid out-of-range values from easing curves
	t := util.FastCoerceIn(fraction, 0.0, 1.0)

	interpolated := UncheckedColor(
		lerpFloat(startL, endL, t),
		lerpFloat(startA, endA, t),
		lerpFloat(startB, endB, t),
		lerpFloat(startAlpha, endAlpha, t),
		oklab,
	)
	return interpolated.Convert(stop.ColorSpace())
}

// lerpFloat linearly interpolates between two float32 values.
func lerpFloat(start, stop, fraction float32) float32 {
	return start + (stop-start)*fraction
}

// Hsv creates a color from HSV (Hue, Saturation, Value) representation.
// hue: 0..360, saturation: 0..1, value: 0..1, alpha: 0..1
func Hsv(hue, saturation, value, alpha float32, colorSpace *colorspace.Rgb) Color {
	if hue < 0 || hue > 360 || saturation < 0 || saturation > 1 || value < 0 || value > 1 {
		panic(fmt.Sprintf("HSV (%f, %f, %f) must be in range (0..360, 0..1, 0..1)", hue, saturation, value))
	}
	if colorSpace == nil {
		colorSpace = colorspace.Srgb
	}
	red := hsvToRgbComponent(5, hue, saturation, value)
	green := hsvToRgbComponent(3, hue, saturation, value)
	blue := hsvToRgbComponent(1, hue, saturation, value)
	return NewColor(red, green, blue, alpha, colorSpace)
}

func hsvToRgbComponent(n int, h, s, v float32) float32 {
	k := float32(int(float32(n)+h/60.0) % 6)
	return v - (v * s * maxf(0, minf(k, minf(4-k, 1))))
}

// Hsl creates a color from HSL (Hue, Saturation, Lightness) representation.
// hue: 0..360, saturation: 0..1, lightness: 0..1, alpha: 0..1
func Hsl(hue, saturation, lightness, alpha float32, colorSpace *colorspace.Rgb) Color {
	if hue < 0 || hue > 360 || saturation < 0 || saturation > 1 || lightness < 0 || lightness > 1 {
		panic(fmt.Sprintf("HSL (%f, %f, %f) must be in range (0..360, 0..1, 0..1)", hue, saturation, lightness))
	}
	if colorSpace == nil {
		colorSpace = colorspace.Srgb
	}
	red := hslToRgbComponent(0, hue, saturation, lightness)
	green := hslToRgbComponent(8, hue, saturation, lightness)
	blue := hslToRgbComponent(4, hue, saturation, lightness)
	return NewColor(red, green, blue, alpha, colorSpace)
}

func hslToRgbComponent(n int, h, s, l float32) float32 {
	k := float32(int(float32(n)+h/30.0) % 12)
	a := s * minf(l, 1.0-l)
	return l - a*maxf(-1, minf(k-3, minf(9-k, 1)))
}

// minf returns the minimum of two float32 values.
func minf(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

// maxf returns the maximum of two float32 values.
func maxf(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

// ColorProducer is a functional interface that produces a Color.
// Useful for avoiding boxing in performance-critical code.
type ColorProducer func() Color

func ColorToNRGBA(c Color) color.NRGBA {
	return color.NRGBA{
		R: uint8(c.Red() * 255),
		G: uint8(c.Green() * 255),
		B: uint8(c.Blue() * 255),
		A: uint8(c.Alpha() * 255),
	}
}
