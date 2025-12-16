package main

import (
	"github.com/zodimo/go-compose/compose/foundation/icon"
	"github.com/zodimo/go-compose/compose/foundation/layout/overlay"
	"github.com/zodimo/go-compose/compose/foundation/lazy"
	"github.com/zodimo/go-compose/compose/foundation/material3/appbar"
	"github.com/zodimo/go-compose/compose/foundation/material3/dialog"
	"github.com/zodimo/go-compose/compose/foundation/material3/navigationbar"
	"github.com/zodimo/go-compose/compose/foundation/material3/scaffold"
	"github.com/zodimo/go-compose/compose/foundation/material3/snackbar"
	m3text "github.com/zodimo/go-compose/compose/foundation/material3/text"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"

	mdicons "golang.org/x/exp/shiny/materialdesign/icons"
)

// Navigation categories
const (
	CategoryActions   = 0
	CategorySelection = 1
	CategoryFeedback  = 2
	CategoryInputs    = 3
)

func UI(c api.Composer) api.LayoutNode {
	// Navigation state
	selectedCategory := c.State("nav_category", func() any { return CategoryActions })
	currentCategory := selectedCategory.Get().(int)

	// Dialog visibility state
	showDialog := c.State("showDialog", func() any { return false })

	// Snackbar state
	snackbarHostState := c.State("snackbarHostState", func() any { return snackbar.NewSnackbarHostState() }).Get().(*snackbar.SnackbarHostState)

	navItems := []struct {
		Label string
		Icon  []byte
	}{
		{"Actions", mdicons.ActionTouchApp},
		{"Selection", mdicons.ToggleCheckBox},
		{"Feedback", mdicons.ActionFeedback},
		{"Inputs", mdicons.ActionInput},
	}

	c = c.Sequence(
		func(c api.Composer) api.Composer {
			// Scaffold with navigation
			scaffold.Scaffold(
				// Content area based on selected category
				lazy.LazyColumn(
					func(scope lazy.LazyListScope) {
						scope.Item(nil, func(c api.Composer) api.Composer {
							return c.Sequence(
								c.When(currentCategory == CategoryActions, ActionsScreen(c)),
								c.When(currentCategory == CategorySelection, SelectionScreen(c)),
								c.When(currentCategory == CategoryFeedback, FeedbackScreen(c, showDialog, snackbarHostState)),
								c.When(currentCategory == CategoryInputs, InputsScreen(c)),
							)(c)
						})
					},
					lazy.WithModifier(weight.Weight(1)),
					lazy.WithModifier(size.FillMaxWidth()),
				),
				scaffold.WithTopBar(
					appbar.TopAppBar(
						m3text.Text("Component Showcase", m3text.TypestyleTitleLarge),
					),
				),
				scaffold.WithBottomBar(
					navigationbar.NavigationBar(
						func(c api.Composer) api.Composer {
							for i, item := range navItems {
								idx := i
								navigationbar.NavigationBarItem(
									currentCategory == idx,
									func() { selectedCategory.Set(idx) },
									func(c api.Composer) api.Composer {
										return icon.Icon(item.Icon)(c)
									},
									func(c api.Composer) api.Composer {
										return m3text.Text(item.Label, m3text.TypestyleLabelMedium)(c)
									},
								)(c)
							}
							return c
						},
					),
				),
				scaffold.WithModifier(size.FillMax()),
			)(c)

			// Dialog overlay
			c.When(showDialog.Get().(bool), overlay.Overlay(
				dialog.AlertDialog(
					func() { showDialog.Set(false) },
					func() { showDialog.Set(false) },
					"Confirm",
					dialog.WithTitle("Example Dialog"),
					dialog.WithText("This is an example AlertDialog demonstrating the Feedback category."),
					dialog.WithDismissButton("Cancel", func() {
						showDialog.Set(false)
					}),
				),
				overlay.WithOnDismiss(func() {
					showDialog.Set(false)
				}),
			))

			// Snackbar host overlay
			snackbar.SnackbarHost(snackbarHostState)(c)

			return c
		},
	)(c)

	return c.Build()
}

// SectionTitle is a helper for section headers
func SectionTitle(title string) api.Composable {
	return func(c api.Composer) api.Composer {
		m3text.Text(title, m3text.TypestyleTitleMedium)(c)
		return c
	}
}

// DialogState is passed to FeedbackScreen
type DialogState = state.MutableValue
