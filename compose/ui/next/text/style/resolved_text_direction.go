package style

// ResolvedTextDirection describes the resolved/actual directionality of text.
//
// Unlike TextDirection which specifies the algorithm to determine direction,
// ResolvedTextDirection represents the final resolved direction.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/next/text/style/ResolvedTextDirection.kt
type ResolvedTextDirection int

const (
	// ResolvedTextDirectionLtr represents text that is left-to-right.
	ResolvedTextDirectionLtr ResolvedTextDirection = iota

	// ResolvedTextDirectionRtl represents text that is right-to-left.
	ResolvedTextDirectionRtl
)

// String returns the string representation of the ResolvedTextDirection.
func (d ResolvedTextDirection) String() string {
	switch d {
	case ResolvedTextDirectionLtr:
		return "Ltr"
	case ResolvedTextDirectionRtl:
		return "Rtl"
	default:
		return "Unknown"
	}
}
