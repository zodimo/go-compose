package Sentinel

// Do not implement this interface on complex objects
// value classes inline to primitives, and Unspecified is just the zero bit pattern.
type Sentinel interface {
	IsSpecified() bool
	neverImplement()
}

// Generic TakeOrElse for any type with IsSpecified() method
// golang interface or value receiver
// DO NOT PASS NIL as arguments
func TakeOrElseValue[T interface{ IsSpecified() bool }](a, b T) T {
	if a.IsSpecified() {
		return a
	}
	return b
}

// Generic TakeOrElse for any type with IsSpecified() method
// golang interface or pointer receiver
func TakeOrElsePointer[T any](a, b *T) *T {
	if a == nil {
		return b
	}
	return a
}
