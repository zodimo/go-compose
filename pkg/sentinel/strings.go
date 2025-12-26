package sentinel

var StringUnspecified = "\x00unspecified"

func TakeOrElseString(a, b string) string {
	if a != StringUnspecified {
		return a
	}
	return b
}
