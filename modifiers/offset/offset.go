package offset

import (
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/modifier"
)

// OffsetData holds the X and Y offset values
type OffsetData struct {
	X unit.Dp
	Y unit.Dp
}

// Offset creates a modifier that translates the element by the given X and Y offset.
// This is useful for creating overlapping layouts like profile cards with avatars.
func Offset(x, y unit.Dp) modifier.Modifier {
	return modifier.NewInspectableModifier(
		modifier.NewModifier(
			&OffsetElement{
				data: OffsetData{
					X: x,
					Y: y,
				},
			},
		),
		modifier.NewInspectorInfo(
			"offset",
			map[string]any{
				"x": x,
				"y": y,
			},
		),
	)
}

// OffsetX creates a modifier that translates only horizontally.
func OffsetX(x unit.Dp) modifier.Modifier {
	return Offset(x, 0)
}

// OffsetY creates a modifier that translates only vertically.
// Negative values move the element up, positive values move it down.
func OffsetY(y unit.Dp) modifier.Modifier {
	return Offset(0, y)
}
