package main

import (
	_ "github.com/zodimo/gio-skia/skia"
	"github.com/zodimo/go-compose/compose/foundation/canvas"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/internal/skiagraphics"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
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

						s.DrawIntoCanvas(func(c graphics.Canvas) {
							c.Save()
							c.Translate(w*0.15, h*0.15)
							p := graphics.NewPaint()
							p.ApplyStyle(graphics.NewStroke(2))
							p.Color = graphics.ColorBlue
							p.Alpha = 0.25

							c.DrawRect(-40, -30, 40, 30, p)
							c.Restore()
						})

					},
					canvas.WithModifier(weight.Weight(1)),
				),
				text.HeadlineSmall("Path Demo - AddPath"),
				canvas.Canvas(
					func(s graphics.DrawScope) {
						s.DrawIntoCanvas(func(c graphics.Canvas) {

							path1 := skiagraphics.NewPath()
							path1.MoveTo(20, 20)
							path1.LineTo(20, 40)
							path1.LineTo(40, 20)

							path2 := skiagraphics.NewPath()
							path2.MoveTo(60, 60)
							path2.LineTo(80, 60)
							path2.LineTo(80, 40)

							strokeStyle := graphics.NewStroke(2)

							c.Save()
							for i := range 2 {
								testPath := skiagraphics.NewPath()
								testPath.MoveTo(20, 20)
								testPath.LineTo(20, 40)
								testPath.LineTo(40, 20)
								if i == 1 {
									testPath.Close()
								}

								testPath.AddPath(path2, geometry.OffsetZero)

								c.DrawPath(testPath, &graphics.Paint{
									Color:       graphics.ColorRed,
									StrokeWidth: 2,
									Style:       strokeStyle,
									Alpha:       1,
								})

								c.Translate(120, 0)
							}
							c.Restore()
						})
					},
					canvas.WithModifier(weight.Weight(1))),
			),
			column.WithSpacing(column.SpaceSides),
			column.WithAlignment(column.Middle),
			column.WithModifier(size.FillMax()),
		)(c)
	}
}
