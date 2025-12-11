package button

import (
	"fmt"
	"go-compose-dev/internal/layoutnode"
	"go-compose-dev/internal/state"

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

		key := c.GenerateID()
		path := c.GetPath()

		clickablePath := fmt.Sprintf("%d/%s/clickable", key, path)
		clickable := c.State(clickablePath, func() any { return &widget.Clickable{} })

		constructorArgs := ButtonConstructorArgs{
			Clickable:    clickable,
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
	Clickable    state.MutableValue
	OnClick      func()
	LabelContent string
}

func buttonWidgetConstructor(options ButtonOptions, constructorArgs ButtonConstructorArgs) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {

			clickable := constructorArgs.Clickable.Get().(*widget.Clickable)
			if clickable.Clicked(gtx) {
				constructorArgs.OnClick()
			}

			return material.Button(options.Theme, constructorArgs.Clickable.Get().(*widget.Clickable), constructorArgs.LabelContent).Layout(gtx)

		}
	})
}
