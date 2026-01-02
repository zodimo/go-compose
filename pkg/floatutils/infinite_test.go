package floatutils

import (
	"math"
	"testing"
)

func TestFloat32Infinite(t *testing.T) {
	if !math.IsInf(float64(Float32Infinite), 1) {
		t.Errorf("Float32Infinite should be positive infinity, got %v", Float32Infinite)
	}
}

func TestFloatInfinite(t *testing.T) {
	if !math.IsInf(FloatInfinite, 1) {
		t.Errorf("FloatInfinite should be positive infinity, got %v", FloatInfinite)
	}
}
