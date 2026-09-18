package fform

import (
	"fmt"

	formengine "github.com/zodimo/go-compose/compose/foundation/form/internal"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
)

// FormState holds the remembered root of a form tree. The root *formengine.Group
// backs the whole form and persists across recompositions.
type FormState struct {
	root *formengine.Group
}

// RememberFormState builds and remembers the form tree root.
//
// build receives the composer so callers can create remembered per-field stores
// via state.MustRemember and wire cross-field logic at build time. The returned
// FormState is remembered across compositions.
func RememberFormState(c api.Composer, build func(c api.Composer) *formengine.Group) *FormState {
	key := fmt.Sprintf("%d/%s/formRoot", c.GenerateID(), c.GetPath())

	formState := state.MustRemember[*FormState](c, key, func() *FormState {
		return &FormState{root: build(c)}
	})

	return formState.Get()
}
