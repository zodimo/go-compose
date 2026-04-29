package pointer

// https://vscode.dev/github/JetBrains/compose-multiplatform-core/blob/jb-main/compose/ui/ui/src/commonMain/kotlin/androidx/compose/ui/input/pointer/PointerEvent.kt#L422
type PointerInputChange struct{}

type InternalPointerEvent struct{}

type PointerEventType int

type PointerButtons struct{}

type PointerKeyboardModifiers struct{}

type PointerEvent struct {
	changes []PointerInputChange

	/** The state of buttons (e.g. mouse or stylus buttons) during this event. */
	Buttons PointerButtons

	/** The state of modifier keys during this event. */
	KeyboardModifiers PointerKeyboardModifiers

	/** The primary reason the [PointerEvent] was sent. */
	Type PointerEventType
}
