package floatutils

import (
	"math"
	"testing"
)

func TestFloat32Equals(t *testing.T) {
	tests := []struct {
		name    string
		a       float32
		b       float32
		epsilon float32
		want    bool
	}{
		{"exact equal", 1.0, 1.0, Float32EqualityThreshold, true},
		{"within epsilon", 1.0, 1.0 + 1e-7, Float32EqualityThreshold, true},
		{"outside epsilon", 1.0, 1.0 + 1e-5, Float32EqualityThreshold, false},
		{"negative values equal", -1.0, -1.0, Float32EqualityThreshold, true},
		{"zero values", 0.0, 0.0, Float32EqualityThreshold, true},
		{"zero and small positive", 0.0, 1e-7, Float32EqualityThreshold, true},
		{"zero and small negative", 0.0, -1e-7, Float32EqualityThreshold, true},
		{"large values within epsilon", 1000000.0, 1000000.0 + 1e-7, Float32EqualityThreshold, true},
		{"custom epsilon small", 1.0, 1.1, 0.01, false},
		{"custom epsilon large", 1.0, 1.1, 0.2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float32Equals(tt.a, tt.b, tt.epsilon); got != tt.want {
				t.Errorf("Float32Equals(%v, %v, %v) = %v, want %v", tt.a, tt.b, tt.epsilon, got, tt.want)
			}
		})
	}
}

func TestFloat64Equals(t *testing.T) {
	tests := []struct {
		name    string
		a       float64
		b       float64
		epsilon float64
		want    bool
	}{
		{"exact equal", 1.0, 1.0, Float64EqualityThreshold, true},
		{"within epsilon", 1.0, 1.0 + 1e-10, Float64EqualityThreshold, true},
		{"outside epsilon", 1.0, 1.0 + 1e-8, Float64EqualityThreshold, false},
		{"negative values equal", -1.0, -1.0, Float64EqualityThreshold, true},
		{"zero values", 0.0, 0.0, Float64EqualityThreshold, true},
		{"zero and small positive", 0.0, 1e-10, Float64EqualityThreshold, true},
		{"precise comparison", 0.1 + 0.2, 0.3, Float64EqualityThreshold, true},
		{"custom epsilon small", 1.0, 1.1, 0.01, false},
		{"custom epsilon large", 1.0, 1.1, 0.2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float64Equals(tt.a, tt.b, tt.epsilon); got != tt.want {
				t.Errorf("Float64Equals(%v, %v, %v) = %v, want %v", tt.a, tt.b, tt.epsilon, got, tt.want)
			}
		})
	}
}

func TestIsInfinite(t *testing.T) {
	tests := []struct {
		name string
		f    float64
		want bool
	}{
		{"positive infinity", math.Inf(1), true},
		{"negative infinity", math.Inf(-1), true},
		{"positive number", 1.0, false},
		{"negative number", -1.0, false},
		{"zero", 0.0, false},
		{"NaN", math.NaN(), false},
		{"max float64", math.MaxFloat64, false},
		{"smallest positive", math.SmallestNonzeroFloat64, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInfinite(tt.f); got != tt.want {
				t.Errorf("IsInfinite(%v) = %v, want %v", tt.f, got, tt.want)
			}
		})
	}
}

func TestIsInfinite_Float32(t *testing.T) {
	tests := []struct {
		name string
		f    float32
		want bool
	}{
		{"positive infinity", float32(math.Inf(1)), true},
		{"negative infinity", float32(math.Inf(-1)), true},
		{"positive number", float32(1.0), false},
		{"NaN", float32(math.NaN()), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInfinite(tt.f); got != tt.want {
				t.Errorf("IsInfinite(%v) = %v, want %v", tt.f, got, tt.want)
			}
		})
	}
}

func TestIsSpecified(t *testing.T) {
	tests := []struct {
		name string
		f    float64
		want bool
	}{
		{"positive number", 1.0, true},
		{"negative number", -1.0, true},
		{"zero", 0.0, true},
		{"positive infinity", math.Inf(1), true},
		{"negative infinity", math.Inf(-1), true},
		{"NaN", math.NaN(), false},
		{"max float64", math.MaxFloat64, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSpecified(tt.f); got != tt.want {
				t.Errorf("IsSpecified(%v) = %v, want %v", tt.f, got, tt.want)
			}
		})
	}
}

func TestIsSpecified_Float32(t *testing.T) {
	tests := []struct {
		name string
		f    float32
		want bool
	}{
		{"positive number", float32(1.0), true},
		{"NaN", float32(math.NaN()), false},
		{"infinity", float32(math.Inf(1)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSpecified(tt.f); got != tt.want {
				t.Errorf("IsSpecified(%v) = %v, want %v", tt.f, got, tt.want)
			}
		})
	}
}

func TestIsUnspecified(t *testing.T) {
	tests := []struct {
		name string
		f    float64
		want bool
	}{
		{"NaN", math.NaN(), true},
		{"positive number", 1.0, false},
		{"negative number", -1.0, false},
		{"zero", 0.0, false},
		{"positive infinity", math.Inf(1), false},
		{"negative infinity", math.Inf(-1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUnspecified(tt.f); got != tt.want {
				t.Errorf("IsUnspecified(%v) = %v, want %v", tt.f, got, tt.want)
			}
		})
	}
}

func TestIsNaN(t *testing.T) {
	tests := []struct {
		name string
		f    float64
		want bool
	}{
		{"NaN", math.NaN(), true},
		{"positive number", 1.0, false},
		{"zero", 0.0, false},
		{"infinity", math.Inf(1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNaN(tt.f); got != tt.want {
				t.Errorf("IsNaN(%v) = %v, want %v", tt.f, got, tt.want)
			}
		})
	}
}
