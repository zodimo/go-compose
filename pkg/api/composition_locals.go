package api

// ProvidedValue pairs a CompositionLocal with a value.
type ProvidedValue struct {
	CompositionLocal interface{}
	Value            interface{}
}
