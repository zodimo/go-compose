package divider

import (
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/internal/state"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/theme"

	"git.sr.ht/~schnwalter/gio-mw/widget/divider"
)

type Modifier = modifier.Modifier

var EmptyModifier = modifier.EmptyModifier

type Composable = api.Composable
type Composer = api.Composer

type ThemeManager = theme.ThemeManager

type MutableValue = state.MutableValue

var M3Divider = divider.Divider
