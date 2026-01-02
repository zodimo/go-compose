package lerp

import (
	"math"
	"testing"
)

func TestBetween(t *testing.T) {
	tests := []struct {
		name     string
		start    float64
		stop     float64
		fraction float64
		want     float64
	}{
		{"fraction 0 returns start", 0.0, 100.0, 0.0, 0.0},
		{"fraction 1 returns stop", 0.0, 100.0, 1.0, 100.0},
		{"fraction 0.5 returns midpoint", 0.0, 100.0, 0.5, 50.0},
		{"fraction 0.25", 0.0, 100.0, 0.25, 25.0},
		{"negative start", -100.0, 100.0, 0.5, 0.0},
		{"both negative", -100.0, -50.0, 0.5, -75.0},
		{"reverse direction", 100.0, 0.0, 0.5, 50.0},
		{"same start and stop", 50.0, 50.0, 0.5, 50.0},
		{"fraction beyond 1", 0.0, 100.0, 1.5, 150.0},
		{"negative fraction", 0.0, 100.0, -0.5, -50.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Between(tt.start, tt.stop, tt.fraction)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Between(%v, %v, %v) = %v, want %v", tt.start, tt.stop, tt.fraction, got, tt.want)
			}
		})
	}
}

func TestBetween_NaN(t *testing.T) {
	nan := math.NaN()

	// NaN start returns NaN
	result := Between(nan, 100.0, 0.5)
	if !math.IsNaN(result) {
		t.Errorf("Between(NaN, 100.0, 0.5) should return NaN, got %v", result)
	}

	// NaN stop returns NaN
	result = Between(0.0, nan, 0.5)
	if !math.IsNaN(result) {
		t.Errorf("Between(0.0, NaN, 0.5) should return NaN, got %v", result)
	}

	// Both NaN
	result = Between(nan, nan, 0.5)
	if !math.IsNaN(result) {
		t.Errorf("Between(NaN, NaN, 0.5) should return NaN, got %v", result)
	}
}

func TestBetween32(t *testing.T) {
	tests := []struct {
		name     string
		start    float32
		stop     float32
		fraction float32
		want     float32
	}{
		{"fraction 0 returns start", 0.0, 100.0, 0.0, 0.0},
		{"fraction 1 returns stop", 0.0, 100.0, 1.0, 100.0},
		{"fraction 0.5 returns midpoint", 0.0, 100.0, 0.5, 50.0},
		{"negative values", -100.0, 100.0, 0.5, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Between32(tt.start, tt.stop, tt.fraction)
			if math.Abs(float64(got-tt.want)) > 1e-6 {
				t.Errorf("Between32(%v, %v, %v) = %v, want %v", tt.start, tt.stop, tt.fraction, got, tt.want)
			}
		})
	}
}

func TestBetween32_NaN(t *testing.T) {
	nan := float32(math.NaN())

	result := Between32(nan, 100.0, 0.5)
	if !math.IsNaN(float64(result)) {
		t.Errorf("Between32(NaN, 100.0, 0.5) should return NaN, got %v", result)
	}

	result = Between32(0.0, nan, 0.5)
	if !math.IsNaN(float64(result)) {
		t.Errorf("Between32(0.0, NaN, 0.5) should return NaN, got %v", result)
	}
}

func TestFloat32Lerp(t *testing.T) {
	tests := []struct {
		name     string
		start    float32
		stop     float32
		fraction float32
		want     float32
	}{
		{"fraction 0 returns start", 0.0, 100.0, 0.0, 0.0},
		{"fraction 1 returns stop", 0.0, 100.0, 1.0, 100.0},
		{"fraction 0.5 returns midpoint", 0.0, 100.0, 0.5, 50.0},
		{"fraction 0.25", 0.0, 100.0, 0.25, 25.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Float32(tt.start, tt.stop, tt.fraction)
			if math.Abs(float64(got-tt.want)) > 1e-6 {
				t.Errorf("Float32(%v, %v, %v) = %v, want %v", tt.start, tt.stop, tt.fraction, got, tt.want)
			}
		})
	}
}

