package text

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// textStyleUnspecified is the sentinel for TextStyle. It is deliberately
// package-private: the only public way to obtain it is TextStyleCompanion.Unspecified().
var textStyleUnspecified *TextStyle = &TextStyle{
	spanStyle:      spanStyleUnspecified,
	paragraphStyle: paragraphStyleUnspecified,
}

var _ TextStyleInterface = (*TextStyle)(nil)

type TextStyle struct {
	spanStyle      *SpanStyle
	paragraphStyle *ParagraphStyle
}

// TextStyleOptions mirrors TextStyle's span and paragraph properties, plus the
// convenience With* constructors below. A field left at its sentinel means
// "leave the source value alone".
type TextStyleOptions struct {
	Span      SpanStyleOptions
	Paragraph ParagraphStyleOptions
}

func newTextStyleOptions() TextStyleOptions {
	return TextStyleOptions{
		Span:      newSpanStyleOptions(),
		Paragraph: newParagraphStyleOptions(),
	}
}

type TextStyleOption = func(*TextStyleOptions)

// SPAN STYLE OPTIONS

func WithColor(color graphics.Color) TextStyleOption {
	return func(opts *TextStyleOptions) {
		opts.Span.Color = color
	}
}

func WithFontSize(size unit.TextUnit) TextStyleOption {
	return func(opts *TextStyleOptions) {
		opts.Span.FontSize = size
	}
}

func WithFontWeight(weight font.FontWeight) TextStyleOption {
	return func(opts *TextStyleOptions) {
		opts.Span.FontWeight = weight
	}
}

func WithFontStyle(fontStyle font.FontStyle) TextStyleOption {
	return func(opts *TextStyleOptions) {
		opts.Span.FontStyle = fontStyle
	}
}

func WithFontFamily(family font.FontFamily) TextStyleOption {
	return func(opts *TextStyleOptions) {
		opts.Span.FontFamily = family
	}
}

func WithLetterSpacing(spacing unit.TextUnit) TextStyleOption {
	return func(opts *TextStyleOptions) {
		opts.Span.LetterSpacing = spacing
	}
}

func WithBackground(color graphics.Color) TextStyleOption {
	return func(opts *TextStyleOptions) {
		opts.Span.Background = color
	}
}

func WithTextDecoration(decoration *style.TextDecoration) TextStyleOption {
	return func(opts *TextStyleOptions) {
		opts.Span.TextDecoration = decoration
	}
}

func WithShadow(shadow *graphics.Shadow) TextStyleOption {
	return func(opts *TextStyleOptions) {
		opts.Span.Shadow = shadow
	}
}

// PARAGRAPH STYLE OPTIONS

func WithTextAlign(textAlign style.TextAlign) TextStyleOption {
	return func(opts *TextStyleOptions) {
		opts.Paragraph.TextAlign = textAlign
	}
}

func WithTextDirection(textDirection style.TextDirection) TextStyleOption {
	return func(opts *TextStyleOptions) {
		opts.Paragraph.TextDirection = textDirection
	}
}

func WithLineHeight(lineHeight unit.TextUnit) TextStyleOption {
	return func(opts *TextStyleOptions) {
		opts.Paragraph.LineHeight = lineHeight
	}
}

func WithLineBreak(lineBreak style.LineBreak) TextStyleOption {
	return func(opts *TextStyleOptions) {
		opts.Paragraph.LineBreak = lineBreak
	}
}

func (ts TextStyle) Alpha() float32 {
	return ts.spanStyle.Color().Alpha()
}
func (ts TextStyle) Background() graphics.Color {
	return ts.spanStyle.Background()
}
func (ts TextStyle) Color() graphics.Color {
	return ts.spanStyle.Color()
}
func (ts TextStyle) FontFamily() font.FontFamily {
	return ts.spanStyle.FontFamily()
}
func (ts TextStyle) FontSize() unit.TextUnit {
	return ts.spanStyle.FontSize()
}
func (ts TextStyle) FontStyle() font.FontStyle {
	return ts.spanStyle.FontStyle()
}
func (ts TextStyle) FontWeight() font.FontWeight {
	return ts.spanStyle.FontWeight()
}
func (ts TextStyle) LetterSpacing() unit.TextUnit {
	return ts.spanStyle.LetterSpacing()
}

