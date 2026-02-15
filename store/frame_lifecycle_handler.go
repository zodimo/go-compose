package store

import (
	"context"
	"log"
	"sync"

	"github.com/zodimo/go-compose/lifecycle"
)

type FrameLifecycleHandler struct {
	mu           sync.Mutex
	scopedValues map[string]ScopedValue
	accessedKeys map[string]struct{}
	frameActive  bool
	debugMode    bool
}

func NewFrameLifecycleHandler(options ...FrameLifecycleHandlerOption) *FrameLifecycleHandler {
	flh := &FrameLifecycleHandler{
		scopedValues: make(map[string]ScopedValue),
		accessedKeys: make(map[string]struct{}),
		debugMode:    false,
	}

	startFrameTrigger := flh.StartFrame
	endFrameTrigger := flh.EndFrame

	opts := FrameLifecycleHandlerOptions{
		StartFrameTriggerReceiver: func(trigger func()) {},
		EndFrameTriggerReceiver:   func(trigger func() []string) {},
	}

	for _, option := range options {
		if option != nil {
			option(&opts)
		}
	}

	flh.debugMode = opts.DebugMode

	//provide triggers to receivers
	opts.StartFrameTriggerReceiver(startFrameTrigger)
	opts.EndFrameTriggerReceiver(endFrameTrigger)

	return flh
}

func (flh *FrameLifecycleHandler) StartFrame() {
	if flh.debugMode {
		log.Println("FrameLifecycleHandler: StartFrame")
	}
	flh.mu.Lock()
	defer flh.mu.Unlock()
	if flh.frameActive {
		panic("Startframe can only be called on an inactive frame")
	}
	flh.frameActive = true
	flh.accessedKeys = make(map[string]struct{})
}

func (flh *FrameLifecycleHandler) AttachValueToLifecycle(key string, value any) {
	if flh.debugMode {
		log.Printf("FrameLifecycleHandler: AttachValueToLifecycle %s\n", key)
	}
	flh.mu.Lock()
	defer flh.mu.Unlock()
	// Mark this key as accessed for GC tracking
	flh.accessedKeys[key] = struct{}{}
	flh.scopedValues[key] = initScopedValue(value)
}

func (flh *FrameLifecycleHandler) Ping(key string) {
	if flh.debugMode {
		log.Printf("FrameLifecycleHandler: Ping %s\n", key)
	}
	flh.mu.Lock()
	defer flh.mu.Unlock()
	flh.accessedKeys[key] = struct{}{}
}

// return keys to remove
func (flh *FrameLifecycleHandler) EndFrame() []string {
	if flh.debugMode {
		log.Println("FrameLifecycleHandler: EndFrame")
	}
	flh.mu.Lock()
	if !flh.frameActive {
		flh.mu.Unlock()
		panic("Endframe can only be called on an active frame")
	}
	flh.frameActive = false

	keysToRemove := []string{}

	for key := range flh.scopedValues {
		if _, accessed := flh.accessedKeys[key]; !accessed {
			keysToRemove = append(keysToRemove, key)
		}
	}

	toGC := make(map[string]ScopedValue, len(keysToRemove))

	for _, key := range keysToRemove {
		scopedValue := flh.scopedValues[key]

		delete(flh.scopedValues, key)
		toGC[key] = scopedValue
	}
	flh.mu.Unlock()

	// Call lifecycle callbacks outside the lock to avoid deadlocks
	for key, scopedValue := range toGC {
		if flh.debugMode {
			log.Printf("FrameLifecycleHandler: Cleanup %s\n", key)
		}
		scopedValue.Cleanup()
	}

	return keysToRemove
}

func initScopedValue(initialValue any) ScopedValue {

	scopedCancelFunc := func() {}
	onForgotten := func() {}
	onCleared := func() {}

	if observer, ok := initialValue.(lifecycle.RememberObserver); ok {
		onForgotten = observer.OnForgotten
	}

	// If implements HasViewModelScope, create and provide a cancellable context
	if scopeHolder, ok := initialValue.(lifecycle.HasViewModelScope); ok {
		ctx, cancelFunc := context.WithCancel(context.Background())
		scopeHolder.SetViewModelScope(ctx)
		scopedCancelFunc = cancelFunc
	}

	if vm, ok := initialValue.(lifecycle.ViewModel); ok {
		onCleared = vm.OnCleared
	}

	scopedValue := NewScopedValue(
		scopedCancelFunc,
		onForgotten,
		onCleared,
	)

	return scopedValue
}
