package button

import (
	"go-compose-dev/internal/modifier"
	"go-compose-dev/internal/state"
	"go-compose-dev/internal/theme"
	"go-compose-dev/pkg/api"
)

type Modifier = modifier.Modifier

var EmptyModifier = modifier.EmptyModifier

type Composable = api.Composable
type Composer = api.Composer

type ThemeManager = theme.ThemeManager

type MutableValue = state.MutableValue
