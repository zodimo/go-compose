package navigationdrawer

import (
	"time"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/material3/surface"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/animation"
	"github.com/zodimo/go-compose/internal/modifier"
	animMod "github.com/zodimo/go-compose/internal/modifiers/animation"
	"github.com/zodimo/go-compose/internal/modifiers/clickable"
	"github.com/zodimo/go-compose/internal/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/unit"
	"git.sr.ht/~schnwalter/gio-mw/token"
)

type Modifier = modifier.Modifier
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
		tm := theme.GetThemeManager()
		m3 := tm.GetMaterial3Theme()

		drawerContainerColor := m3.Scheme.SurfaceContainerLow.AsNRGBA()

		// Animation state
		anim := c.State(c.GenerateID().String()+"/anim", func() any {
			return &animation.VisibilityAnimation{
				Duration: time.Millisecond * 300,
				State:    animation.Invisible,
			}
		}).Get().(*animation.VisibilityAnimation)

		// Sync animation state
		if opts.IsOpen {
			anim.Appear(time.Now())
		} else {
			anim.Disappear(time.Now())
		}

		baseScrim := m3.Scheme.Scrim.SetOpacity(token.OpacityLevel8).AsNRGBA()
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
							modifier.EmptyModifier.
								Then(size.FillMax()).
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
						surface.WithShape(shape.RoundedCornerShape{Radius: unit.Dp(16)}),
						surface.WithModifier(
							modifier.EmptyModifier.
								Then(animMod.AnimatedWidth(anim, 360)). // Animate width
								Then(size.FillMaxHeight()),
						),
					)(c)
				}

				return c
			},
		)(c)
	}
}
