package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/material3/button"
	"github.com/zodimo/go-compose/compose/foundation/material3/scaffold"
	"github.com/zodimo/go-compose/compose/foundation/material3/text"
	"github.com/zodimo/go-compose/compose/navigation"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/runtime"
	"github.com/zodimo/go-compose/state"
	"github.com/zodimo/go-compose/store"
	"github.com/zodimo/go-compose/theme"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Navigation Demo"))
		w.Option(app.Size(unit.Dp(800), unit.Dp(600)))

		if err := Run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

func Run(window *app.Window) error {
	enLocale := system.Locale{Language: "en", Direction: system.LTR}
	var ops op.Ops

	store := store.NewPersistentState(map[string]state.MutableValue{})
	store.SetOnStateChange(func() {
		window.Invalidate()
	})
	runtime := runtime.NewRuntime()
	themeManager := theme.GetThemeManager()

	for {
		switch frameEvent := window.Event().(type) {
		case app.DestroyEvent:
			return frameEvent.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, frameEvent)
			gtx.Locale = enLocale
			gtx = themeManager.Material3ThemeInit(gtx)

			composer := compose.NewComposer(store)
			composer = DemoUI()(composer)
			layoutNode := composer.Build()

			callOp := runtime.Run(gtx, layoutNode)
			callOp.Add(gtx.Ops)
			frameEvent.Frame(gtx.Ops)
		}
	}
}

func DemoUI() api.Composable {
	return func(c api.Composer) api.Composer {
		navController := navigation.RememberNavController(c)

		return scaffold.Scaffold(
			func(c api.Composer) api.Composer {
				// Content
				return navigation.NavHost(navController, "home", func(b *navigation.NavGraphBuilder) {
					b.Composable("home", HomeScreen(navController))
					b.Composable("details", DetailsScreen(navController))
				})(c)
			},
		)(c)
	}
}

func HomeScreen(navController *navigation.NavController) api.Composable {
	return func(c api.Composer) api.Composer {
		return column.Column(
			c.Sequence(
				text.Text("Home Screen", text.TypestyleBodyLarge),
				button.Filled(
					func() {
						navController.Navigate("details")
					},
					"Go to Details",
				),
			),
			column.WithSpacing(layout.SpaceAround),
			column.WithAlignment(layout.Middle),
		)(c)
	}
}

func DetailsScreen(navController *navigation.NavController) api.Composable {
	return func(c api.Composer) api.Composer {
		return column.Column(
			c.Sequence(
				text.Text("Details Screen", text.TypestyleBodyLarge),
				button.Filled(
					func() {
						navController.PopBackStack()
					},
					"Go Back",
				),
			),
			column.WithSpacing(layout.SpaceAround),
			column.WithAlignment(layout.Middle),
		)(c)
	}
}
