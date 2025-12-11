package button

import (
	"fmt"
	"go-compose-dev/internal/layoutnode"

	"git.sr.ht/~schnwalter/gio-mw/widget/button"
)

const Material3ButtonNodeID = "Material3Button"

func Text(onClick func(), label string, options ...ButtonOption) Composable {
	return buttonComposable(button.Text(), onClick, label, options...)
}

func Outlined(onClick func(), label string, options ...ButtonOption) Composable {
	return buttonComposable(button.Outlined(), onClick, label, options...)
}

func Filled(onClick func(), label string, options ...ButtonOption) Composable {
	return buttonComposable(button.Filled(), onClick, label, options...)
}

func FilledTonal(onClick func(), label string, options ...ButtonOption) Composable {
	return buttonComposable(button.FilledTonal(), onClick, label, options...)
}

func Elevated(onClick func(), label string, options ...ButtonOption) Composable {
	return buttonComposable(button.Elevated(), onClick, label, options...)
}

func buttonComposable(material3Button *button.Button, onClick func(), label string, options ...ButtonOption) Composable {
	return func(c Composer) Composer {
		opts := DefaultButtonOptions()
		for _, option := range options {
			option(&opts)
		}

		if opts.Button == nil {
			key := c.GenerateID()
			path := c.GetPath()

			buttonStatePath := fmt.Sprintf("%d/%s/button", key, path)
			buttonValue := c.State(buttonStatePath, func() any { return material3Button })
			opts.Button = buttonValue.Get().(*button.Button)
		}

		constructorArgs := ButtonConstructorArgs{
			Button:       opts.Button,
			OnClick:      onClick,
			LabelContent: label,
		}

		c.StartBlock(Material3ButtonNodeID)
		c.Modifier(func(modifier Modifier) Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.SetWidgetConstructor(buttonWidgetConstructor(opts, constructorArgs))

		return c.EndBlock()
	}
}

type ButtonConstructorArgs struct {
	Button       *button.Button
	OnClick      func()
	LabelContent string
}

func buttonWidgetConstructor(_ ButtonOptions, constructorArgs ButtonConstructorArgs) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {

			button := constructorArgs.Button
			onClick := constructorArgs.OnClick
			if button.Clicked(gtx) {
				onClick()
			}

			return button.Layout(gtx, constructorArgs.LabelContent)

		}
	})
}
