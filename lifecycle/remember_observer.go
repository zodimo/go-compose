package lifecycle

// RememberObserver is an interface for objects that need lifecycle callbacks
// when remembered in composition. This is the Go equivalent of
// androidx.compose.runtime.RememberObserver.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/runtime/runtime/src/commonMain/kotlin/androidx/compose/runtime/RememberObserver.kt
type RememberObserver interface {
	// OnForgotten is called when the object is no longer being remembered.
	// This happens when the remember call is removed from composition.
	OnForgotten()
}
