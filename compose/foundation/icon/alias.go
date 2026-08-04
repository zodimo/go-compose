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
