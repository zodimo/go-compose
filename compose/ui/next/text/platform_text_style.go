package text

import (
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// platformTextStyleUnspecified is the sentinel for PlatformTextStyle. It is
// deliberately package-private: the only public way to obtain it is
// PlatformTextStyleCompanion.Unspecified().
var platformTextStyleUnspecified = &PlatformTextStyle{
	SpanStyle:      platformSpanStyleUnspecified,
	ParagraphStyle: platformParagraphStyleUnspecified,
}

// platformSpanStyleUnspecified is the sentinel for PlatformSpanStyle.
var platformSpanStyleUnspecified = &PlatformSpanStyle{}

// platformParagraphStyleUnspecified is the sentinel for PlatformParagraphStyle.
var platformParagraphStyleUnspecified = &PlatformParagraphStyle{}

// PlatformSpanStyle contains platform-specific text styling.
type PlatformSpanStyle struct {
}

// PlatformParagraphStyle contains platform-specific paragraph styling.
type PlatformParagraphStyle struct {
}

// PlatformTextStyle contains platform-specific text styling combining span and paragraph styles.
type PlatformTextStyle struct {
	SpanStyle      *PlatformSpanStyle
	ParagraphStyle *PlatformParagraphStyle
}

func IsSpecifiedPlatformTextStyle(s *PlatformTextStyle) bool {
	return isSpecifiedPlatformTextStyle(s)
}

func IsSpecifiedPlatformSpanStyle(s *PlatformSpanStyle) bool {
	return isSpecifiedPlatformSpanStyle(s)
}

func IsSpecifiedPlatformParagraphStyle(s *PlatformParagraphStyle) bool {
	return isSpecifiedPlatformParagraphStyle(s)
}

func TakeOrElsePlatformTextStyle(s, def *PlatformTextStyle) *PlatformTextStyle {
	return takeOrElsePlatformTextStyle(s, def)
}

func TakeOrElsePlatformSpanStyle(s, def *PlatformSpanStyle) *PlatformSpanStyle {
	return takeOrElsePlatformSpanStyle(s, def)
}

func TakeOrElsePlatformParagraphStyle(s, def *PlatformParagraphStyle) *PlatformParagraphStyle {
	return takeOrElsePlatformParagraphStyle(s, def)
}

func EqualPlatformTextStyle(a, b *PlatformTextStyle) bool {
	return equalPlatformTextStyle(a, b)
}

func EqualPlatformSpanStyle(a, b *PlatformSpanStyle) bool {
	return equalPlatformSpanStyle(a, b)
}

func EqualPlatformParagraphStyle(a, b *PlatformParagraphStyle) bool {
	return equalPlatformParagraphStyle(a, b)
}

func StringPlatformParagraphStyle(s *PlatformParagraphStyle) string {
	if !isSpecifiedPlatformParagraphStyle(s) {
		return "PlatformParagraphStyle.Unspecified"
	}
	// Currently empty struct
	return "PlatformParagraphStyle()"
}

// === Sentinel operations ===

func isSpecifiedPlatformTextStyle(s *PlatformTextStyle) bool {
	return s != nil && s != platformTextStyleUnspecified
}

func isUnspecifiedPlatformTextStyle(s *PlatformTextStyle) bool {
	return !isSpecifiedPlatformTextStyle(s)
}

func isSpecifiedPlatformSpanStyle(s *PlatformSpanStyle) bool {
	return s != nil && s != platformSpanStyleUnspecified
}

func isUnspecifiedPlatformSpanStyle(s *PlatformSpanStyle) bool {
	return !isSpecifiedPlatformSpanStyle(s)
}

func isSpecifiedPlatformParagraphStyle(s *PlatformParagraphStyle) bool {
	return s != nil && s != platformParagraphStyleUnspecified
}

func isUnspecifiedPlatformParagraphStyle(s *PlatformParagraphStyle) bool {
	return !isSpecifiedPlatformParagraphStyle(s)
}

func takeOrElsePlatformTextStyle(s, def *PlatformTextStyle) *PlatformTextStyle {
	if !isSpecifiedPlatformTextStyle(s) {
		return def
	}
	return s
}

func takeOrElsePlatformSpanStyle(s, def *PlatformSpanStyle) *PlatformSpanStyle {
	if !isSpecifiedPlatformSpanStyle(s) {
		return def
	}
	return s
}

func takeOrElsePlatformParagraphStyle(s, def *PlatformParagraphStyle) *PlatformParagraphStyle {
	if !isSpecifiedPlatformParagraphStyle(s) {
		return def
	}
	return s
}

func coalescePlatformTextStyle(ptr, def *PlatformTextStyle) *PlatformTextStyle {
	if ptr == nil {
		return def
	}
	return ptr
}

func coalescePlatformSpanStyle(ptr, def *PlatformSpanStyle) *PlatformSpanStyle {
	if ptr == nil {
		return def
	}
	return ptr
}

func coalescePlatformParagraphStyle(ptr, def *PlatformParagraphStyle) *PlatformParagraphStyle {
	if ptr == nil {
		return def
	}
	return ptr
}

func mergePlatformTextStyle(a, b *PlatformTextStyle) *PlatformTextStyle {
	if isUnspecifiedPlatformTextStyle(b) {
		return coalescePlatformTextStyle(a, platformTextStyleUnspecified)
	}
	return b
}

func mergePlatformSpanStyle(a, b *PlatformSpanStyle) *PlatformSpanStyle {
	if isUnspecifiedPlatformSpanStyle(b) {
		return coalescePlatformSpanStyle(a, platformSpanStyleUnspecified)
	}
	return b
}

func mergePlatformParagraphStyle(a, b *PlatformParagraphStyle) *PlatformParagraphStyle {
	if isUnspecifiedPlatformParagraphStyle(b) {
		return coalescePlatformParagraphStyle(a, platformParagraphStyleUnspecified)
	}
	return b
}

func samePlatformTextStyle(a, b *PlatformTextStyle) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == platformTextStyleUnspecified
	}
	if b == nil {
		return a == platformTextStyleUnspecified
	}
	return a == b
}

