package theme

func ColorLerp(start, stop ColorDescriptor, fraction float32) ColorDescriptor {
	return start.AppendUpdate(LerpColor(stop, fraction))
}
