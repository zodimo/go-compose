package badge

import (
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/pkg/api"
)

type Modifier = modifier.Modifier

var EmptyModifier = modifier.EmptyModifier

type Composable = api.Composable
type Composer = api.Composer

// Text aliases for convenience
var badgeText = text.LabelSmall