func samePlatformSpanStyle(a, b *PlatformSpanStyle) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == platformSpanStyleUnspecified
	}
	if b == nil {
		return a == platformSpanStyleUnspecified
	}
	return a == b
}

func samePlatformParagraphStyle(a, b *PlatformParagraphStyle) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == platformParagraphStyleUnspecified
	}
	if b == nil {
		return a == platformParagraphStyleUnspecified
	}
	return a == b
}

func semanticEqualPlatformTextStyle(a, b *PlatformTextStyle) bool {
	a = coalescePlatformTextStyle(a, platformTextStyleUnspecified)
	b = coalescePlatformTextStyle(b, platformTextStyleUnspecified)
	if a == platformTextStyleUnspecified && b == platformTextStyleUnspecified {
		return true
	}
	return true
}

func semanticEqualPlatformSpanStyle(a, b *PlatformSpanStyle) bool {
	a = coalescePlatformSpanStyle(a, platformSpanStyleUnspecified)
	b = coalescePlatformSpanStyle(b, platformSpanStyleUnspecified)
	if a == platformSpanStyleUnspecified && b == platformSpanStyleUnspecified {
		return true
	}
	// For now it's an empty struct, so just true if both specified.
	return true
}

func semanticEqualPlatformParagraphStyle(a, b *PlatformParagraphStyle) bool {
	a = coalescePlatformParagraphStyle(a, platformParagraphStyleUnspecified)
	b = coalescePlatformParagraphStyle(b, platformParagraphStyleUnspecified)
	if a == platformParagraphStyleUnspecified && b == platformParagraphStyleUnspecified {
		return true
	}
	return true
}

func equalPlatformTextStyle(a, b *PlatformTextStyle) bool {
	if !samePlatformTextStyle(a, b) {
		return semanticEqualPlatformTextStyle(a, b)
	}
	return true
}

