package main

import (
	"fmt"
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/icon"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/material3/checkbox"
	"github.com/zodimo/go-compose/compose/foundation/material3/chip"
	"github.com/zodimo/go-compose/compose/foundation/material3/radiobutton"
	"github.com/zodimo/go-compose/compose/foundation/material3/segmentedbutton"
	"github.com/zodimo/go-compose/compose/foundation/material3/slider"
	mswitch "github.com/zodimo/go-compose/compose/foundation/material3/switch"
	m3text "github.com/zodimo/go-compose/compose/foundation/material3/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/theme"

	mdicons "golang.org/x/exp/shiny/materialdesign/icons"
)

// SelectionScreen shows checkbox, radio, switch, slider
func SelectionScreen(c api.Composer) api.Composable {
	isChecked := c.State("sel_checked", func() any { return false })
	isSwitched := c.State("sel_switched", func() any { return true })
	radioOption := c.State("sel_radio", func() any { return 0 })
	sliderValue := c.State("sel_slider", func() any { return float32(0.5) })

	return column.Column(
		c.Sequence(
			SectionTitle("Checkbox"),
			spacer.Height(8),
			row.Row(c.Sequence(
				checkbox.Checkbox(isChecked.Get().(bool), func(b bool) {
					isChecked.Set(b)
				}),
				spacer.Width(8),
				m3text.Text(fmt.Sprintf("Checked: %v", isChecked.Get().(bool)), m3text.TypestyleBodyMedium),
			), row.WithAlignment(row.Middle)),

			spacer.Height(16),
			SectionTitle("Switch"),
			spacer.Height(8),
			row.Row(c.Sequence(
				mswitch.Switch(isSwitched.Get().(bool), func(b bool) {
					isSwitched.Set(b)
				}),
				spacer.Width(8),
				m3text.Text(fmt.Sprintf("On: %v", isSwitched.Get().(bool)), m3text.TypestyleBodyMedium),
			), row.WithAlignment(row.Middle)),

			spacer.Height(16),
			SectionTitle("Radio Buttons"),
			spacer.Height(8),
			column.Column(c.Sequence(
				row.Row(c.Sequence(
					radiobutton.RadioButton(radioOption.Get().(int) == 0, func() { radioOption.Set(0) }),
					spacer.Width(8),
					m3text.Text("Option A", m3text.TypestyleBodyMedium),
				), row.WithAlignment(row.Middle)),
				row.Row(c.Sequence(
					radiobutton.RadioButton(radioOption.Get().(int) == 1, func() { radioOption.Set(1) }),
					spacer.Width(8),
					m3text.Text("Option B", m3text.TypestyleBodyMedium),
				), row.WithAlignment(row.Middle)),
				row.Row(c.Sequence(
					radiobutton.RadioButton(radioOption.Get().(int) == 2, func() { radioOption.Set(2) }),
					spacer.Width(8),
					m3text.Text("Option C", m3text.TypestyleBodyMedium),
				), row.WithAlignment(row.Middle)),
			)),

			spacer.Height(16),
			SectionTitle("Slider"),
			spacer.Height(8),
			slider.Slider(
				sliderValue.Get().(float32),
				func(v float32) { sliderValue.Set(v) },
			),
			m3text.Text(fmt.Sprintf("Value: %.2f", sliderValue.Get().(float32)), m3text.TypestyleBodySmall),

			spacer.Height(24),
			SectionTitle("Segmented Button"),
			spacer.Height(8),
			func(c api.Composer) api.Composer {
				selectedIndex := c.State("seg_index", func() any { return 0 }).Get().(int)
				options := []string{"Day", "Week", "Month"}
				checkIcon := icon.Icon(mdicons.NavigationCheck)

				return segmentedbutton.SingleChoiceSegmentedButtonRow(
					c.Range(len(options), func(i int) api.Composable {
						return segmentedbutton.SegmentedButton(
							selectedIndex == i,
							func(checked bool) {
								if checked {
									c.State("seg_index", nil).Set(i)
								}
							},
							options[i],
							getShapeForIndex(i, len(options)),
							segmentedbutton.WithSelectedColor(theme.ColorHelper.SpecificColor(color.NRGBA{R: 0xE8, G: 0xDE, B: 0xF8, A: 0xFF})),
							segmentedbutton.WithSelectedIcon(checkIcon),
						)
					}),
				)(c)
			},

			spacer.Height(24),
			SectionTitle("Chips"),
			spacer.Height(8),
			row.Row(c.Sequence(
				chip.AssistChip(func() { fmt.Println("Assist") }, "Assist"),
				spacer.Width(8),
				func(c api.Composer) api.Composer {
					selected := c.State("chip_sel", func() any { return false })
					label := "Filter"
					if selected.Get().(bool) {
						label = "Selected"
					}
					return chip.FilterChip(
						func() { selected.Set(!selected.Get().(bool)) },
						label,
						chip.WithSelected(selected.Get().(bool)),
					)(c)
				},
				spacer.Width(8),
				chip.InputChip(func() {}, "Input"),
			)),
		),
		column.WithModifier(padding.All(16)),
	)
}

func getShapeForIndex(index, total int) segmentedbutton.SegmentShape {
	if total == 1 {
		return segmentedbutton.SegmentShapeOnly
	}
	if index == 0 {
		return segmentedbutton.SegmentShapeStart
	}
	if index == total-1 {
		return segmentedbutton.SegmentShapeEnd
	}
	return segmentedbutton.SegmentShapeMiddle
}
