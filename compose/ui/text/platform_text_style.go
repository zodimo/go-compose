package text

// PlatformSpanStyle contains platform-specific text styling.
type PlatformSpanStyle struct {
}

func (s *PlatformSpanStyle) Merge(other *PlatformSpanStyle) *PlatformSpanStyle {
	if s == nil {
		return other
	}
	if other == nil {
		return s
	}
	// No fields to merge yet
	return s
}

func (s *PlatformSpanStyle) Equals(other *PlatformSpanStyle) bool {
	if s == other {
		return true
	}
	if s == nil || other == nil {
		return false
	}
	return true
}

// PlatformParagraphStyle contains platform-specific paragraph styling.
type PlatformParagraphStyle struct {
}

// PlatformTextStyle contains platform-specific text styling combining span and paragraph styles.
type PlatformTextStyle struct {
	SpanStyle      *PlatformSpanStyle
	ParagraphStyle *PlatformParagraphStyle
}

func (s *PlatformTextStyle) Merge(other *PlatformTextStyle) *PlatformTextStyle {
	if s == nil {
		return other
	}
	if other == nil {
		return s
	}
	return &PlatformTextStyle{
		SpanStyle:      s.SpanStyle.Merge(other.SpanStyle),
		ParagraphStyle: s.ParagraphStyle.Merge(other.ParagraphStyle),
	}
}

func (s *PlatformTextStyle) Equals(other *PlatformTextStyle) bool {
	if s == other {
		return true
	}
	if s == nil || other == nil {
		return false
	}
	if !s.SpanStyle.Equals(other.SpanStyle) {
		return false
	}
	if !s.ParagraphStyle.Equals(other.ParagraphStyle) {
		return false
	}
	return true
}

func (s *PlatformParagraphStyle) Merge(other *PlatformParagraphStyle) *PlatformParagraphStyle {
	if s == nil {
		return other
	}
	if other == nil {
		return s
	}
	// No fields to merge yet
	return s
}

func (s *PlatformParagraphStyle) Equals(other *PlatformParagraphStyle) bool {
	if s == other {
		return true
	}
	if s == nil || other == nil {
		return false
	}
	return true
}

func lerpPlatformSpanStyle(start, stop *PlatformSpanStyle, fraction float32) *PlatformSpanStyle {
	if start == nil && stop == nil {
		return nil
	}
	if start == nil {
		return stop
	}
	if stop == nil {
		return start
	}
	// Discrete interpolation
	if fraction < 0.5 {
		return start
	}
	return stop
}

func lerpPlatformParagraphStyle(start, stop *PlatformParagraphStyle, fraction float32) *PlatformParagraphStyle {
	if start == nil && stop == nil {
		return nil
	}
	if start == nil {
		return stop
	}
	if stop == nil {
		return start
	}
	// Discrete interpolation
	if fraction < 0.5 {
		return start
	}
	return stop
}

// LerpDiscrete interpolates between two values discretely.
func LerpDiscrete[T any](start, stop T, fraction float32) T {
	if fraction < 0.5 {
		return start
	}
	return stop
}
