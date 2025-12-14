package text

import (
	"fmt"

	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/state"

	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

const BasicTextNodeID = "BasicText"

// Text renders a text leaf. The value is stored in a slot for tooling.
func Text(value string, options ...TextOption) Composable {

	opts := DefaultTextOptions()
	for _, option := range options {
		if option == nil {
			continue
		}
		option(&opts)
	}

	return func(c Composer) Composer {

		path := c.GetPath()
		key := c.GenerateID()

		selectablePath := fmt.Sprintf("%d/%s/selectable", key, path)

		var selectable state.MutableValue
		if opts.Selectable {
			selectable = c.State(selectablePath, func() any { return &widget.Selectable{} })
		}

		constructorArgs := BasicTextConstructorArgs{
			Value:      value,
			Options:    opts,
			selectable: selectable,
		}

		c.StartBlock(BasicTextNodeID)
		c.Modifier(func(modifier Modifier) Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.SetWidgetConstructor(textWidgetConstructor(opts, constructorArgs))
		return c.EndBlock()
	}
}

type BasicTextConstructorArgs struct {
	Value      string
	Options    TextOptions
	selectable state.MutableValue
}

func textWidgetConstructor(options TextOptions, constructorArgs BasicTextConstructorArgs) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {

			// where to find the theme?
			theme := GetThemeManager().MaterialTheme()

			text := constructorArgs.Value
			textOptions := constructorArgs.Options

			textColorMacro := op.Record(gtx.Ops)
			paint.ColorOp{Color: textOptions.TextStyleOptions.Color}.Add(gtx.Ops)
			textColor := textColorMacro.Stop()

			selectColorMacro := op.Record(gtx.Ops)
			paint.ColorOp{Color: textOptions.TextStyleOptions.SelectionColor}.Add(gtx.Ops)
			selectColor := selectColorMacro.Stop()

			if textOptions.Selectable {
				selectable := constructorArgs.selectable.Get().(*widget.Selectable)
				state := selectable
				if state.Text() != constructorArgs.Value {
					state.SetText(constructorArgs.Value)
				}
				state.Alignment = constructorArgs.Options.Alignment
				state.MaxLines = constructorArgs.Options.MaxLines
				state.Truncator = constructorArgs.Options.Truncator
				state.WrapPolicy = constructorArgs.Options.WrapPolicy
				state.LineHeight = constructorArgs.Options.LineHeight
				state.LineHeightScale = constructorArgs.Options.LineHeightScale
				return state.Layout(
					gtx,
					theme.Shaper,
					constructorArgs.Options.TextStyleOptions.Font,
					constructorArgs.Options.TextStyleOptions.TextSize,
					textColor,
					selectColor,
				)

			}

			return widget.Label{
				Alignment:       textOptions.Alignment,
				MaxLines:        textOptions.MaxLines,
				Truncator:       textOptions.Truncator,
				WrapPolicy:      textOptions.WrapPolicy,
				LineHeight:      textOptions.LineHeight,
				LineHeightScale: textOptions.LineHeightScale,
			}.Layout(gtx, theme.Shaper, textOptions.TextStyleOptions.Font, textOptions.TextStyleOptions.TextSize, text, textColor)
		}
	})

}
