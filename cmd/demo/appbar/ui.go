package main

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/material3/appbar"
	"github.com/zodimo/go-compose/compose/material3/iconbutton"
	"github.com/zodimo/go-compose/compose/material3/scaffold"
	mswitch "github.com/zodimo/go-compose/compose/material3/switch"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

func UI() api.Composable {

	return func(c compose.Composer) compose.Composer {
		colorscheme := state.MustRemember(c, "color-scheme", func() *material3.ColorSchemeNext {
			return material3.DarkColorScheme()
		})

		ColorSchemeLightSwitch := state.MustRemember(c, "switch_state_1", func() bool { return false })

		return compose.CompositionLocalProvider(
			[]api.ProvidedValue{material3.LocalColorSchemeNext.Provides(colorscheme.Get())},
			scaffold.Scaffold(
				func(c compose.Composer) compose.Composer {
					return column.Column(
						c.Sequence(
							//ColorScheme Toggle
							row.Row(
								c.Sequence(
									text.BodyMedium("Toggle ColorScheme"),
									spacer.Width(20),
									mswitch.Switch(
										ColorSchemeLightSwitch.Get(),
										func(b bool) {
											ColorSchemeLightSwitch.Set(b)
											if b {
												colorscheme.Set(material3.LightColorScheme())
											} else {
												colorscheme.Set(material3.DarkColorScheme())
											}
										},
									),
								),
								row.WithSpacing(row.SpaceStart),
								row.WithAlignment(row.Middle),
							),

							// 1. Simple TopAppBar
							appbar.TopAppBar(
								text.TitleLarge("Simple TopAppBar"),
							),
							spacer.Height(16),

							// 2. TopAppBar with Navigation Icon
							appbar.TopAppBar(
								text.TitleLarge("With Nav Icon"),
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
								text.TitleLarge("With Actions"),
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
								text.TitleLarge("Center Aligned"),
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
								text.HeadlineSmall("Medium TopAppBar"),
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
								text.HeadlineMedium("Large TopAppBar"),
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

							func(c api.Composer) api.Composer {
								theme := material3.Theme(c)
								// 7. Custom Colors - Primary Theme
								return appbar.TopAppBar(
									text.HeadlineMedium("Custom Colors"),
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
										ContainerColor:             theme.ColorScheme().Primary,   //theme.ColorHelper.ColorSelector().PrimaryRoles.Primary,
										NavigationIconContentColor: theme.ColorScheme().OnPrimary, //theme.ColorHelper.ColorSelector().PrimaryRoles.OnPrimary,
										TitleContentColor:          theme.ColorScheme().OnPrimary, //theme.ColorHelper.ColorSelector().PrimaryRoles.OnPrimary,
										ActionIconContentColor:     theme.ColorScheme().OnPrimary, //theme.ColorHelper.ColorSelector().PrimaryRoles.OnPrimary,
									}),
								)(c)
							},
						),
						column.WithModifier(size.FillMax().
							Then(padding.All(16)), // Add some padding around the column
						),
					)(c)
				},
			),
		)(c)
	}
}
