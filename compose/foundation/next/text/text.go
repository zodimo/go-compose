package text

import (
	"fmt"

	"gioui.org/op"
	"gioui.org/op/paint"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/next/text/input"
	"github.com/zodimo/go-compose/compose/foundation/next/text/modifiers"
	"github.com/zodimo/go-compose/compose/foundation/next/text/selection"
	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/platform"
	"github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/theme"
	"github.com/zodimo/go-ternary"

	gioText "gioui.org/text"
)

const BasicTextNodeID = "BasicText"

func BasicText(
	annotatedString text.AnnotatedString,
	options ...TextOption,
) Composable {

	opts := DefaultTextOptions()
	for _, option := range options {
		if option == nil {
			continue
		}
		option(&opts)
	}

	validateMinMaxLines(opts.MinLines, opts.MaxLines)
	hasInlineContent := annotatedString.HasInlineContent()
	hasLinks := annotatedString.HasLinks()

	return func(c compose.Composer) compose.Composer {

		c.StartBlock(BasicTextNodeID)

		selectionRegistrar := selection.LocalSelectionRegistrar.Current(c)
		var selectionController *modifiers.SelectionController = nil
		if selectionRegistrar != nil {
			backgroundColor := selection.LocalTextSelectionColors.Current(c).BackgroundColor
			selectionController = modifiers.NewSelectionController(
				selectionRegistrar.NextSelectableId(),
				selectionRegistrar,
				backgroundColor,
			)
		}
		textShaper := compose.LocalTextShaper.Current(c)
		familyResolver := platform.LocalFontFamilyResolver.Current(c)
		layoutDirection := platform.LocalLayoutDirection.Current(c)

		//content color....

		if !hasInlineContent && !hasLinks {

			c.Modifier(func(modifier modifier.Modifier) modifier.Modifier {
				return modifier.
					Then(
						textModifier().
							Then(opts.Modifier),
					)
			})

			c.SetWidgetConstructor(textWidgetConstructor(BasicTextConstructorArgs{
				annotatedString:     annotatedString,
				style:               opts.TextStyle,
				onTextLayout:        opts.OnTextLayout,
				overflow:            opts.OverFlow,
				softWrap:            opts.SoftWrap,
				maxLines:            opts.MaxLines,
				minLines:            opts.MinLines,
				fontFamilyResolver:  familyResolver,
				placeholders:        nil,
				onPlaceholderLayout: nil,
				selectionController: selectionController,
				color:               opts.Color,
				onShowTranslation:   nil,
				autoSize:            opts.AutoSize,
				layoutDirection:     layoutDirection,

				textShaper: textShaper,
			}))

		} else {
			c.Modifier(func(modifier modifier.Modifier) modifier.Modifier {
				return modifier.
					Then(
						textModifier().
							Then(opts.Modifier),
					)
			})

			key := c.GenerateID()
			path := c.GetPath()

			displayTextPath := fmt.Sprintf("%d/%s/displayText", key, path)

			displayText := c.State(displayTextPath, func() any {
				return annotatedString
			})

			onShowTranslation := func(value modifiers.TextSubstitutionValue) {
				displayText.Set(ternary.Ternary(value.IsShowingSubstitution, value.Substitution, value.Original))
			}

			c.SetWidgetConstructor(textWithLinksAndInlineContentConstructor(BasicTextConstructorArgs{
				annotatedString:     displayText.Get().(text.AnnotatedString),
				onTextLayout:        opts.OnTextLayout,
				hasInlineContent:    true,
				inlineContent:       opts.InlineContent,
				style:               opts.TextStyle,
				overflow:            opts.OverFlow,
				softWrap:            opts.SoftWrap,
				maxLines:            opts.MaxLines,
				minLines:            opts.MinLines,
				fontFamilyResolver:  familyResolver,
				placeholders:        nil,
				onPlaceholderLayout: nil,
				selectionController: selectionController,
				color:               opts.Color,
				onShowTranslation:   onShowTranslation,
				autoSize:            opts.AutoSize,
				layoutDirection:     layoutDirection,

				textShaper: textShaper,
			}))
		}

		return c.EndBlock()
	}

}

