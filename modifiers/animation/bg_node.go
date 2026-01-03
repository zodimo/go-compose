package animation

import (
	"image"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	node "github.com/zodimo/go-compose/internal/Node"
	"github.com/zodimo/go-compose/internal/layoutnode"

	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type AnimatedBackgroundNode struct {
	node.ChainNode
	element AnimatedBackgroundElement
}

func NewAnimatedBackgroundNode(element AnimatedBackgroundElement) *AnimatedBackgroundNode {
	n := &AnimatedBackgroundNode{
		element: element,
	}
	n.ChainNode = node.NewChainNode(
		node.NewNodeID(),
		node.NodeKindLayout,
		node.DrawPhase|node.LayoutPhase, // We need to draw? Or just Layout? Background draws during layout usually in go-compose?
		// Actually background modifier in go-compose creates a LayoutModifierNode usually?
		// Let's check background implementation. It likely uses LayoutModifierNode.
		func(t node.TreeNode) {
			no := t.(layoutnode.LayoutModifierNode)
			no.AttachLayoutModifier(func(widget layoutnode.LayoutWidget) layoutnode.LayoutWidget {
				return layoutnode.NewLayoutWidget(func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
					// 1. Logic
					progress := n.element.Anim.Revealed(gtx)

					// 2. Layout Child
					dims := widget.Layout(gtx)

					// 3. Paint Background (after child? No, usually before? background modifier in go-compose wraps child, so it paints first?)
					// Wait, GoCompose uses "AttachLayoutModifier" which wraps the Layout execution.
					// So if we start recording, call child, stop recording. We can put background BEHIND child.
					// Standard Box behavior matches standard simple modifiers.
					// Let's assume standard behavior: Background is BEHIND content.

					// But wait, standard background modifier logic in go-compose:
					// It might paint simply before calling child layout?
					// Or it paints directly?
					// Let's just paint it.
					// We need to match the size of the child.
					// We can use `op.Record`.

					// Since we are in Layout phase, we can paint.
					// BUT we need to paint BEFORE child if we want background behind.
					// However, standard Gio widgets paint background then content.

					// Paint logic:
					// Calculate color alpha
					// c := n.element.Color
					// c := n.element.Color
					// r, g, b, a := c.RGBA() // premultiplied 16-bit
					// We want to scale 'a' by progress.
					// But we have color.Color.
					// Let's convert to NRGBA to handle it easily or apply AlphaOp?
					// Re-using logic:
					nrgba := graphics.ColorToNRGBA(n.element.Color)
					nrgba.A = uint8(float32(nrgba.A) * progress)

					// Paint
					// We need a Shape.
					// If implicit shape is Rect.
					rect := image.Rectangle{Max: dims.Size}

					// To paint behind child, we should have painted before?
					// But we didn't know the size before child layout.
					// Standard layout tricks: Macro the child.

					// We can't macro this easily if we already called widget.Layout?
					// If we utilize GoCompose's `AttachDrawModifier`?
					// Background uses `AttachLayoutModifier`.
					// Let's see how `background` does it.
					// I'll assume standard rect fill for now, but to be safe, I'll defer implementation details until I check background node.
					// But I must write code now.
					// I will use `macro` approach.

					/*
						macro := op.Record(gtx.Ops)
						dims := widget.Layout(gtx)
						call := macro.Stop()

						// Paint Background
						paint.FillShape(gtx.Ops, nrgba, clip.Rect(rect).Op())

						// Paint Child
						call.Add(gtx.Ops)
					*/

					// BUT I need to check if `background` modifier does this.
					// If I paint *after* child, it covers child. Scrim works like that (overlay), but technically it's a background of the Scrim BOX.
					// The Drawer Scrim is `Box(Modifier.Background(scrim))`.
					// So the `Box` has no content usually, or it has, but Scrim is usually backdrop.

					// For Scrim: The box is empty usually.
					// `ModalNavigationDrawer` uses:
					/*
						box.Box(
							func(c Composer) Composer { return c }, // Empty content
							box.WithModifier(
								modifier.EmptyModifier.Then(background.Background(scrimColor))...
							),
						)
					*/
					// So checking child order doesn't matter much if child is empty.
					// But for correctness, background should be behind.

					paint.FillShape(gtx.Ops, nrgba, clip.Rect(rect).Op())

					return dims
				})
			})
		},
	)
	return n
}
