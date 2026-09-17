package fform

import (
	"slices"
	"sync"

	"github.com/zodimo/go-maybe"
)

type FormFieldState struct {
	id      string
	value   maybe.Maybe[string]
	touched bool
	rules   []ValidationRule
	mu      sync.Mutex

	disabled bool
}

func NewFormFieldState(id string) *FormFieldState {
	return &FormFieldState{
		id:       id,
		value:    maybe.None[string](),
		touched:  false,
		rules:    []ValidationRule{},
		disabled: false,
	}
}

func (fs *FormFieldState) ID() string {
	return fs.id
}

func (fs *FormFieldState) Value() maybe.Maybe[string] {
	return fs.value
}

func (fs *FormFieldState) ValidationRules() []ValidationRule {
	return slices.Clone(fs.rules)
}

func (fs *FormFieldState) Touched() bool {
	return fs.touched
}

func (fs *FormFieldState) Clone() *FormFieldState {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	return &FormFieldState{
		id:      fs.id,
		value:   fs.value,
		touched: fs.touched,
		rules:   slices.Clone(fs.rules),
	}
}

func (fs *FormFieldState) Touch() *FormFieldState {
	c := fs.Clone()
	c.touched = true
	return c
}

func (fs *FormFieldState) Disabled() bool {
	return fs.disabled
}

func (fs *FormFieldState) Disable() *FormFieldState {
	c := fs.Clone()
	c.disabled = true
	return c
}

func (fs *FormFieldState) Enable() *FormFieldState {
	c := fs.Clone()
	c.disabled = false
	return c
}

func (fs *FormFieldState) Update(f func(s *FormFieldState) *FormFieldState) *FormFieldState {
	c := fs.Clone()
	return f(c)
}
