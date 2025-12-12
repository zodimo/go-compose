package store

import (
	"fmt"
	"go-compose-dev/internal/state"
	"reflect"
	"sync"
)

type MutableValueInterface = state.MutableValue
type PersistentStateInterface = state.PersistentState

var _ MutableValueInterface = &MutableValue{}
var _ MutableValueInterface = &MutableValueTypedWrapper[any]{}

type TypedMutableValueInterface[T any] interface {
	Get() T
	Set(value T)

	Unwrap() MutableValueInterface
}

var _ TypedMutableValueInterface[any] = &MutableValueTypedWrapper[any]{}

// MutableValue is a trivial state container.
type MutableValue struct {
	cell           any
	changeNotifier func(any)
	mutex          sync.Mutex
	compare        func(any, any) bool
}

func NewMutableValue(initial any, changeNotifier func(any), compare func(any, any) bool) *MutableValue {
	return &MutableValue{
		cell:           initial,
		changeNotifier: changeNotifier,
		compare:        compare,
	}
}

func (mv *MutableValue) Get() any {
	return mv.cell
}

func (mv *MutableValue) Set(value any) {
	changed := !mv.compare(mv.cell, value)
	if changed {
		mv.mutex.Lock()
		defer mv.mutex.Unlock()
		mv.cell = value
		if mv.changeNotifier != nil {
			mv.changeNotifier(value)
		}
	}
}

type PersistentState struct {
	scopes        map[string]MutableValueInterface
	onStateChange func()
}

func NewPersistentState(scopes map[string]MutableValueInterface) PersistentStateInterface {
	return &PersistentState{scopes: scopes}
}

func (ps *PersistentState) SetOnStateChange(callback func()) {
	ps.onStateChange = callback
}

func (ps *PersistentState) GetState(id string, initial func() any) MutableValueInterface {
	if v, ok := ps.scopes[id]; ok {
		return v
	}
	ps.scopes[id] = &MutableValue{
		cell: initial(),
		changeNotifier: func(any) {
			if ps.onStateChange != nil {
				ps.onStateChange()
			}
		},
		compare: reflect.DeepEqual,
	}
	return ps.scopes[id]
}

type MutableValueTyped[T any] struct {
	cell           T
	changeNotifier func(T)
	compare        func(T, T) bool
}

func NewMutableValueTyped[T any](initial T, changeNotifier func(T), compare func(T, T) bool) *MutableValueTyped[T] {
	return &MutableValueTyped[T]{
		cell:           initial,
		changeNotifier: changeNotifier,
		compare:        compare,
	}
}

type MutableValueTypedWrapper[T any] struct {
	mv *MutableValue
}

func (w *MutableValueTypedWrapper[T]) Get() T {
	return w.mv.cell.(T)
}

func (w *MutableValueTypedWrapper[T]) Set(value T) {
	w.mv.Set(value)
}

func WrapMutableValue[T any](mv *MutableValue) (TypedMutableValueInterface[T], error) {

	_, ok := mv.cell.(T)
	if !ok {
		var zero T
		return nil, fmt.Errorf("cell is not of type %T", zero)
	}

	return &MutableValueTypedWrapper[T]{
		mv: mv,
	}, nil
}

func (w *MutableValueTypedWrapper[T]) Unwrap() MutableValueInterface {
	return w.mv
}
