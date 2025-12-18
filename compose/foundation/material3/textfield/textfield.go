package textfield

const Material3TextFieldNodeID = "Material3TextField"

type HandlerWrapper struct {
	Func func(string)
}

type OnSubmitWrapper struct {
	Func func()
}

type TextFieldStateTracker struct {
	LastValue string
}

// TextField implements a Material Design 3 text field.
// It defaults to the Filled variant for backward compatibility,
// but users should prefer explicit Filled or Outlined calls.
func TextField(
	value string,
	onValueChange func(string),
	label string,
	options ...TextFieldOption,
) Composable {
	return Filled(value, onValueChange, label, options...)
}
