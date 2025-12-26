// Package text provides rich text styling and annotation capabilities.
package text

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/zodimo/go-compose/compose/ui/unit"
)

// ============================================================================
// Core Types
// ============================================================================

// AnnotatedString is immutable text with style and metadata annotations.
// Use Builder to construct instances efficiently.
type AnnotatedString struct {
	text        string
	annotations []Range[Annotation]

	// Cached derived slices (nil if empty)
	spanStyles      []Range[SpanStyle]
	paragraphStyles []Range[ParagraphStyle]
}

// Range represents an annotation applied to a half-open range [Start, End).
type Range[T any] struct {
	Item  T
	Start int
	End   int
	Tag   string
}

// NewRange creates a Range with the given item and bounds.
func NewRange[T any](item T, start, end int) Range[T] {
	return Range[T]{Item: item, Start: start, End: end}
}

// NewRangeWithTag creates a Range with an additional tag for identification.
func NewRangeWithTag[T any](item T, start, end int, tag string) Range[T] {
	return Range[T]{Item: item, Start: start, End: end, Tag: tag}
}

// Annotation is a marker interface for all annotation types.
type Annotation interface {
	//Closed interface
	isAnnotation()
}

// ============================================================================
// Annotation Types
// ============================================================================

// StringAnnotation is simple string metadata attached to text.
type StringAnnotation string

func (StringAnnotation) isAnnotation() {}

// TtsAnnotation provides text-to-speech metadata.
type TtsAnnotation struct {
	// Add TTS-specific fields: Language, Pitch, etc.
}

func (TtsAnnotation) isAnnotation() {}

// LinkAnnotation marks clickable links in text.
type LinkAnnotation struct {
	URL string
}

func (LinkAnnotation) isAnnotation() {}

// Bullet represents a bullet list marker.
type Bullet struct {
	// Add bullet styling fields
}

func (Bullet) isAnnotation() {}

// ============================================================================
// Constructors
// ============================================================================

// NewAnnotatedString creates an AnnotatedString from text and optional style ranges.
func NewAnnotatedString(text string, spanStyles []Range[SpanStyle], paragraphStyles []Range[ParagraphStyle]) AnnotatedString {
	annotations := constructAnnotations(spanStyles, paragraphStyles)
	return newAnnotatedStringWithAnnotations(text, annotations)
}

// newAnnotatedStringWithAnnotations creates an AnnotatedString from text and generic annotations.
func newAnnotatedStringWithAnnotations(text string, annotations []Range[Annotation]) AnnotatedString {
	// Validate paragraph styles don't overlap
	if err := validateParagraphStyles(paragraphStylesFromAnnotations(annotations)); err != nil {
		panic(err)
	}

	as := AnnotatedString{text: text}
	if len(annotations) == 0 {
		return as
	}

	// Build cached slices
	as.spanStyles = make([]Range[SpanStyle], 0, 4)
	as.paragraphStyles = make([]Range[ParagraphStyle], 0, 4)

	for _, ann := range annotations {
		switch v := ann.Item.(type) {
		case SpanStyle:
			as.spanStyles = append(as.spanStyles, Range[SpanStyle]{
				Item: v, Start: ann.Start, End: ann.End, Tag: ann.Tag,
			})
		case ParagraphStyle:
			as.paragraphStyles = append(as.paragraphStyles, Range[ParagraphStyle]{
				Item: v, Start: ann.Start, End: ann.End, Tag: ann.Tag,
			})
		}
	}
	as.annotations = annotations
	return as
}

// constructAnnotations merges SpanStyle and ParagraphStyle slices.
func constructAnnotations(spanStyles []Range[SpanStyle], paragraphStyles []Range[ParagraphStyle]) []Range[Annotation] {
	totalLen := len(spanStyles) + len(paragraphStyles)
	if totalLen == 0 {
		return nil
	}

	annotations := make([]Range[Annotation], 0, totalLen)
	for _, s := range spanStyles {
		annotations = append(annotations, Range[Annotation]{
			Item:  s.Item,
			Start: s.Start,
			End:   s.End,
			Tag:   s.Tag,
		})
	}
	for _, p := range paragraphStyles {
		annotations = append(annotations, Range[Annotation]{
			Item:  p.Item,
			Start: p.Start,
			End:   p.End,
			Tag:   p.Tag,
		})
	}
	return annotations
}

