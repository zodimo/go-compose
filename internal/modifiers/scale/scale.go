package scale

import (
	"go-compose-dev/internal/modifier"
)

type ScaleData struct {
	ScaleX float32
	ScaleY float32
}

func Scale(scale float32) modifier.Modifier {
	return modifier.NewModifier(
		&ScaleElement{
			data: ScaleData{
				ScaleX: scale,
				ScaleY: scale,
			},
		},
	)
}

func ScaleXY(x, y float32) modifier.Modifier {
	return modifier.NewModifier(
		&ScaleElement{
			data: ScaleData{
				ScaleX: x,
				ScaleY: y,
			},
		},
	)
}
