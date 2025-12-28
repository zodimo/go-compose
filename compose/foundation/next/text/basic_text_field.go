package text

import (
	"gioui.org/op"
	"gioui.org/op/paint"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/next/text/input"
	"github.com/zodimo/go-compose/compose/ui/platform"
	"github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/theme"
)

const BasicTextFieldNodeID = "BasicTextField"

// BasicTextField is an interactive text input composable that accepts text input
// through software or hardware keyboard, but provides no decorations like hint
// or placeholder.
//
// All editing state is hoisted through the TextFieldState parameter. Whenever
// the contents change via user input or semantics, the state is updated.
// Similarly, all programmatic updates to the state reflect in the composable.
//
// To add decorations (icons, labels, helper text), use the WithDecorator option.
//
// To filter or transform input, use WithInputTransformation.
//
// To control line limits and scrolling, use WithLineLimits.
//
// This is a port of androidx.compose.foundation.text.BasicTextField.
func BasicTextField(
	state *input.TextFieldState,
	options ...TextFieldOption,
) Composable {

	opts := DefaultTextFieldOptions()
	for _, option := range options {
		if option == nil {
			continue
		}
		option(&opts)
	}

	return func(c compose.Composer) compose.Composer {

		c.StartBlock(BasicTextFieldNodeID)

		familyResolver := platform.LocalFontFamilyResolver.Current(c)
		layoutDirection := platform.LocalLayoutDirection.Current(c)

		// Determine single-line mode from line limits
		singleLine := input.IsSingleLine(opts.LineLimits)

		// Calculate min/max lines from limits
		minLines := 1
		maxLines := 1
		if multiLine, ok := opts.LineLimits.(input.MultiLine); ok {
			minLines = multiLine.MinHeightInLines
			maxLines = multiLine.MaxHeightInLines
		}

		c.Modifier(func(m modifier.Modifier) modifier.Modifier {
			return m.Then(textFieldModifier().Then(opts.Modifier))
		})

		c.SetWidgetConstructor(textFieldWidgetConstructor(BasicTextFieldConstructorArgs{
			state:                state,
			textStyle:            opts.TextStyle,
			onTextLayout:         opts.OnTextLayout,
			singleLine:           singleLine,
			maxLines:             maxLines,
			minLines:             minLines,
			enabled:              opts.Enabled,
			readOnly:             opts.ReadOnly,
			inputTransformation:  opts.InputTransformation,
			outputTransformation: opts.OutputTransformation,
			fontFamilyResolver:   familyResolver,
			cursorColor:          opts.CursorColor,
			layoutDirection:      layoutDirection,
		}))

		return c.EndBlock()
	}
}

// BasicTextFieldConstructorArgs holds the arguments for the text field widget constructor.
type BasicTextFieldConstructorArgs struct {
	state                *input.TextFieldState
	textStyle            *text.TextStyle
	onTextLayout         func(text.TextLayoutResult)
	singleLine           bool
	maxLines             int
	minLines             int
	enabled              bool
	readOnly             bool
	inputTransformation  input.InputTransformation
	outputTransformation input.OutputTransformation
	fontFamilyResolver   interface{} // font.FontFamilyResolver
	cursorColor          interface{} // graphics.ColorProducer
	layoutDirection      unit.LayoutDirection
}

// textFieldWidgetConstructor creates the widget constructor for BasicTextField.
func textFieldWidgetConstructor(args BasicTextFieldConstructorArgs) layoutnode.LayoutNodeWidgetConstructor {
	// Create the bridge adapter from the text field state
	textString := args.state.Text()
	sourceAdapter := input.NewTextSourceAdapterFromString(textString)

	// Create the layout controller
	controller := input.NewTextLayoutController(sourceAdapter)

	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			// Get theme and shaper
			materialTheme := GetThemeManager().MaterialTheme()
			tm := theme.GetThemeManager()

			// Update the source adapter with current text from state
			currentText := args.state.Text()
			if sourceAdapter.Text() != currentText {
				sourceAdapter = input.NewTextSourceAdapterFromString(currentText)
				controller = input.NewTextLayoutController(sourceAdapter)
			}

			// Resolve text style with defaults
			textStyle := text.TextStyleResolveDefaults(args.textStyle, args.layoutDirection)

			// Configure the controller from the text style
			controller.SetTextStyle(textStyle)
			controller.ConfigureFromTextStyle(textStyle)
			controller.SetMaxLines(args.maxLines)
			controller.SetSingleLine(args.singleLine)
			controller.SetTruncator("")
			controller.SetLineHeightScale(1)

			// Resolve text color
			var textColorDescriptor theme.ColorDescriptor
			textColorDescriptor = theme.ColorHelper.SpecificColor(textStyle.Color())
			resolvedTextColor := tm.ResolveColorDescriptor(textColorDescriptor).AsNRGBA()

			// Create text color material
			textColorMacro := op.Record(gtx.Ops)
			paint.ColorOp{Color: resolvedTextColor}.Add(gtx.Ops)
			textColor := textColorMacro.Stop()

			// Use the controller to layout and paint the text
			dims := controller.LayoutAndPaint(gtx, materialTheme.Shaper, textColor)

			return dims
		}
	})
}

// textFieldModifier returns the base modifier for text fields.
func textFieldModifier() modifier.Modifier {
	return modifier.EmptyModifier
}
