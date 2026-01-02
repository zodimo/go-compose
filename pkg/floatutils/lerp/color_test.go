package lerp

import (
	"math"
	"testing"

	"github.com/zodimo/go-compose/pkg/floatutils"
)

type RGBA struct {
	R, G, B, A float32
}

func rgbaEqual(a, b struct{ R, G, B, A float32 }, epsilon float32) bool {
	return math.Abs(float64(a.R-b.R)) <= float64(epsilon) &&
		math.Abs(float64(a.G-b.G)) <= float64(epsilon) &&
		math.Abs(float64(a.B-b.B)) <= float64(epsilon) &&
		math.Abs(float64(a.A-b.A)) <= float64(epsilon)
}

func TestLerpColor(t *testing.T) {
	tests := []struct {
		name string
		a    struct{ R, G, B, A float32 }
		b    struct{ R, G, B, A float32 }
		p    float32
		want struct{ R, G, B, A float32 }
	}{
		{
			name: "fraction 0 returns first color",
			a:    struct{ R, G, B, A float32 }{1.0, 0.0, 0.0, 1.0},
			b:    struct{ R, G, B, A float32 }{0.0, 0.0, 1.0, 1.0},
			p:    0.0,
			want: struct{ R, G, B, A float32 }{1.0, 0.0, 0.0, 1.0},
		},
		{
			name: "fraction 1 returns second color",
			a:    struct{ R, G, B, A float32 }{1.0, 0.0, 0.0, 1.0},
			b:    struct{ R, G, B, A float32 }{0.0, 0.0, 1.0, 1.0},
			p:    1.0,
			want: struct{ R, G, B, A float32 }{0.0, 0.0, 1.0, 1.0},
		},
		{
			name: "fraction 0.5 returns midpoint",
			a:    struct{ R, G, B, A float32 }{1.0, 0.0, 0.0, 1.0},
			b:    struct{ R, G, B, A float32 }{0.0, 0.0, 1.0, 1.0},
			p:    0.5,
			want: struct{ R, G, B, A float32 }{0.5, 0.0, 0.5, 1.0},
		},
		{
			name: "alpha interpolation",
			a:    struct{ R, G, B, A float32 }{1.0, 1.0, 1.0, 0.0},
			b:    struct{ R, G, B, A float32 }{1.0, 1.0, 1.0, 1.0},
			p:    0.5,
			want: struct{ R, G, B, A float32 }{1.0, 1.0, 1.0, 0.5},
		},
		{
			name: "black to white",
			a:    struct{ R, G, B, A float32 }{0.0, 0.0, 0.0, 1.0},
			b:    struct{ R, G, B, A float32 }{1.0, 1.0, 1.0, 1.0},
			p:    0.5,
			want: struct{ R, G, B, A float32 }{0.5, 0.5, 0.5, 1.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LerpColor(tt.a, tt.b, tt.p)
			if !rgbaEqual(got, tt.want, floatutils.Float32EqualityThreshold) {
				t.Errorf("Color(%v, %v, %v) = %v, want %v", tt.a, tt.b, tt.p, got, tt.want)
			}
		})
	}
}

func TestLerpColorPrecise(t *testing.T) {
	tests := []struct {
		name string
		a    struct{ R, G, B, A float32 }
		b    struct{ R, G, B, A float32 }
		p    float32
	}{
		{
			name: "fraction 0 returns first color",
			a:    struct{ R, G, B, A float32 }{1.0, 0.0, 0.0, 1.0},
			b:    struct{ R, G, B, A float32 }{0.0, 0.0, 1.0, 1.0},
			p:    0.0,
		},
		{
			name: "fraction 1 returns second color",
			a:    struct{ R, G, B, A float32 }{1.0, 0.0, 0.0, 1.0},
			b:    struct{ R, G, B, A float32 }{0.0, 0.0, 1.0, 1.0},
			p:    1.0,
		},
		{
			name: "gamma correct blending",
			a:    struct{ R, G, B, A float32 }{0.0, 0.0, 0.0, 1.0},
			b:    struct{ R, G, B, A float32 }{1.0, 1.0, 1.0, 1.0},
			p:    0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LerpColorPrecise(tt.a, tt.b, tt.p)

			// At p=0, should equal a
			if tt.p == 0 {
				if !rgbaEqual(got, tt.a, floatutils.Float32EqualityThreshold) {
					t.Errorf("ColorLerpPrecise at p=0 should return a, got %v", got)
				}
			}

			// At p=1, should equal b
			if tt.p == 1 {
				if !rgbaEqual(got, tt.b, floatutils.Float32EqualityThreshold) {
					t.Errorf("ColorLerpPrecise at p=1 should return b, got %v", got)
				}
			}

			// Gamma-correct blending produces different midpoint than linear
			if tt.p == 0.5 && tt.a.R == 0 && tt.b.R == 1 {
				// sqrt(0.5) ≈ 0.707, not 0.5 like linear
				expectedR := float32(math.Sqrt(0.5))
				if math.Abs(float64(got.R-expectedR)) > 1e-5 {
					t.Errorf("ColorLerpPrecise gamma-correct R = %v, want ~%v", got.R, expectedR)
				}
			}
		})
	}
}

func TestLerpColorList(t *testing.T) {
	t.Run("equal length lists", func(t *testing.T) {
		a := []struct{ R, G, B, A float32 }{
			{0, 0, 0, 1},
			{1, 0, 0, 1},
		}
		b := []struct{ R, G, B, A float32 }{
			{1, 1, 1, 1},
			{0, 1, 0, 1},
		}

		result := LerpColorList(a, b, 0.5)

		if len(result) != 2 {
			t.Fatalf("Expected length 2, got %d", len(result))
		}

		// First color: (0,0,0,1) to (1,1,1,1) at 0.5 = (0.5,0.5,0.5,1)
		expected0 := struct{ R, G, B, A float32 }{0.5, 0.5, 0.5, 1.0}
		if !rgbaEqual(result[0], expected0, floatutils.Float32EqualityThreshold) {
			t.Errorf("result[0] = %v, want %v", result[0], expected0)
		}
	})

	t.Run("different length lists", func(t *testing.T) {
		a := []struct{ R, G, B, A float32 }{
			{0, 0, 0, 1},
		}
		b := []struct{ R, G, B, A float32 }{
			{1, 0, 0, 1},
			{0, 1, 0, 1},
		}

		result := LerpColorList(a, b, 0.5)

		if len(result) != 2 {
			t.Fatalf("Expected length 2, got %d", len(result))
		}
	})
}

func TestLerpColorListPrecice(t *testing.T) {
	t.Run("equal length lists", func(t *testing.T) {
		a := []struct{ R, G, B, A float32 }{
			{0, 0, 0, 1},
			{1, 0, 0, 1},
		}
		b := []struct{ R, G, B, A float32 }{
			{1, 1, 1, 1},
			{0, 1, 0, 1},
		}

		result := LerpColorListPrecice(a, b, 0.5)

		if len(result) != 2 {
			t.Fatalf("Expected length 2, got %d", len(result))
		}
	})

	t.Run("different length lists", func(t *testing.T) {
		a := []struct{ R, G, B, A float32 }{
			{0, 0, 0, 1},
		}
		b := []struct{ R, G, B, A float32 }{
			{1, 0, 0, 1},
			{0, 1, 0, 1},
		}

		result := LerpColorListPrecice(a, b, 0.5)

		if len(result) != 2 {
			t.Fatalf("Expected length 2, got %d", len(result))
		}
	})
}
