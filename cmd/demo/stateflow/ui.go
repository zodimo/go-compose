package main

import (
	"fmt"
	"time"

	"gioui.org/layout"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/button"
	m3text "github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/pkg/flow"
)

// Global state flow for demonstration
var counterFlow = flow.NewMutableStateFlow(0)
var counterFlow2 = flow.NewMutableStateFlow(0)

func init() {
	// Start a background ticker to update the flow
	go func() {
		for {
			time.Sleep(1 * time.Second)
			counterFlow.Update(func(current int) int {
				return current + 1
			})
		}
	}()

	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			counterFlow2.Update(func(current int) int {
				return current + 1
			})
		}
	}()
}

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		return box.Box(
			func(c api.Composer) api.Composer {
				// Collect the flow as state inside the layout context (Box)
				// This ensures the LaunchedEffect node is a child of Box, not the Root.
				countState := flow.CollectStateFlowAsState(c, "counter1", counterFlow)
				count := countState.Get()

				countState2 := flow.CollectStateFlowAsState(c, "counter2", counterFlow2)
				count2 := countState2.Get()

				return column.Column(
					c.Sequence(
						m3text.HeadlineMedium("StateFlow Integration Demo"),
						m3text.BodyLarge(fmt.Sprintf("Current Count (from flow: 1s interval): %d", count)),
						row.Row(
							c.Sequence(
								button.Filled(func() {
									counterFlow.Emit(0)
								}, "Reset"),
								spacer.Width(5),
								button.Filled(func() {
									counterFlow.Update(func(curr int) int { return curr + 5 })
								}, "Add 5"),
							),
							row.WithModifier(padding.Vertical(20, 0)),
						),
						spacer.Height(20),
						m3text.BodyLarge(fmt.Sprintf("Current Count (from flow2: 500ms interval): %d", count2)),
						row.Row(
							c.Sequence(
								button.Filled(func() {
									counterFlow2.Emit(0)
								}, "Reset"),
								spacer.Width(5),
								button.Filled(func() {
									counterFlow2.Update(func(curr int) int { return curr + 5 })
								}, "Add 5"),
							),
							row.WithModifier(padding.Vertical(20, 0)),
						),
					),
					column.WithModifier(
						padding.All(20),
					),
					column.WithAlignment(layout.Middle),
				)(c)
			},
			box.WithModifier(size.FillMax()),
			box.WithAlignment(layout.Center),
		)(c)
	}
}