func (ts TextStyle) LineBreak() style.LineBreak {
	return ts.paragraphStyle.LineBreak()
}
func (ts TextStyle) LineHeight() unit.TextUnit {
	return ts.paragraphStyle.LineHeight()
}
func (ts TextStyle) TextAlign() style.TextAlign {
	return ts.paragraphStyle.TextAlign()
}
func (ts TextStyle) TextDecoration() *style.TextDecoration {
	return ts.spanStyle.TextDecoration()
}
func (ts TextStyle) TextDirection() style.TextDirection {
	return ts.paragraphStyle.TextDirection()
}

// asOption lets an already-populated TextStyleOptions seed a TextStyleOptions
// target, so the whole struct can be passed as a single TextStyleOption.
func (o TextStyleOptions) asOption(opts *TextStyleOptions) {
	*opts = o
}

func stringTextStyle(ts *TextStyle) string {
	ts = coalesceTextStyle(ts, textStyleUnspecified)
	if !isSpecifiedTextStyle(ts) {
		return "TextStyle{Unspecified}"
	}
	return fmt.Sprintf("TextStyle{spanStyle: %s, paragraphStyle: %s}",
		stringSpanStyle(ts.spanStyle),
		stringParagraphStyle(ts.paragraphStyle),
	)
}

// ensureMutableTextStyle panics if ts is a sentinel value.
// Call this at the start of any function that mutates a TextStyle to fail-fast on misuse.
func ensureMutableTextStyle(ts *TextStyle) {
	if !isSpecifiedTextStyle(ts) {
		panic("attempt to mutate sentinel TextStyleUnspecified; use CopyTextStyle first")
	}
}

func textStyleResolveDefaults(ts *TextStyle, direction unit.LayoutDirection) *TextStyle {
	ts = coalesceTextStyle(ts, textStyleUnspecified)
	return &TextStyle{
		spanStyle:      spanStyleResolveDefaults(ts.spanStyle),
		paragraphStyle: paragraphStyleResolveDefaults(ts.paragraphStyle, direction),
	}
}

func isSpecifiedTextStyle(s *TextStyle) bool {
	return s != nil && s != textStyleUnspecified
}

func isUnspecifiedTextStyle(s *TextStyle) bool {
	return !isSpecifiedTextStyle(s)
}

func takeOrElseTextStyle(s, def *TextStyle) *TextStyle {
	if !isSpecifiedTextStyle(s) {
		return def
	}
	return s
}

// mergeTextStyle layers b onto a: a value from b wins only where b actually set it.
func mergeTextStyle(a, b *TextStyle) *TextStyle {
	a = coalesceTextStyle(a, textStyleUnspecified)
	b = coalesceTextStyle(b, textStyleUnspecified)

	if a == textStyleUnspecified {
		return b
	}
	if b == textStyleUnspecified {
		return a
	}

	// Both are custom: allocate new merged style
	return &TextStyle{
		spanStyle:      mergeSpanStyle(a.spanStyle, b.spanStyle),
		paragraphStyle: mergeParagraphStyle(a.paragraphStyle, b.paragraphStyle),
	}
}

func coalesceTextStyle(ptr, def *TextStyle) *TextStyle {
	if ptr == nil {
		return def
	}
	return ptr
}

// copyTextStyle returns a copy of ts with the given bare options applied on top.
// Options override; every field they did not touch keeps ts's value. The source
// is never mutated, so a sentinel source is safe.
func copyTextStyle(ts *TextStyle, options ...func(*TextStyle)) *TextStyle {
	ts = coalesceTextStyle(ts, textStyleUnspecified)
	work := *ts
	for _, option := range options {
		option(&work)
	}
	return &work
}

// copyTextStyleWithOptions is the option-constructor form used by callers.
// Span and paragraph options are assembled separately, then overlaid onto ts so
// that an option setting a single field leaves every other field untouched.
func copyTextStyleWithOptions(ts *TextStyle, options ...TextStyleOption) *TextStyle {
	ts = coalesceTextStyle(ts, textStyleUnspecified)
	opts := newTextStyleOptions()
	for _, option := range options {
		option(&opts)
	}
	return &TextStyle{
		spanStyle:      copySpanStyleWithOptions(ts.spanStyle, opts.Span.asOption),
		paragraphStyle: copyParagraphStyleWithOptions(ts.paragraphStyle, opts.Paragraph.asOption),
	}
}

