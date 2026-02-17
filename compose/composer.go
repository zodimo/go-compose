package compose

import (
	"github.com/zodimo/go-compose/internal/composer/zipper"
	"github.com/zodimo/go-compose/internal/sequence"
	"github.com/zodimo/go-compose/pkg/api"
)

type Composable = api.Composable
type Composer = api.Composer

func NewComposer(options ...api.ComposerOption) Composer {
	return zipper.NewComposer(options...)
}

// Use This Sequence When not inside of a composable but composing composables
var Sequence = sequence.Sequence

// Identity Composable
func Id() Composable {
	return func(c Composer) Composer {
		return c
	}
}
