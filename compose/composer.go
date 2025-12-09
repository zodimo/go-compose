package compose

import (
	"go-compose-dev/internal/composer/zipper"
	"go-compose-dev/internal/state"
	"go-compose-dev/pkg/api"
)

func NewComposer(store state.PersistentState) api.Composer {
	return zipper.NewComposer(store)
}
