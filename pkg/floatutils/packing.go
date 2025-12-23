package floatutils

import "math"

// PackFloats packs two float32 values into a single int64.
// val1 occupies the higher 32 bits, and val2 occupies the lower 32 bits.
func PackFloats(val1, val2 float32) int64 {
	v1 := uint64(math.Float32bits(val1))
	v2 := uint64(math.Float32bits(val2))
	return int64((v1 << 32) | v2)
}

// UnpackFloat1 extracts the float32 value from the higher 32 bits of the packed int64.
func UnpackFloat1(packed int64) float32 {
	return math.Float32frombits(uint32(uint64(packed) >> 32))

}

// UnpackFloat2 extracts the float32 value from the lower 32 bits of the packed int64.
func UnpackFloat2(packed int64) float32 {
	return math.Float32frombits(uint32(uint64(packed) & 0xFFFFFFFF))
}
