package stopclickthru

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	mBox "github.com/zodimo/go-compose/modifiers/box"
	"github.com/zodimo/go-compose/modifiers/pointer"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
)

func StopClickThru(content api.Composable) api.Composable {
	return func(c api.Composer) api.Composer {
		return box.Box(
			c.Sequence(
				//box to block pointer events passing though the content
				box.Box(
					compose.Id(),
					box.WithModifier(
						mBox.MatchParentSize().Then(pointer.BlockPointer()),
					),
				),
				content,
			),
			box.WithModifier(
				size.WrapContentSize(),
			),
		)(c)
	}
}
