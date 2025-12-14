package icon

import (
	"image/color"

	"github.com/zodimo/go-compose/internal/color/colorhelper"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/theme"

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

type ColorDescriptor = theme.ThemeColorDescriptor

var themeManager = theme.GetThemeManager()

var specificColor = theme.SpecificColor