func TestFloatList32(t *testing.T) {
	t.Run("both nil returns nil", func(t *testing.T) {
		result := FloatList32(nil, nil, 0.5)
		if result != nil {
			t.Errorf("FloatList32(nil, nil, 0.5) should return nil, got %v", result)
		}
	})

	t.Run("a nil returns nil", func(t *testing.T) {
		result := FloatList32(nil, []float32{1, 2, 3}, 0.5)
		if result != nil {
			t.Errorf("FloatList32(nil, [...], 0.5) should return nil, got %v", result)
		}
	})

	t.Run("b nil returns nil", func(t *testing.T) {
		result := FloatList32([]float32{1, 2, 3}, nil, 0.5)
		if result != nil {
			t.Errorf("FloatList32([...], nil, 0.5) should return nil, got %v", result)
		}
	})

	t.Run("equal length lists", func(t *testing.T) {
		a := []float32{0, 10, 20}
		b := []float32{100, 110, 120}
		result := FloatList32(a, b, 0.5)

		expected := []float32{50, 60, 70}
		if len(result) != len(expected) {
			t.Fatalf("Expected length %d, got %d", len(expected), len(result))
		}
		for i, want := range expected {
			if math.Abs(float64(result[i]-want)) > 1e-6 {
				t.Errorf("result[%d] = %v, want %v", i, result[i], want)
			}
		}
	})

	t.Run("different length lists uses max length", func(t *testing.T) {
		a := []float32{0, 10}
		b := []float32{100, 110, 120}
		result := FloatList32(a, b, 0.5)

		if len(result) != 3 {
			t.Fatalf("Expected length 3, got %d", len(result))
		}
		// Third element should lerp between a[1] (10) and b[2] (120)
		expected := Float32(10, 120, 0.5)
		if math.Abs(float64(result[2]-expected)) > 1e-6 {
			t.Errorf("result[2] = %v, want %v", result[2], expected)
		}
	})
}

func TestInt(t *testing.T) {
	tests := []struct {
		name     string
		start    int
		stop     int
		fraction float32
		want     int
	}{
		{"fraction 0", 0, 100, 0.0, 0},
		{"fraction 1", 0, 100, 1.0, 100},
		{"fraction 0.5", 0, 100, 0.5, 50},
		{"fraction 0.25", 0, 100, 0.25, 25},
		{"truncates down", 0, 10, 0.55, 5}, // 5.5 truncates to 5
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Int(tt.start, tt.stop, tt.fraction)
			if got != tt.want {
				t.Errorf("Int(%v, %v, %v) = %v, want %v", tt.start, tt.stop, tt.fraction, got, tt.want)
			}
		})
	}
}

func TestIntFixed(t *testing.T) {
	tests := []struct {
		name     string
		start    int
		stop     int
		fraction int // 0-256 where 256 = 1.0
		want     int
	}{
		{"fraction 0 (0)", 0, 100, 0, 0},
		{"fraction 256 (1.0)", 0, 100, 256, 100},
		{"fraction 128 (0.5)", 0, 100, 128, 50},
		{"fraction 64 (0.25)", 0, 100, 64, 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IntFixed(tt.start, tt.stop, tt.fraction)
			if got != tt.want {
				t.Errorf("IntFixed(%v, %v, %v) = %v, want %v", tt.start, tt.stop, tt.fraction, got, tt.want)
			}
		})
	}
}

func TestIntPrecise(t *testing.T) {
	tests := []struct {
		name     string
		start    int
		stop     int
		fraction float32
		want     int
	}{
		{"fraction 0", 0, 100, 0.0, 0},
		{"fraction 1", 0, 100, 1.0, 100},
		{"fraction 0.5", 0, 100, 0.5, 50},
		{"rounds correctly", 0, 10, 0.55, 6},      // 5.5 rounds to 6 (nearest)
		{"rounds down below 0.5", 0, 10, 0.44, 4}, // 4.4 rounds to 4
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IntPrecise(tt.start, tt.stop, tt.fraction)
			if got != tt.want {
				t.Errorf("IntPrecise(%v, %v, %v) = %v, want %v", tt.start, tt.stop, tt.fraction, got, tt.want)
			}
		})
	}
}
