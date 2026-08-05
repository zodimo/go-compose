package runtime

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/ui/platform"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/render"
	_ "github.com/zodimo/go-compose/internal/render/gio"
	"github.com/zodimo/go-compose/internal/unitconvert"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/store"
	"github.com/zodimo/go-compose/theme"
)

func gioAppMain() {
	app.Main()
}

func (a *Application) runWindow(cfg *windowConfig) {
	defer a.wg.Done()

	win := new(app.Window)
	var options []app.Option
	if cfg.Title != "" {
		options = append(options, app.Title(cfg.Title))
	}
	if cfg.Size != unit.DpSizeZero && cfg.Size != unit.DpSizeUnspecified {
		options = append(options, app.Size(
			unitconvert.DpToGioUnitUnsafe(cfg.Size.Width),
			unitconvert.DpToGioUnitUnsafe(cfg.Size.Height),
		))
	}
	if len(options) > 0 {
		win.Option(options...)
	}

	st := store.NewPersistentState(store.WithRootContext(a.ctx))
	sub := st.Subscribe(func() { win.Invalidate() })
	defer sub.Unsubscribe()

	themeManager := theme.GetThemeManager()
	rt := NewRuntime()

	var ops op.Ops
	for {
		switch e := win.Event().(type) {
		case app.DestroyEvent:
			if e.Err != nil {
				select {
				case a.windowExit <- e.Err:
				default:
				}
			}
			return
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			gtx.Locale = system.Locale{Language: "en", Direction: system.LTR}
			gtx = themeManager.Material3ThemeInit(gtx).(layout.Context)
			composer := compose.NewComposer(api.ComposerWithStore(st))
			ui := cfg.Content
			if ui == nil {
				ui = func(c api.Composer) api.Composer { return c }
			}
			if cfg.useLocalWindow {
				ui = compose.CompositionLocalProvider1(
					platform.LocalWindow,
					platform.NewWindow(win),
					ui,
				)
			}
			cmd := rt.Run(gtx, composer, ui)
			render.ApplyToGio(cmd, gtx.Ops)
			e.Frame(gtx.Ops)
		}
	}
}
