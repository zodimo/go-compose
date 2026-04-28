package path

import (
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/internal/skiagraphics"
)

func New() graphics.Path {
	return skiagraphics.NewPath()
}
