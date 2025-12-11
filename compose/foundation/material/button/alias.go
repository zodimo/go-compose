package button

import (
	"go-compose-dev/internal/modifier"
	"go-compose-dev/internal/theme"
	"go-compose-dev/pkg/api"

	"gioui.org/widget"
)

type Modifier = modifier.Modifier

var EmptyModifier = modifier.EmptyModifier

type Composable = api.Composable
type Composer = api.Composer

type ThemeManager = theme.ThemeManager

var GetThemeManager = theme.GetThemeManager

type GioClickable = widget.Clickable
