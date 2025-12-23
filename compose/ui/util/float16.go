package util

import (
	"math"
)

// FloatToHalf converts a float32 to a half-precision float (float16) represented as uint16.
func FloatToHalf(f float32) uint16 {
	bits := math.Float32bits(f)
	s := (bits >> 31) & 0x1
	e := (bits >> 23) & 0xFF
	m := bits & 0x7FFFFF

	var outE uint32
	var outM uint32

	if e == 255 {
		// NaN or Infinity
		outE = 31
		if m != 0 {
			outM = 0x200 // Keep it a NaN, preserving sign not strictly required for color
		} else {
			outM = 0
		}
	} else {
		// Normalized or Denormalized
		eInt := int(e) - 127 + 15
		if eInt >= 31 {
			// Overflow
			outE = 31
			outM = 0
		} else if eInt <= 0 {
			// Underflow
			if eInt < -10 {
				// Too small, becomes zero
				outE = 0
				outM = 0
			} else {
				// Denormalized
				m = (m | 0x800000) >> (1 - eInt)
				outE = 0
				outM = m >> 13
			}
		} else {
			outE = uint32(eInt)
			outM = m >> 13
		}
	}

	return uint16((s << 15) | (outE << 10) | outM)
}

// HalfToFloat converts a half-precision float (float16) represented as uint16 to float32.
func HalfToFloat(h uint16) float32 {
	s := (uint32(h) >> 15) & 0x1
	e := (uint32(h) >> 10) & 0x1F
	m := uint32(h) & 0x3FF

	var outE uint32
	var outM uint32

	if e == 0 {
		if m == 0 {
			// Zero
			outE = 0
			outM = 0
		} else {
			// Denormalized
			// Normalize it
			for (m & 0x400) == 0 {
				m <<= 1
				outE--
			}
			outE += 127 - 15 + 1 // +1 because we shifted explicitly
			outM = (m & 0x3FF) << 13
		}
	} else if e == 31 {
		// NaN or Infinity
		outE = 255
		if m != 0 {
			outM = 0x400000 // NaN
		} else {
			outM = 0 // Infinity
		}
	} else {
		// Normalized
		outE = e + 127 - 15
		outM = m << 13
	}

	bits := (s << 31) | (outE << 23) | outM
	return math.Float32frombits(bits)
}
