package surface

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/border"
	"github.com/zodimo/go-compose/modifiers/clip"
	"github.com/zodimo/go-compose/modifiers/shadow"
	"github.com/zodimo/go-compose/theme"
)

// Surface is a layout composable that represents a Material surface.
// It handles clipping, background color, and elevation.
func Surface(
	content Composable,
	options ...SurfaceOption,
) Composable {

	opts := DefaultSurfaceOptions()
	for _, option := range options {
		if option == nil {
			continue
		}
		option(&opts)
	}

	opts.Color = theme.TakeOrElseColor(opts.Color, theme.ColorHelper.ColorSelector().SurfaceRoles.Surface)
	opts.ContentColor = theme.TakeOrElseColor(opts.ContentColor, theme.ColorHelper.ColorSelector().SurfaceRoles.OnSurface)
	opts.BorderColor = theme.TakeOrElseColor(opts.BorderColor, theme.ColorHelper.SpecificColor(graphics.ColorTransparent))

	return func(c Composer) Composer {

		// Apply modifiers: Clip then Background.
		// Shadow should be behind everything.

		// Apply modifiers:
		// 1. Shadow (doesn't clip, sits behind)
		// 2. Clip (applies to everything following)
		// 3. Background (respects clip)
		// 4. Border (respects clip)
		// 5. Custom Modifiers (clickable, etc. - should respect clip)

		surfaceModifier := modifier.EmptyModifier.
			Then(shadow.Simple(opts.ShadowElevation, opts.Shape)).
			Then(clip.Clip(opts.Shape)).
			Then(background.Background(opts.Color, func(o *background.BackgroundOptions) { o.Shape = opts.Shape })).
			Then(border.Border(opts.BorderWidth, opts.BorderColor, opts.Shape)).
			Then(opts.Modifier)

		// Use Box to hold content and apply modifiers
		return box.Box(content, box.WithModifier(surfaceModifier), box.WithAlignment(opts.Alignment))(c)
	}
}
