package size

const NotSet = -1

type SizeData struct {
	Width  int
	Height int

	FillMaxWidth  bool
	FillMaxHeight bool
	FillMax       bool

	// WrapContent options
	WrapWidth  bool
	WrapHeight bool
	Alignment  Alignment
	Unbounded  bool

	Required bool
}
