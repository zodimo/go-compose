package main

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/icon"
	"github.com/zodimo/go-compose/compose/foundation/lazy"
	"github.com/zodimo/go-compose/compose/material3/appbar"
	"github.com/zodimo/go-compose/compose/material3/button"
	"github.com/zodimo/go-compose/compose/material3/dialog"
	"github.com/zodimo/go-compose/compose/material3/navigationbar"
	"github.com/zodimo/go-compose/compose/material3/scaffold"
	"github.com/zodimo/go-compose/compose/material3/snackbar"
	m3text "github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"

	mdicons "golang.org/x/exp/shiny/materialdesign/icons"
)

// Navigation categories
const (
	CategoryActions    = 0
	CategorySelection  = 1
	CategoryFeedback   = 2
	CategoryInputs     = 3
	CategoryTypography = 4
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
		{"Typography", mdicons.ActionLabel},
	}

	c = c.Sequence(
		// Scaffold with navigation
		scaffold.Scaffold(
			// Content area based on selected category
			lazy.LazyColumn(
				func(scope lazy.LazyListScope) {
					scope.Item(nil, func(c api.Composer) api.Composer {
						return c.Sequence(
							//use Lazy to not invoke the composlable and keep it in memory
							c.WhenLazy(currentCategory == CategoryActions, func() api.Composable { return ActionsScreen(c) }),
							c.WhenLazy(currentCategory == CategorySelection, func() api.Composable { return SelectionScreen(c) }),
							c.WhenLazy(currentCategory == CategoryFeedback, func() api.Composable { return FeedbackScreen(c, showDialog, snackbarHostState) }),
							c.WhenLazy(currentCategory == CategoryInputs, func() api.Composable { return InputsScreen(c) }),
							c.WhenLazy(currentCategory == CategoryTypography, func() api.Composable { return TypographyScreen(c) }),
						)(c)
					})
				},
				lazy.WithModifier(weight.Weight(1).Then(size.FillMaxWidth())),
			),
			scaffold.WithTopBar(
				appbar.TopAppBar(
					m3text.TextWithStyle(
						"Component Showcase",
						m3text.TypestyleTitleLarge,
					),
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
									return m3text.TextWithStyle(item.Label, m3text.TypestyleLabelMedium)(c)
								},
							)(c)
						}
						return c
					},
				),
			),
			scaffold.WithModifier(size.FillMax()),
		),
		// Dialog overlay
		c.When(showDialog.Get().(bool),
			dialog.AlertDialog(
				func() {
					fmt.Println("Dialog Dismiss requested")
					showDialog.Set(false)
				},
				button.Text(func() {
					fmt.Println("Dialog Confirm button clicked")
					showDialog.Set(false)
				}, "Confirm"),
				dialog.TextContent("This is an example AlertDialog demonstrating the Feedback category."),
				dialog.WithTitleText("Example Dialog"),
				dialog.WithDismissButton(button.Text(
					func() {
						fmt.Println("Dialog Dismis button clicked")
						showDialog.Set(false)
					}, "Cancel")),
			),
		),
		// Snackbar host overlay
		snackbar.SnackbarHost(snackbarHostState),
	)(c)

	return c.Build()
}

// SectionTitle is a helper for section headers
func SectionTitle(title string) api.Composable {
	return m3text.TitleLarge(title)
}

// DialogState is passed to FeedbackScreen
type DialogState = state.MutableValue
