package intl

import (
	"fmt"

	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/next/text/intl/LocaleList.kt;drc=4970f6e96cdb06089723da0ab8ec93ae3f067c7a;l=31

// localeListUnspecified is the sentinel for LocaleList. It is deliberately
// package-private: the only public way to obtain it is
// LocaleListCompanion.Unspecified().
var localeListUnspecified = &LocaleList{}

type LocaleList struct {
	List []Locale
}

var LocaleListEmpty = &LocaleList{
	List: []Locale{},
}

var LocalListCurrent = &LocaleList{
	List: []Locale{},
}

func IsSpecifiedLocaleList(s *LocaleList) bool {
	return LocaleListCompanion.IsSpecified(s)
}
func TakeOrElseLocaleList(s, def *LocaleList) *LocaleList {
	return takeOrElseLocaleList(s, def)
}

// Identity (2 ns)
func SameLocaleList(a, b *LocaleList) bool {
	return sameLocaleList(a, b)
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualLocaleList(a, b *LocaleList) bool {
	return semanticEqualLocaleList(a, b)
}

// EqualLocaleList returns true if both LocaleLists have the same locales.
func EqualLocaleList(a, b *LocaleList) bool {
	return equalLocaleList(a, b)
}

func CoalesceLocaleList(ptr, def *LocaleList) *LocaleList {
	return coalesceLocaleList(ptr, def)
}

// === Sentinel operations ===

func isSpecifiedLocaleList(s *LocaleList) bool {
	return s != nil && s != localeListUnspecified
}

func isUnspecifiedLocaleList(s *LocaleList) bool {
	return !isSpecifiedLocaleList(s)
}

func takeOrElseLocaleList(s, def *LocaleList) *LocaleList {
	if !isSpecifiedLocaleList(s) {
		return def
	}
	return s
}

func coalesceLocaleList(ptr, def *LocaleList) *LocaleList {
	if ptr == nil {
		return def
	}
	return ptr
}

// sameLocaleList is identity (2 ns).
func sameLocaleList(a, b *LocaleList) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == localeListUnspecified
	}
	if b == nil {
		return a == localeListUnspecified
	}
	return a == b
}

// semanticEqualLocaleList is field-by-field equality (20 ns). Two unspecified
// lists are equal.
func semanticEqualLocaleList(a, b *LocaleList) bool {

	a = coalesceLocaleList(a, localeListUnspecified)
	b = coalesceLocaleList(b, localeListUnspecified)

	if a == localeListUnspecified && b == localeListUnspecified {
		return true
	}
	if len(a.List) != len(b.List) {
		return false
	}
	for i, locale := range a.List {
		if locale != b.List[i] {
			return false
		}
	}
	return true
}

func equalLocaleList(a, b *LocaleList) bool {
	if !sameLocaleList(a, b) {
		return semanticEqualLocaleList(a, b)
	}
	return true
}

func mergeLocaleList(a, b *LocaleList) *LocaleList {
	if isUnspecifiedLocaleList(b) {
		return coalesceLocaleList(a, localeListUnspecified)
	}
	return b
}

func stringLocaleList(s *LocaleList) string {
	if !isSpecifiedLocaleList(s) {
		return "LocaleList{Unspecified}"
	}
	return "LocaleList{" + localeListItems(s) + "}"
}

func localeListItems(s *LocaleList) string {
	out := ""
	for i, locale := range s.List {
		if i > 0 {
			out += ", "
		}
		out += fmt.Sprintf("%v", locale)
	}
	return out
}

// copyLocaleList returns a copy of s. The source is never mutated.
func copyLocaleList(s *LocaleList, options ...func(*LocaleList)) *LocaleList {
	cp := &LocaleList{}
	if s != nil {
		cp.List = append([]Locale(nil), s.List...)
	}
	for _, option := range options {
		option(cp)
	}
	return cp
}

// LocaleListCompanion is the single exported sentinel-family symbol for LocaleList.
var LocaleListCompanion = companions.NewStructCompanion(
	func() *LocaleList { return localeListUnspecified },
	stringLocaleList,
	isSpecifiedLocaleList,
	isUnspecifiedLocaleList,
	takeOrElseLocaleList,
	mergeLocaleList,
	sameLocaleList,
	semanticEqualLocaleList,
	equalLocaleList,
	coalesceLocaleList,
	copyLocaleList,
)
