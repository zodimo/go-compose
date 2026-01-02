package floatutils

import (
	"math"
	"testing"
)

func TestPackFloats(t *testing.T) {
	tests := []struct {
		name string
		val1 float32
		val2 float32
	}{
		{"positive values", 1.5, 2.5},
		{"zero values", 0.0, 0.0},
		{"negative values", -1.5, -2.5},
		{"mixed signs", -1.0, 1.0},
		{"large values", 1000000.0, 2000000.0},
		{"small values", 0.000001, 0.000002},
		{"infinity", float32(math.Inf(1)), float32(math.Inf(-1))},
		{"max float32", math.MaxFloat32, -math.MaxFloat32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packed := PackFloats(tt.val1, tt.val2)
			got1 := UnpackFloat1(packed)
			got2 := UnpackFloat2(packed)

			if got1 != tt.val1 {
				t.Errorf("UnpackFloat1(PackFloats(%v, %v)) = %v, want %v", tt.val1, tt.val2, got1, tt.val1)
			}
			if got2 != tt.val2 {
				t.Errorf("UnpackFloat2(PackFloats(%v, %v)) = %v, want %v", tt.val1, tt.val2, got2, tt.val2)
			}
		})
	}
}

func TestPackFloats_NaN(t *testing.T) {
	nan := float32(math.NaN())
	normal := float32(1.0)

	// Pack NaN in first position
	packed1 := PackFloats(nan, normal)
	got1 := UnpackFloat1(packed1)
	got2 := UnpackFloat2(packed1)

	if !math.IsNaN(float64(got1)) {
		t.Errorf("UnpackFloat1 should return NaN, got %v", got1)
	}
	if got2 != normal {
		t.Errorf("UnpackFloat2 should return %v, got %v", normal, got2)
	}

	// Pack NaN in second position
	packed2 := PackFloats(normal, nan)
	got3 := UnpackFloat1(packed2)
	got4 := UnpackFloat2(packed2)

	if got3 != normal {
		t.Errorf("UnpackFloat1 should return %v, got %v", normal, got3)
	}
	if !math.IsNaN(float64(got4)) {
		t.Errorf("UnpackFloat2 should return NaN, got %v", got4)
	}
}

func TestUnpackFloat1(t *testing.T) {
	// Specific bit pattern tests
	tests := []struct {
		name   string
		packed int64
		want   float32
	}{
		{"zero packed", 0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnpackFloat1(tt.packed); got != tt.want {
				t.Errorf("UnpackFloat1(%v) = %v, want %v", tt.packed, got, tt.want)
			}
		})
	}
}

func TestUnpackFloat2(t *testing.T) {
	// Specific bit pattern tests
	tests := []struct {
		name   string
		packed int64
		want   float32
	}{
		{"zero packed", 0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnpackFloat2(tt.packed); got != tt.want {
				t.Errorf("UnpackFloat2(%v) = %v, want %v", tt.packed, got, tt.want)
			}
		})
	}
}

func TestPackFloats_RoundTrip(t *testing.T) {
	// Test many random-ish values
	values := []float32{
		0.0, 1.0, -1.0, 0.5, -0.5,
		100.0, -100.0, 0.001, -0.001,
		math.MaxFloat32, -math.MaxFloat32,
		math.SmallestNonzeroFloat32, -math.SmallestNonzeroFloat32,
	}

	for _, v1 := range values {
		for _, v2 := range values {
			packed := PackFloats(v1, v2)
			got1 := UnpackFloat1(packed)
			got2 := UnpackFloat2(packed)

			if got1 != v1 {
				t.Errorf("Round-trip failed: UnpackFloat1(PackFloats(%v, %v)) = %v", v1, v2, got1)
			}
			if got2 != v2 {
				t.Errorf("Round-trip failed: UnpackFloat2(PackFloats(%v, %v)) = %v", v1, v2, got2)
			}
		}
	}
}
