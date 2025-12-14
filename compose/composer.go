package compose

import (
	"github.com/zodimo/go-compose/internal/composer/zipper"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
)

type Composable = api.Composable
type Composer = api.Composer

func NewComposer(store state.PersistentState) Composer {
	return zipper.NewComposer(store)
}

type composableSequence struct {
	contents []Composable
}

func NewComposableSequence(composables []Composable) *composableSequence {
	return &composableSequence{contents: composables}
}

func (s *composableSequence) Compose(c Composer) Composer {
	for _, content := range s.contents {
		c = content(c)
	}
	return c
}

func Sequence(contents ...Composable) Composable {
	return func(c Composer) Composer {
		return NewComposableSequence(contents).Compose(c)
	}
}

func Id() Composable {
	return func(c Composer) Composer {
		return c
	}
}
