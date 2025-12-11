package clip

import (
	"go-compose-dev/internal/modifier"
)

func Clip(shape Shape) Modifier {

	return modifier.NewInspectableModifier(
		modifier.NewModifier(
			&ClipElement{
				clipData: ClipData{
					Shape: shape,
				},
			},
		),
		modifier.NewInspectorInfo(
			"clip",
			map[string]any{
				"shape": shape,
			},
		),
	)
}
