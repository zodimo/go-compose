package main

import (
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI(c api.Composer) api.LayoutNode {

	c = column.Column(
		c.Sequence(
			// Row with explicitly large size
			row.Row(
				c.Sequence(
					// Box 1: Wrapped Content (Text)
					column.Column(
						c.Sequence(
							text.Text("Wrap Content (Default)",
								text.WithTextStyleOptions(
									text.StyleWithColor(color.NRGBA{R: 255, G: 255, B: 255, A: 255}),
								),
							),
						),
						column.WithModifier(background.Background(color.NRGBA{R: 0, G: 0, B: 200, A: 255})),
						// This should wrap the text size
						column.WithModifier(size.WrapContentSize()),
					),

					// Box 2: Wrapped Content (Top Start)
					column.Column(
						c.Sequence(
							text.Text("Top Start",
								text.WithTextStyleOptions(
									text.StyleWithColor(color.NRGBA{R: 255, G: 255, B: 255, A: 255}),
								),
							),
						),
						column.WithModifier(background.Background(color.NRGBA{R: 0, G: 100, B: 0, A: 255})),
						// Wrapped but parent forces size? No, we are in a row.
						// Let's create a fixed size container and put a wrapped item inside it to test alignment.
						column.WithModifier(size.WrapContentSize(size.TopStart)),
					),
				),
				row.WithModifier(size.FillMaxWidth()),
				row.WithModifier(size.Height(100)), // Fixed height for row
				row.WithModifier(background.Background(color.NRGBA{R: 50, G: 50, B: 50, A: 255})),
			),

			// A Fixed Size Box containing an aligned Wrapped Child
			column.Column(
				c.Sequence(
					// Child is wrapped and aligned BottomEnd
					column.Column(
						c.Sequence(
							text.Text("Bottom End Aligned",
								text.WithTextStyleOptions(
									text.StyleWithColor(color.NRGBA{R: 0, G: 0, B: 0, A: 255}),
								),
							),
						),
						column.WithModifier(background.Background(color.NRGBA{R: 200, G: 200, B: 0, A: 255})),
						// This container wraps its text content
						column.WithModifier(size.WrapContentSize()),
					),
				),
				// Parent is fixed size
				column.WithModifier(size.Size(300, 200)),
				column.WithModifier(background.Background(color.NRGBA{R: 100, G: 0, B: 0, A: 255})),
				// We want the inner column to be aligned BottomEnd within this parent?
				// Actually, size.WrapContentSize(Align) on the CHILD modifier essentially says:
				// "My size should be my content size, but fail: WrapContentSize on a node dictates how IT behaves."

				// Wait, if I want to align a child within a parent, conventionally in Compose:
				// Parent (Box) -> modifiers.size(300.dp)
				// Child -> modifiers.align(Alignment.BottomEnd) [BoxScope]
				// OR
				// Child -> modifiers.wrapContentSize(Alignment.BottomEnd) ???

				// In Jetpack Compose:
				// .wrapContentSize(Alignment.BottomEnd) on a Modifier means:
				// The element will satisfy its own content size.
				// BUT if the incoming constraints are larger, it will effectively occupy that larger space *logically* in the chain (or allow pass through),
				// but visually position the content at BottomEnd.
				//
				// Let's verify our implementation:
				// NewSizeNode:
				// 1. Calculate child constraints (Unbounded/0 min).
				// 2. Measure child -> childDims.
				// 3. Determine 'mySize'.
				//    If wrapped: mySize = Clamp(childDims, constraints.Min, constraints.Max).
				//    Wait, if constraints.Min are from Parent (e.g. fixed size parent forces exact constraints?),
				//    then Clamp(child, fixed, fixed) = fixed.
				//    So 'mySize' becomes the parent size.
				// 4. Align: Offset(child in mySize).

				// So if parent is Size(300, 200), it passes Min=300, Max=300.
				// WrapContentSize on child:
				// 1. Child constraints Min=0, Max=300. Child measures small (e.g. 50x20).
				// 2. mySize = Clamp(50, 300, 300) = 300.
				// 3. Offset(50 inside 300) via Alignment.
				// This matches Compose behavior where wrapContentSize allows a small element to be placed within a larger enforced constraint!

				// So applying WrapContentSize(BottomEnd) to the INNER node is correct to simulate alignment within the forced constraints of the parent chain.

				column.WithModifier(size.WrapContentSize(size.BottomEnd)),
			),
		),
		column.WithModifier(size.FillMax()),
		column.WithModifier(background.Background(color.NRGBA{R: 30, G: 30, B: 30, A: 255})),
	)(c)

	return c.Build()
}
