package main

import (
	"fmt"
	"image/color"

	"github.com/zodimo/go-compose/compose/foundation/icon"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/material3/segmentedbutton"
	"github.com/zodimo/go-compose/compose/foundation/material3/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/pkg/api"

	mdicons "golang.org/x/exp/shiny/materialdesign/icons"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		// State for single-select (0, 1, or 2)
		singleSelectState := c.State("single_select_index", func() any { return 0 })
		selectedIndex := singleSelectState.Get().(int)

		// State for multi-select (map of checked states)
		multiSelectState := c.State("multi_select_map", func() any {
			return map[int]bool{0: true, 1: false, 2: false, 3: false}
		})
		checkedMap := multiSelectState.Get().(map[int]bool)

		// Single-select options
		singleOptions := []string{"Day", "Week", "Month"}

		// Multi-select options
		multiOptions := []string{"XS", "S", "M", "L"}

		// Checkmark icon for selected state
		checkIcon := icon.Icon(mdicons.NavigationCheck)

		return column.Column(
			c.Sequence(
				// Title
				text.Text("Segmented Button Demo", text.TypestyleTitleLarge),

				spacer.Height(16),

				// Section: Single-Select
				text.Text("Single Select", text.TypestyleTitleMedium),
				text.Text(fmt.Sprintf("Selected: %s", singleOptions[selectedIndex]), text.TypestyleBodyMedium),

				spacer.Height(8),

				segmentedbutton.SingleChoiceSegmentedButtonRow(
					c.Range(len(singleOptions), func(i int) api.Composable {
						label := singleOptions[i]
						shape := getShapeForIndex(i, len(singleOptions))
						return segmentedbutton.SegmentedButton(
							selectedIndex == i,
							func(checked bool) {
								if checked {
									singleSelectState.Set(i)
								}
							},
							label,
							shape,
							segmentedbutton.WithSelectedColor(color.NRGBA{R: 0xE8, G: 0xDE, B: 0xF8, A: 0xFF}),
							segmentedbutton.WithSelectedIcon(checkIcon),
						)
					}),
				),

				spacer.Height(24),

				// Section: Multi-Select
				text.Text("Multi Select", text.TypestyleTitleMedium),

				spacer.Height(8),

				segmentedbutton.MultiChoiceSegmentedButtonRow(
					c.Range(len(multiOptions), func(i int) api.Composable {
						label := multiOptions[i]
						shape := getShapeForIndex(i, len(multiOptions))
						isChecked := checkedMap[i]
						return segmentedbutton.SegmentedButton(
							isChecked,
							func(checked bool) {
								newMap := make(map[int]bool)
								for k, v := range checkedMap {
									newMap[k] = v
								}
								newMap[i] = checked
								multiSelectState.Set(newMap)
							},
							label,
							shape,
							segmentedbutton.WithSelectedColor(color.NRGBA{R: 0xE8, G: 0xDE, B: 0xF8, A: 0xFF}),
							segmentedbutton.WithSelectedIcon(checkIcon),
						)
					}),
				),
			),
			column.WithModifier(padding.All(24)),
		)(c)
	}
}

// getShapeForIndex returns the appropriate segment shape based on position.
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