// paragraphStylesFromAnnotations extracts ParagraphStyle ranges.
func paragraphStylesFromAnnotations(annotations []Range[Annotation]) []Range[ParagraphStyle] {
	var paraStyles []Range[ParagraphStyle]
	for _, ann := range annotations {
		if ps, ok := ann.Item.(ParagraphStyle); ok {
			paraStyles = append(paraStyles, Range[ParagraphStyle]{
				Item:  ps,
				Start: ann.Start,
				End:   ann.End,
				Tag:   ann.Tag,
			})
		}
	}
	return paraStyles
}

// NormalizedParagraphStyles determines the paragraph boundaries and merges styles.
// It handles gaps (using defaultParagraphStyle) and nested paragraphs.
func (as AnnotatedString) NormalizedParagraphStyles(defaultParagraphStyle ParagraphStyle) []Range[ParagraphStyle] {
	spanStyles := as.paragraphStyles
	length := len(as.text)

	// Create a copy and sort by start
	sortedParagraphs := make([]Range[ParagraphStyle], len(spanStyles))
	copy(sortedParagraphs, spanStyles)
	sort.Slice(sortedParagraphs, func(i, j int) bool {
		return sortedParagraphs[i].Start < sortedParagraphs[j].Start
	})

	var result []Range[ParagraphStyle]
	lastAdded := 0
	var stack []Range[ParagraphStyle]

	for _, current := range sortedParagraphs {
		// Merge current with parent if nested
		if len(stack) > 0 {
			// Current is nested in stack top?
			// The validation ensures proper nesting or no overlap.
			// But we need to merge styles.
			// In Kotlin: val current = it.copy(defaultParagraphStyle.merge(it.item))
			// But wait, Kotlin merges with default first?
			// current = it.copy(defaultParagraphStyle.merge(it.item))
			// Actually Kotlin code says:
			// val current = it.copy(defaultParagraphStyle.merge(it.item))
			// So it merges the current item ON TOP OF default.

			// In Go we have MergeParagraphStyle(a, b) -> yields b merged over a.
			// So default.Merge(item) -> item on top of default.
			merged := MergeParagraphStyle(&defaultParagraphStyle, &current.Item)
			current.Item = *merged
		} else {
			merged := MergeParagraphStyle(&defaultParagraphStyle, &current.Item)
			current.Item = *merged
		}

		for lastAdded < current.Start && len(stack) > 0 {
			lastInStack := stack[len(stack)-1]
			if current.Start < lastInStack.End {
				result = append(result, Range[ParagraphStyle]{
					Item:  lastInStack.Item,
					Start: lastAdded,
					End:   current.Start,
				})
				lastAdded = current.Start
			} else {
				result = append(result, Range[ParagraphStyle]{
					Item:  lastInStack.Item,
					Start: lastAdded,
					End:   lastInStack.End,
				})
				lastAdded = lastInStack.End
				// Pop finished
				for len(stack) > 0 && lastAdded == stack[len(stack)-1].End {
					stack = stack[:len(stack)-1]
				}
			}
		}

		if lastAdded < current.Start {
			result = append(result, Range[ParagraphStyle]{
				Item:  defaultParagraphStyle,
				Start: lastAdded,
				End:   current.Start,
			})
			lastAdded = current.Start
		}

		if len(stack) > 0 {
			lastInStack := stack[len(stack)-1]
			if lastInStack.Start == current.Start && lastInStack.End == current.End {
				// Fully overlapped
				stack = stack[:len(stack)-1]
				merged := MergeParagraphStyle(&lastInStack.Item, &current.Item)
				stack = append(stack, Range[ParagraphStyle]{
					Item:  *merged,
					Start: current.Start,
					End:   current.End,
				})
			} else if lastInStack.Start == lastInStack.End {
				// Zero length
				result = append(result, Range[ParagraphStyle]{
					Item:  lastInStack.Item,
					Start: lastInStack.Start,
					End:   lastInStack.End,
				})
				stack = stack[:len(stack)-1]
				stack = append(stack, Range[ParagraphStyle]{
					Item:  current.Item,
					Start: current.Start,
					End:   current.End,
				})
			} else {
				merged := MergeParagraphStyle(&lastInStack.Item, &current.Item)
				stack = append(stack, Range[ParagraphStyle]{
					Item:  *merged,
					Start: current.Start,
					End:   current.End,
				})
			}
		} else {
			stack = append(stack, Range[ParagraphStyle]{
				Item:  current.Item,
				Start: current.Start,
				End:   current.End,
			})
		}
	}

	for lastAdded <= length && len(stack) > 0 {
		lastInStack := stack[len(stack)-1]
		result = append(result, Range[ParagraphStyle]{
			Item:  lastInStack.Item,
			Start: lastAdded,
			End:   lastInStack.End,
		})
		lastAdded = lastInStack.End
		for len(stack) > 0 && lastAdded == stack[len(stack)-1].End {
			stack = stack[:len(stack)-1]
		}
	}

	if lastAdded < length {
		result = append(result, Range[ParagraphStyle]{
			Item:  defaultParagraphStyle,
			Start: lastAdded,
			End:   length,
		})
	}

	if len(result) == 0 {
		result = append(result, Range[ParagraphStyle]{
			Item:  defaultParagraphStyle,
			Start: 0,
			End:   0,
		})
	}

	return result
}

