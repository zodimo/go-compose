package text

import (
	"github.com/zodimo/go-compose/compose/ui/next/text"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/theme"
)

type Modifier = modifier.Modifier

var EmptyModifier = modifier.EmptyModifier

type Composable = api.Composable
type Composer = api.Composer

type ThemeManager = theme.ThemeManager

var GetThemeManager = theme.GetThemeManager

type TextStyle = text.TextStyle

type uiTextLayoutResult = text.TextLayoutResult
type uiPlaceholder = text.Placeholder
type uiTextUnit = unit.TextUnit

var uiTextUnitSp = unit.Sp

type uiAnnotatedString = text.AnnotatedString

type uiConstraints = unit.Constraints
