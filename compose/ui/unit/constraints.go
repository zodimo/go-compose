package unit

import (
	"fmt"
	"math"
)

// Constraints represents immutable constraints for measuring layouts.
// It packs minWidth, maxWidth, minHeight, maxHeight, and a 2-bit focus indicator into a uint64.
// The bit allocation is dynamic: one dimension gets 13-18 bits, the other gets 13-16 bits.
//
// Focus modes (2 bits):
//   - 0: MaxFocusHeight - height gets 18 bits, width gets 13 bits
//   - 1: MinFocusHeight - height gets 16 bits, width gets 15 bits
//   - 2: MinFocusWidth  - width gets 16 bits, height gets 15 bits
//   - 3: MaxFocusWidth  - width gets 18 bits, height gets 13 bits
type Constraints uint64

// Infinity represents unbounded constraints (max value).
const Infinity = math.MaxInt32

const (
	focusMask uint64 = 0x3

	// Bit allocation constants
	minFocusBits    = 16
	minNonFocusBits = 15
	maxFocusBits    = 18
	maxNonFocusBits = 13

	minFocusMask    = 0xFFFF  // 2^16 - 1
	minNonFocusMask = 0x7FFF  // 2^15 - 1
	maxFocusMask    = 0x3FFFF // 2^18 - 1
	maxNonFocusMask = 0x1FFF  // 2^13 - 1

	// Max values for the *other* dimension
	maxAllowedForMinFocusBits    = (1 << (31 - minFocusBits)) - 2    // 32766
	maxAllowedForMinNonFocusBits = (1 << (31 - minNonFocusBits)) - 2 // 65534
	maxAllowedForMaxFocusBits    = (1 << (31 - maxFocusBits)) - 2    // 8190
	maxAllowedForMaxNonFocusBits = (1 << (31 - maxNonFocusBits)) - 2 // 262142

	// Mask to preserve focus bits and max dimensions, zero min dimensions
	maxDimensionsAndFocusMask uint64 = 0xFFFFFFFE00000003
)

// NewConstraints creates a new Constraints object with validation.
func NewConstraints(minWidth, maxWidth, minHeight, maxHeight int) Constraints {
	if !(maxWidth >= minWidth && maxHeight >= minHeight && minWidth >= 0 && minHeight >= 0) {
		panic(fmt.Sprintf("Invalid constraints: w(%d-%d), h(%d-%d)", minWidth, maxWidth, minHeight, maxHeight))
	}
	return createConstraints(minWidth, maxWidth, minHeight, maxHeight)
}

// Fixed creates constraints for a fixed size in both dimensions.
func Fixed(width, height int) Constraints {
	if width < 0 || height < 0 {
		panic("width and height must be >= 0")
	}
	return createConstraints(width, width, height, height)
}

// FixedWidth creates constraints with fixed width and unbounded height.
func FixedWidth(width int) Constraints {
	if width < 0 {
		panic("width must be >= 0")
	}
	return createConstraints(width, width, 0, Infinity)
}

// FixedHeight creates constraints with fixed height and unbounded width.
func FixedHeight(height int) Constraints {
	if height < 0 {
		panic("height must be >= 0")
	}
	return createConstraints(0, Infinity, height, height)
}

// --- Getters ---

// focusIndex returns the 2-bit focus mode.
func (c Constraints) focusIndex() int {
	return int(uint64(c) & focusMask)
}

// MinWidth returns the minimum width in pixels.
func (c Constraints) MinWidth() int {
	mask := widthMask(indexToBitOffset(c.focusIndex()))
	return int((uint64(c) >> 2) & mask)
}

// MaxWidth returns the maximum width in pixels, or Infinity.
func (c Constraints) MaxWidth() int {
	mask := widthMask(indexToBitOffset(c.focusIndex()))
	width := int((uint64(c) >> 33) & mask)
	if width == 0 {
		return Infinity
	}
	return width - 1
}

// MinHeight returns the minimum height in pixels.
func (c Constraints) MinHeight() int {
	bitOffset := indexToBitOffset(c.focusIndex())
	mask := heightMask(bitOffset)
	offset := minHeightOffsets(bitOffset)
	return int((uint64(c) >> offset) & mask)
}

// MaxHeight returns the maximum height in pixels, or Infinity.
func (c Constraints) MaxHeight() int {
	bitOffset := indexToBitOffset(c.focusIndex())
	mask := heightMask(bitOffset)
	offset := minHeightOffsets(bitOffset) + 31
	height := int((uint64(c) >> offset) & mask)
	if height == 0 {
		return Infinity
	}
	return height - 1
}

