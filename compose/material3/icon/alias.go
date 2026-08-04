package icon

import (
	"image/color"

	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/pkg/api"
)

type Composable = api.Composable
type Composer = api.Composer

type layoutContext = layoutnode.LayoutContext
type layoutDimensions = layoutnode.LayoutDimensions
type IconWidget = func(gtx layoutContext, foreground color.NRGBA) layoutDimensions

// IconSource represents a source for icon data.
// It can be either IconBytes (raw SVG/icon data) or SymbolName (font-based icon).
type IconSource interface {
	isIconSource()
}

// IconBytes wraps a []byte for use as an IconSource.
// Use this with icons from golang.org/x/exp/shiny/materialdesign/icons.
type IconBytes []byte

func (IconBytes) isIconSource() {}
