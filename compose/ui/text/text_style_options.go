package text

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
)

type TextStyleOption func(ts *TextStyle)

func WithColor(color graphics.Color) TextStyleOption {
	return func(ts *TextStyle) {
		ts = CoalesceTextStyle(ts, TextStyleUnspecified)
		ts.spanStyle = CopySpanStyle(ts.spanStyle, func(s *SpanStyle) {
			s.color = color
		})
	}
}
