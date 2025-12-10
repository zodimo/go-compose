package weight

import "go-compose-dev/internal/modifier"

func Weight(weight int) Modifier {

	if weight < 0 {
		panic("weight cannot be negative")
	}

	return modifier.NewInspectableModifier(
		modifier.NewModifier(
			&WeightElement{
				weight: WeightData{
					Weight: float32(weight),
				},
			},
		),
		modifier.NewInspectorInfo(
			"weight",
			map[string]any{
				"weight": weight,
			},
		),
	)
}
