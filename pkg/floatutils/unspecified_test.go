package floatutils

import (
	"testing"
)

func TestFloat64Unspecified(t *testing.T) {
	if IsSpecified(Float64Unspecified) {
		t.Errorf("Float64Unspecified should be unspecified, got %v", Float64Unspecified)
	}
}

func TestFloat32Unspecified(t *testing.T) {
	if IsSpecified(Float32Unspecified) {
		t.Errorf("Float32Unspecified should be unspecified, got %v", Float32Unspecified)
	}
}

func TestTakeOrElse(t *testing.T) {
	tests := []struct {
		name         string
		value        float64
		defaultValue float64
		want         float64
	}{
		{"specified value returned", 5.0, 10.0, 5.0},
		{"unspecified float64 returns default", Float64Unspecified, 10.0, 10.0},
		{"unspecified float32 returns default", float64(Float32Unspecified), 10.0, 10.0},

		{"zero is specified", 0.0, 10.0, 0.0},
		{"negative is specified", -5.0, 10.0, -5.0},
		{"infinity is specified", FloatInfinite, 10.0, FloatInfinite},
		{"negative infinity is specified", -FloatInfinite, 10.0, -FloatInfinite},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TakeOrElse(tt.value, tt.defaultValue)
			if IsUnspecified(tt.want) {
				if IsSpecified(got) {
					t.Errorf("TakeOrElse(%v, %v) = %v, want unspecified", tt.value, tt.defaultValue, got)
				}
			} else if got != tt.want {
				t.Errorf("TakeOrElse(%v, %v) = %v, want %v", tt.value, tt.defaultValue, got, tt.want)
			}
		})
	}
}

func TestTakeOrElse_Float32(t *testing.T) {
	tests := []struct {
		name         string
		value        float32
		defaultValue float32
		want         float32
	}{
		{"specified value returned", 5.0, 10.0, 5.0},
		{"unspecified returns default", Float32Unspecified, 10.0, 10.0},
		{"zero is specified", 0.0, 10.0, 0.0},
		{"negative is specified", -5.0, 10.0, -5.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TakeOrElse(tt.value, tt.defaultValue)
			if IsUnspecified(tt.want) {
				if IsSpecified(got) {
					t.Errorf("TakeOrElse(%v, %v) = %v, want unspecified", tt.value, tt.defaultValue, got)
				}
			} else if got != tt.want {
				t.Errorf("TakeOrElse(%v, %v) = %v, want %v", tt.value, tt.defaultValue, got, tt.want)
			}
		})
	}
}
