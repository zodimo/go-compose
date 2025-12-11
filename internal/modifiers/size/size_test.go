package size

import (
	"image"
	"testing"

	"gioui.org/layout"
)

func TestApplySizeDataToConstraints(t *testing.T) {
	initial := layout.Constraints{
		Min: image.Point{X: 100, Y: 100},
		Max: image.Point{X: 200, Y: 200},
	}

	tests := []struct {
		name     string
		data     SizeData
		expected layout.Constraints
	}{
		{
			name: "WrapWidth",
			data: SizeData{
				Width: NotSet, Height: NotSet, // Important initialization
				WrapWidth: true,
			},
			expected: layout.Constraints{
				Min: image.Point{X: 0, Y: 100},
				Max: image.Point{X: 200, Y: 200},
			},
		},
		{
			name: "WrapHeight",
			data: SizeData{
				Width: NotSet, Height: NotSet,
				WrapHeight: true,
			},
			expected: layout.Constraints{
				Min: image.Point{X: 100, Y: 0},
				Max: image.Point{X: 200, Y: 200},
			},
		},
		{
			name: "Unbounded WrapWidth",
			data: SizeData{
				Width: NotSet, Height: NotSet,
				WrapWidth: true, Unbounded: true,
			},
			expected: layout.Constraints{
				Min: image.Point{X: 0, Y: 100},
				Max: image.Point{X: 1000000, Y: 200},
			},
		},
		{
			name: "FillWidth override WrapWidth logic check",
			// If Fill is set, Wrap logic shouldn't matter for Min because Fill sets Min=Max
			// But our implementation sets Wrap LAST.
			// Check implementation: Fill sets Min=Max. Then Wrap sets Min=0.
			// Result: Min=0, Max=200.
			// Is this correct?
			// Compose: Modifiers are chained.
			// If we have one SizeData with BOTH Fill and Wrap?
			// Usually impossible via constructors (they are separate nodes).
			// But if constructed manually?
			// If I have FillMaxWidth && WrapWidth:
			// Fill sets Min=Max=200.
			// Wrap sets Min=0.
			// Result: Min=0, Max=200.
			// Effectively Wrap wins on Min.
			// Fill sets Min=Max. Wrap sets Min=0.
			data: SizeData{
				Width: NotSet, Height: NotSet,
				FillMaxWidth: true, WrapWidth: true,
			},
			expected: layout.Constraints{
				Min: image.Point{X: 0, Y: 100},
				Max: image.Point{X: 200, Y: 200},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplySizeDataToConstraints(initial, tt.data)
			if got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}
