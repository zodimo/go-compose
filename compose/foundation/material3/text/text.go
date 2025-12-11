package text

import (
	"go-compose-dev/compose/foundation/text"

	"go-compose-dev/internal/layoutnode"
	"go-compose-dev/internal/modifier"
	"go-compose-dev/pkg/api"

	"git.sr.ht/~schnwalter/gio-mw/token"
	"git.sr.ht/~schnwalter/gio-mw/wdk"
)

// Text displays text with the given style from the Material Theme.
// It retrieves the theme from the layout context at runtime.
func Text(value string, style Typestyle, options ...text.TextOption) api.Composable {
	return func(c api.Composer) api.Composer {
		// Resolve options
		opts := text.DefaultTextOptions()
		for _, option := range options {
			option(&opts)
		}

		c.StartBlock("Material3Text")
		c.Modifier(func(modifier modifier.Modifier) modifier.Modifier {
			return modifier.Then(opts.Modifier)
		})

		c.SetWidgetConstructor(layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
			return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {

				// Map options to LabelStyle
				labelStyle := wdk.LabelStyle{
					Typestyle:  style,
					Alignment:  opts.Alignment,
					MaxLines:   opts.MaxLines,
					WrapPolicy: opts.WrapPolicy,
					// Color is tricky: foundation text options uses NRGBA, wdk uses token.MatColor.
					// If the user specified a color, use it. Otherwise wdk.LayoutLabel defaults to theme color.
				}

				if opts.TextStyleOptions != nil && opts.TextStyleOptions.Color.A > 0 {
					labelStyle.Color = token.MatColor(opts.TextStyleOptions.Color)
				}

				return wdk.LayoutLabel(gtx, labelStyle, value)
			}
		}))
		return c.EndBlock()
	}
}
