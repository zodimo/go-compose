package effect

import (
	"context"
	"fmt"
	"reflect"

	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
)

// LaunchedEffect runs a side-effect in a goroutine.
// The effect is restarted if any of the keys change.
func LaunchedEffect(block func(context.Context), keys ...any) api.Composable {
	return func(c api.Composer) api.Composer {
		c.StartBlock("LaunchedEffect")

		// Key for the persistent state, ensuring uniqueness per node using its ID.
		// Since State() might be global in this implementation, we need to scope it manually.
		stateKey := fmt.Sprintf("launched_effect_%d", c.GetID().Value())

		effectState := state.MustRemember(c, stateKey, func() *launchEffectState {
			return &launchEffectState{}
		})

		state := effectState.Get()

		// Check if keys changed
		keysChanged := !keysEqual(state.lastKeys, keys)

		if keysChanged {
			// Cancel previous
			if state.cancel != nil {
				state.cancel()
			}

			// Start new
			ctx, cancel := context.WithCancel(context.Background())
			state.cancel = cancel
			// Copy keys to ensure we store a snapshot (though variadic slice is usually fresh)
			keysCopy := make([]any, len(keys))
			copy(keysCopy, keys)
			state.lastKeys = keysCopy

			go func() {
				block(ctx)
			}()
		}

		// Set a dummy widget constructor that does nothing (zero size)
		c.SetWidgetConstructor(layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
			return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
				return layoutnode.LayoutDimensions{}
			}
		}))

		return c.EndBlock()
	}
}

type launchEffectState struct {
	cancel   context.CancelFunc
	lastKeys []any
}

// keysEqual compares two key slices using identity (==) for pointer and interface
// types, and reflect.DeepEqual for value types. This avoids deep-comparing mutable
// internal state of objects like StateFlow (which contain sync.Mutex, atomic.Value,
// channels, etc.) that would cause false negatives with reflect.DeepEqual.
func keysEqual(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !keyEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

// keyEqual compares two key values. For pointer, interface, channel, map, and func
// types it uses identity comparison (==). For value types it falls back to
// reflect.DeepEqual for structural comparison.
func keyEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	// Different types means different keys
	if va.Type() != vb.Type() {
		return false
	}

	switch va.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Map, reflect.Func:
		// Use pointer/identity comparison — avoids deep-comparing mutable internals
		return va.Pointer() == vb.Pointer()
	default:
		// Value types (int, string, struct, etc.) — use structural equality
		return reflect.DeepEqual(a, b)
	}
}