// --- Properties ---

// HasBoundedWidth returns false if maxWidth is Infinity.
func (c Constraints) HasBoundedWidth() bool {
	mask := widthMask(indexToBitOffset(c.focusIndex()))
	return (int(uint64(c)>>33) & int(mask)) != 0
}

// HasBoundedHeight returns false if maxHeight is Infinity.
func (c Constraints) HasBoundedHeight() bool {
	bitOffset := indexToBitOffset(c.focusIndex())
	mask := heightMask(bitOffset)
	offset := minHeightOffsets(bitOffset) + 31
	return (int(uint64(c)>>offset) & int(mask)) != 0
}

// HasFixedWidth returns true if there's exactly one valid width.
func (c Constraints) HasFixedWidth() bool {
	return c.MinWidth() == c.MaxWidth()
}

// HasFixedHeight returns true if there's exactly one valid height.
func (c Constraints) HasFixedHeight() bool {
	return c.MinHeight() == c.MaxHeight()
}

// IsZero returns true if maxWidth or maxHeight is 0.
func (c Constraints) IsZero() bool {
	return c.MaxWidth() == 0 || c.MaxHeight() == 0
}

// --- Utilities ---

// Copy creates new constraints with optional overrides.
func (c Constraints) Copy(minWidth, maxWidth, minHeight, maxHeight int) Constraints {
	return NewConstraints(minWidth, maxWidth, minHeight, maxHeight)
}

// CopyMaxDimensions returns constraints with min dimensions zeroed, preserving max and focus.
func (c Constraints) CopyMaxDimensions() Constraints {
	return Constraints(uint64(c) & maxDimensionsAndFocusMask)
}

// ConstrainWidth clamps the width to these constraints.
func (c Constraints) ConstrainWidth(width int) int {
	return fastCoerceIn(width, c.MinWidth(), c.MaxWidth())
}

// ConstrainHeight clamps the height to these constraints.
func (c Constraints) ConstrainHeight(height int) int {
	return fastCoerceIn(height, c.MinHeight(), c.MaxHeight())
}

// ConstrainSize clamps an IntSize to these constraints.
func (c Constraints) ConstrainSize(size IntSize) IntSize {
	return IntSize{
		Width:  c.ConstrainWidth(size.Width),
		Height: c.ConstrainHeight(size.Height),
	}
}

// Constrain clamps another Constraints to these constraints.
// The result will satisfy these constraints, but may not satisfy the input constraints.
func (c Constraints) Constrain(other Constraints) Constraints {
	return NewConstraints(
		fastCoerceIn(other.MinWidth(), c.MinWidth(), c.MaxWidth()),
		fastCoerceIn(other.MaxWidth(), c.MinWidth(), c.MaxWidth()),
		fastCoerceIn(other.MinHeight(), c.MinHeight(), c.MaxHeight()),
		fastCoerceIn(other.MaxHeight(), c.MinHeight(), c.MaxHeight()),
	)
}

// IsSatisfiedBy checks if a size satisfies these constraints.
func (c Constraints) IsSatisfiedBy(size IntSize) bool {
	return size.Width >= c.MinWidth() && size.Width <= c.MaxWidth() &&
		size.Height >= c.MinHeight() && size.Height <= c.MaxHeight()
}

// Offset expands constraints by the given deltas.
func (c Constraints) Offset(horizontal, vertical int) Constraints {
	minW := maxInt(0, c.MinWidth()+horizontal)
	maxW := addMaxWithMin(c.MaxWidth(), horizontal)
	minH := maxInt(0, c.MinHeight()+vertical)
	maxH := addMaxWithMin(c.MaxHeight(), vertical)
	return NewConstraints(minW, maxW, minH, maxH)
}

func addMaxWithMin(max, delta int) int {
	if max == Infinity {
		return Infinity
	}
	return maxInt(0, max+delta)
}

// String returns a readable representation.
func (c Constraints) String() string {
	maxW := c.MaxWidth()
	maxWS := "Infinity"
	if maxW != Infinity {
		maxWS = fmt.Sprint(maxW)
	}
	maxH := c.MaxHeight()
	maxHS := "Infinity"
	if maxH != Infinity {
		maxHS = fmt.Sprint(maxH)
	}
	return fmt.Sprintf("Constraints(minWidth=%d, maxWidth=%s, minHeight=%d, maxHeight=%s)",
		c.MinWidth(), maxWS, c.MinHeight(), maxHS)
}

