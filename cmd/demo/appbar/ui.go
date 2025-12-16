package main

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/material3/appbar"
	"github.com/zodimo/go-compose/compose/foundation/material3/iconbutton"
	"github.com/zodimo/go-compose/compose/foundation/material3/scaffold"
	"github.com/zodimo/go-compose/compose/foundation/material3/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

func UI() api.Composable {
	return scaffold.Scaffold(
		func(c compose.Composer) compose.Composer {
			return column.Column(
				c.Sequence(
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
					Spacer(16),

					// 4. Center Aligned TopAppBar
					appbar.CenterAlignedTopAppBar(
						text.Text("Center Aligned", text.TypestyleTitleLarge),
						appbar.WithNavigationIcon(
							iconbutton.Standard(
								func() {},
								icons.NavigationMenu, // Using NavigationMenu as placeholder
								"Menu",
							),
						),
						appbar.WithActions(
							iconbutton.Standard(
								func() {},
								icons.SocialPerson, // Assuming SocialPerson exists or use ActionAccountCircle if available
								"Profile",
							),
						),
					),
					Spacer(16),

					// 5. Medium TopAppBar
					appbar.MediumTopAppBar(
						text.Text("Medium TopAppBar", text.TypestyleHeadlineSmall),
						appbar.WithNavigationIcon(
							iconbutton.Standard(
								func() {},
								icons.NavigationArrowBack,
								"Back",
							),
						),
						appbar.WithActions(
							iconbutton.Standard(
								func() {},
								icons.ActionSearch,
								"Search",
							),
						),
					),
					Spacer(16),

					// 6. Large TopAppBar
					appbar.LargeTopAppBar(
						text.Text("Large TopAppBar", text.TypestyleHeadlineMedium),
						appbar.WithNavigationIcon(
							iconbutton.Standard(
								func() {},
								icons.NavigationArrowBack,
								"Back",
							),
						),
						appbar.WithActions(
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
					Spacer(16),
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
