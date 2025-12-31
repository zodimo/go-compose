package text

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/next/text"
	"github.com/zodimo/go-compose/compose/ui/next/text/style"
)

type TextOptions struct {
	Modifier Modifier

	TextStyle *text.TextStyle

	OnTextLayout func(text.TextLayoutResult)
	OverFlow     style.TextOverFlow
	SoftWrap     bool

	MaxLines int
	MinLines int

	InlineContent map[string]InlineTextContent

	Color    graphics.ColorProducer
	AutoSize TextAutoSize
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

func WithOnTextLayout(onTextLayout func(text.TextLayoutResult)) TextOption {
	return func(o *TextOptions) {
		o.OnTextLayout = onTextLayout
	}
}

func WithOverFlow(overFlow style.TextOverFlow) TextOption {
	return func(o *TextOptions) {
		o.OverFlow = overFlow
	}
}
func WithSoftWrap(softWrap bool) TextOption {
	return func(o *TextOptions) {
		o.SoftWrap = softWrap
	}
}

func WithMaxLines(maxLines int) TextOption {
	return func(o *TextOptions) {
		o.MaxLines = maxLines
	}
}

func WithMinLines(minLines int) TextOption {
	return func(o *TextOptions) {
		o.MinLines = minLines
	}
}

func WithColor(color graphics.ColorProducer) TextOption {
	return func(o *TextOptions) {
		o.Color = color
	}
}
func WithAutoSize(autoSize TextAutoSize) TextOption {
	return func(o *TextOptions) {
		o.AutoSize = autoSize
	}
}
