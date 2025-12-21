package text

import "fmt"

type AnnotatedString interface {
	Options() AnnotatedStringOptions
	WithOptions(options AnnotatedStringOptions) AnnotatedString
	fmt.Stringer
}

type AnnotatedStringOptions struct {
	SpanStyles      []SpanStyle
	ParagraphStyles []ParagraphStyle
}

type NewAnnotatedStringOption func(options *AnnotatedStringOptions)

func defaultAnnotatedStringOptions() AnnotatedStringOptions {
	return AnnotatedStringOptions{
		SpanStyles:      []SpanStyle{},
		ParagraphStyles: []ParagraphStyle{},
	}
}

type AnnotatedStringImpl struct {
	text    string
	options AnnotatedStringOptions
}

func NewAnnotatedString(text string, options ...NewAnnotatedStringOption) AnnotatedString {
	opts := defaultAnnotatedStringOptions()
	for _, option := range options {
		if option == nil {
			continue
		}
		option(&opts)
	}
	return AnnotatedStringImpl{
		text:    text,
		options: opts,
	}
}

func (a AnnotatedStringImpl) String() string {
	return a.text
}

func (a AnnotatedStringImpl) Options() AnnotatedStringOptions {
	return a.options
}

func (a AnnotatedStringImpl) WithOptions(options AnnotatedStringOptions) AnnotatedString {
	return AnnotatedStringImpl{
		text:    a.text,
		options: options,
	}
}

var _ AnnotatedString = (*AnnotatedText)(nil)

type AnnotatedText string

func (a AnnotatedText) String() string {
	return string(a)
}

func (a AnnotatedText) Options() AnnotatedStringOptions {
	return defaultAnnotatedStringOptions()
}

func (a AnnotatedText) WithOptions(options AnnotatedStringOptions) AnnotatedString {
	return AnnotatedStringImpl{
		text:    string(a),
		options: options,
	}
}
