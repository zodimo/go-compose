package stateful

import (
	"context"
	"fmt"
)

// type Stateful[T any] = struct {
// 	State T
// 	Clone func(T) T
// }

type Stateful[State any] = struct {
	state State
	clone func(State) State
}

func NewStateful[T any](initialState T, clone func(T) T) Stateful[T] {
	return Stateful[T]{
		state: initialState,
		clone: clone,
	}
}

func Duplicate[T any](s Stateful[T]) Stateful[T] {
	return Stateful[T]{
		state: s.clone(s.state),
		clone: s.clone,
	}
}

type ActorOptions struct {
	inboxSize int
}

type ActorOption = func(o *ActorOptions)

func WithInboxSize(size int) ActorOption {
	return func(o *ActorOptions) {
		o.inboxSize = size
	}
}

type ActorRequest[State any, S ~Stateful[State]] interface {
	isActorRequest()
}

type ActorUpdateOptions struct {
	Wait bool
}

type ActorUpdateOption = func(o *ActorUpdateOptions)

func UpdateWithWait(wait bool) ActorUpdateOption {
	return func(o *ActorUpdateOptions) {
		o.Wait = wait
	}
}

type ActorUpdateRequest[State any, S ~Stateful[State]] struct {
	ActorRequest[State, S]
	options *ActorUpdateOptions
	update  func(s State) State
	done    chan struct{}
}

func NewActorUpdateRequest[State any, S ~Stateful[State]](update func(s State) State, options ...ActorUpdateOption) ActorUpdateRequest[State, S] {
	opts := ActorUpdateOptions{}

	for _, opt := range options {
		if opt != nil {
			opt(&opts)
		}
	}
	return ActorUpdateRequest[State, S]{
		options: &opts,
		update:  update,
		done:    make(chan struct{}),
	}
}

type ActorReadRequest[State any, S ~Stateful[State]] struct {
	ActorRequest[State, S]
	response chan State
}

type Actor[State any, S ~Stateful[State]] struct {
	ctx      context.Context
	stateful S
	inbox    chan ActorRequest[State, S]
}

func (a *Actor[State, S]) Update(ctx context.Context, request ActorUpdateRequest[State, S]) error {
	if err := a.ctx.Err(); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case a.inbox <- request:

	}
	if request.options.Wait {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-request.done:
			return nil
		}
	}
	return nil

}

// async
func (a *Actor[State, S]) Send(ctx context.Context, req ActorRequest[State, S]) error {
	if err := a.ctx.Err(); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case a.inbox <- req:
		return nil
	}
}

func (a *Actor[State, S]) Request(ctx context.Context) (state State, err error) {

	if err := a.ctx.Err(); err != nil {
		var zero State
		return zero, err
	}

	request := ActorReadRequest[State, S]{
		response: make(chan State, 1),
	}

	select {
	case <-ctx.Done():
		var zero State
		return zero, ctx.Err()
	case a.inbox <- request:
		select {
		case <-ctx.Done():
			var zero State
			return zero, ctx.Err()
		case result, ok := <-request.response:
			if !ok {
				var zero State
				return zero, fmt.Errorf("response channel closed")
			}
			return result, nil
		}
	}
}

func (a *Actor[State, S]) run() {
	for {
		select {
		case <-a.ctx.Done():
			return
		case in := <-a.inbox:
			switch req := in.(type) {
			case ActorReadRequest[State, S]:
				stateful := Stateful[State](a.stateful)
				cloned := stateful.clone(stateful.state)
				go func() {
					req.response <- cloned
				}()

			case ActorUpdateRequest[State, S]:
				stateful := Stateful[State](a.stateful)
				cloned := stateful.clone(stateful.state)
				stateful.state = req.update(cloned)
				a.stateful = stateful
				if req.options.Wait {
					go func() {
						req.done <- struct{}{}
					}()
				}

			}
		}
	}
}

func NewActor[State any, S ~Stateful[State]](ctx context.Context, initialState S, options ...ActorOption) *Actor[State, S] {
	opts := ActorOptions{
		inboxSize: 1,
	}
	for _, opt := range options {
		if opt != nil {
			opt(&opts)
		}
	}
	a := Actor[State, S]{
		ctx:      ctx,
		stateful: initialState,
		inbox:    make(chan ActorRequest[State, S], opts.inboxSize),
	}
	go func() {
		a.run()
	}()
	return &a
}
