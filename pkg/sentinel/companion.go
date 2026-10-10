package sentinel

type Companion[T any] interface {
	Unspecified() T
	String(T) string
	IsSpecified(T) bool
	TakeOrElse(src, def T) T
	Merge(a, b T) T
	Same(a, b T) bool
	SemanticEqual(a, b T) bool
	Equal(a, b T) bool
}

type ScalarCompanion[T any] interface {
	Companion[T]
	Copy(s T) T
}

type StructCompoanion[T any] interface {
	Companion[T]
	Coalesce(ptr, def T) T
	Copy(s T, options ...func(T)) T
}
