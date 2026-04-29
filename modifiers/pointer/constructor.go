package pointer

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/input/pointer"
	"github.com/zodimo/go-compose/internal/modifier"
)

func BlockPointer() ui.Modifier {
	return modifier.NewModifier(InputBlockerElement{})
}

func PointerInput(key any, block func(scope pointer.PointerInputScope)) ui.Modifier {

	return modifier.NewInspectableModifier(
		modifier.NewModifier(
			PointerInputElement{
				key:   key,
				block: block,
			},
		),
		modifier.NewInspectorInfo(
			"PointerInput",
			map[string]any{
				"key":   key,
				"block": block,
			},
		),
	)
}
