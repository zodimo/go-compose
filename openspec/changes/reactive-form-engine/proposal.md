## Why

The current `compose/foundation/form` package is a UI shell without a state engine: `FormState` is an empty struct, `FormFieldState[T]` mutates `touched`/`disabled`/`err` as plain fields that never trigger recomposition, and cross-field logic (e.g. disabling one group when a checkbox changes) has no home. There is no tree/graph data structure anywhere in the repo to build on. A production form system needs a retained, traversable, type-safe node graph (Angular Reactive Forms-style) that is *pure Go* so it is unit-testable, plus a thin binding that connects node mutations to the framework's recomposition loop.

## What Changes

- **NEW** `compose/foundation/form/internal/` — a pure-Go form tree engine: `FormNode` interface, `Control[T]`, `Group`, `Array` nodes with parent pointers, recursive status propagation (`Valid`/`Invalid`/`Pending`/`Disabled`), `touched`/`dirty`/`pristine`/`enabled` lifecycle, typed validators, `RawValue()` export, and a path resolver supporting `Get("a.b.c")` and `Get("items[0].name")`.
- **NEW** `ValueStore[T]` seam — engine nodes read/write values through a tiny interface. A plain mutex holder is used for pure-Go tests; a `state.MutableValueTyped[T]` adapter is provided for UI use.
- **NEW** change notification — the engine uses `state.SubscriptionManager` for `OnValueChange` / `OnStatusChange` / `OnTouchChange`; group-level listeners actually fire when children change (fixing the no-op `updateAncestors` from the design sketch).
- **NEW** `Form()` composable integration — `RememberFormState(c)` builds the tree; the composable subscribes to tree notifications and bumps a remembered version `MutableValue`, which drives recomposition through the existing store → `window.Invalidate()` loop.
- **REPLACE** `FormFieldState[T]` — removed in favor of engine nodes + a `FormFieldBinding[T]` adapter. The public entry points used by `cmd/demo/form/ui.go` and `components/textfield.go` (`RememberFormState`, `RememberFormFieldState`, `Form`, `FormScope.Field/Form`, `TextFieldComponent`) keep working.
- **ENHANCE** error UX — errors are gated on touched state: `WithError(control.HasErrors() && control.IsTouched())` + `WithSupportingText(firstError)`, replacing the hand-rolled red `BodySmall` below the field.
- **DEAD-CODE MARKERS** — header comments added to `pkg/cforms/forms.go`, `pkg/cforms/errors.go`, `cmd/demo/form/form.go`, and the dead declarations in `types.go`/`options.go`/`components/textfield_options.go` (no deletion, per decision).

## Capabilities

### New Capabilities
- `form-engine`: The pure-Go form tree/graph engine — `FormNode` interface, `Control[T]`, `Group`, `Array` nodes, parent pointers, recursive status/touched/dirty propagation, `Pending` status, validators, `RawValue()` export, DFS/BFS traversal, path resolver (`a.b.c`, `items[0].name`), `ValueStore[T]` seam, and `SubscriptionManager`-based change notification with correct ancestor propagation.
- `form-binding`: The UI integration layer — `Form()` composable that remembers the tree and drives recomposition via a version tick, `FormScope.Field/Form` binding named nodes, `FormFieldBinding[T]` adapter replacing `FormFieldState[T]`, and touched-gated Material 3 error display (`WithError` + `WithSupportingText`).

### Modified Capabilities
<!-- No existing specs (openspec/specs/ is empty); all behavior here is net-new. -->

## Impact

- **New code**: `compose/foundation/form/internal/` (engine, ~5-7 files) and updates to `compose/foundation/form/` public surface (`form_state.go`, `scope.go`, `alias.go`, `components/textfield.go`, `components/textfield_options.go`).
- **Removed**: `FormFieldState[T]` and its dead methods (`ID`, `ClearError`, `Touched`, `Disable`, `Enable`, `Update`, `Clone`).
- **Dead-code markers**: `pkg/cforms/forms.go`, `pkg/cforms/errors.go`, `cmd/demo/form/form.go`, plus dead declarations in `compose/foundation/form/types.go`, `options.go`, `components/textfield_options.go`.
- **Dependencies**: none new — reuses `state.SubscriptionManager`, `state.MutableValueTyped`, `state.MustRemember` (all existing). The engine itself depends only on `state.SubscriptionManager`.
- **Demo**: `cmd/demo/form/ui.go` updated to exercise the tree engine (nested groups, array, cross-field disable, touched-gated errors); stale window title "Package CForms Demo" fixed in `main.go`.