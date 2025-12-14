package main

import (
	"fmt"
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/material3/scaffold"
	"github.com/zodimo/go-compose/compose/foundation/material3/slider"
	m3text "github.com/zodimo/go-compose/compose/foundation/material3/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/pkg/api"
)

type Composable = api.Composable
type Composer = api.Composer
type Modifier = api.Modifier

func UI() Composable {
	return func(c Composer) Composer {
		// State for sliders
		basicValue := c.State("basic", func() any { return float32(0.5) })
		steppedValue := c.State("stepped", func() any { return float32(20.0) })
		rangeMappedValue := c.State("range", func() any { return float32(50.0) })
		customColorValue := c.State("custom", func() any { return float32(0.3) })

		return scaffold.Scaffold(
			func(c Composer) Composer {
				return column.Column(
					func(c Composer) Composer {
						// 1. Basic Continuous Slider
						Label("Continuous Slider (0.0 - 1.0)")(c)
						slider.Slider(
							basicValue.Get().(float32),
							func(v float32) { basicValue.Set(v) },
							slider.WithOnValueChangeFinished(func() {
								fmt.Printf("Basic Slider Finished: %.2f\n", basicValue.Get().(float32))
							}),
						)(c)
						ValueText(basicValue.Get().(float32))(c)

						spacer.SpacerHeight(24)(c)

						// 2. Stepped Slider (0-100, 5 steps)
						Label("Stepped Slider (0-100, 5 steps)")(c)
						slider.Slider(
							steppedValue.Get().(float32),
							func(v float32) { steppedValue.Set(v) },
							slider.WithSteps(4),
							slider.WithValueRange(0, 100),
						)(c)
						ValueText(steppedValue.Get().(float32))(c)

						spacer.SpacerHeight(24)(c)

						// 3. Slider with mapped range (0-100) continuous
						Label("Mapped Range (0-100)")(c)
						slider.Slider(
							rangeMappedValue.Get().(float32),
							func(v float32) { rangeMappedValue.Set(v) },
							slider.WithValueRange(0, 100),
						)(c)
						ValueText(rangeMappedValue.Get().(float32))(c)

						spacer.SpacerHeight(24)(c)

						// 4. Custom Colors
						Label("Custom Colors")(c)
						slider.Slider(
							customColorValue.Get().(float32),
							func(v float32) { customColorValue.Set(v) },
							slider.WithColors(slider.SliderColors{
								ThumbColor:            color.NRGBA{R: 255, A: 255},
								ActiveTrackColor:      color.NRGBA{R: 200, G: 50, B: 50, A: 255},
								InactiveTrackColor:    color.NRGBA{R: 200, G: 200, B: 200, A: 255},
								ActiveTickColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
								InactiveTickColor:     color.NRGBA{R: 0, G: 0, B: 0, A: 255},
								DisabledThumbColor:    color.NRGBA{R: 100, G: 100, B: 100, A: 255},
								DisabledActiveTrack:   color.NRGBA{R: 100, G: 100, B: 100, A: 255},
								DisabledActiveTick:    color.NRGBA{R: 100, G: 100, B: 100, A: 255},
								DisabledInactiveTrack: color.NRGBA{R: 220, G: 220, B: 220, A: 255},
								DisabledInactiveTick:  color.NRGBA{R: 220, G: 220, B: 220, A: 255},
							}),
						)(c)

						spacer.SpacerHeight(24)(c)

						// 5. Disabled Slider
						Label("Disabled Slider")(c)
						slider.Slider(
							0.5,
							func(v float32) {},
							slider.WithEnabled(false),
						)(c)

						return c
					},
					column.WithModifier(padding.All(16)),
					column.WithAlignment(column.Start),
				)(c)
			},
		)(c)
	}
}

func Label(txt string) Composable {
	return m3text.Text(txt, m3text.TypestyleBodyLarge)
}

func ValueText(val float32) Composable {
	return m3text.Text(fmt.Sprintf("Value: %.2f", val), m3text.TypestyleBodyMedium)
}
