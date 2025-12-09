package zipper

import (
	idModels "go-compose-dev/pkg/compose-identifier/models"
)

func NewComposer(state PersistentState) Composer {

	idManager := idModels.GeOrCreateIdentityManager("composer")
	idManager.ResetKeyCounter()

	return &composer{
		focus:     nil,
		path:      []pathItem{},
		memo:      EmptyMemo,
		state:     state,
		idManager: idManager,
	}
}
