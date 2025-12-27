package text

import (
	"fmt"
	"image"

	"github.com/zodimo/go-compose/compose/foundation/text/selection"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/state"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/op"
	"gioui.org/op/clip"
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

		textSelectionColors := selection.LocalTextSelectionColors.Current(c)

		constructorArgs := BasicTextConstructorArgs{
			Value:      value,
			Options:    opts,
			selectable: selectable,
			textColor: theme.TakeOrElseColor(
				opts.TextStyleOptions.Color,
				theme.ColorHelper.SpecificColor(graphics.ColorBlack),
			),
			textSelectionColor: theme.TakeOrElseColor(
				opts.TextStyleOptions.SelectionColor,
				theme.ColorHelper.SpecificColor(textSelectionColors.BackgroundColor)),
		}

		c.StartBlock(BasicTextNodeID)
		c.Modifier(func(modifier Modifier) Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.SetWidgetConstructor(textWidgetConstructor(constructorArgs))
		return c.EndBlock()
	}
}

type BasicTextConstructorArgs struct {
	Value              string
	Options            TextOptions
	selectable         state.MutableValue
	textColor          theme.ColorDescriptor
	textSelectionColor theme.ColorDescriptor
}

func textWidgetConstructor(constructorArgs BasicTextConstructorArgs) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {

			// where to find the theme?
			materialTheme := GetThemeManager().MaterialTheme()
			tm := theme.GetThemeManager()

			text := constructorArgs.Value
			textOptions := constructorArgs.Options

			// Resolve ColorDescriptors to NRGBA
			resolvedTextColor := tm.ResolveColorDescriptor(constructorArgs.textColor).AsNRGBA()
			resolvedSelectColor := tm.ResolveColorDescriptor(constructorArgs.textSelectionColor).AsNRGBA()

			textColorMacro := op.Record(gtx.Ops)
			paint.ColorOp{Color: resolvedTextColor}.Add(gtx.Ops)
			textColor := textColorMacro.Stop()

			selectColorMacro := op.Record(gtx.Ops)
			paint.ColorOp{Color: resolvedSelectColor}.Add(gtx.Ops)
			selectColor := selectColorMacro.Stop()

			var dims layoutnode.LayoutDimensions

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
				dims = state.Layout(
					gtx,
					materialTheme.Shaper,
					constructorArgs.Options.TextStyleOptions.Font,
					constructorArgs.Options.TextStyleOptions.TextSize,
					textColor,
					selectColor,
				)

			} else {

				dims = widget.Label{
					Alignment:       textOptions.Alignment,
					MaxLines:        textOptions.MaxLines,
					Truncator:       textOptions.Truncator,
					WrapPolicy:      textOptions.WrapPolicy,
					LineHeight:      textOptions.LineHeight,
					LineHeightScale: textOptions.LineHeightScale,
				}.Layout(gtx, materialTheme.Shaper, textOptions.TextStyleOptions.Font, textOptions.TextStyleOptions.TextSize, text, textColor)
			}

			// Draw strikethrough if enabled
			if textOptions.TextStyleOptions.Strikethrough {
				// Draw a line through the middle of the text
				lineHeight := 1
				y := dims.Size.Y / 2
				rect := image.Rect(0, y, dims.Size.X, y+lineHeight)
				paint.FillShape(gtx.Ops, resolvedTextColor, clip.Rect(rect).Op())
			}

			return dims
		}
	})

}
