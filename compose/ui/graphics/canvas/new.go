package canvas

import (
	"gioui.org/layout"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/internal/skiagraphics"
)

func New(gtx layout.Context) graphics.Canvas {
	return skiagraphics.NewCanvas(gtx)
}
