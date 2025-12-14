package card

import (
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/widget"
	"git.sr.ht/~schnwalter/gio-mw/widget/card"
)

type Modifier = modifier.Modifier

var EmptyModifier = modifier.EmptyModifier

type Composable = api.Composable
type Composer = api.Composer

type ThemeManager = theme.ThemeManager

type MutableValue = state.MutableValue

type cardKind = card.Kind

const (
	cardElevated cardKind = card.Elevated
	cardFilled   cardKind = card.Filled
	cardOutlined cardKind = card.Outlined
)

type m3CardChild = card.Child

var m3CardImage = card.Image
var m3CardContent = card.Content
var m3CardContentCover = card.ContentCover

type GioImage = widget.Image
