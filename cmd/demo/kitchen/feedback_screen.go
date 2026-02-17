package main

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/icon"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/badge"
	"github.com/zodimo/go-compose/compose/material3/button"
	"github.com/zodimo/go-compose/compose/material3/progress"
	"github.com/zodimo/go-compose/compose/material3/snackbar"
	m3text "github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"

	mdicons "golang.org/x/exp/shiny/materialdesign/icons"
)

// FeedbackScreen shows dialog trigger, progress, badges, and snackbar
func FeedbackScreen(c api.Composer, showDialog DialogState, snackbarHostState snackbar.SnackbarHostState) api.Composable {
	progressVal := state.MustRemember(c, "fb_progress", func() float32 { return float32(0.6) })
	snackbarCount := state.MustRemember(c, "snackbar_count", func() int { return 0 })

	return func(c api.Composer) api.Composer {
		return column.Column(
			c.Sequence(
				SectionTitle("Dialog"),
				spacer.Height(8),
				row.Row(c.Sequence(
					button.Filled(func() {
						showDialog.Set(true)
					}, "Show Dialog"),
				)),

				spacer.Height(24),
				SectionTitle("Progress Indicators"),
				spacer.Height(8),
				row.Row(c.Sequence(
					progress.CircularProgressIndicator(progressVal.Get()),
					spacer.Width(16),
					progress.LinearProgressIndicator(
						progressVal.Get(),
						progress.WithModifier(size.Width(150)),
					),
				), row.WithAlignment(row.Middle)),
				spacer.Height(8),
				row.Row(c.Sequence(
					button.Text(func() {
						p := progressVal.Get() + 0.1
						if p > 1 {
							p = 0
						}
						progressVal.Set(p)
					}, "+10%"),
					button.Text(func() {
						progressVal.Set(float32(0))
					}, "Reset"),
				)),

				spacer.Height(24),
				SectionTitle("Snackbar"),
				m3text.BodySmall("Queue: messages shown one at a time (FIFO)"),
				spacer.Height(8),

				// Row 1: Short and Long duration
				row.Row(c.Sequence(
					button.Filled(func() {
						count := snackbarCount.Get() + 1
						snackbarCount.Set(count)
						snackbarHostState.ShowSnackbar(
							fmt.Sprintf("Short snackbar #%d (4s)", count),
						)
					}, "Short"),
					spacer.Width(8),
					button.Outlined(func() {
						count := snackbarCount.Get() + 1
						snackbarCount.Set(count)
						snackbarHostState.ShowSnackbar(
							fmt.Sprintf("Long snackbar #%d (10s)", count),
							snackbar.WithDuration(snackbar.SnackbarDurationLong),
						)
					}, "Long"),
				)),
				spacer.Height(8),

				// Row 2: With Action and Queue 3
				row.Row(c.Sequence(
					button.FilledTonal(func() {
						count := snackbarCount.Get() + 1
						snackbarCount.Set(count)
						snackbarHostState.ShowSnackbar(
							fmt.Sprintf("Action snackbar #%d", count),
							snackbar.WithActionLabel("Undo"),
							snackbar.WithOnResult(func(result snackbar.SnackbarResult) {
								if result == snackbar.SnackbarActionPerformed {
									fmt.Println("Undo action performed!")
								}
							}),
						)
					}, "With Action"),
					spacer.Width(8),
					button.Text(func() {
						for i := 1; i <= 3; i++ {
							count := snackbarCount.Get() + 1
							snackbarCount.Set(count)
							snackbarHostState.ShowSnackbar(
								fmt.Sprintf("Queued %d of 3 (msg #%d)", i, count),
							)
						}
					}, "Queue 3"),
				)),
				spacer.Height(8),

				// Row 3: With Dismiss and Action + Dismiss
				row.Row(c.Sequence(
					button.Outlined(func() {
						count := snackbarCount.Get() + 1
						snackbarCount.Set(count)
						snackbarHostState.ShowSnackbar(
							fmt.Sprintf("Dismiss snackbar #%d", count),
							snackbar.WithDismissAction(),
							snackbar.WithDuration(snackbar.SnackbarDurationIndefinite),
							snackbar.WithOnResult(func(result snackbar.SnackbarResult) {
								fmt.Printf("Dismiss snackbar result: %d\n", result)
							}),
						)
					}, "With Dismiss"),
					spacer.Width(8),
					button.Filled(func() {
						count := snackbarCount.Get() + 1
						snackbarCount.Set(count)
						snackbarHostState.ShowSnackbar(
							fmt.Sprintf("Full snackbar #%d", count),
							snackbar.WithActionLabel("Retry"),
							snackbar.WithDismissAction(),
							snackbar.WithOnResult(func(result snackbar.SnackbarResult) {
								if result == snackbar.SnackbarActionPerformed {
									fmt.Println("Retry action performed!")
								} else {
									fmt.Println("Full snackbar dismissed")
								}
							}),
						)
					}, "Action + Dismiss"),
				)),

				spacer.Height(24),
				SectionTitle("Badges"),
				spacer.Height(8),
				row.Row(c.Sequence(
					badge.BadgedBox(
						badge.Badge(badge.WithContent(m3text.LabelSmall("5"))),
						icon.Icon(mdicons.SocialNotifications),
					),
					spacer.Width(24),
					badge.BadgedBox(
						badge.Badge(badge.WithContent(m3text.LabelSmall("99"))),
						icon.Icon(mdicons.CommunicationEmail),
					),
					spacer.Width(24),
					badge.BadgedBox(
						badge.Badge(), // Small dot badge
						icon.Icon(mdicons.ActionShoppingCart),
					),
				), row.WithAlignment(row.Middle)),
			),
			column.WithModifier(padding.All(16)),
		)(c)
	}
}
