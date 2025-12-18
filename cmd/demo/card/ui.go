package main

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/paint"
	gioWidget "gioui.org/widget"

	fImage "github.com/zodimo/go-compose/compose/foundation/image"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/lazy"
	"github.com/zodimo/go-compose/compose/foundation/material3/button"
	"github.com/zodimo/go-compose/compose/foundation/material3/card"
	m3text "github.com/zodimo/go-compose/compose/foundation/material3/text"
	"github.com/zodimo/go-compose/compose/foundation/material3/textfield"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	uilayout "github.com/zodimo/go-compose/compose/ui/layout"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/clickable"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/theme"
)

func UI(c api.Composer) api.LayoutNode {
	// State for interactive card demo
	inputValue := c.State("card_input", func() any { return "" })

	c = lazy.LazyColumn(
		func(scope lazy.LazyListScope) {
			// Title
			scope.Item("title", column.Column(
				c.Sequence(
					m3text.Text("Card Component Demo", m3text.TypestyleHeadlineMedium),
					spacer.Height(24),
				),
			))

			// Elevated Card
			scope.Item("elevated", column.Column(
				c.Sequence(
					SectionTitle("Elevated Card"),
					spacer.Height(8),
					card.Elevated(
						card.CardContents(
							card.Content(
								CardContent(
									"Elevated Card",
									"Elevated cards use shadow to create depth. They're best for\nscanning and organizing content.",
								),
							),
						),
						card.WithModifier(size.Width(340, size.SizeRequired())),
					),
					spacer.Height(24),
				),
			))

			// Filled Card
			scope.Item("filled", column.Column(
				c.Sequence(
					SectionTitle("Filled Card"),
					spacer.Height(8),
					card.Filled(
						card.CardContents(
							card.Content(
								CardContent(
									"Filled Card",
									"Filled cards use color to create depth. They use a subtle\nbackground color distinct from the surface.",
								),
							),
						),
						card.WithModifier(size.Width(340, size.SizeRequired())),
					),
					spacer.Height(24),
				),
			))

			// Outlined Card
			scope.Item("outlined", column.Column(
				c.Sequence(
					SectionTitle("Outlined Card"),
					spacer.Height(8),
					card.Outlined(
						card.CardContents(
							card.Content(
								CardContent(
									"Outlined Card",
									"Outlined cards use a border instead of shadow. They're best\nfor grouping related content.",
								),
							),
						),
						card.WithModifier(size.Width(340, size.SizeRequired())),
					),
					spacer.Height(24),
				),
			))

			// Interactive Card - demonstrates Tab navigation works
			scope.Item("interactive", column.Column(
				c.Sequence(
					SectionTitle("Interactive Card (Tab Navigation Demo)"),
					spacer.Height(8),
					card.Elevated(
						card.CardContents(
							card.Content(
								InteractiveCardContent(inputValue),
							),
						),
						card.WithModifier(size.Width(340, size.SizeRequired())),
					),
					spacer.Height(24),
				),
			))

			// ContentCover Example - full-bleed header
			scope.Item("cover", column.Column(
				c.Sequence(
					SectionTitle("Card with ContentCover"),
					m3text.Text("ContentCover provides full-bleed content without padding", m3text.TypestyleBodySmall),
					spacer.Height(8),
					card.Elevated(
						card.CardContents(
							// Cover at top - no internal padding
							card.ContentCover(
								CoverBanner("Featured", "Primary color header"),
							),
							// Regular content below
							card.Content(
								CardContent(
									"Main Content",
									"This content has padding while the cover above extends edge-to-edge.",
								),
							),
						),
						card.WithModifier(size.Width(340, size.SizeRequired())),
					),
					spacer.Height(24),
				),
			))

			// Multiple Content Sections
			scope.Item("multiple", column.Column(
				c.Sequence(
					SectionTitle("Card with Multiple Sections"),
					m3text.Text("Cards can have multiple Content sections", m3text.TypestyleBodySmall),
					spacer.Height(8),
					card.Filled(
						card.CardContents(
							card.Content(
								CardContent("Section 1", "First content section"),
							),
							card.Content(
								CardContent("Section 2", "Second content section stacked below."),
							),
							card.Content(
								ActionButtons(),
							),
						),
						card.WithModifier(size.Width(340, size.SizeRequired())),
					),
					spacer.Height(24),
				),
			))

			// Card with Image content type
			scope.Item("image", column.Column(
				c.Sequence(
					SectionTitle("Card with Image"),
					m3text.Text("Using go-compose Image within ContentCover", m3text.TypestyleBodySmall),
					spacer.Height(8),
					card.Elevated(
						card.CardContents(
							// Use ContentCover with go-compose Image composable
							card.ContentCover(
								fImage.Image(
									CreateImageResource(),
									fImage.WithContentScale(uilayout.ContentScaleCrop),
									fImage.WithModifier(size.Height(120)),
								),
							),
							// Content below
							card.Content(
								CardContent(
									"Media Card",
									"Image content is displayed edge-to-edge at the card's width.",
								),
							),
						),
						card.WithModifier(size.Width(340, size.SizeRequired())),
					),
					spacer.Height(24),
				),
			))

			// Clickable Cards - demonstrates hover effect with onClick
			scope.Item("clickable", column.Column(
				c.Sequence(
					SectionTitle("Clickable Cards (Hover Effect Demo)"),
					m3text.Text("Cards with onClick show hover effect clipped to rounded corners", m3text.TypestyleBodySmall),
					spacer.Height(8),
					row.Row(
						c.Sequence(
							card.Elevated(
								card.CardContents(
									card.Content(SmallCardContent("Click Me")),
								),
								card.WithModifier(
									size.Width(120).Then(clickable.OnClick(func() {
										// Handle click
									})),
								),
							),
							spacer.Width(16),
							card.Filled(
								card.CardContents(
									card.Content(SmallCardContent("Click Me")),
								),
								card.WithModifier(
									size.Width(120).Then(clickable.OnClick(func() {
										// Handle click
									})),
								),
							),
							spacer.Width(16),
							card.Outlined(
								card.CardContents(
									card.Content(SmallCardContent("Click Me")),
								),
								card.WithModifier(
									size.Width(120).Then(clickable.OnClick(func() {
										// Handle click
									})),
								),
							),
						),
					),
					spacer.Height(24),
				),
			))

			// Cards in a Row
			scope.Item("row", column.Column(
				c.Sequence(
					SectionTitle("Cards in a Row"),
					spacer.Height(8),
					row.Row(
						c.Sequence(
							card.Elevated(
								card.CardContents(
									card.Content(SmallCardContent("Card 1")),
								),
								card.WithModifier(size.Width(150)),
							),
							spacer.Width(16),
							card.Filled(
								card.CardContents(
									card.Content(SmallCardContent("Card 2")),
								),
								card.WithModifier(size.Width(150)),
							),
						),
					),
					spacer.Height(24),
				),
			))
		},
		lazy.WithModifier(padding.All(24).Then(size.FillMax())),
	)(c)

	return c.Build()
}