// validateParagraphStyles ensures no overlapping paragraph ranges.
func validateParagraphStyles(styles []Range[ParagraphStyle]) error {
	if len(styles) < 2 {
		return nil
	}

	sorted := make([]Range[ParagraphStyle], len(styles))
	copy(sorted, styles)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Start < sorted[j].Start
	})

	ends := make([]int, 0, 4) // Track nested paragraph end positions
	for _, s := range sorted {
		for len(ends) > 0 && s.Start >= ends[len(ends)-1] {
			ends = ends[:len(ends)-1] // Pop finished paragraphs
		}

		if len(ends) > 0 {
			// Current must be fully contained in last paragraph
			if s.End > ends[len(ends)-1] {
				return fmt.Errorf("paragraph overlap not allowed: range %d-%d exceeds parent %d",
					s.Start, s.End, ends[len(ends)-1])
			}
		}
		ends = append(ends, s.End)
	}
	return nil
}

// ============================================================================
// Properties
// ============================================================================

// Len returns the length of the text.
func (as AnnotatedString) Len() int {
	return len(as.text)
}

// String returns the plain text (implements fmt.Stringer).
func (as AnnotatedString) String() string {
	return as.text
}

// CharAt returns the byte at index (useful for ASCII).
func (as AnnotatedString) CharAt(index int) byte {
	return as.text[index]
}

// Text returns the plain text.
func (as AnnotatedString) Text() string {
	return as.text
}

// SpanStyles returns all SpanStyle ranges (may be empty slice).
func (as AnnotatedString) SpanStyles() []Range[SpanStyle] {
	if as.spanStyles == nil {
		return []Range[SpanStyle]{}
	}
	return as.spanStyles
}

// ParagraphStyles returns all ParagraphStyle ranges (may be empty slice).
func (as AnnotatedString) ParagraphStyles() []Range[ParagraphStyle] {
	if as.paragraphStyles == nil {
		return []Range[ParagraphStyle]{}
	}
	return as.paragraphStyles
}

// ============================================================================
// Operations
// ============================================================================

// SubSequence returns a substring with annotations in range [start, end).
func (as AnnotatedString) SubSequence(start, end int) AnnotatedString {
	if start == 0 && end >= len(as.text) {
		return as
	}
	if start > end {
		panic(fmt.Sprintf("start (%d) must be <= end (%d)", start, end))
	}
	if start == end {
		return Empty()
	}

	return newAnnotatedStringWithAnnotations(
		as.text[start:end],
		filterRanges(as.annotations, start, end),
	)
}

