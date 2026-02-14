package store

import (
	"reflect"

	"github.com/zodimo/go-compose/state"
)

type PersistentStateInterface = state.PersistentState

type PersistentState struct {
	scopes      map[string]state.MutableValue
	subscribers *state.SubscriptionManager
}

func NewPersistentState(scopes map[string]state.MutableValue) PersistentStateInterface {
	return &PersistentState{
		scopes:      scopes,
		subscribers: state.NewSubscriptionManager(),
	}
}

// Deprecated: use Subscribe instead
func (ps *PersistentState) SetOnStateChange(callback func()) {
	ps.subscribers.Subscribe(callback)
}

func (ps *PersistentState) Subscribe(callback func()) state.Subscription {
	return ps.subscribers.Subscribe(callback)
}

func (ps *PersistentState) GetState(id string, initial func() any, options ...state.StateOption) state.MutableValue {

	opts := state.StateOptions{
		Compare: reflect.DeepEqual,
	}
	for _, option := range options {
		option(&opts)
	}

	if v, ok := ps.scopes[id]; ok {
		return v
	}

	ps.scopes[id] = state.NewMutableValue(initial(), func(any) {
		ps.subscribers.NotifyAll()
	}, opts.Compare)
	return ps.scopes[id]
}
