package main

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/material3/appbar"
	"github.com/zodimo/go-compose/compose/material3/iconbutton"
	"github.com/zodimo/go-compose/compose/material3/scaffold"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

func UI() api.Composable {
	return scaffold.Scaffold(
		func(c compose.Composer) compose.Composer {
			theme := material3.Theme(c)
			return column.Column(
				c.Sequence(
					// 1. Simple TopAppBar
					appbar.TopAppBar(
						text.TextWithStyle("Simple TopAppBar", text.TypestyleTitleLarge),
					),
					spacer.Height(16),

					// 2. TopAppBar with Navigation Icon
					appbar.TopAppBar(
						text.TextWithStyle("With Nav Icon", text.TypestyleTitleLarge),
						appbar.WithNavigationIcon(
							iconbutton.Standard(
								func() {},
								icons.NavigationMenu,
								"Menu",
							),
						),
					),
					spacer.Height(16),

					// 3. TopAppBar with Actions
					appbar.TopAppBar(
						text.TextWithStyle("With Actions", text.TypestyleTitleLarge),
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
					spacer.Height(16),

					// 4. Center Aligned TopAppBar
					appbar.CenterAlignedTopAppBar(
						text.TextWithStyle("Center Aligned", text.TypestyleTitleLarge),
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
					spacer.Height(16),

					// 5. Medium TopAppBar
					appbar.MediumTopAppBar(
						text.TextWithStyle("Medium TopAppBar", text.TypestyleHeadlineSmall),
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
					spacer.Height(16),

					// 6. Large TopAppBar
					appbar.LargeTopAppBar(
						text.TextWithStyle("Large TopAppBar", text.TypestyleHeadlineMedium),
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
					spacer.Height(16),

					// 7. Custom Colors - Primary Theme
					appbar.TopAppBar(
						text.TextWithStyle("Custom Colors", text.TypestyleTitleLarge),
						appbar.WithNavigationIcon(
							iconbutton.Standard(
								func() {},
								icons.NavigationMenu,
								"Menu",
							),
						),
						appbar.WithActions(
							iconbutton.Standard(
								func() {},
								icons.ActionSearch,
								"Search",
							),
						),
						appbar.WithColors(appbar.TopAppBarColors{
							ContainerColor:             theme.ColorScheme().Primary.Color,   //theme.ColorHelper.ColorSelector().PrimaryRoles.Primary,
							NavigationIconContentColor: theme.ColorScheme().Primary.OnColor, //theme.ColorHelper.ColorSelector().PrimaryRoles.OnPrimary,
							TitleContentColor:          theme.ColorScheme().Primary.OnColor, //theme.ColorHelper.ColorSelector().PrimaryRoles.OnPrimary,
							ActionIconContentColor:     theme.ColorScheme().Primary.OnColor, //theme.ColorHelper.ColorSelector().PrimaryRoles.OnPrimary,
						}),
					),
					spacer.Height(16),
				),
				column.WithModifier(size.FillMax().
					Then(padding.All(16)), // Add some padding around the column
				),
			)(c)
		},
	)
}
