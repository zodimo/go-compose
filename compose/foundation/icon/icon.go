package icon

import (
	"go-compose-dev/internal/layoutnode"
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"
)

var FallbackColorDescriptor = ThemeManager.ColorRoleDescriptors().BasicColorRoleDescriptors.BasicFg

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

			colorDescriptor := FallbackColorDescriptor
			if options.Color.IsSome() {
				colorDescriptor = options.Color.UnwrapUnsafe()
			}
			if options.LazyColor.IsSome() {
				colorDescriptor = options.LazyColor.UnwrapUnsafe()()
			}

			themeColor := ThemeManager.ResolveColorDescriptor(colorDescriptor)

			return iconWidget(gtx, themeColor.AsNRGBA())
		}
	})
}

func requireIconWidget(data []byte) IconWidget {
	iconWidget, err := widget.NewIcon(data)
	if err != nil {
		panic(err)
	}
	return func(gtx layout.Context, foreground color.Color) layout.Dimensions {
		if nrgba, ok := foreground.(color.NRGBA); ok {
			return iconWidget.Layout(gtx, nrgba)
		}
		return iconWidget.Layout(gtx, ToNRGBA(foreground))
	}
}