// Text renders a text leaf. The value is stored in a slot for tooling.
func Text(value string, options ...TextOption) Composable {
	// return BasicText(stringAnnotation(value), options...)
	return compose.Id()
}

type BasicTextConstructorArgs struct {
	annotatedString     text.AnnotatedString
	style               *text.TextStyle
	onTextLayout        func(text.TextLayoutResult)
	hasInlineContent    bool
	inlineContent       map[string]InlineTextContent
	overflow            style.TextOverFlow
	softWrap            bool
	maxLines            int
	minLines            int
	fontFamilyResolver  font.FontFamilyResolver
	placeholders        []text.Range[text.Placeholder]
	onPlaceholderLayout func([]geometry.Rect)
	selectionController *modifiers.SelectionController
	color               graphics.ColorProducer
	onShowTranslation   func(modifiers.TextSubstitutionValue)
	autoSize            TextAutoSize
	layoutDirection     unit.LayoutDirection

	textShaper *text.TextShaper
}

func textWithLinksAndInlineContentConstructor(constructorArgs BasicTextConstructorArgs) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			return layoutnode.LayoutDimensions{}
		}
	})
}

func textWidgetConstructor(constructorArgs BasicTextConstructorArgs) layoutnode.LayoutNodeWidgetConstructor {
	// Create the bridge adapter from the annotated string
	textString := constructorArgs.annotatedString.Text()
	sourceAdapter := input.NewTextSourceAdapterFromString(textString)

	// Create the layout controller
	controller := input.NewTextLayoutController(sourceAdapter)

	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			// Get theme and shaper

			// Resolve text style with defaults
			textStyle := text.TextStyleResolveDefaults(constructorArgs.style, constructorArgs.layoutDirection)

			// Configure the controller from the text style
			controller.SetTextStyle(textStyle)
			controller.ConfigureFromTextStyle(textStyle)
			controller.SetMaxLines(constructorArgs.maxLines)
			controller.SetTruncator("...")
			controller.SetLineHeightScale(1)

			// Resolve text color
			var textColorDescriptor theme.ColorDescriptor
			if constructorArgs.color != nil {
				textColorDescriptor = theme.ColorHelper.SpecificColor(constructorArgs.color())
			} else {
				textColorDescriptor = theme.ColorHelper.SpecificColor(textStyle.Color())
			}
			resolvedTextColor := theme.GetThemeManager().ResolveColorDescriptor(textColorDescriptor).AsNRGBA()

			// Create text color material
			textColorMacro := op.Record(gtx.Ops)
			paint.ColorOp{Color: resolvedTextColor}.Add(gtx.Ops)
			textColor := textColorMacro.Stop()

			// Use the controller to layout and paint the text
			dims := controller.LayoutAndPaint(gtx, constructorArgs.textShaper.Shaper, textColor)

			return dims
		}
	})
}

func validateMinMaxLines(minLines int, maxLines int) {
	if minLines <= 0 || maxLines <= 0 {
		panic("both minLines and maxLines must be greater than zero")
	}
	if minLines > maxLines {
		panic("minLines must be less than or equal to maxLines")
	}
}

func textModifier() modifier.Modifier {
	return modifier.EmptyModifier
}

func textAlignToGioAlign(textAlign style.TextAlign) gioText.Alignment {
	switch textAlign {
	case style.TextAlignLeft:
		return gioText.Start
	case style.TextAlignCenter:
		return gioText.Middle
	case style.TextAlignRight:
		return gioText.End
	case style.TextAlignJustify:
		panic("TextAlignJustify not supported")
	default:
		panic(fmt.Sprintf("unknown TextAlign %s", textAlign))
	}
}

func lineBreakToGioWrapPolicy(lineBreak style.LineBreak) gioText.WrapPolicy {
	switch lineBreak {
	case style.LineBreakSimple:
		return gioText.WrapGraphemes
	case style.LineBreakHeading:
		return gioText.WrapWords
	case style.LineBreakParagraph:
		return gioText.WrapHeuristically
	default:
		panic(fmt.Sprintf("unknown TextLineBreak %s", lineBreak))
	}
}
