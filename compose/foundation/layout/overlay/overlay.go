package overlay

import (
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/internal/layoutnode"

	"gioui.org/gesture"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
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
		c.Modifier(func(modifier ui.Modifier) ui.Modifier {
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
			// Persistent gesture for scrim click detection
			scrimClick := node.State("scrimClick", func() any { return &gesture.Click{} }).Get().(*gesture.Click)

			// Persistent tag for input blocking (catches events not handled by scrim or content)
			inputBlockerTag := node.State("inputBlockerTag", func() any { return new(int) }).Get().(*int)

			// Resolve ScrimColor to NRGBA
			scrimColor := graphics.ColorToNRGBA(options.ScrimColor)

			parentSize := gtx.ToGio().Constraints.Max

			// Capture backend before entering layout.Stack callback
			backend := gtx.DrawBackend()

			// Layout with Stack - center content
			dims := layout.Stack{Alignment: layout.Center}.Layout(*gtx.ToGio(),
				// Layer 0: Scrim background
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					paint.Fill(gtx.Ops, scrimColor)

					// Register full-screen scrim click handler for dismiss AND input blocking
					// This is registered AFTER content, so content handlers receive events first.
					// Events not consumed by content will hit this handler.
					func() {
						area := clip.Rect{Max: parentSize}.Push(gtx.Ops)
						defer area.Pop()

						// Register scrim click gesture for dismiss
						scrimClick.Add(gtx.Ops)
						event.Op(gtx.Ops, scrimClick)

						// Also register input blocker to catch any remaining pointer events
						event.Op(gtx.Ops, inputBlockerTag)
					}()

					// Process scrim click events for dismiss
					if options.OnClick != nil {
						for {
							e, ok := scrimClick.Update(gtx.Source)
							if !ok {
								break
							}
							if e.Kind == gesture.KindClick {
								options.OnClick()
							}
						}
					}

					// Drain any remaining pointer events to block them from reaching background
					for {
						_, ok := gtx.Event(pointer.Filter{
							Target: inputBlockerTag,
							Kinds:  pointer.Press | pointer.Release | pointer.Move | pointer.Drag | pointer.Scroll | pointer.Enter | pointer.Leave | pointer.Cancel,
						})
						if !ok {
							break
						}
					}
					return layout.Dimensions{Size: gtx.Constraints.Max}
				}),
				// Layer 1: Content (drawn on top)
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					children := node.Children()
					if len(children) > 0 {
						child := children[0].(layoutnode.NodeCoordinator)
						return layoutnode.ToGioDimensions(child.Layout(layoutnode.NewLayoutContextWithBackend(&gtx, backend)))
					}
					return layout.Dimensions{}
				}),
			)

			return layoutnode.FromGioDimensions(dims)
		}
	})
}
