package navigationdrawer

import (
	"time"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/material3/surface"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/animation"
	"github.com/zodimo/go-compose/internal/modifier"
	animMod "github.com/zodimo/go-compose/internal/modifiers/animation"
	"github.com/zodimo/go-compose/internal/modifiers/size"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/unit"
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
		tm := theme.GetThemeManager()
		m3 := tm.GetMaterial3Theme()

		drawerContainerColor := m3.Scheme.SurfaceContainerLow.AsNRGBA()

		// Animation state
		// We use a persistent pointer for the animation state
		anim := c.State(c.GenerateID().String()+"/anim", func() any {
			return &animation.VisibilityAnimation{
				Duration: time.Millisecond * 250,
				State:    animation.Invisible,
			}
		}).Get().(*animation.VisibilityAnimation)

		// Sync animation state with props
		if opts.IsOpen {
			anim.Appear(time.Now())
		} else {
			anim.Disappear(time.Now())
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
							surface.WithShape(shape.RoundedCornerShape{Radius: unit.Dp(0)}),
							surface.WithModifier(
								modifier.EmptyModifier.
									Then(size.Width(360)). // Inner content fits 360
									Then(size.FillMaxHeight()),
							),
						)(c)
					},
					// The outer box constrains the width based on animation
					box.WithModifier(
						modifier.EmptyModifier.
							Then(animMod.AnimatedWidth(anim, 360)).
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
