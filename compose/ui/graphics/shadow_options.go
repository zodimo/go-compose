package graphics

import (
	"github.com/zodimo/go-compose/compose/ui/geometry"
)

type ShadowOption func(*Shadow)

func WithColor(color Color) ShadowOption {
	return func(s *Shadow) {
		s.Color = color
	}
}

func WithOffset(offset geometry.Offset) ShadowOption {
	return func(s *Shadow) {
		s.Offset = offset
	}
}

func WithBlurRadius(blurRadius float32) ShadowOption {
	return func(s *Shadow) {
		s.BlurRadius = blurRadius
	}
}
