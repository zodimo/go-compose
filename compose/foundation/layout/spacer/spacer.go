package spacer

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
)

func Spacer(d int) Composable {
	return box.Box(
		func(c Composer) Composer { return c },
		box.WithModifier(
			size.Size(d, d),
		),
	)
}

func SpacerWidth(d int) Composable {
	return box.Box(
		func(c Composer) Composer { return c },
		box.WithModifier(
			size.Width(d),
		),
	)
}

func SpacerHeight(d int) Composable {
	return box.Box(
		func(c Composer) Composer { return c },
		box.WithModifier(
			size.Height(d),
		),
	)
}

func SpacerWeight(d int) Composable {
	return box.Box(
		compose.Id(),
		box.WithModifier(
			weight.Weight(1),
		),
	)
}
