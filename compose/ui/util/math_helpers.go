package util

import "math"

// Lerp returns the linear interpolation between start and stop by fraction.
func Lerp(start, stop, fraction float32) float32 {
	return (1-fraction)*start + fraction*stop
}

// FastCoerceIn clamps the value between min and max.
func FastCoerceIn(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// PackInts packs two int32 values into a single int64.
func PackInts(val1, val2 int32) int64 {
	return (int64(val1) << 32) | (int64(val2) & 0xFFFFFFFF)
}

// UnpackInt1 returns the first packed int32 from an int64.
func UnpackInt1(value int64) int32 {
	return int32(value >> 32)
}

// UnpackInt2 returns the second packed int32 from an int64.
func UnpackInt2(value int64) int32 {
	return int32(value)
}

// Pow returns x**y.
func Pow(x, y float32) float32 {
	return float32(math.Pow(float64(x), float64(y)))
}

// Abs returns the absolute value of x.
func Abs(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

// FastCbrt returns the cube root of x.
func FastCbrt(x float32) float32 {
	return float32(math.Cbrt(float64(x)))
}

// PackFloats packs two float32 values into a single int64.
func PackFloats(v1, v2 float32) int64 {
	bits1 := uint64(math.Float32bits(v1))
	bits2 := uint64(math.Float32bits(v2))
	return int64((bits1 << 32) | (bits2 & 0xFFFFFFFF))
}

// FastCoerceInInt clamps integer value between min and max.
func FastCoerceInInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
