package lerp

import "testing"

func TestLerpDiscrete(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		fraction float64
		want     string
	}{
		{"fraction 0 returns a", "first", "second", 0.0, "first"},
		{"fraction 0.49 returns a", "first", "second", 0.49, "first"},
		{"fraction 0.5 returns b", "first", "second", 0.5, "second"},
		{"fraction 0.51 returns b", "first", "second", 0.51, "second"},
		{"fraction 1 returns b", "first", "second", 1.0, "second"},
		{"negative fraction returns a", "first", "second", -0.5, "first"},
		{"fraction > 1 returns b", "first", "second", 1.5, "second"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LerpDiscrete(tt.a, tt.b, tt.fraction)
			if got != tt.want {
				t.Errorf("LerpDiscrete(%q, %q, %v) = %q, want %q", tt.a, tt.b, tt.fraction, got, tt.want)
			}
		})
	}
}

func TestLerpDiscrete_Int(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		fraction float32
		want     int
	}{
		{"fraction 0 returns a", 10, 20, 0.0, 10},
		{"fraction 0.49 returns a", 10, 20, 0.49, 10},
		{"fraction 0.5 returns b", 10, 20, 0.5, 20},
		{"fraction 1 returns b", 10, 20, 1.0, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LerpDiscrete(tt.a, tt.b, tt.fraction)
			if got != tt.want {
				t.Errorf("LerpDiscrete(%v, %v, %v) = %v, want %v", tt.a, tt.b, tt.fraction, got, tt.want)
			}
		})
	}
}

type MyStruct struct {
	Value string
}

func TestLerpDiscrete_Struct(t *testing.T) {
	a := MyStruct{Value: "first"}
	b := MyStruct{Value: "second"}

	t.Run("returns a below 0.5", func(t *testing.T) {
		got := LerpDiscrete(a, b, 0.3)
		if got != a {
			t.Errorf("Expected a, got %v", got)
		}
	})

	t.Run("returns b at 0.5", func(t *testing.T) {
		got := LerpDiscrete(a, b, 0.5)
		if got != b {
			t.Errorf("Expected b, got %v", got)
		}
	})
}
