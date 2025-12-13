package icon

import (
	"go-compose-dev/internal/color/colorhelper"
	"go-compose-dev/internal/modifier"
	"go-compose-dev/internal/theme"
	"go-compose-dev/pkg/api"
	"image/color"

	"gioui.org/layout"
	// . "golang.org/x/exp/shiny/materialdesign/icons"
)

type Modifier = modifier.Modifier

var EmptyModifier = modifier.EmptyModifier

type Composable = api.Composable
type Composer = api.Composer

type layoutContext = layout.Context
type layoutDimensions = layout.Dimensions
type IconWidget = func(gtx layoutContext, foreground color.Color) layoutDimensions

var ToNRGBA = colorhelper.ToNRGBA

type ThemeColorSet = theme.ThemeColorSetter
type ColorReader = theme.ColorReader
