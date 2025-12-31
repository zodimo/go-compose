package unit

// LayoutDirection represents the horizontal layout direction of content.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-unit/src/commonMain/kotlin/androidx/compose/ui/unit/LayoutDirection.kt
type LayoutDirection int

const (
	// LayoutDirectionLtr represents left-to-right layout direction.
	LayoutDirectionLtr LayoutDirection = iota

	// LayoutDirectionRtl represents right-to-left layout direction.
	LayoutDirectionRtl
)

// String returns the string representation of the LayoutDirection.
func (d LayoutDirection) String() string {
	switch d {
	case LayoutDirectionLtr:
		return "Ltr"
	case LayoutDirectionRtl:
		return "Rtl"
	default:
		return "Unknown"
	}
}
