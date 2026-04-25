package main

import (
	_ "github.com/zodimo/gio-skia/skia"
	"github.com/zodimo/go-compose/compose/foundation/canvas"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {

		return column.Column(
			c.Sequence(
				text.HeadlineMedium("Canvas Demo"),
				canvas.Canvas(
					func(s graphics.DrawScope) {

						size := s.Size()
						w := size.Width()
						h := size.Height()

						// ─────────────────────────────────────────────────────────
						// Row 1: Rectangles
						// ─────────────────────────────────────────────────────────

						s.DrawIntoCanvas(func(c graphics.Canvas) {
							// DrawRect
							c.Save()
							c.Translate(w*0.15, h*0.15)
							//rect := models.Rect{Left: -40, Top: -30, Right: 40, Bottom: 30}
							// p := skia.NewPaintFill(color.NRGBA{R: 100, G: 200, B: 255, A: 255})
							p := graphics.NewPaint()
							p.ApplyStyle(graphics.NewStroke(2))
							p.Color = graphics.ColorBlue
							p.Alpha = 0.25
							// p.SetColor(color.NRGBA{R: 100, G: 200, B: 255, A: 255})

							c.DrawRect(-40, -30, 40, 30, p)
							c.Restore()
						})

					},
					canvas.WithModifier(size.FillMax()),
				),
			),
			column.WithSpacing(column.SpaceSides),
			column.WithAlignment(column.Middle),
			column.WithModifier(size.FillMax()),
		)(c)
	}
}
