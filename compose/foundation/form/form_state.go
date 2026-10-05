package fform

import (
	"fmt"

	formengine "github.com/zodimo/go-compose/compose/foundation/form/internal"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
)

// FormState holds the remembered root of a form tree. The root *formengine.Group
// backs the whole form and persists across recompositions, and the exported
// methods below form the form-level lifecycle surface: validate the whole tree,
// read the serialized value, reset, and bulk lifecycle mutations.
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

// Root returns the tree root node. It is exposed for advanced wiring (custom
// traversal, path resolution) that the FormScope helpers do not cover.
func (s *FormState) Root() *Group {
	return s.root
}

// Status returns the aggregated status of the whole form: INVALID if any
// descendant is invalid, DISABLED if the entire form is disabled, PENDING while
// asynchronous work is outstanding, otherwise VALID.
func (s *FormState) Status() Status {
	return s.root.Status()
}

// Validate validates every node in the tree and reports whether the form is
// valid. A disabled control always validates true.
func (s *FormState) Validate() bool {
	return s.root.Validate()
}

// Errors returns every validation error in the tree keyed by dotted path, or nil
// when the form is valid or disabled.
func (s *FormState) Errors() map[string]string {
	return s.root.Errors()
}

// Value returns the form serialized as plain Go values: a map[string]any keyed
// by the root group's child names, with nested groups as maps, arrays as ordered
// slices, and effectively-disabled nodes omitted.
func (s *FormState) Value() map[string]any {
	raw, _ := s.root.RawValue().(map[string]any)
	if raw == nil {
		return map[string]any{}
	}
	return raw
}

// RawValue returns the form serialized as plain Go values without the
// map[string]any narrowing of Value. Callers dealing with the root group can
// use either.
func (s *FormState) RawValue() any {
	return s.root.RawValue()
}

// IsTouched reports whether any node in the tree has been touched.
func (s *FormState) IsTouched() bool {
	return s.root.IsTouched()
}

// IsDirty reports whether any node in the tree differs from its initial value.
func (s *FormState) IsDirty() bool {
	return s.root.IsDirty()
}

// IsPristine is the negation of IsDirty.
func (s *FormState) IsPristine() bool {
	return s.root.IsPristine()
}

// IsEnabled reports whether the form root is enabled.
func (s *FormState) IsEnabled() bool {
	return s.root.IsEnabled()
}

// MarkAllTouched marks every control in the tree touched, so touched-gated error
// display surfaces validation failures the user has not yet visited. This is the
// conventional "attempted submit" step before reading Validate().
func (s *FormState) MarkAllTouched() {
	s.root.MarkAsTouched()
}

// MarkAllUntouched clears the touched flag across the tree.
func (s *FormState) MarkAllUntouched() {
	s.root.MarkAsUntouched()
}

// MarkAllPristine clears the dirty flag across the tree without changing values.
func (s *FormState) MarkAllPristine() {
	s.root.MarkAsPristine()
}

// Reset restores every control to its initial value, clears errors, and marks
// the whole tree pristine and untouched.
func (s *FormState) Reset() {
	s.root.Reset()
}

// SetDisabled disables or enables the entire form. Children keep their own
// disabled settings, so re-enabling restores them.
func (s *FormState) SetDisabled(disabled bool) {
	s.root.SetDisabled(disabled)
}

// Submit validates the form, marks every control touched so invalid fields
// render their errors, and invokes onValid with the serialized value only when
// the form is valid. It reports whether submission proceeded.
//
// Submit is the conventional button handler: it encodes the
// validate -> touch-all -> proceed-only-if-valid flow so callers do not have to
// order those steps themselves.
func (s *FormState) Submit(onValid func(value map[string]any)) bool {
	valid := s.Validate()
	s.MarkAllTouched()
	if !valid {
		return false
	}
	if onValid != nil {
		onValid(s.Value())
	}
	return true
}
