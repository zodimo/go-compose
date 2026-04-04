package lifecycle

import "context"

type CoroutineLauncherOptions struct {
	RecoverFromPanic func(panicMessage any)
	PanicHandler     func(panicMessage any, recoverFromPanic func(panicMessage any))
}
type CoroutineLauncherOption func(*CoroutineLauncherOptions)

func WithRecoverFromPanic(recoverFromPanic func(panicMessage any)) CoroutineLauncherOption {
	return func(options *CoroutineLauncherOptions) {
		if recoverFromPanic != nil {
			options.RecoverFromPanic = recoverFromPanic
		}
	}
}

func WithPanicHandler(panicHandler func(panicMessage any, recoverFromPanic func(panicMessage any))) CoroutineLauncherOption {
	return func(options *CoroutineLauncherOptions) {
		if panicHandler != nil {
			options.PanicHandler = panicHandler
		}
	}
}

func CoroutineLauncherOptionsDefaults() *CoroutineLauncherOptions {
	return &CoroutineLauncherOptions{
		RecoverFromPanic: func(panicMessage any) {
			//default dont silence the error
			panic(panicMessage)
		},
		PanicHandler: func(panicMessage any, recoverFromPanic func(panicMessage any)) {
			if recoverFromPanic != nil {
				recoverFromPanic(panicMessage)
			}
		},
	}
}

// HasCoroutineScope is an optional interface for ViewModels that need a managed context/scope.
// When a Remembered value implements this interface, it will receive a context that is
// cancelled when the leaving the composition.
type HasCoroutineScope interface {
	// SetCoroutineScope provides a context that is cancelled when the remembered value is leaving the composition.
	// This is analogous to viewModelScope in Kotlin, which is a CoroutineScope.
	SetCoroutineScope(ctx context.Context)
}

type CoroutineScope interface {
	HasCoroutineScope
	CoroutineLauncher
	Context() (context.Context, bool)
	MustContext() context.Context
}

type CoroutineLauncher interface {
	Launch(block func(ctx context.Context), options ...CoroutineLauncherOption)
}
