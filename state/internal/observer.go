package internal

import (
	"fmt"
	"sync"

	"github.com/zodimo/go-compose/state/core"
)

// singleton
var SingletonObserverManager core.ObserverManager

func init() {
	SingletonObserverManager = &ObserverManager{}
}

type ObserverManager struct {
	readObservers []core.ReadObserver
	mu            sync.RWMutex
	isLocked      bool
}

func (m *ObserverManager) lock() {
	if m.isLocked {
		fmt.Println("observerManager: is locked")
	}
	m.mu.Lock()
	m.isLocked = true
}
func (m *ObserverManager) unlock() {
	if !m.isLocked {
		fmt.Println("observerManager: not locked")
		return
	}
	m.mu.Unlock()
	m.isLocked = false
}

func (m *ObserverManager) pushObserver(observer core.ReadObserver) {
	m.lock()
	defer m.unlock()
	m.readObservers = append(m.readObservers, observer)
}

func (m *ObserverManager) popObserver() core.ReadObserver {
	if len(m.readObservers) == 0 {
		panic("observerManager: no observers")
	}
	observer := m.readObservers[len(m.readObservers)-1]
	m.readObservers = m.readObservers[:len(m.readObservers)-1]
	return observer
}

func (m *ObserverManager) WithReadObserver(observer core.ReadObserver, block func()) {
	if observer == nil {
		panic("observerManager: observer cannot be nil")
	}
	if block == nil {
		panic("observerManager: block cannot be nil")
	}
	m.pushObserver(observer)
	defer m.popObserver()
	block()
}

func (m *ObserverManager) NotifyRead(source core.StateChangeNotifier) {
	m.mu.RLock()
	if len(m.readObservers) > 0 {
		observer := m.readObservers[len(m.readObservers)-1]
		m.mu.RUnlock()
		observer(source)
		return
	}
	m.mu.RUnlock()
}

// NotifyRead notifies the active read observer that a state object has been read.
// Source is typically the state object itself (e.g. *MutableValue).
func NotifyRead(source core.StateChangeNotifier) {
	SingletonObserverManager.NotifyRead(source)
}

// WithReadObserver executes the block with the given read observer.
// It restores the previous observer after the block finishes.
func WithReadObserver(observer core.ReadObserver, block func()) {
	SingletonObserverManager.WithReadObserver(observer, block)
}
