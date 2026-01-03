package text

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-maybe"

	gioFont "gioui.org/font"
	gioUnit "gioui.org/unit"
)

type TextOptions struct {
	Modifier Modifier

	TextStyle *text.TextStyle

	MaxLines int
	// Truncator is the text that will be shown at the end of the final
	// line if MaxLines is exceeded. Defaults to "…" if empty.
	Truncator string

	//BACKWARDS COMPATIBILITY

	// Alignment specifies the text alignment.
	Alignment maybe.Maybe[Alignment]

	// WrapPolicy configures how displayed text will be broken into lines.
	WrapPolicy maybe.Maybe[WrapPolicy]
	// LineHeight controls the distance between the baselines of lines of text.
	// If zero, a sensible default will be used.
	LineHeight maybe.Maybe[gioUnit.Sp]
	// LineHeightScale applies a scaling factor to the LineHeight. If zero, a
	// sensible default will be used.
	LineHeightScale maybe.Maybe[float32]

	// Selectable provides text selection state for the label. If not set, the label cannot
	// be selected or copied interactively.
	Selectable maybe.Maybe[bool]

	// Face defines the text style.
	Font maybe.Maybe[gioFont.Font]
	// SelectionColor is the color of the background for selected text.
	SelectionColor graphics.Color
	// TextSize determines the size of the text glyphs.
	TextSize maybe.Maybe[gioUnit.Sp]
}

type TextOption func(*TextOptions)

func WithModifier(m Modifier) TextOption {
	return func(o *TextOptions) {
		o.Modifier = m
	}
}

func WithTextStyle(ts *text.TextStyle) TextOption {
	return func(o *TextOptions) {
		o.TextStyle = ts
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
	return WithAlignment(style.FromGioTextAlign(alignment))
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

// Deprecated use TextStyle
func WithWrapPolicy(wrapPolicy WrapPolicy) TextOption {
	return func(o *TextOptions) {
		o.WrapPolicy = maybe.Some(wrapPolicy)
	}
}

// Deprecated use TextStyle
func WithLineHeight(lineHeightInSP float32) TextOption {
	return func(o *TextOptions) {
		o.LineHeight = maybe.Some(gioUnit.Sp(lineHeightInSP))
	}
}

func WithLineHeightScale(lineHeightScale float32) TextOption {

	return func(o *TextOptions) {
		o.LineHeightScale = maybe.Some(lineHeightScale)
	}
}

func WithColor(color graphics.Color) TextOption {
	return func(o *TextOptions) {

	}
}

// TextSelectable sets the text to be selectable.
// stores *widget.Selectable in runtime memoization
func Selectable() TextOption {
	return func(o *TextOptions) {
		o.Selectable = maybe.Some(true)
	}
}

func StyleWithFont(font gioFont.Font) TextOption {
	return func(o *TextOptions) {
		o.Font = maybe.Some(font)
	}
}

// Deprecated use TextStyle
func StyleWithColor(color graphics.Color) TextOption {
	return WithTextStyleOptions(text.WithColor(color))
}

// StyleWithSelectionColor sets the selection highlight color.
// Deprecated use TextStyle
func StyleWithSelectionColor(color graphics.Color) TextOption {
	return func(o *TextOptions) {
		o.SelectionColor = color
	}
}

// Deprecated use TextStyle
func StyleWithTextSize(sizeInSP float32) TextOption {
	return func(o *TextOptions) {
		o.TextSize = maybe.Some(gioUnit.Sp(sizeInSP))
	}
}

// StyleWithStrikethrough enables strikethrough text decoration.
// Deprecated use TextStyle
func StyleWithStrikethrough() TextOption {
	return WithTextStyleOptions(text.WithTextDecoration(style.TextDecorationLineThrough))
}
