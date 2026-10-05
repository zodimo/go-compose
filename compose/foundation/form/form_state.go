package fform

import (
	"errors"
	"fmt"
	"sort"

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

// CodedErrors returns every coded validation failure in the tree as a slice of
// *ValidationError, each stamped with the failing control's dotted path. Errors
// without a code are still included with an empty Code. It returns nil when the
// form is valid or disabled.
//
// Like Errors(), it reflects the last validation run: call Validate() first (or
// Submit, which validates) so a fresh tree's validators have executed.
//
// This is the structured counterpart to Errors(): use it to branch on ErrorCode
// (localization, focusing the first invalid field, mapping to an API error body)
// instead of parsing message strings.
func (s *FormState) CodedErrors() []*ValidationError {
	var out []*ValidationError

	Walk(s.root, func(node FormNode) {
		control, ok := node.(interface{ ValidationErrors() []error })
		if !ok {
			return
		}
		for _, err := range control.ValidationErrors() {
			var ve *ValidationError
			if errors.As(err, &ve) {
				out = append(out, ve.WithPath(node.Path()))
				continue
			}
			// A non-coded validator failure: surface it with its message only.
			out = append(out, &ValidationError{Message: err.Error(), Path: node.Path()})
		}
	})

	// Deterministic order: by path, then code, then message.
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Path != out[j].Path {
			return out[i].Path < out[j].Path
		}
		if out[i].Code != out[j].Code {
			return out[i].Code < out[j].Code
		}
		return out[i].Message < out[j].Message
	})

	return out
}

// Codes returns the distinct ErrorCodes raised anywhere in the form, sorted.
func (s *FormState) Codes() []ErrorCode {
	coded := s.CodedErrors()
	seen := make(map[ErrorCode]struct{}, len(coded))
	for _, ve := range coded {
		if ve.Code == "" {
			continue
		}
		seen[ve.Code] = struct{}{}
	}
	codes := make([]ErrorCode, 0, len(seen))
	for code := range seen {
		codes = append(codes, code)
	}
	sort.Slice(codes, func(i, j int) bool { return codes[i] < codes[j] })
	return codes
}

// FirstInvalidPath returns the first failing control's dotted path, or ("",
// false) when the form is valid. It is the usual hook for "focus the first
// error".
//
// Group children are stored in a map, so there is no stable insertion order to
// walk; "first" therefore means the lexically smallest dotted path among the
// failing controls, which is deterministic. Arrays keep positional order within
// a path (for example phones[0].number sorts before phones[1].number).
func (s *FormState) FirstInvalidPath() (string, bool) {
	coded := s.CodedErrors()
	if len(coded) == 0 {
		return "", false
	}
	// CodedErrors is sorted by path, so the first entry is the lexically
	// smallest failing path.
	return coded[0].Path, true
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
