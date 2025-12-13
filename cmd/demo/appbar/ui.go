package main

import (
	"go-compose-dev/compose"
	"go-compose-dev/compose/foundation/layout/column"
	"go-compose-dev/compose/foundation/layout/row"
	"go-compose-dev/compose/foundation/material3/appbar"
	"go-compose-dev/compose/foundation/material3/iconbutton"
	"go-compose-dev/compose/foundation/material3/scaffold"
	"go-compose-dev/compose/foundation/material3/text"
	"go-compose-dev/internal/modifiers/padding"
	"go-compose-dev/internal/modifiers/size"
	"go-compose-dev/pkg/api"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

func UI() api.Composable {
	return scaffold.Scaffold(
		func(c compose.Composer) compose.Composer {
			return column.Column(
				compose.Sequence(
					// 1. Simple TopAppBar
					appbar.TopAppBar(
						text.Text("Simple TopAppBar", text.TypestyleTitleLarge),
					),
					Spacer(16),

					// 2. TopAppBar with Navigation Icon
					appbar.TopAppBar(
						text.Text("With Nav Icon", text.TypestyleTitleLarge),
						appbar.WithNavigationIcon(
							iconbutton.Standard(
								func() {},
								icons.NavigationMenu,
								"Menu",
							),
						),
					),
					Spacer(16),

					// 3. TopAppBar with Actions
					appbar.TopAppBar(
						text.Text("With Actions", text.TypestyleTitleLarge),
						appbar.WithActions(
							row.Row(
								compose.Sequence(
									iconbutton.Standard(
										func() {},
										icons.ActionFavorite,
										"Favorite",
									),
									iconbutton.Standard(
										func() {},
										icons.ActionSearch,
										"Search",
									),
									iconbutton.Standard(
										func() {},
										icons.NavigationMoreVert,
										"More",
									),
								),
							),
						),
					),
					Spacer(16),

					// 4. Complete TopAppBar (Nav + Actions)
					appbar.TopAppBar(
						text.Text("Complete TopAppBar", text.TypestyleTitleLarge),
						appbar.WithNavigationIcon(
							iconbutton.Standard(
								func() {},
								icons.NavigationArrowBack,
								"Back",
							),
						),
						appbar.WithActions(
							row.Row(
								compose.Sequence(
									iconbutton.Standard(
										func() {},
										icons.ContentContentCopy,
										"Copy",
									),
								),
							),
						),
					),
				),
				column.WithModifier(size.FillMax()),
				column.WithModifier(padding.All(16)), // Add some padding around the column
			)(c)
		},
	)
}

func Spacer(dp int) compose.Composable {
	return func(c compose.Composer) compose.Composer {
		return row.Row(
			func(c compose.Composer) compose.Composer { return c },
			row.WithModifier(size.Height(dp)),
		)(c)
	}
}
