package text

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

type TextStyleOption = func(ts *TextStyle)

// SPAN STYLE

func WithColor(color graphics.Color) TextStyleOption {
	return func(ts *TextStyle) {
		ensureMutableTextStyle(ts)
		ts.spanStyle = SpanStyleCompanion.Copy(ts.spanStyle, func(s *SpanStyle) {
			s.color = color
		})
	}
}

func WithFontSize(size unit.TextUnit) TextStyleOption {
	return func(ts *TextStyle) {
		ensureMutableTextStyle(ts)
		ts.spanStyle = SpanStyleCompanion.Copy(ts.spanStyle, func(s *SpanStyle) {
			s.fontSize = size
		})
	}
}

func WithFontWeight(weight font.FontWeight) TextStyleOption {
	return func(ts *TextStyle) {
		ensureMutableTextStyle(ts)
		ts.spanStyle = SpanStyleCompanion.Copy(ts.spanStyle, func(s *SpanStyle) {
			s.fontWeight = weight
		})
	}
}

func WithFontStyle(style font.FontStyle) TextStyleOption {
	return func(ts *TextStyle) {
		ensureMutableTextStyle(ts)
		ts.spanStyle = SpanStyleCompanion.Copy(ts.spanStyle, func(s *SpanStyle) {
			s.fontStyle = style
		})
	}
}

func WithFontFamily(family font.FontFamily) TextStyleOption {
	return func(ts *TextStyle) {
		ensureMutableTextStyle(ts)
		ts.spanStyle = SpanStyleCompanion.Copy(ts.spanStyle, func(s *SpanStyle) {
			s.fontFamily = family
		})
	}
}

func WithLetterSpacing(spacing unit.TextUnit) TextStyleOption {
	return func(ts *TextStyle) {
		ensureMutableTextStyle(ts)
		ts.spanStyle = SpanStyleCompanion.Copy(ts.spanStyle, func(s *SpanStyle) {
			s.letterSpacing = spacing
		})
	}
}

func WithBackground(color graphics.Color) TextStyleOption {
	return func(ts *TextStyle) {
		ensureMutableTextStyle(ts)
		ts.spanStyle = SpanStyleCompanion.Copy(ts.spanStyle, func(s *SpanStyle) {
			s.background = color
		})
	}
}

func WithTextDecoration(decoration *style.TextDecoration) TextStyleOption {
	return func(ts *TextStyle) {
		ensureMutableTextStyle(ts)
		ts.spanStyle = SpanStyleCompanion.Copy(ts.spanStyle, func(s *SpanStyle) {
			s.textDecoration = decoration
		})
	}
}

func WithShadow(shadow *graphics.Shadow) TextStyleOption {
	return func(ts *TextStyle) {
		ensureMutableTextStyle(ts)
		ts.spanStyle = SpanStyleCompanion.Copy(ts.spanStyle, func(s *SpanStyle) {
			s.shadow = shadow
		})
	}
}

// PARAGRAPH STYLE

func WithTextAlign(textAlign style.TextAlign) TextStyleOption {
	return func(ts *TextStyle) {
		ensureMutableTextStyle(ts)
		ts.paragraphStyle = ParagraphStyleCompanion.Copy(ts.paragraphStyle, func(s *ParagraphStyle) {
			s.textAlign = textAlign
		})
	}
}

func WithTextDirection(textDirection style.TextDirection) TextStyleOption {
	return func(ts *TextStyle) {
		ensureMutableTextStyle(ts)
		ts.paragraphStyle = ParagraphStyleCompanion.Copy(ts.paragraphStyle, func(s *ParagraphStyle) {
			s.textDirection = textDirection
		})
	}
}

func WithLineHeight(lineHeight unit.TextUnit) TextStyleOption {
	return func(ts *TextStyle) {
		ensureMutableTextStyle(ts)
		ts.paragraphStyle = ParagraphStyleCompanion.Copy(ts.paragraphStyle, func(s *ParagraphStyle) {
			s.lineHeight = lineHeight
		})
	}
}

func WithLineBreak(lineBreak style.LineBreak) TextStyleOption {
	return func(ts *TextStyle) {
		ensureMutableTextStyle(ts)
		ts.paragraphStyle = ParagraphStyleCompanion.Copy(ts.paragraphStyle, func(s *ParagraphStyle) {
			s.lineBreak = lineBreak
		})
	}
}
