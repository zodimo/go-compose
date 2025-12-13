package badge

import (
	"image"
	"image/color"

	"go-compose-dev/internal/layoutnode"
	"go-compose-dev/internal/modifier"
	"go-compose-dev/pkg/api"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"git.sr.ht/~schnwalter/gio-mw/wdk"
)

// Badge is a Material 3 badge component.
func Badge(options ...BadgeOption) api.Composable {
	return func(c api.Composer) api.Composer {
		opts := DefaultBadgeOptions()
		for _, option := range options {
			option(&opts)
		}

		c.StartBlock("Badge")
		c.Modifier(func(modifier Modifier) Modifier {
			return modifier.Then(opts.Modifier)
		})

		if opts.Content != nil {
			c.WithComposable(opts.Content)
		}

		c.SetWidgetConstructor(badgeWidgetConstructor(opts))

		return c.EndBlock()
	}
}

func badgeWidgetConstructor(opts BadgeOptions) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			theme := wdk.GetMaterialTheme(gtx)

			// Resolve Defaults
			containerColor := opts.ContainerColor
			if containerColor == (color.NRGBA{}) {
				containerColor = theme.Scheme.Error.Color.AsNRGBA()
			}
			contentColor := opts.ContentColor
			if contentColor == (color.NRGBA{}) {
				contentColor = theme.Scheme.Error.OnColor.AsNRGBA()
			}

			children := node.Children()
			hasContent := len(children) > 0

			// M3 Specs
			var (
				heightDp     unit.Dp
				maxWidthDp   unit.Dp
				cornerRadius unit.Dp
			)

			if hasContent {
				// Large badge
				heightDp = 16
				maxWidthDp = 34  // Max character count size
				cornerRadius = 8 // Large badge shape: 8dp corner radius
			} else {
				// Small badge: fixed 6x6dp
				heightDp = 6
				maxWidthDp = 6
				cornerRadius = 3 // Small badge shape: 3dp corner radius
			}

			// Set constraints
			gtx.Constraints.Min.Y = gtx.Dp(heightDp)
			gtx.Constraints.Max.X = gtx.Dp(maxWidthDp)
			if !hasContent {
				// Small badge is always 6x6dp (square)
				gtx.Constraints.Min.X = gtx.Dp(6)
			}

			// Measure logic
			macro := op.Record(gtx.Ops)

			// Layout content if any
			var dims layout.Dimensions
			if hasContent {
				// Inset layout
				dims = layout.Inset{
					Left: unit.Dp(4), Right: unit.Dp(4),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					childCoords := children[0].(layoutnode.NodeCoordinator)
					return childCoords.Layout(gtx)
				})
			} else {
				// Empty small badge (always 6x6dp)
				dims = layout.Dimensions{Size: image.Pt(gtx.Dp(6), gtx.Dp(6))}
			}

			call := macro.Stop()

			// Draw Background
			w := dims.Size.X
			h := dims.Size.Y
			if h < gtx.Dp(heightDp) {
				h = gtx.Dp(heightDp)
			}

			size := image.Pt(w, h)

			// M3 Spec: Fixed corner radius (3dp for small, 8dp for large)
			radius := gtx.Dp(cornerRadius)
			rr := clip.RRect{
				Rect: image.Rectangle{Max: size},
				SE:   radius, SW: radius, NW: radius, NE: radius,
			}

			paint.FillShape(gtx.Ops, containerColor, clip.Outline{Path: rr.Path(gtx.Ops)}.Op())

			// Draw Content
			if hasContent {
				paint.ColorOp{Color: contentColor}.Add(gtx.Ops)
			}
			call.Add(gtx.Ops)

			return layout.Dimensions{Size: size}
		}
	})
}

// BadgedBox is a layout that places a badge at the top-end corner of the content.
func BadgedBox(
	badge api.Composable,
	content api.Composable,
	modifiers ...modifier.Modifier,
) api.Composable {
	return func(c api.Composer) api.Composer {
		c.StartBlock("BadgedBox")

		// Apply modifiers
		for _, m := range modifiers {
			c.Modifier(func(old modifier.Modifier) modifier.Modifier {
				return old.Then(m)
			})
		}

		// Child 0: Content
		c.WithComposable(content)
		// Child 1: Badge
		c.WithComposable(badge)

		c.SetWidgetConstructor(badgedBoxWidgetConstructor())

		return c.EndBlock()
	}
}

func badgedBoxWidgetConstructor() layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			children := node.Children()
			if len(children) < 2 {
				if len(children) == 1 {
					return children[0].(layoutnode.NodeCoordinator).Layout(gtx)
				}
				return layout.Dimensions{}
			}

			contentNode := children[0].(layoutnode.NodeCoordinator)
			badgeNode := children[1].(layoutnode.NodeCoordinator)

			// 1. Measure Content
			macroContent := op.Record(gtx.Ops)
			dimsContent := contentNode.Layout(gtx)
			callContent := macroContent.Stop()

			// 2. Measure Badge
			macroBadge := op.Record(gtx.Ops)
			dimsBadge := badgeNode.Layout(gtx)
			callBadge := macroBadge.Stop()

			// 3. Draw Content
			callContent.Add(gtx.Ops)

			// 4. Draw Badge Offset
			// M3 Spec: Distance from top trailing icon corner to bottom leading badge corner
			// Small badge: 6x6dp, Large badge: 14x12dp (HxV)
			var offsetX, offsetY int
			if dimsBadge.Size.Y <= gtx.Dp(6) {
				// Small badge
				offsetX = gtx.Dp(6)
				offsetY = gtx.Dp(6)
			} else {
				// Large badge
				offsetX = gtx.Dp(12)
				offsetY = gtx.Dp(14)
			}
			// Position badge: top-end corner of content minus offsets
			x := dimsContent.Size.X - dimsBadge.Size.X + offsetX
			y := -offsetY

			stack := op.Offset(image.Pt(x, y)).Push(gtx.Ops)
			callBadge.Add(gtx.Ops)
			stack.Pop()

			return dimsContent
		}
	})
}
