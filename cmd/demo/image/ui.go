package main

import (
	"embed"
	"image/color"
	_ "image/png"

	fImage "github.com/zodimo/go-compose/compose/foundation/image"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	uilayout "github.com/zodimo/go-compose/compose/ui/layout"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/theme"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/surface"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/modifiers/border"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
)

//go:embed gopher.png
var assets embed.FS

func UI() api.Composable {

	imageResource := graphics.NewResourceFromImageFS(assets, "gopher.png")

	return func(c api.Composer) api.Composer {
		return column.Column(
			func(c api.Composer) api.Composer {
				// 1. Basic Image (Fit)
				text.TextWithStyle("ContentScale: Fit", text.TypestyleBodyLarge)(c)
				spacer.Height(8)(c)
				surface.Surface(
					fImage.Image(
						imageResource,
						fImage.WithContentScale(uilayout.ContentScaleFit),
					),
					surface.WithModifier(
						size.Size(150, 100, size.SizeRequired()).
							Then(
								border.Border(
									unit.Dp(1),
									theme.ColorHelper.SpecificColor(graphics.FromNRGBA(color.NRGBA{0, 0, 0, 255})),
									shape.ShapeRectangle,
								),
							),
					),
				)(c)
				spacer.Height(16)(c)

				// 2. Crop
				text.TextWithStyle("ContentScale: Crop", text.TypestyleBodyLarge)(c)
				spacer.Height(8)(c)
				surface.Surface(
					fImage.Image(
						imageResource,
						fImage.WithContentScale(uilayout.ContentScaleCrop),
					),
					surface.WithModifier(
						size.Size(150, 100, size.SizeRequired()).Then(
							border.Border(
								unit.Dp(1),
								theme.ColorHelper.SpecificColor(graphics.FromNRGBA(color.NRGBA{0, 0, 0, 255})),
								shape.ShapeRectangle,
							),
						),
					),
				)(c)
				spacer.Height(16)(c)

				// 3. Alignment (BottomEnd) + Alpha
				text.TextWithStyle("Alignment: BottomEnd + Alpha 0.5", text.TypestyleBodyLarge)(c)
				spacer.Height(8)(c)
				surface.Surface(
					fImage.Image(
						imageResource,
						fImage.WithContentScale(uilayout.ContentScaleNone),
						fImage.WithAlignment(size.BottomEnd),
						fImage.WithAlpha(0.5),
						fImage.WithModifier(size.Size(50, 50)),
					),
					surface.WithModifier(
						size.Size(150, 100, size.SizeRequired()).Then(
							border.Border(
								unit.Dp(1),
								theme.ColorHelper.SpecificColor(graphics.FromNRGBA(color.NRGBA{0, 0, 0, 255})),
								shape.ShapeRectangle,
							),
						),
					),
				)(c)

				return c
			},
			column.WithModifier(
				size.FillMax().Then(padding.All(16)),
			),
		)(c)
	}
}
