package text

import (
	"fmt"
	"image"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/text/selection"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/platform"
	"github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/state"

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

		textShaper := compose.LocalTextShaper.Current(c)
		layoutDirection := platform.LocalLayoutDirection.Current(c)
		// Resolve text style with defaults
		opts.TextStyle = text.TextStyleResolveDefaults(opts.TextStyle, layoutDirection)

		path := c.GetPath()
		key := c.GenerateID()

		selectablePath := fmt.Sprintf("%d/%s/selectable", key, path)

		// @TODO selection container present

		var selectable state.MutableValue
		if opts.Selectable.IsSome() {
			if opts.Selectable.UnwrapUnsafe() {
				selectable = c.State(selectablePath, func() any { return &widget.Selectable{} })
			}
		}

		textSelectionColors := selection.LocalTextSelectionColors.Current(c)

		constructorArgs := BasicTextConstructorArgs{
			Value:              value,
			Options:            opts,
			selectable:         selectable,
			textSelectionColor: textSelectionColors.BackgroundColor,

			textShaper: textShaper,
			// layoutDirection:    layoutDirection,
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
	textSelectionColor graphics.Color

	textShaper *text.TextShaper
	// layoutDirection    unit.LayoutDirection
}

func textWidgetConstructor(constructorArgs BasicTextConstructorArgs) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {

			textValue := constructorArgs.Value
			textOptions := constructorArgs.Options

			// local alias
			textStyle := textOptions.TextStyle

			resolvedTextColor := graphics.ColorToNRGBA(textStyle.Color())
			resolvedSelectColor := graphics.ColorToNRGBA(constructorArgs.textSelectionColor)

			textColorMacro := op.Record(gtx.Ops)
			paint.ColorOp{Color: resolvedTextColor}.Add(gtx.Ops)
			textColor := textColorMacro.Stop()

			selectColorMacro := op.Record(gtx.Ops)
			paint.ColorOp{Color: resolvedSelectColor}.Add(gtx.Ops)
			_ = selectColorMacro.Stop()

			var dims layoutnode.LayoutDimensions

			// if textOptions.Selectable {
			// 	selectable := constructorArgs.selectable.Get().(*widget.Selectable)
			// 	state := selectable
			// 	if state.Text() != constructorArgs.Value {
			// 		state.SetText(constructorArgs.Value)
			// 	}
			// 	state.Alignment = constructorArgs.Options.Alignment
			// 	state.MaxLines = constructorArgs.Options.MaxLines
			// 	state.Truncator = constructorArgs.Options.Truncator
			// 	state.WrapPolicy = constructorArgs.Options.WrapPolicy
			// 	state.LineHeight = constructorArgs.Options.LineHeight
			// 	state.LineHeightScale = constructorArgs.Options.LineHeightScale
			// 	dims = state.Layout(
			// 		gtx,
			// 		materialTheme.Shaper,
			// 		constructorArgs.Options.TextStyleOptions.Font,
			// 		constructorArgs.Options.TextStyleOptions.TextSize,
			// 		textColor,
			// 		selectColor,
			// 	)

			// } else {

			lineHeight := textStyle.LineHeight()
			if lineHeight.IsSpecified() {
				if !lineHeight.IsSp() {
					panic("lineHeight, only sp is supported")
				}
			}

			// fmt.Printf("textStyle [%s]: %s\n", textValue, text.StringTextStyle(textStyle))
			dims = widget.Label{
				Alignment:       style.TextAlignToGioTextAlignment(textStyle.TextAlign()),
				MaxLines:        textOptions.MaxLines,
				Truncator:       textOptions.Truncator,
				WrapPolicy:      style.LineBreakToGioWrapPolicy(textStyle.LineBreak()),
				LineHeight:      textStyle.LineHeight().AsGioSp(),
				LineHeightScale: 0, // TODO how should this be handled?
			}.Layout(
				gtx,
				constructorArgs.textShaper.Shaper,
				font.ToGioFont(
					textStyle.FontFamily(),
					textStyle.FontWeight(),
					textStyle.FontStyle(),
				),
				textStyle.FontSize().AsGioSp(),
				textValue,
				textColor,
			)
			// }

			textDecoration := style.TakeOrElseTextDecoration(textOptions.TextStyle.TextDecoration(), style.TextDecorationNone)
			if style.IsSpecifiedTextDecoration(textDecoration) {
				// Draw strikethrough if enabled
				if textDecoration.Contains(style.TextDecorationLineThrough) {
					// Draw a line through the middle of the text
					lineHeight := 1
					y := dims.Size.Y / 2
					rect := image.Rect(0, y, dims.Size.X, y+lineHeight)
					paint.FillShape(gtx.Ops, resolvedTextColor, clip.Rect(rect).Op())
				}
			}

			return dims
		}
	})

}
