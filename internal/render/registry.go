package render

import "sync"

var (
	registryMu sync.Mutex
	registry   = map[string]BackendFactory{}
	activeName = "gio"
)

// Register registers a render backend factory under the given name.
// It is typically called from an init() function in the backend's package,
// which is blank-imported by the app shell (design D6).
//
// Example:
//
//	import _ "github.com/zodimo/go-compose/internal/render/gio"
func Register(name string, f BackendFactory) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[name] = f
}

// ActiveBackend returns a new instance of the currently active backend.
// Panics if the named backend has not been registered.
func ActiveBackend() Backend {
	registryMu.Lock()
	defer registryMu.Unlock()
	f, ok := registry[activeName]
	if !ok {
		panic("render: no backend registered for " + activeName)
	}
	return f()
}

// SetActiveBackend selects the active backend by name.
// Panics if the named backend has not been registered.
func SetActiveBackend(name string) {
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, ok := registry[name]; !ok {
		panic("render: unknown backend: " + name)
	}
	activeName = name
}

// ActiveBackendName returns the name of the currently active backend.
func ActiveBackendName() string {
	registryMu.Lock()
	defer registryMu.Unlock()
	return activeName
}

// resetForTesting clears the registry and resets the active backend.
// Only for use in tests to ensure isolation between test cases.
func resetForTesting() {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry = map[string]BackendFactory{}
	activeName = "gio"
}
