package lerp

func LerpDiscrete[T any, F Float](a, b T, fraction F) T {
	if fraction < 0.5 {
		return a
	}
	return b
}
