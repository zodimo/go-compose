package main

import (
	"github.com/zodimo/go-compose/compose/foundation/icon"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/material3/badge"
	"github.com/zodimo/go-compose/compose/foundation/material3/button"
	"github.com/zodimo/go-compose/compose/foundation/material3/progress"
	"github.com/zodimo/go-compose/compose/foundation/material3/snackbar"
	m3text "github.com/zodimo/go-compose/compose/foundation/material3/text"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"

	mdicons "golang.org/x/exp/shiny/materialdesign/icons"
)

// SnackbarHostState for managing snackbar display
type SnackbarHostState = snackbar.SnackbarHostState

// FeedbackScreen shows dialog trigger, progress, badges, and snackbar
func FeedbackScreen(c api.Composer, showDialog DialogState, snackbarHostState *SnackbarHostState) api.Composable {
	progressVal := c.State("fb_progress", func() any { return float32(0.6) })

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
					progress.CircularProgressIndicator(progressVal.Get().(float32)),
					spacer.Width(16),
					progress.LinearProgressIndicator(
						progressVal.Get().(float32),
						progress.WithModifier(size.Width(150)),
					),
				), row.WithAlignment(row.Middle)),
				spacer.Height(8),
				row.Row(c.Sequence(
					button.Text(func() {
						p := progressVal.Get().(float32) + 0.1
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
				spacer.Height(8),
				row.Row(c.Sequence(
					button.Filled(func() {
						snackbarHostState.ShowSnackbar("Hello from Snackbar!")
					}, "Show Snackbar"),
				)),

				spacer.Height(24),
				SectionTitle("Badges"),
				spacer.Height(8),
				row.Row(c.Sequence(
					badge.BadgedBox(
						badge.Badge(badge.WithContent(m3text.Text("5", m3text.TypestyleLabelSmall))),
						icon.Icon(mdicons.SocialNotifications),
					),
					spacer.Width(24),
					badge.BadgedBox(
						badge.Badge(badge.WithContent(m3text.Text("99", m3text.TypestyleLabelSmall))),
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
