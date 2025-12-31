package main

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	foundationText "github.com/zodimo/go-compose/compose/foundation/next/text"
	composeUIText "github.com/zodimo/go-compose/compose/ui/next/text"
	"github.com/zodimo/go-compose/pkg/api"
)

func UI() api.Composable {
	return func(c api.Composer) api.Composer {

		return column.Column(
			c.Sequence(
				// Title
				foundationText.Text("BasicText Demo (Next)"),
				spacer.Height(24),

				// Section: Basic Text
				foundationText.Text("Basic Text"),
				spacer.Height(8),
				foundationText.BasicText(
					composeUIText.NewAnnotatedString("Hello, World! This is BasicText using the new bridge architecture.", nil, nil),
				),
				spacer.Height(16),

				// Section: Single Line
				foundationText.Text("Single Line (Truncated)"),
				spacer.Height(8),
				foundationText.BasicText(
					composeUIText.NewAnnotatedString("This is a very long text that should be truncated to a single line because we set maxLines to 1. The remaining text should show an ellipsis.", nil, nil),
					foundationText.WithMaxLines(1),
					foundationText.WithSoftWrap(true),
				),
				spacer.Height(16),

				// Section: Multi-line
				foundationText.Text("Multi-line (Max 2 Lines)"),
				spacer.Height(8),
				foundationText.BasicText(
					composeUIText.NewAnnotatedString("This is a long paragraph that demonstrates multi-line text support. It should wrap to multiple lines but be limited to 2 lines maximum, with truncation applied after that.", nil, nil),
					foundationText.WithMaxLines(2),
					foundationText.WithSoftWrap(true),
				),
				spacer.Height(16),

				// Section: Unlimited Lines
				foundationText.Text("Unlimited Lines"),
				spacer.Height(8),
				foundationText.BasicText(
					composeUIText.NewAnnotatedString("This text has no line limit and should wrap naturally. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.", nil, nil),
				),
				spacer.Height(24),

				// Footer
				foundationText.Text("✓ Using TextLayoutController bridge"),
			),
		)(c)
	}

}
