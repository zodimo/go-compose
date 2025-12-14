package icon

import (
	"go-compose-dev/internal/layoutnode"
	"go-compose-dev/internal/theme"
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"
)

// icons from golang.org/x/exp/shiny/materialdesign/icons
func Icon(iconByte []byte, options ...IconOption) Composable {
	opts := DefaultIconOptions()
	for _, option := range options {
		if option == nil {
			continue
		}
		option(&opts)
	}

	return func(c Composer) Composer {
		c.StartBlock("Icon")
		c.Modifier(func(modifier Modifier) Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.SetWidgetConstructor(iconWidgetConstructor(opts, iconByte))

		return c.EndBlock()
	}
}

func iconWidgetConstructor(options IconOptions, iconByte []byte) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			iconWidget := requireIconWidget(iconByte)

			if options.ThemeColor.IsSome() {
				themeManager := theme.GetThemeManager()
				return iconWidget(gtx, options.ThemeColor.UnwrapUnsafe().ThemeColor(themeManager.ThemeColorResolver()))
			}

			return iconWidget(gtx, options.Color)
		}
	})
}

func requireIconWidget(data []byte) IconWidget {
	iconWidget, err := widget.NewIcon(data)
	if err != nil {
		panic(err)
	}
	return func(gtx layout.Context, foreground color.Color) layout.Dimensions {
		return iconWidget.Layout(gtx, ToNRGBA(foreground))
	}
}
