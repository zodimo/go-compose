package fform

import (
	"sync"

	"github.com/zodimo/go-compose/state"
)

type FormFieldState[T any] struct {
	id      string
	value   state.MutableValueTyped[T]
	touched bool
	mu      sync.Mutex

	disabled bool
	err      error
}

func NewFormFieldState[T any](id string, value state.MutableValueTyped[T]) *FormFieldState[T] {
	return &FormFieldState[T]{
		id:       id,
		value:    value,
		touched:  false,
		disabled: false,
	}
}

func (fs *FormFieldState[T]) ID() string {
	return fs.id
}

func (fs *FormFieldState[T]) Value() T {
	return fs.value.Get()
}

func (fs *FormFieldState[T]) SetValue(v T) *FormFieldState[T] {
	fs.value.Set(v)
	return fs
}

func (fs *FormFieldState[T]) Error() error {
	return fs.err
}

func (fs *FormFieldState[T]) SetError(err error) *FormFieldState[T] {
	fs.err = err
	return fs
}

func (fs *FormFieldState[T]) HasError() bool {
	return fs.err != nil
}

func (fs *FormFieldState[T]) ClearError() *FormFieldState[T] {
	fs.err = nil
	return fs
}

func (fs *FormFieldState[T]) Touched() bool {
	return fs.touched
}

func (fs *FormFieldState[T]) Touch() *FormFieldState[T] {
	fs.touched = true
	return fs
}

func (fs *FormFieldState[T]) Disabled() bool {
	return fs.disabled
}

func (fs *FormFieldState[T]) Disable() *FormFieldState[T] {
	fs.disabled = true
	return fs
}

func (fs *FormFieldState[T]) Enable() *FormFieldState[T] {
	fs.disabled = false
	return fs
}

func (fs *FormFieldState[T]) Update(f func(s *FormFieldState[T]) *FormFieldState[T]) *FormFieldState[T] {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	c := fs.Clone()
	fs = f(c)
	return fs
}

func (fs *FormFieldState[T]) Clone() *FormFieldState[T] {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	return &FormFieldState[T]{
		id:       fs.id,
		value:    fs.value,
		touched:  fs.touched,
		disabled: fs.disabled,
		err:      fs.err,
	}
}
