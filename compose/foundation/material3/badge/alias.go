package badge

import (
	"go-compose-dev/compose/foundation/material3/text"
	"go-compose-dev/internal/modifier"
	"go-compose-dev/pkg/api"
)

type Modifier = modifier.Modifier

var EmptyModifier = modifier.EmptyModifier

type Composable = api.Composable
type Composer = api.Composer

// Text aliases for convenience
var Text = text.Text
var TypestyleLabelSmall = text.TypestyleLabelSmall