// --- Internal Helpers ---

func createConstraints(minWidth, maxWidth, minHeight, maxHeight int) Constraints {
	// Calculate required bits for each dimension
	heightVal := maxHeight
	if maxHeight == Infinity {
		heightVal = minHeight
	}
	heightBits := bitsNeededForSizeUnchecked(heightVal)

	widthVal := maxWidth
	if maxWidth == Infinity {
		widthVal = minWidth
	}
	widthBits := bitsNeededForSizeUnchecked(widthVal)

	if widthBits+heightBits > 31 {
		panic(fmt.Sprintf("Can't represent width %d and height %d in Constraints", widthVal, heightVal))
	}

	// Branchless conversion of Infinity to 0 for max values
	maxWidthValue := int32(maxWidth) + 1
	maxWidthValue = maxWidthValue & ^(maxWidthValue >> 31)

	maxHeightValue := int32(maxHeight) + 1
	maxHeightValue = maxHeightValue & ^(maxHeightValue >> 31)

	bitOffset := widthBits - 13
	focus := bitOffsetToIndex(bitOffset)

	minHeightOffset := minHeightOffsets(bitOffset)
	maxHeightOffset := minHeightOffset + 31

	value := uint64(focus) |
		(uint64(minWidth) << 2) |
		(uint64(maxWidthValue) << 33) |
		(uint64(minHeight) << minHeightOffset) |
		(uint64(maxHeightValue) << maxHeightOffset)

	return Constraints(value)
}

func bitsNeededForSizeUnchecked(size int) int {
	switch {
	case size < maxNonFocusMask:
		return maxNonFocusBits // 13
	case size < minNonFocusMask:
		return minNonFocusBits // 15
	case size < minFocusMask:
		return minFocusBits // 16
	case size < maxFocusMask:
		return maxFocusBits // 18
	default:
		return 255 // Error value
	}
}

func indexToBitOffset(index int) int {
	// (index & 0x1) << 1 + ((index & 0x2) >> 1) * 3
	return (index&0x1)<<1 + ((index&0x2)>>1)*3
}

func bitOffsetToIndex(bits int) int {
	return (bits >> 1) + (bits & 0x1)
}

func minHeightOffsets(bitOffset int) int {
	return 15 + bitOffset
}

func widthMask(bitOffset int) uint64 {
	return (1 << (13 + bitOffset)) - 1
}

func heightMask(bitOffset int) uint64 {
	return (1 << (18 - bitOffset)) - 1
}

func fastCoerceIn(value, minVal, maxVal int) int {
	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return value
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// FitPrioritizingWidth creates constraints favoring width bit allocation.
func FitPrioritizingWidth(minWidth, maxWidth, minHeight, maxHeight int) Constraints {
	minW := min(minWidth, maxFocusMask-1)
	maxW := minW
	if maxWidth != Infinity {
		maxW = min(maxWidth, maxFocusMask-1)
	}

	consumed := minW
	if maxW != Infinity {
		consumed = maxW
	}
	maxAllowed := maxAllowedForSize(consumed)

	maxH := maxAllowed
	if maxHeight != Infinity {
		maxH = min(maxAllowed, maxHeight)
	}
	minH := min(maxAllowed, minHeight)

	return NewConstraints(minW, maxW, minH, maxH)
}

// FitPrioritizingHeight creates constraints favoring height bit allocation.
func FitPrioritizingHeight(minWidth, maxWidth, minHeight, maxHeight int) Constraints {
	minH := min(minHeight, maxFocusMask-1)
	maxH := minH
	if maxHeight != Infinity {
		maxH = min(maxHeight, maxFocusMask-1)
	}

	consumed := minH
	if maxH != Infinity {
		consumed = maxH
	}
	maxAllowed := maxAllowedForSize(consumed)

	maxW := maxAllowed
	if maxWidth != Infinity {
		maxW = min(maxAllowed, maxWidth)
	}
	minW := min(maxAllowed, minWidth)

	return NewConstraints(minW, maxW, minH, maxH)
}

func maxAllowedForSize(size int) int {
	switch {
	case size < maxNonFocusMask:
		return maxAllowedForMaxNonFocusBits // 262142
	case size < minNonFocusMask:
		return maxAllowedForMinNonFocusBits // 65534
	case size < minFocusMask:
		return maxAllowedForMinFocusBits // 32766
	case size < maxFocusMask:
		return maxAllowedForMaxFocusBits // 8190
	default:
		panic(fmt.Sprintf("Can't represent a size of %d in Constraints", size))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
