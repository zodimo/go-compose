package iconbutton

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/internal/layoutnode"

	"gioui.org/layout"
	"gioui.org/widget"
	"git.sr.ht/~schnwalter/gio-mw/token"
	"git.sr.ht/~schnwalter/gio-mw/widget/button"
)

const Material3IconButtonNodeID = "Material3IconButton"

// Standard represents a standard text-based icon button (no background).
func Standard(onClick func(), icon []byte, description string, options ...IconButtonOption) Composable {

	return iconButtonWithContentColorComposable(button.Text(), onClick, icon, description, options...)
}

// Filled represents a filled icon button (high emphasis).
func Filled(onClick func(), icon []byte, description string, options ...IconButtonOption) Composable {
	return iconButtonComposable(button.Filled(), onClick, icon, description, options...)
}

// FilledTonal represents a filled tonal icon button (medium emphasis).
func FilledTonal(onClick func(), icon []byte, description string, options ...IconButtonOption) Composable {
	return iconButtonComposable(button.FilledTonal(), onClick, icon, description, options...)
}

// Outlined represents an outlined icon button.
func Outlined(onClick func(), icon []byte, description string, options ...IconButtonOption) Composable {
	return iconButtonComposable(button.Outlined(), onClick, icon, description, options...)
}

func iconButtonWithContentColorComposable(material3Button *button.Button, onClick func(), icon []byte, description string, options ...IconButtonOption) Composable {
	return func(c Composer) Composer {
		opts := DefaultIconButtonOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		contentColor := material3.LocalContentColor.Current(c)

		if opts.Button == nil {
			key := c.GenerateID()
			path := c.GetPath()

			buttonStatePath := fmt.Sprintf("%d/%s/iconbutton", key, path)
			buttonValue := c.State(buttonStatePath, func() any { return material3Button })
			opts.Button = buttonValue.Get().(*button.Button)
		}

		opts.Button.WithColorScheme(&button.AlternativeColorScheme{
			EnabledLabelColor: token.MatColor(graphics.ColorToNRGBA(contentColor)),
			EnabledIconColor:  token.MatColor(graphics.ColorToNRGBA(contentColor)),
		})

		constructorArgs := IconButtonConstructorArgs{
			Button:      opts.Button,
			OnClick:     onClick,
			Icon:        icon,
			Description: description,
		}

		c.StartBlock(Material3IconButtonNodeID)
		c.Modifier(func(modifier ui.Modifier) ui.Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.SetWidgetConstructor(iconButtonWidgetConstructor(opts, constructorArgs))

		return c.EndBlock()
	}
}

func iconButtonComposable(material3Button *button.Button, onClick func(), icon []byte, description string, options ...IconButtonOption) Composable {
	return func(c Composer) Composer {
		opts := DefaultIconButtonOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		if opts.Button == nil {
			key := c.GenerateID()
			path := c.GetPath()

			buttonStatePath := fmt.Sprintf("%d/%s/iconbutton", key, path)
			buttonValue := c.State(buttonStatePath, func() any { return material3Button })
			opts.Button = buttonValue.Get().(*button.Button)
		}

		constructorArgs := IconButtonConstructorArgs{
			Button:      opts.Button,
			OnClick:     onClick,
			Icon:        icon,
			Description: description,
		}

		c.StartBlock(Material3IconButtonNodeID)
		c.Modifier(func(modifier ui.Modifier) ui.Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.SetWidgetConstructor(iconButtonWidgetConstructor(opts, constructorArgs))

		return c.EndBlock()
	}
}

type IconButtonConstructorArgs struct {
	Button      *button.Button
	OnClick     func()
	Icon        []byte
	Description string
}

func iconButtonWidgetConstructor(_ IconButtonOptions, constructorArgs IconButtonConstructorArgs) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {

			button := constructorArgs.Button
			onClick := constructorArgs.OnClick
			if button.Clicked(*gtx.ToGio()) {
				onClick()
			}

			iconWidget, err := widget.NewIcon(constructorArgs.Icon)
			if err != nil {
				// In a real app we might want to handle this gracefully,
				// but failing fast for invalid icon data is acceptable.
				panic(fmt.Sprintf("failed to create icon: %v", err))
			}

			mwIconWidget := func(gtx layout.Context, c token.MatColor) layout.Dimensions {
				return iconWidget.Layout(gtx, c.AsNRGBA())
			}

			return layoutnode.FromGioDimensions(button.LayoutIconOnly(*gtx.ToGio(), constructorArgs.Description, mwIconWidget))
		}
	})
}