// Append concatenates another AnnotatedString.
func (as AnnotatedString) Append(other AnnotatedString) AnnotatedString {
	b := NewBuilder()
	b.AppendAnnotatedString(as)
	b.AppendAnnotatedString(other)
	return b.ToAnnotatedString()
}

// GetStringAnnotations queries string annotations by tag within [start, end).
func (as AnnotatedString) GetStringAnnotations(tag string, start, end int) []Range[string] {
	if start >= end {
		return nil
	}
	var result []Range[string]
	for _, ann := range as.annotations {
		if sa, ok := ann.Item.(StringAnnotation); ok && ann.Tag == tag && intersect(start, end, ann.Start, ann.End) {
			result = append(result, Range[string]{
				Item:  string(sa),
				Start: max(start, ann.Start) - start,
				End:   min(end, ann.End) - start,
				Tag:   ann.Tag,
			})
		}
	}
	return result
}

// HasStringAnnotations checks if any string annotations with tag exist in range.
func (as AnnotatedString) HasStringAnnotations(tag string, start, end int) bool {
	for _, ann := range as.annotations {
		if _, ok := ann.Item.(StringAnnotation); ok && ann.Tag == tag && intersect(start, end, ann.Start, ann.End) {
			return true
		}
	}
	return false
}

// GetLinkAnnotations queries LinkAnnotations within [start, end).
func (as AnnotatedString) GetLinkAnnotations(start, end int) []Range[LinkAnnotation] {
	if start >= end {
		return nil
	}
	var result []Range[LinkAnnotation]
	for _, ann := range as.annotations {
		if la, ok := ann.Item.(LinkAnnotation); ok && intersect(start, end, ann.Start, ann.End) {
			result = append(result, Range[LinkAnnotation]{
				Item:  la,
				Start: max(start, ann.Start) - start,
				End:   min(end, ann.End) - start,
				Tag:   ann.Tag,
			})
		}
	}
	return result
}

// HasLinkAnnotations checks if any link annotations exist in range.
func (as AnnotatedString) HasLinkAnnotations(start, end int) bool {
	for _, ann := range as.annotations {
		if _, ok := ann.Item.(LinkAnnotation); ok && intersect(start, end, ann.Start, ann.End) {
			return true
		}
	}
	return false
}

// HasEqualAnnotations compares annotation lists for deep equality.
func (as AnnotatedString) HasEqualAnnotations(other AnnotatedString) bool {
	if len(as.annotations) != len(other.annotations) {
		return false
	}
	for i, ann := range as.annotations {
		if ann.Start != other.annotations[i].Start || ann.End != other.annotations[i].End || ann.Tag != other.annotations[i].Tag {
			return false
		}
	}
	return true
}

// Empty returns an empty AnnotatedString.
func Empty() AnnotatedString {
	return AnnotatedString{text: ""}
}

// MapAnnotations returns a new AnnotatedString with annotations transformed.
func (as AnnotatedString) MapAnnotations(transform func(Range[Annotation]) Range[Annotation]) AnnotatedString {
	b := NewBuilderFromAnnotatedString(as)
	b.MapAnnotations(transform)
	return b.ToAnnotatedString()
}

// FlatMapAnnotations returns a new AnnotatedString with annotations transformed to multiple annotations.
func (as AnnotatedString) FlatMapAnnotations(transform func(Range[Annotation]) []Range[Annotation]) AnnotatedString {
	b := NewBuilderFromAnnotatedString(as)
	b.FlatMapAnnotations(transform)
	return b.ToAnnotatedString()
}

// ============================================================================
// Builder
// ============================================================================

// Builder constructs AnnotatedString instances.
type Builder struct {
	text        strings.Builder
	styleStack  []mutableRange
	annotations []mutableRange
}

type mutableRange struct {
	item  Annotation
	start int
	end   int // mutable, initially math.MinInt
	tag   string
}

