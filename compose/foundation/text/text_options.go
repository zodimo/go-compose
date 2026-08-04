package text

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/textconvert"
	"github.com/zodimo/go-maybe"

	gioFont "gioui.org/font"
)

type TextOptions struct {
	Modifier ui.Modifier

	TextStyle *text.TextStyle

	MaxLines int
	// Truncator is the text that will be shown at the end of the final
	// line if MaxLines is exceeded. Defaults to "…" if empty.
	Truncator string

	//BACKWARDS COMPATIBILITY

	// LineHeightScale applies a scaling factor to the LineHeight. If zero, a
	// sensible default will be used.
	LineHeightScale maybe.Maybe[float32]

	// Selectable provides text selection state for the label. If not set, the label cannot
	// be selected or copied interactively.
	Selectable maybe.Maybe[bool]

	// SelectionColor is the color of the background for selected text.
	SelectionColor graphics.Color
}

type TextOption func(*TextOptions)

func WithModifier(m ui.Modifier) TextOption {
	return func(o *TextOptions) {
		o.Modifier = m
	}
}

// replace TextStyle
func WithTextStyle(ts *text.TextStyle) TextOption {
	return func(o *TextOptions) {
		o.TextStyle = ts
	}
}

// merge TextStyle
func WithAdditionalTextStyle(ts *text.TextStyle) TextOption {
	return func(o *TextOptions) {
		textStyle := text.MergeTextStyle(
			o.TextStyle,
			ts,
		)
		o.TextStyle = textStyle
	}
}

func WithTextStyleOptions(options ...text.TextStyleOption) TextOption {
	return func(o *TextOptions) {
		// Use CopyTextStyle to avoid mutating the shared TextStyleUnspecified singleton
		textStyle := text.CopyTextStyle(
			text.CoalesceTextStyle(o.TextStyle, text.TextStyleUnspecified),
			options...,
		)
		o.TextStyle = textStyle
	}
}

func WithGioAlignment(alignment Alignment) TextOption {
	return WithAlignment(textconvert.FromGioTextAlign(alignment))
}

func WithAlignment(alignment style.TextAlign) TextOption {
	return WithTextStyleOptions(text.WithTextAlign(alignment))
}

func WithMaxLines(maxLines int) TextOption {
	return func(o *TextOptions) {
		o.MaxLines = maxLines
	}
}

func WithTruncator(truncator string) TextOption {
	return func(o *TextOptions) {
		o.Truncator = truncator
	}
}

// Deprecated: use WithTextStyle, WithAdditionalTextStyle or WithTextStyleOptions
func WithWrapPolicy(wrapPolicy WrapPolicy) TextOption {
	return WithTextStyleOptions(text.WithLineBreak(textconvert.GioWrapPolicyToLineBreak(wrapPolicy)))
}

// Deprecated: use WithTextStyle, WithAdditionalTextStyle or WithTextStyleOptions
func WithLineHeight(lineHeightInSP float32) TextOption {
	return WithTextStyleOptions(text.WithLineHeight(unit.Sp(lineHeightInSP)))

}

// Deprecated
func WithLineHeightScale(lineHeightScale float32) TextOption {
	return func(o *TextOptions) {
		o.LineHeightScale = maybe.Some(lineHeightScale)
	}
}

func WithColor(color graphics.Color) TextOption {
	return WithTextStyleOptions(text.WithColor(color))
}

// TextSelectable sets the text to be selectable.
// stores *widget.Selectable in runtime memoization
func Selectable() TextOption {
	return func(o *TextOptions) {
		o.Selectable = maybe.Some(true)
	}
}

// Deprecated: use WithTextStyle, WithAdditionalTextStyle or WithTextStyleOptions
func StyleWithFont(font gioFont.Font) TextOption {
	return WithAdditionalTextStyle(text.TextStyleFromGioFont(font))
}

// Deprecated: use WithTextStyle, WithAdditionalTextStyle or WithTextStyleOptions
func StyleWithColor(color graphics.Color) TextOption {
	return WithTextStyleOptions(text.WithColor(color))
}

// StyleWithSelectionColor sets the selection highlight color.
// Deprecated: use WithTextStyle, WithAdditionalTextStyle or WithTextStyleOptions
func StyleWithSelectionColor(color graphics.Color) TextOption {
	return func(o *TextOptions) {
		o.SelectionColor = color
	}
}

// Deprecated: use WithTextStyle, WithAdditionalTextStyle or WithTextStyleOptions
func StyleWithTextSize(sizeInSP float32) TextOption {
	return WithTextStyleOptions(text.WithFontSize(unit.Sp(sizeInSP)))
}

// StyleWithStrikethrough enables strikethrough text decoration.
// Deprecated: use WithTextStyle, WithAdditionalTextStyle or WithTextStyleOptions
func StyleWithStrikethrough() TextOption {
	return WithTextStyleOptions(text.WithTextDecoration(style.TextDecorationLineThrough))
}