func equalPlatformSpanStyle(a, b *PlatformSpanStyle) bool {
	if !samePlatformSpanStyle(a, b) {
		return semanticEqualPlatformSpanStyle(a, b)
	}
	return true
}

func equalPlatformParagraphStyle(a, b *PlatformParagraphStyle) bool {
	if !samePlatformParagraphStyle(a, b) {
		return semanticEqualPlatformParagraphStyle(a, b)
	}
	return true
}

func stringPlatformTextStyle(s *PlatformTextStyle) string {
	if !isSpecifiedPlatformTextStyle(s) {
		return "PlatformTextStyle.Unspecified"
	}
	return "PlatformTextStyle()"
}

func stringPlatformSpanStyle(s *PlatformSpanStyle) string {
	if !isSpecifiedPlatformSpanStyle(s) {
		return "PlatformSpanStyle.Unspecified"
	}
	return "PlatformSpanStyle()"
}

func stringPlatformParagraphStyle(s *PlatformParagraphStyle) string {
	return StringPlatformParagraphStyle(s)
}

func copyPlatformTextStyle(s *PlatformTextStyle, options ...func(*PlatformTextStyle)) *PlatformTextStyle {
	cp := &PlatformTextStyle{}
	if isSpecifiedPlatformTextStyle(s) {
		cp.SpanStyle = s.SpanStyle
		cp.ParagraphStyle = s.ParagraphStyle
	}
	for _, option := range options {
		option(cp)
	}
	return cp
}

func copyPlatformSpanStyle(s *PlatformSpanStyle, options ...func(*PlatformSpanStyle)) *PlatformSpanStyle {
	cp := &PlatformSpanStyle{}
	for _, option := range options {
		option(cp)
	}
	return cp
}

func copyPlatformParagraphStyle(s *PlatformParagraphStyle, options ...func(*PlatformParagraphStyle)) *PlatformParagraphStyle {
	cp := &PlatformParagraphStyle{}
	for _, option := range options {
		option(cp)
	}
	return cp
}

// PlatformTextStyleCompanion is the single exported sentinel-family symbol for PlatformTextStyle.
var PlatformTextStyleCompanion = companions.NewStructCompanion(
	func() *PlatformTextStyle { return platformTextStyleUnspecified },
	stringPlatformTextStyle,
	isSpecifiedPlatformTextStyle,
	isUnspecifiedPlatformTextStyle,
	takeOrElsePlatformTextStyle,
	mergePlatformTextStyle,
	samePlatformTextStyle,
	semanticEqualPlatformTextStyle,
	equalPlatformTextStyle,
	coalescePlatformTextStyle,
	copyPlatformTextStyle,
)

// PlatformSpanStyleCompanion is the single exported sentinel-family symbol for PlatformSpanStyle.
var PlatformSpanStyleCompanion = companions.NewStructCompanion(
	func() *PlatformSpanStyle { return platformSpanStyleUnspecified },
	stringPlatformSpanStyle,
	isSpecifiedPlatformSpanStyle,
	isUnspecifiedPlatformSpanStyle,
	takeOrElsePlatformSpanStyle,
	mergePlatformSpanStyle,
	samePlatformSpanStyle,
	semanticEqualPlatformSpanStyle,
	equalPlatformSpanStyle,
	coalescePlatformSpanStyle,
	copyPlatformSpanStyle,
)

// PlatformParagraphStyleCompanion is the single exported sentinel-family symbol for PlatformParagraphStyle.
var PlatformParagraphStyleCompanion = companions.NewStructCompanion(
	func() *PlatformParagraphStyle { return platformParagraphStyleUnspecified },
	stringPlatformParagraphStyle,
	isSpecifiedPlatformParagraphStyle,
	isUnspecifiedPlatformParagraphStyle,
	takeOrElsePlatformParagraphStyle,
	mergePlatformParagraphStyle,
	samePlatformParagraphStyle,
	semanticEqualPlatformParagraphStyle,
	equalPlatformParagraphStyle,
	coalescePlatformParagraphStyle,
	copyPlatformParagraphStyle,
)