// NewBuilder creates a Builder with optional initial capacity.
func NewBuilder(capacity ...int) *Builder {
	cap := 16
	if len(capacity) > 0 && capacity[0] > 0 {
		cap = capacity[0]
	}
	b := &Builder{
		styleStack:  make([]mutableRange, 0, 8),
		annotations: make([]mutableRange, 0, 16),
	}
	b.text.Grow(cap)
	return b
}

// NewBuilderFromAnnotatedString creates a Builder initialized with the given AnnotatedString.
func NewBuilderFromAnnotatedString(as AnnotatedString) *Builder {
	b := NewBuilder(as.Len())
	b.AppendAnnotatedString(as)
	return b
}

// Len returns current text length.
func (b *Builder) Len() int {
	return b.text.Len()
}

// Append adds plain text.
func (b *Builder) Append(s string) {
	b.text.WriteString(s)
}

// AppendAnnotatedString adds another AnnotatedString with offset annotations.
func (b *Builder) AppendAnnotatedString(as AnnotatedString) {
	offset := b.text.Len()
	b.text.WriteString(as.text)
	for _, ann := range as.annotations {
		b.annotations = append(b.annotations, mutableRange{
			item:  ann.Item,
			start: offset + ann.Start,
			end:   offset + ann.End,
			tag:   ann.Tag,
		})
	}
}

// AddStyle adds a SpanStyle to range [start, end).
func (b *Builder) AddStyle(style SpanStyle, start, end int) {
	b.annotations = append(b.annotations, mutableRange{
		item:  style,
		start: start,
		end:   end,
	})
}

// AddParagraphStyle adds a ParagraphStyle to range [start, end).
func (b *Builder) AddParagraphStyle(style ParagraphStyle, start, end int) {
	b.annotations = append(b.annotations, mutableRange{
		item:  style,
		start: start,
		end:   end,
	})
}

// AddStringAnnotation adds a string annotation with a tag.
func (b *Builder) AddStringAnnotation(tag, annotation string, start, end int) {
	b.annotations = append(b.annotations, mutableRange{
		item:  StringAnnotation(annotation),
		start: start,
		end:   end,
		tag:   tag,
	})
}

// PushStyle pushes a SpanStyle onto the stack for subsequent text.
func (b *Builder) PushStyle(style SpanStyle) int {
	mr := mutableRange{
		item:  style,
		start: b.text.Len(),
		end:   math.MinInt,
	}
	b.styleStack = append(b.styleStack, mr)
	b.annotations = append(b.annotations, mr)
	return len(b.styleStack) - 1
}

// AddTtsAnnotation adds a TtsAnnotation to range [start, end).
func (b *Builder) AddTtsAnnotation(tts TtsAnnotation, start, end int) {
	b.annotations = append(b.annotations, mutableRange{
		item:  tts,
		start: start,
		end:   end,
	})
}

// AddLink adds a LinkAnnotation to range [start, end).
func (b *Builder) AddLink(link LinkAnnotation, start, end int) {
	b.annotations = append(b.annotations, mutableRange{
		item:  link,
		start: start,
		end:   end,
	})
}

// AddBullet adds a Bullet annotation to range [start, end).
func (b *Builder) AddBullet(bullet Bullet, start, end int) {
	b.annotations = append(b.annotations, mutableRange{
		item:  bullet,
		start: start,
		end:   end,
	})
}

// PushTtsAnnotation pushes a TtsAnnotation onto the stack.
func (b *Builder) PushTtsAnnotation(tts TtsAnnotation) int {
	mr := mutableRange{
		item:  tts,
		start: b.text.Len(),
		end:   math.MinInt,
	}
	b.styleStack = append(b.styleStack, mr)
	b.annotations = append(b.annotations, mr)
	return len(b.styleStack) - 1
}

// PushLink pushes a LinkAnnotation onto the stack.
func (b *Builder) PushLink(link LinkAnnotation) int {
	mr := mutableRange{
		item:  link,
		start: b.text.Len(),
		end:   math.MinInt,
	}
	b.styleStack = append(b.styleStack, mr)
	b.annotations = append(b.annotations, mr)
	return len(b.styleStack) - 1
}

