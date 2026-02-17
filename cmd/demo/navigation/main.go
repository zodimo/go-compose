package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/material3/button"
	"github.com/zodimo/go-compose/compose/material3/scaffold"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/navigation"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/runtime"
	"github.com/zodimo/go-compose/store"
	"github.com/zodimo/go-compose/theme"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Navigation Demo with Arguments"))
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

	store := store.NewPersistentState()
	store.Subscribe(func() {
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

			composer := compose.NewComposer(api.ComposerWithStore(store))

			callOp := runtime.Run(gtx, composer, DemoUI())
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
				return navigation.NavHost(
					navController,
					"home",
					func(b *navigation.NavGraphBuilder) {
						// Home screen without arguments
						b.Composable("home", HomeScreen(navController))

						// Details screen with required path argument {itemId}
						b.ComposableWithArgs("details/{itemId}", func(entry *navigation.BackStackEntry) api.Composable {
							return DetailsScreen(navController, entry)
						})
					},
				)(c)
			},
		)(c)
	}
}

func HomeScreen(navController *navigation.NavController) api.Composable {
	return func(c api.Composer) api.Composer {
		return column.Column(
			c.Sequence(
				text.TitleLarge("Home Screen"),
				text.BodyLarge("Select an item to view details:"),

				// Navigate to different items with different IDs
				button.Filled(
					func() {
						navController.Navigate("details/101")
					},
					"View Item 101",
				),
				button.Filled(
					func() {
						navController.Navigate("details/202")
					},
					"View Item 202",
				),
				button.Filled(
					func() {
						navController.Navigate("details/303")
					},
					"View Item 303",
				),
			),
			column.WithSpacing(layout.SpaceAround),
			column.WithAlignment(layout.Middle),
		)(c)
	}
}

func DetailsScreen(navController *navigation.NavController, entry *navigation.BackStackEntry) api.Composable {
	return func(c api.Composer) api.Composer {
		// Extract itemId from arguments
		itemId := "unknown"
		if entry.Arguments.IsSome() {
			args := entry.Arguments.UnwrapUnsafe()
			if id, ok := args.GetString("itemId"); ok {
				itemId = id
			}
		}

		return column.Column(
			c.Sequence(
				text.TitleLarge("Details Screen"),
				text.HeadlineMedium(fmt.Sprintf("Viewing Item: %s", itemId)),
				text.BodyLarge(fmt.Sprintf("This is the detail view for item ID: %s", itemId)),
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
