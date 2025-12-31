package style

// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-text/src/commonMain/kotlin/androidx/compose/ui/next/text/style/TextOverflow.kt;drc=4970f6e96cdb06089723da0ab8ec93ae3f067c7a;l=23
type TextOverFlow int

const (
	// Clip the overflowing text to fix its container.
	OverFlowClip TextOverFlow = 1

	// Use an ellipsis at the end of the string to indicate that the text has overflowed.
	// For example, [This is a ...].
	OverFlowEllipsis TextOverFlow = 2

	// Display all text, even if there is not enough space in the specified bounds. When
	// overflow is visible, text may be rendered outside the bounds of the composable displaying
	// the text. This ensures that all text is displayed to the user, and is typically the right
	// choice for most text display. It does mean that the text may visually occupy a region
	// larger than the bounds of its composable. This can lead to situations where text
	// displays outside the bounds of the background and clickable on a Text composable with a
	// fixed height and width.
	//
	// To make the background and click region expand to match the size of the text, allow it to
	// expand vertically/horizontally using Modifier.heightIn/Modifier.widthIn or similar.
	//
	// Note: text that expands past its bounds using Visible may be clipped by other modifiers
	// such as Modifier.clipToBounds.
	OverFlowVisible TextOverFlow = 3

	// Use an ellipsis at the start of the string to indicate that the text has overflowed.
	// For example, [... is a text].
	//
	// Note that not all platforms support the ellipsis at the start. For example, on Android
	// the start ellipsis is only available for a single line text (i.e. when either a soft wrap
	// is disabled or a maximum number of lines maxLines set to 1). In case of multiline text it
	// will fallback to Clip.
	OverFlowStartEllipsis TextOverFlow = 4

	// Use an ellipsis in the middle of the string to indicate that the text has overflowed.
	// For example, [This ... text].
	//
	// Note that not all platforms support the ellipsis in the middle. For example, on Android
	// the middle ellipsis is only available for a single line text (i.e. when either a soft
	// wrap is disabled or a maximum number of lines maxLines set to 1). In case of multiline
	// text it will fallback to Clip.
	OverFlowMiddleEllipsis TextOverFlow = 5
)

func (t TextOverFlow) String() string {
	switch t {
	case OverFlowClip:
		return "Clip"
	case OverFlowEllipsis:
		return "Ellipsis"
	case OverFlowMiddleEllipsis:
		return "MiddleEllipsis"
	case OverFlowVisible:
		return "Visible"
	case OverFlowStartEllipsis:
		return "StartEllipsis"
	default:
		return "Invalid"
	}
}