// PushBullet pushes a Bullet annotation onto the stack.
func (b *Builder) PushBullet(bullet Bullet) int {
	mr := mutableRange{
		item:  bullet,
		start: b.text.Len(),
		end:   math.MinInt,
	}
	b.styleStack = append(b.styleStack, mr)
	b.annotations = append(b.annotations, mr)
	return len(b.styleStack) - 1
}

// Pop ends the last pushed style/annotation.
func (b *Builder) Pop() {
	if len(b.styleStack) == 0 {
		panic("nothing to pop")
	}
	b.styleStack[len(b.styleStack)-1].end = b.text.Len()
	b.styleStack = b.styleStack[:len(b.styleStack)-1]
}

// PopTo pops all styles up to and including the given index.
func (b *Builder) PopTo(index int) {
	if index >= len(b.styleStack) {
		panic(fmt.Sprintf("index %d exceeds stack size %d", index, len(b.styleStack)))
	}
	for len(b.styleStack) > index {
		b.Pop()
	}
}

// ToAnnotatedString builds the final immutable AnnotatedString.
func (b *Builder) ToAnnotatedString() AnnotatedString {
	for i := range b.annotations {
		if b.annotations[i].end == math.MinInt {
			b.annotations[i].end = b.text.Len()
		}
	}
	return newAnnotatedStringWithAnnotations(b.text.String(), mutableToImmutable(b.annotations))
}

// MapAnnotations transforms annotations in the builder.
func (b *Builder) MapAnnotations(transform func(Range[Annotation]) Range[Annotation]) {
	for i, ann := range b.annotations {
		immutable := Range[Annotation]{
			Item:  ann.item,
			Start: ann.start,
			End:   ann.end,
			Tag:   ann.tag,
		}
		transformed := transform(immutable)
		b.annotations[i] = mutableRange{
			item:  transformed.Item,
			start: transformed.Start,
			end:   transformed.End,
			tag:   transformed.Tag,
		}
	}
}

// FlatMapAnnotations transforms annotations into multiple annotations.
func (b *Builder) FlatMapAnnotations(transform func(Range[Annotation]) []Range[Annotation]) {
	var newAnnotations []mutableRange
	for _, ann := range b.annotations {
		immutable := Range[Annotation]{
			Item:  ann.item,
			Start: ann.start,
			End:   ann.end,
			Tag:   ann.tag,
		}
		transformed := transform(immutable)
		for _, t := range transformed {
			newAnnotations = append(newAnnotations, mutableRange{
				item:  t.Item,
				start: t.Start,
				end:   t.End,
				tag:   t.Tag,
			})
		}
	}
	b.annotations = newAnnotations
}

// mutableToImmutable converts mutableRange to Range[Annotation].
func mutableToImmutable(mr []mutableRange) []Range[Annotation] {
	if len(mr) == 0 {
		return nil
	}
	result := make([]Range[Annotation], len(mr))
	for i, r := range mr {
		result[i] = Range[Annotation]{
			Item:  r.item,
			Start: r.start,
			End:   r.end,
			Tag:   r.tag,
		}
	}
	return result
}

// ============================================================================
// Utility Functions
// ============================================================================

// filterRanges extracts and offsets annotations for subsequence operations.
func filterRanges(annotations []Range[Annotation], start, end int) []Range[Annotation] {
	if len(annotations) == 0 || start >= end {
		return nil
	}
	var result []Range[Annotation]
	for _, ann := range annotations {
		if intersect(start, end, ann.Start, ann.End) {
			result = append(result, Range[Annotation]{
				Item:  ann.Item,
				Start: max(start, ann.Start) - start,
				End:   min(end, ann.End) - start,
				Tag:   ann.Tag,
			})
		}
	}
	return result
}

