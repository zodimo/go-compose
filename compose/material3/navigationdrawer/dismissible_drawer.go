package navigationdrawer

import (
	"time"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/material3/surface"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/animation"
	animMod "github.com/zodimo/go-compose/modifiers/animation"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/state"

	"github.com/zodimo/go-compose/compose/ui/unit"
)

// DismissibleNavigationDrawer uses a drawer that is usually visible but can be dismissed.
// When open, it sits side-by-side with the content.
// When closed, the content takes the full width.
func DismissibleNavigationDrawer(
	drawerContent Composable,
	content Composable,
	options ...ModalNavigationDrawerOption, // Reusing options pattern for consistency
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
			return animation.NewVisibilityAnimation(time.Millisecond*250, animation.Invisible)
		}).Get()

		// Sync animation state with props
		if opts.IsOpen {
			anim.Appear(c.TimeNow())
		} else {
			anim.Disappear(c.TimeNow())
		}

		return row.Row(
			func(c Composer) Composer {
				// 1. Drawer Sheet (Animated Width)
				// We wrap in a Box to clip/control size
				box.Box(
					func(c Composer) Composer {
						return surface.Surface(
							drawerContent,
							surface.WithColor(drawerContainerColor),
							surface.WithShape(&shape.RoundedCornerShape{Radius: unit.Dp(0)}),
							surface.WithModifier(
								size.Width(360). // Inner content fits 360
											Then(size.FillMaxHeight()),
							),
						)(c)
					},
					// The outer box constrains the width based on animation
					box.WithModifier(
						animMod.AnimatedWidth(anim, 360).
							Then(size.FillMaxHeight()),
						// We might need clipping here if the content shouldn't squash
						// But usually side-by-side drawers might squash?
						// Actually standard behavior: Drawer slides out.
						// If we just reduce width of container, the inner content might reflow.
						// To simulate "Draw sliding out", we should Clip the container.
						// Note: `clip.Clip` logic in go-compose might be complex for partial width.
						// For now, let's just animate width.
					),
				)(c)

				// 2. Main Content
				box.Box(
					content,
					box.WithModifier(size.FillMax()),
				)(c)

				return c
			},
			row.WithModifier(size.FillMax()),
		)(c)
	}
}
