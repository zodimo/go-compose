package text

import (
	"github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/platform"
	"github.com/zodimo/go-compose/compose/ui/text/style"

	uiText "github.com/zodimo/go-compose/compose/ui/text"

	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/pkg/api"

	"git.sr.ht/~schnwalter/gio-mw/token"
	"git.sr.ht/~schnwalter/gio-mw/wdk"
)

// TextWithStyle displays text with the given style from the Material Theme.
// It retrieves the theme from the layout context at runtime.
func TextWithStyle(value string, tokenStyle Typestyle, options ...text.TextOption) api.Composable {
	return func(c api.Composer) api.Composer {
		// Resolve options
		opts := text.DefaultTextOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		contentColor := material3.LocalContentColor.Current(c)

		// textShaper := compose.LocalTextShaper.Current(c)
		// familyResolver := platform.LocalFontFamilyResolver.Current(c)
		layoutDirection := platform.LocalLayoutDirection.Current(c)

		textColor := graphics.ColorToNRGBA(opts.TextStyle.Color().TakeOrElse(contentColor))

		// Resolve text style with defaults
		opts.TextStyle = uiText.TextStyleResolveDefaults(opts.TextStyle, layoutDirection)

		c.StartBlock("Material3Text")
		c.Modifier(func(modifier modifier.Modifier) modifier.Modifier {
			return modifier.Then(opts.Modifier)
		})

		c.SetWidgetConstructor(layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
			return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {

				// Map options to LabelStyle
				labelStyle := wdk.LabelStyle{
					Typestyle:  tokenStyle,
					Alignment:  style.TextAlignToGioTextAlignment(opts.TextStyle.TextAlign()),
					MaxLines:   opts.MaxLines,
					WrapPolicy: style.LineBreakToGioWrapPolicy(opts.TextStyle.LineBreak()),
				}

				labelStyle.Color = token.MatColor(textColor)

				return wdk.LayoutLabel(gtx, labelStyle, value)
			}
		}))
		return c.EndBlock()
	}
}