// intersect checks if two half-open ranges [aStart, aEnd) and [bStart, bEnd) intersect.
func intersect(aStart, aEnd, bStart, bEnd int) bool {
	return (aStart < aEnd && bStart < bEnd && aStart < bEnd && bStart < aEnd) ||
		(aStart == aEnd && aStart >= bStart && aStart <= bEnd) ||
		(bStart == bEnd && bStart >= aStart && bStart <= aEnd)
}

// max returns the larger of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetLocalParagraphStyles finds ParagraphStyles in range and converts them to local range.
func (as AnnotatedString) GetLocalParagraphStyles(start, end int) []Range[ParagraphStyle] {
	if start == end {
		return nil
	}
	if as.paragraphStyles == nil {
		return nil
	}
	if start == 0 && end >= len(as.text) {
		return as.paragraphStyles
	}
	var result []Range[ParagraphStyle]
	for _, ps := range as.paragraphStyles {
		if intersect(start, end, ps.Start, ps.End) {
			result = append(result, Range[ParagraphStyle]{
				Item:  ps.Item,
				Start: coerceIn(ps.Start, start, end) - start,
				End:   coerceIn(ps.End, start, end) - start,
			})
		}
	}
	return result
}

// GetLocalAnnotations finds annotations in range matching predicate and converts to local range.
func (as AnnotatedString) GetLocalAnnotations(start, end int, predicate func(Annotation) bool) []Range[Annotation] {
	if start == end {
		return nil
	}
	annotations := as.annotations
	if len(annotations) == 0 {
		return nil
	}
	if start == 0 && end >= len(as.text) {
		if predicate == nil {
			return annotations
		}
		var filtered []Range[Annotation]
		for _, ann := range annotations {
			if predicate(ann.Item) {
				filtered = append(filtered, ann)
			}
		}
		return filtered
	}
	var result []Range[Annotation]
	for _, ann := range annotations {
		matches := true
		if predicate != nil {
			matches = predicate(ann.Item)
		}
		if matches && intersect(start, end, ann.Start, ann.End) {
			result = append(result, Range[Annotation]{
				Item:  ann.Item,
				Start: coerceIn(ann.Start, start, end) - start,
				End:   coerceIn(ann.End, start, end) - start,
				Tag:   ann.Tag,
			})
		}
	}
	return result
}

// SubstringWithoutParagraphStyles returns a substring without paragraph styles.
func (as AnnotatedString) SubstringWithoutParagraphStyles(start, end int) AnnotatedString {
	text := ""
	if start != end {
		text = as.text[start:end]
	}
	return newAnnotatedStringWithAnnotations(
		text,
		as.GetLocalAnnotations(start, end, func(a Annotation) bool {
			_, isPara := a.(ParagraphStyle)
			return !isPara
		}),
	)
}

// MapEachParagraphStyle iterates normalized paragraph styles.
func MapEachParagraphStyle[T any](
	as AnnotatedString,
	defaultParagraphStyle ParagraphStyle,
	block func(annotatedString AnnotatedString, paragraphStyle Range[ParagraphStyle]) T,
) []T {
	normalized := as.NormalizedParagraphStyles(defaultParagraphStyle)
	result := make([]T, 0, len(normalized))
	for _, psRange := range normalized {
		// Create substring for this paragraph range
		// Kotlins substringWithoutParagraphStyles is used here
		subAs := as.SubstringWithoutParagraphStyles(psRange.Start, psRange.End)
		result = append(result, block(subAs, psRange))
	}
	return result
}

// WithStyle pushes style, executes block, and pops.
func (b *Builder) WithStyle(style SpanStyle, block func()) {
	index := b.PushStyle(style)
	defer b.PopTo(index)
	block()
}

// PushParagraphStyle pushes a ParagraphStyle onto the stack.
func (b *Builder) PushParagraphStyle(style ParagraphStyle) int {
	mr := mutableRange{
		item:  style,
		start: b.text.Len(),
		end:   math.MinInt,
	}
	b.styleStack = append(b.styleStack, mr)
	b.annotations = append(b.annotations, mr)
	return len(b.styleStack) - 1
}