// SectionTitle creates a section heading
func SectionTitle(title string) api.Composable {
	return m3text.Text(title, m3text.TypestyleTitleMedium)
}

// CardContent creates a standard card content with title and description
func CardContent(title, description string) api.Composable {
	return func(c api.Composer) api.Composer {
		return column.Column(
			c.Sequence(
				m3text.Text(title, m3text.TypestyleTitleMedium),
				spacer.Height(8),
				m3text.Text(description, m3text.TypestyleBodyMedium),
			),
		)(c)
	}
}

// SmallCardContent creates compact card content
func SmallCardContent(title string) api.Composable {
	return func(c api.Composer) api.Composer {
		return box.Box(
			m3text.Text(title, m3text.TypestyleTitleSmall),
			box.WithAlignment(layout.Center),
		)(c)
	}
}

// InteractiveCardContent creates a card with interactive elements
func InteractiveCardContent(inputValue api.MutableValue) api.Composable {
	return func(c api.Composer) api.Composer {
		return column.Column(
			c.Sequence(
				m3text.Text("Interactive Form", m3text.TypestyleTitleMedium),
				spacer.Height(12),
				m3text.Text("Tab between fields works correctly:", m3text.TypestyleBodySmall),
				spacer.Height(12),
				textfield.TextField(
					inputValue.Get().(string),
					func(s string) { inputValue.Set(s) },
					"Enter something",
					textfield.WithModifier(size.FillMaxWidth()),
				),
				spacer.Height(12),
				textfield.TextField(
					"",
					func(s string) {},
					"Another field",
					textfield.WithModifier(size.FillMaxWidth()),
				),
				spacer.Height(16),
				button.Filled(func() {}, "Submit"),
			),
		)(c)
	}
}

// CoverBanner creates a full-bleed colored banner for ContentCover
func CoverBanner(title, subtitle string) api.Composable {
	return func(c api.Composer) api.Composer {
		return box.Box(
			column.Column(
				c.Sequence(
					m3text.Text(title, m3text.TypestyleHeadlineSmall),
					m3text.Text(subtitle, m3text.TypestyleBodySmall),
				),
			),
			box.WithModifier(padding.All(24).
				Then(background.Background(theme.ColorHelper.ColorSelector().PrimaryRoles.Container)).
				Then(size.FillMaxWidth()),
			),
			box.WithAlignment(layout.NW),
		)(c)
	}
}

// ActionButtons creates a row of action buttons for cards
func ActionButtons() api.Composable {
	return func(c api.Composer) api.Composer {
		return box.Box(
			row.Row(
				c.Sequence(
					button.Text(func() {}, "Cancel"),
					spacer.Width(8),
					button.Filled(func() {}, "Confirm"),
				),
			),
			box.WithModifier(padding.Horizontal(16, 16).Then(padding.Vertical(0, 16))),
			box.WithAlignment(layout.E),
		)(c)
	}
}

// CreatePlaceholderImage creates a placeholder colored image for the demo
func CreatePlaceholderImage(width, height int) *gioWidget.Image {
	// Create a gradient-like colored image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with a nice gradient
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Create a purple gradient
			r := uint8(103 + x*40/width)
			g := uint8(80 + y*40/height)
			b := uint8(164)
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	return &gioWidget.Image{
		Src:   paint.NewImageOp(img),
		Fit:   gioWidget.Cover, // Fill the width
		Scale: 1,
	}
}

// CreateImageResource creates a go-compose ImageResource from a programmatic gradient
func CreateImageResource() graphics.ImageResource {
	width, height := 340, 120
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with a nice purple gradient
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := uint8(103 + x*40/width)
			g := uint8(80 + y*40/height)
			b := uint8(164)
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	// graphics.ImageResource expects an ImageOp
	return graphics.ImageResource{
		ImageOp: paint.NewImageOp(img),
	}
}
