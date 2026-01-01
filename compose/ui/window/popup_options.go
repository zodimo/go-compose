package window

import (
	"github.com/zodimo/go-compose/compose/ui/unit"
)

type PopupOptions struct {
	Alignment PopupAlignment
	OffsetX   unit.Dp
	OffsetY   unit.Dp
}

func DefaultPopupOptions() PopupOptions {
	return PopupOptions{
		Alignment: AlignTopLeft,
		OffsetX:   0,
		OffsetY:   0,
	}
}

type PopupOption func(*PopupOptions)

func WithAlignment(alignment PopupAlignment) PopupOption {
	return func(opts *PopupOptions) {
		opts.Alignment = alignment
	}
}

func WithOffset(x, y unit.Dp) PopupOption {
	return func(opts *PopupOptions) {
		opts.OffsetX = x
		opts.OffsetY = y
	}
}