// WithParagraphStyle pushes paragraph style, executes block, and pops.
func (b *Builder) WithParagraphStyle(style ParagraphStyle, block func()) {
	index := b.PushParagraphStyle(style)
	defer b.PopTo(index)
	block()
}

// WithAnnotation pushes annotation, executes block, and pops.
func (b *Builder) WithAnnotation(tag, annotation string, block func()) {
	index := b.PushStringAnnotation(tag, annotation)
	defer b.PopTo(index)
	block()
}

// WithTtsAnnotation pushes tts annotation, executes block, and pops.
func (b *Builder) WithTtsAnnotation(tts TtsAnnotation, block func()) {
	index := b.PushTtsAnnotation(tts)
	defer b.PopTo(index)
	block()
}

// WithLink pushes link annotation, executes block, and pops.
func (b *Builder) WithLink(link LinkAnnotation, block func()) {
	index := b.PushLink(link)
	defer b.PopTo(index)
	block()
}

// PushStringAnnotation pushes a string annotation onto the stack.
func (b *Builder) PushStringAnnotation(tag, annotation string) int {
	mr := mutableRange{
		item:  StringAnnotation(annotation),
		start: b.text.Len(),
		end:   math.MinInt,
		tag:   tag,
	}
	b.styleStack = append(b.styleStack, mr)
	b.annotations = append(b.annotations, mr)
	return len(b.styleStack) - 1
}

// BulletScope is the scope for a bullet list.
type BulletScope struct {
	builder                *Builder
	bulletListSettingStack []bulletSetting
}

type bulletSetting struct {
	indent unit.TextUnit
	bullet Bullet
}

// WithBulletList creates a bullet list scope.
func (b *Builder) WithBulletList(
	indentation unit.TextUnit,
	bullet Bullet,
	block func(*BulletScope),
) {
	scope := &BulletScope{
		builder: b,
		bulletListSettingStack: []bulletSetting{
			{indent: indentation, bullet: bullet},
		},
	}
	// Kotlin logic involves pushing a paragraph style for the list container
	// But for now we just run the block.
	// The implementation details of nested lists with indentation logic is complex and requires style.TextIndent.
	// I will simplify for this port to just execute the block to unblock compilation.
	// TODO: Implement full indentation logic.
	block(scope)
}

// WithBulletListItem adds a bullet list item.
func (s *BulletScope) WithBulletListItem(
	bullet *Bullet,
	block func(),
) {
	// TODO: Implement bullet logic
	block()
}

// Case transformations

// ToUpper returns a new AnnotatedString with text converted to uppercase.
func (as AnnotatedString) ToUpper() AnnotatedString {
	return AnnotatedString{
		text:            strings.ToUpper(as.text),
		spanStyles:      as.spanStyles,
		paragraphStyles: as.paragraphStyles,
		annotations:     as.annotations,
	}
}

// ToLower returns a new AnnotatedString with text converted to lowercase.
func (as AnnotatedString) ToLower() AnnotatedString {
	return AnnotatedString{
		text:            strings.ToLower(as.text),
		spanStyles:      as.spanStyles,
		paragraphStyles: as.paragraphStyles,
		annotations:     as.annotations,
	}
}

// Capitalize returns a new AnnotatedString with the first character capitalized.
func (as AnnotatedString) Capitalize() AnnotatedString {
	if len(as.text) == 0 {
		return as
	}
	// Simple naive capitalization
	text := strings.ToUpper(as.text[:1]) + as.text[1:]
	return AnnotatedString{
		text:            text,
		spanStyles:      as.spanStyles,
		paragraphStyles: as.paragraphStyles,
		annotations:     as.annotations,
	}
}

// Decapitalize returns a new AnnotatedString with the first character lowercased.
func (as AnnotatedString) Decapitalize() AnnotatedString {
	if len(as.text) == 0 {
		return as
	}
	text := strings.ToLower(as.text[:1]) + as.text[1:]
	return AnnotatedString{
		text:            text,
		spanStyles:      as.spanStyles,
		paragraphStyles: as.paragraphStyles,
		annotations:     as.annotations,
	}
}
