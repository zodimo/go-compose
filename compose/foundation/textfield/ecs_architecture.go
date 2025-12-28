package textfield

import (
	"context"
	"sync"

	"gioui.org/widget"
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/compose/effect"
	"gioui.org/text"
	"gioui.org/font"
	"gioui.org/unit"
	"gioui.org/op"
)

// =========================================================================================
// 1. Data-Oriented Component Storage
// =========================================================================================

// TextFieldData acts as the "Archetype" for a TextField.
// It is a single POD (Plain Old Data) struct.
// Note: While we use a pointer to widget.Editor (due to its internal locks),
// keeping the rest of the state in a value struct allows for better cache locality if iterated.
type TextFieldData struct {
	// State
	Text          string
	Editor        *widget.Editor

	// Configuration (Props that might change)
	SingleLine    bool
	OnValueChange func(string)
	OnSubmit      func()

	// Tracking for Input System
	LastSeenText  string
}

// TextFieldSystem acts as the "System" in ECS.
// It manages a dense array of TextFieldData.
type TextFieldSystem struct {
	mu sync.RWMutex

	// Flat buffer of components.
	components []*TextFieldData

	// Registry mapping IDs to indices (Sparse Set pattern)
	registry map[uint32]int
}

// NewTextFieldSystem creates a new system instance.
// This should be created at the application root and passed down via CompositionLocal.
func NewTextFieldSystem() *TextFieldSystem {
	return &TextFieldSystem{
		components: make([]*TextFieldData, 0, 100),
		registry:   make(map[uint32]int),
	}
}

// LocalTextFieldSystem provides access to the TextFieldSystem in the composition tree.
var LocalTextFieldSystem = compose.CompositionLocalOf(func() *TextFieldSystem {
	// Default fallback (though typically the root should provide one)
	return NewTextFieldSystem()
})

// EnsureComponent allocates or retrieves the component for a given ID.
func (sys *TextFieldSystem) EnsureComponent(id uint32) *TextFieldData {
	sys.mu.Lock()
	defer sys.mu.Unlock()

	if idx, ok := sys.registry[id]; ok {
		if sys.components[idx] != nil {
			return sys.components[idx]
		}
	}

	// Create new component
	comp := &TextFieldData{
		Editor: &widget.Editor{},
	}

	// Find empty slot or append
	// (Simple append for now, reuse logic would go here)
	sys.components = append(sys.components, comp)
	sys.registry[id] = len(sys.components) - 1
	return comp
}

// RemoveComponent handles cleanup (e.g. when a node is disposed).
func (sys *TextFieldSystem) RemoveComponent(id uint32) {
	sys.mu.Lock()
	defer sys.mu.Unlock()

	if idx, ok := sys.registry[id]; ok {
		// Nil out the slot. A real ECS would use swap-remove or a free list.
		sys.components[idx] = nil
		delete(sys.registry, id)
	}
}

// =========================================================================================
// 2. The Input System (Separating Event Loop from Layout)
// =========================================================================================

// ProcessEvents would ideally be called *before* the Layout phase in the frame loop.
// It iterates over all active text fields and processes their input events.
func (sys *TextFieldSystem) ProcessEvents(gtx layoutnode.LayoutContext) {
	sys.mu.RLock()
	defer sys.mu.RUnlock()

	for _, comp := range sys.components {
		if comp == nil {
			continue
		}

		// 1. Drive Editor Events
		for {
			event, ok := comp.Editor.Update(gtx)
			if !ok {
				break
			}

			switch e := event.(type) {
			case widget.SubmitEvent:
				if comp.OnSubmit != nil {
					comp.OnSubmit()
				}
				_ = e
			}
		}

		// 2. Detect User Change (State Sync)
		newText := comp.Editor.Text()
		if newText != comp.Text {
			if comp.OnValueChange != nil {
				comp.OnValueChange(newText)
			}
		}
	}
}

// =========================================================================================
// 3. The Composable (View)
// =========================================================================================

type ECSTextFieldOption func(*TextFieldData)

// ECSTextField is the refactored Composable.
// It uses the System to manage state, rather than internal closures.
func ECSTextField(
	value string,
	onValueChange func(string),
	options ...ECSTextFieldOption,
) Composable {
	return func(c Composer) Composer {
		// 1. Identity
		// Use the numeric value of the ID for map lookups (avoids string allocation)
		id := c.GenerateID()
		entityID := id.Value()

		// 2. System Access
		system := LocalTextFieldSystem.Current(c)

		// 3. Component Retrieval
		comp := system.EnsureComponent(entityID)

		// 4. Lifecycle Management
		// We use LaunchedEffect to handle cleanup on disposal.
		// When the key (entityID) changes or the node is removed, the context is cancelled.
		c.WithComposable(effect.LaunchedEffect(func(ctx context.Context) {
			<-ctx.Done()
			system.RemoveComponent(entityID)
		}, entityID))

		// 5. Sync Props to Component
		comp.Text = value
		comp.OnValueChange = onValueChange
		for _, opt := range options {
			opt(comp)
		}

		// 6. Sync External Change (Data Binding)
		if comp.Text != comp.LastSeenText {
			if comp.Editor.Text() != comp.Text {
				comp.Editor.SetText(comp.Text)
			}
			comp.LastSeenText = comp.Text
		}

		// 7. Render
		c.StartBlock("ECSTextField")
		c.SetWidgetConstructor(layoutnode.NewLayoutNodeWidgetConstructor(
			func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
				return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
					// NOTE: We rely on system.ProcessEvents(gtx) being called elsewhere (e.g. root)
					// But for safety in this hybrid environment, we can peek:
					// comp.Editor.Update(gtx) // If needed here

					// Use a basic shaper. In a real app, this should come from a CompositionLocal or Theme.
					shaper := text.NewShaper(text.NoSystemFonts())
					return comp.Editor.Layout(gtx, shaper, font.Font{}, unit.Sp(14), op.CallOp{}, op.CallOp{})
				}
			},
		))

		return c.EndBlock()
	}
}
