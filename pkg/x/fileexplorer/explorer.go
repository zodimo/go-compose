package fileexplorer

import (
	"context"
	"fmt"

	"gioui.org/x/explorer"
	"github.com/zodimo/go-compose/compose/effect"
	"github.com/zodimo/go-compose/compose/ui/platform"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
)

func RememberExplorer(c api.Composer, explorerAction func(expl *explorer.Explorer)) (onClick func(), launchedEffect api.Composable) {

	key := c.GenerateID()
	path := c.GetPath()

	rememberKeyPrefix := fmt.Sprintf("%s-%s", key, path)

	window := platform.LocalWindow.Current(c)

	expl := state.MustRemember(c, rememberKeyPrefix+"explorer",
		func() *explorer.Explorer {
			return explorer.NewExplorer(window)
		},
	).Get()

	running := state.MustRemember(c, rememberKeyPrefix+"running", func() bool {
		return false
	})

	effectCounter := state.MustRemember(c, rememberKeyPrefix+"effectCounter", func() int {
		return 0
	})

	onClick = func() {
		if !running.Get() {
			running.Set(true)
		}
	}

	launchedEffect = effect.LaunchedEffect(
		func(ctx context.Context) {
			defer func() {
				running.Set(false)
				effectCounter.Set(effectCounter.Get() + 1)
			}()

			explorerAction(expl)
		},
		effectCounter.Get(),
	)

	return onClick, c.When(running.Get(), launchedEffect)

}
