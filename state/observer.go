package state

import (
	"github.com/zodimo/go-compose/state/core"
	"github.com/zodimo/go-compose/state/internal"
)

type ReadObserver = core.ReadObserver

// NotifyRead notifies the active read observer that a state object has been read.
// Source is typically the state object itself (e.g. *MutableValue).
func NotifyRead(source StateChangeNotifier) {
	internal.NotifyRead(source)
}

// WithReadObserver executes the block with the given read observer.
// It restores the previous observer after the block finishes.
func WithReadObserver(observer ReadObserver, block func()) {
	internal.WithReadObserver(observer, block)
}

type ObserverManager = core.ObserverManager
