package intl

// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/text/intl/LocaleList.kt;drc=4970f6e96cdb06089723da0ab8ec93ae3f067c7a;l=31

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

// Equals returns true if both LocaleLists have the same locales.
func (l LocaleList) Equals(other *LocaleList) bool {
	if other == nil {
		return false
	}
	if len(l.List) != len(other.List) {
		return false
	}
	for i, locale := range l.List {
		if locale != other.List[i] {
			return false
		}
	}
	return true
}
