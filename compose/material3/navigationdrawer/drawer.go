package navigationdrawer

import (
	"time"

	"git.sr.ht/~schnwalter/gio-mw/token"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/material3/surface"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/animation"
	animMod "github.com/zodimo/go-compose/modifiers/animation"
	"github.com/zodimo/go-compose/modifiers/clickable"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"

	"github.com/zodimo/go-compose/compose/ui/unit"
)

type Composable = api.Composable
type Composer = api.Composer

// ModalNavigationDrawer implements a navigation drawer that overlays the content.
// It uses a generic Box layout to stack a scrim and the drawer sheet over the content.
func ModalNavigationDrawer(
	drawerContent Composable,
	content Composable,
	options ...ModalNavigationDrawerOption,

) Composable {
	opts := DefaultModalNavigationDrawerOptions()
	for _, opt := range options {
		opt(&opts)
	}

	return func(c Composer) Composer {
		theme := material3.Theme(c)
		drawerContainerColor := theme.ColorScheme().SurfaceContainerLow //theme.ColorHelper.ColorSelector().SurfaceRoles.ContainerLow

		// Animation state
		anim := state.MustState(c, c.GenerateID().String()+"/anim", func() *animation.VisibilityAnimation {
			return animation.NewVisibilityAnimation(time.Millisecond*300, animation.Invisible)
		}).Get()

		// Sync animation state
		if opts.IsOpen {
			anim.Appear(c.TimeNow())
		} else {
			anim.Disappear(c.TimeNow())
		}

		baseScrim := graphics.SetOpacity(theme.ColorScheme().Scrim, float32(token.OpacityLevel8))
		// We use AnimatedBackground for scrim opacity

		return box.Box(
			func(c Composer) Composer {
				// 1. Main Content (Background)
				content(c)

				// 2. Scrim (Animated Opacity)
				// We always render scrim, but AnimatedBackground handles transparency.
				// However, if fully invisible, we might want to skip click handling.
				// For now, let's rely on AnimatedBackground to just hide it visually.
				// But click handling: if opacity is 0, click should ideally pass through?
				// Gio Clickable intercepts even if invisible?
				// Yes.
				// So if !anim.Visible(), we should NOT render the scrim box.
				// But anim.Visible() is stale.
				// We can just render it. The click intercept requires Modal to blocking.
				// If drawer is closed, scrim should be GONE.
				// If we rely on stale `anim.Visible()`, it might flicker.
				// But `opts.IsOpen` is the source of truth for "Goal".
				// If `!opts.IsOpen` AND `anim.State == Invisible`, then we can skip.

				shouldRender := opts.IsOpen || anim.Visible()

				if shouldRender {
					box.Box(
						func(c Composer) Composer { return c },
						box.WithModifier(
							size.FillMax().
								Then(animMod.AnimatedBackground(anim, baseScrim, shape.ShapeRectangle)). // RectangleShape singleton
								Then(clickable.OnClick(func() {
									if opts.OnClose != nil {
										opts.OnClose()
									}
								})),
						),
					)(c)

					// 3. Drawer Sheet
					surface.Surface(
						drawerContent,
						surface.WithColor(drawerContainerColor),
						surface.WithShape(&shape.RoundedCornerShape{Radius: unit.Dp(16)}),
						surface.WithModifier(
							animMod.AnimatedWidth(anim, 360). // Animate width
												Then(size.FillMaxHeight()),
						),
					)(c)
				}

				return c
			},
		)(c)
	}
}
