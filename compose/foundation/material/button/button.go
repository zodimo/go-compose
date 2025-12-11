package button

import (
	"fmt"
	"go-compose-dev/internal/layoutnode"

	"gioui.org/widget"
	"gioui.org/widget/material"
)

const MaterialButtonNodeID = "MaterialButton"

func Button(onClick func(), label string, options ...ButtonOption) Composable {
	return func(c Composer) Composer {
		opts := DefaultButtonOptions()
		for _, option := range options {
			option(&opts)
		}

		if opts.Clickable == nil {
			key := c.GenerateID()
			path := c.GetPath()

			clickablePath := fmt.Sprintf("%d/%s/clickable", key, path)
			clickableValue := c.State(clickablePath, func() any { return &GioClickable{} })
			opts.Clickable = clickableValue.Get().(*widget.Clickable)
		}

		constructorArgs := ButtonConstructorArgs{
			Clickable:    opts.Clickable,
			LabelContent: label,
			OnClick:      onClick,
		}

		c.StartBlock(MaterialButtonNodeID)
		c.Modifier(func(modifier Modifier) Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.SetWidgetConstructor(buttonWidgetConstructor(opts, constructorArgs))

		return c.EndBlock()
	}
}

func GetClickablePath(c Composer) string {
	return c.GetPath().String() + "/clickable"
}

type ButtonConstructorArgs struct {
	Clickable    *GioClickable
	OnClick      func()
	LabelContent string
}

func buttonWidgetConstructor(options ButtonOptions, constructorArgs ButtonConstructorArgs) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {

			clickable := constructorArgs.Clickable
			if clickable.Clicked(gtx) {
				constructorArgs.OnClick()
			}

			return material.Button(options.Theme, constructorArgs.Clickable, constructorArgs.LabelContent).Layout(gtx)

		}
	})
}
