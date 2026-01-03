package overlay

import (
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/internal/layoutnode"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

func Overlay(content Composable, options ...OverlayOption) Composable {
	opts := DefaultOverlayOptions()
	for _, option := range options {
		if option == nil {
			continue
		}
		option(&opts)
	}
	return func(c Composer) Composer {

		theme := material3.Theme(c)
		opts.ScrimColor = opts.ScrimColor.TakeOrElse(graphics.SetOpacity(theme.ColorScheme().Scrim, 0.32))

		c.StartBlock("Overlay")
		c.Modifier(func(modifier Modifier) Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.WithComposable(content)
		c.SetWidgetConstructor(overlayWidgetConstructor(opts))

		return c.EndBlock()
	}
}

func overlayWidgetConstructor(options OverlayOptions) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			// Persistent clickable for the scrim
			scrimClickable := node.State("scrimClickable", func() any { return &widget.Clickable{} }).Get().(*widget.Clickable)

			// Handle scrim clicks
			if options.OnDismiss != nil {
				if scrimClickable.Clicked(gtx) {
					options.OnDismiss()
				}
			}

			// Resolve ScrimColor to NRGBA
			scrimColor := graphics.ColorToNRGBA(options.ScrimColor)

			// 1. Draw Scrim
			// We use a Stack to draw scrim behind content
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					// Fill the screen/parent with scrim
					// We use a Clickable to capture clicks, and paint the color
					return scrimClickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						paint.Fill(gtx.Ops, scrimColor)
						return layout.Dimensions{Size: gtx.Constraints.Max}
					})
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					// 2. Draw Content
					// We assume only one child content for Overlay
					children := node.Children()
					if len(children) > 0 {
						child := children[0].(layoutnode.NodeCoordinator)
						return child.Layout(gtx)
					}
					return layout.Dimensions{}
				}),
			)
		}
	})
}
