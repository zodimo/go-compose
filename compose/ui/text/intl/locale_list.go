package intl

// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/next/text/intl/LocaleList.kt;drc=4970f6e96cdb06089723da0ab8ec93ae3f067c7a;l=31

var LocaleListUnspecified = &LocaleList{}

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
	return s != nil && s != LocaleListUnspecified
}
func TakeOrElseLocaleList(s, def *LocaleList) *LocaleList {
	if !IsSpecifiedLocaleList(s) {
		return def
	}
	return s
}

// Identity (2 ns)
func SameLocaleList(a, b *LocaleList) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == LocaleListUnspecified
	}
	if b == nil {
		return a == LocaleListUnspecified
	}
	return a == b
}

// Semantic equality (field-by-field, 20 ns)
func SemanticEqualLocaleList(a, b *LocaleList) bool {

	a = CoalesceLocaleList(a, LocaleListUnspecified)
	b = CoalesceLocaleList(b, LocaleListUnspecified)

	for i, locale := range a.List {
		if locale != b.List[i] {
			return false
		}
	}
	return true
}

// EqualLocaleList returns true if both LocaleLists have the same locales.
func EqualLocaleList(a, b *LocaleList) bool {
	if !SameLocaleList(a, b) {
		return SemanticEqualLocaleList(a, b)
	}
	return true
}

func CoalesceLocaleList(ptr, def *LocaleList) *LocaleList {
	if ptr == nil {
		return def
	}
	return ptr
}
