package floatutils

type Float interface {
	~float32 | ~float64
}

func TakeOrElse[T Float](v T, defaultValue T) T {
	if !IsSpecified(v) {
		return defaultValue
	}
	return v
}