// textStyleFromOptions builds a TextStyle from options applied over the sentinel.
func textStyleFromOptions(options ...TextStyleOption) *TextStyle {
	return copyTextStyleWithOptions(textStyleUnspecified, options...)
}

// The exported wrappers below are the deliberately-documented public surface of this
// sentinel family. Everything else (the sentinel itself and the ten canonical ops)
// stays package-private and is reachable through the Companions.

// CopyTextStyleWithOptions returns a copy of ts with options applied on top.
// Options override; fields they do not touch keep ts's value. Use
// TextStyleCompanion.Unspecified() to start from an unspecified style.
func CopyTextStyleWithOptions(ts *TextStyle, options ...TextStyleOption) *TextStyle {
	return copyTextStyleWithOptions(ts, options...)
}

// TextStyleFromOptions builds a TextStyle from options applied over the unspecified style.
func TextStyleFromOptions(options ...TextStyleOption) *TextStyle {
	return textStyleFromOptions(options...)
}

// ResolveTextStyleDefaults resolves an unspecified TextStyle against the framework defaults.
func ResolveTextStyleDefaults(ts *TextStyle, direction unit.LayoutDirection) *TextStyle {
	return textStyleResolveDefaults(ts, direction)
}

// SpanStyleFromOptions builds a SpanStyle from options applied over the unspecified style.
func SpanStyleFromOptions(options ...SpanStyleOption) *SpanStyle {
	return copySpanStyleWithOptions(spanStyleUnspecified, options...)
}

// CopySpanStyleWithOptions returns a copy of s with options applied on top.
func CopySpanStyleWithOptions(s *SpanStyle, options ...SpanStyleOption) *SpanStyle {
	return copySpanStyleWithOptions(s, options...)
}

// ResolveSpanStyleDefaults resolves an unspecified SpanStyle against the framework defaults.
func ResolveSpanStyleDefaults(s *SpanStyle) *SpanStyle {
	return spanStyleResolveDefaults(s)
}

// ParagraphStyleFromOptions builds a ParagraphStyle from options applied over the unspecified style.
func ParagraphStyleFromOptions(options ...ParagraphStyleOption) *ParagraphStyle {
	return copyParagraphStyleWithOptions(paragraphStyleUnspecified, options...)
}

// CopyParagraphStyleWithOptions returns a copy of s with options applied on top.
func CopyParagraphStyleWithOptions(s *ParagraphStyle, options ...ParagraphStyleOption) *ParagraphStyle {
	return copyParagraphStyleWithOptions(s, options...)
}

// ResolveParagraphStyleDefaults resolves an unspecified ParagraphStyle against the framework defaults.
func ResolveParagraphStyleDefaults(s *ParagraphStyle, direction unit.LayoutDirection) *ParagraphStyle {
	return paragraphStyleResolveDefaults(s, direction)
}

// Identity (2 ns)
func sameTextStyle(a, b *TextStyle) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == textStyleUnspecified
	}
	if b == nil {
		return a == textStyleUnspecified
	}
	return a == b
}

// Semantic equality (field-by-field, 20 ns)
func semanticEqualTextStyle(a, b *TextStyle) bool {

	a = coalesceTextStyle(a, textStyleUnspecified)
	b = coalesceTextStyle(b, textStyleUnspecified)

	return equalSpanStyle(a.spanStyle, b.spanStyle) &&
		equalParagraphStyle(a.paragraphStyle, b.paragraphStyle)
}

func equalTextStyle(a, b *TextStyle) bool {
	if sameTextStyle(a, b) {
		return true
	}
	return semanticEqualTextStyle(a, b)
}

// TextStyleCompanion is the only exported value of this file's sentinel family member.
var TextStyleCompanion = companions.NewStructCompanion(
	func() *TextStyle { return textStyleUnspecified },
	stringTextStyle,
	isSpecifiedTextStyle,
	isUnspecifiedTextStyle,
	takeOrElseTextStyle,
	mergeTextStyle,
	sameTextStyle,
	semanticEqualTextStyle,
	equalTextStyle,
	coalesceTextStyle,
	copyTextStyle,
)
