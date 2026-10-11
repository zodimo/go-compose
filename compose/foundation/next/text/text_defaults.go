package text

import (
	"math"

	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/next/text/style"
)

func DefaultTextOptions() TextOptions {

	return TextOptions{
		Modifier:      ui.EmptyModifier,
		TextStyle:     nil,
		OnTextLayout:  nil,
		OverFlow:      style.OverFlowClip,
		SoftWrap:      true,
		MaxLines:      math.MaxInt32,
		MinLines:      1,
		InlineContent: map[string]InlineTextContent{},
		Color:         graphics.ColorCompanion.Unspecified(),
		AutoSize:      nil,
	}
}
