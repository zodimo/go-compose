## Context

The `compose/foundation/form` package is a UI shell: `Form()` is a `LazyColumn` + `CompositionLocalProvider`, `FormState` is an empty struct, and `FormFieldState[T]` holds `value` (backed by `state.MutableValueTyped[T]`) plus `touched`/`disabled`/`err` as **plain fields that never trigger recomposition**. There is no tree/graph data structure anywhere in the repo — `foundation/layout/tree` is a rendering widget with no node objects or parent pointers, and `internal/node.TreeNode` is internal, parentless, and non-traversable. A design sketch (`form/OiBqXP9A9jiK-2026-09-17-162915.md`) proposes an Angular Reactive Forms-style engine but does not compile and has structural defects.

Framework facts that constrain the design:
- **No snapshot / no fine-grained recomposition.** Every `state.MutableValue.Set` fans out through a single store notifier → `window.Invalidate()` → the *entire* composable tree re-executes. This is the only recomposition mechanism.
- `state.SubscriptionManager` is the production callback registry (thread-safe, idempotent `Unsubscribe`).
- `state.DerivedState` is the only dependency-tracking primitive (not needed for this design — status is computed on read).
- `pkg/flow` is experimental with zero production imports; `MutableStateFlow` implements `state.StateChangeNotifier` and can be adopted later without breaking the seam.
- The Material 3 convention is `WithError(bool)` + `WithSupportingText(string)` where supporting text doubles as the error message.

## Goals / Non-Goals

**Goals:**
- A pure-Go, unit-testable form tree engine: `FormNode`, `Control[T]`, `Group`, `Array`, parent pointers, recursive status/lifecycle propagation, validators, `RawValue()` export, and a real path resolver (`a.b.c`, `items[0].name`).
- Correct change propagation: child mutations notify ancestor listeners (fixing the sketch's no-op `updateAncestors`).
- UI integration that triggers recomposition through the existing store loop, with no new dependencies.
- Touched-gated error display via the existing Material 3 error API.
- Replace `FormFieldState[T]` while keeping the public entry points working (`Form`, `FormScope`, `TextFieldComponent`, `RememberFormState`, `RememberFormFieldState`).

**Non-Goals:**
- Fine-grained recomposition (does not exist in the framework; whole-tree recomposition is accepted).
- `Array` *rendering* (add/remove row UI via `lazy.Range()`/`Key()`); the `Array` node is in the engine, rendering is a follow-up.
- Adopting `pkg/flow` as the notification backbone (reusable later via the same seam).
- Deleting dead code — only header `// DEAD CODE` markers per decision.
- Async validators (only the `Pending` status is reserved; async scheduling is out of scope).

## Decisions

### D1. Engine lives in `compose/foundation/form/internal/`
The tree engine is `internal/` so only `compose/foundation/form` and its subpackages (e.g. `components/`) can import it. Demos and external users reach the engine through `alias.go` re-exports on the public package.

*Alternatives:* `pkg/form` (rejected — non-UI slot, would split the API across the repo and duplicate the existing `pkg/cforms` mistake); `compose/foundation/next/form/` (rejected per user decision — internal keeps the engine private to the form package).

### D2. `ValueStore[T]` seam — engine stays framework-agnostic
```go
type ValueStore[T any] interface {
    Get() T
    Set(T)
}
```
- **Plain holder** (default, zero-dep): `NewPlainValueStore[T](initial T)` — mutex-guarded, for pure-Go tests and non-UI use.
- **UI adapter**: `NewMutableValueValueStore[T](mv state.MutableValueTyped[T])` — delegates to a remembered `MutableValueTyped[T]`, so `Set` triggers recomposition.

*Alternatives:* engine directly holding `state.MutableValueTyped[T]` (rejected — couples the engine to the composer's remember lifecycle and kills pure-Go testability); engine holding a plain `T` with a global notifier (accepted as the plain holder is exactly that).

### D3. Status & lifecycle computed on read; propagation only for listeners
- `Status()` is **computed recursively on read** — no cached status field, no version counter. Because the framework re-renders the whole tree on any mutation, reads during composition are always current.
- Change notification (`notify`) is **eager and propagates**: every mutating method (`Set`, `MarkAsTouched`, `SetDisabled`, `Reset`, `Array.Add/Insert/Remove`) calls `node.notify(...)` then `parent.onChildChanged()` up the chain. Each node's `onChildChanged` recomputes *its own* derived state and notifies *its own* subscribers before recursing up. This fixes the sketch's no-op `updateAncestors`.
- `StatusPending` is defined and returned when a node has pending async work (reserved; no async validators in this change).

*Alternatives:* cached status + push invalidation (rejected — needs invalidation bookkeeping for zero benefit since reads are cheap and recomposition is whole-tree).

### D4. Disabled = own || inherited (fixes the sketch's destructive `SetDisabled`)
Each node tracks `ownDisabled` and `inheritedDisabled` separately. `SetDisabled(x)` sets only `ownDisabled`; effective disabled is `own || parent.InheritedDisabled()`. `Group.SetDisabled(true)` no longer destroys child state — re-enabling the parent restores each child's own setting.

### D5. Validators return `error`, not `string`
`ValidatorFunc[T] func(value T) error`, aggregated with `errors.Join` into `control.Errors()` (all messages joined, so the UI can surface every failure). This matches the existing `ValidationRule[T]` convention and the `components/textfield.go` validate closure. The sketch's `func(value T, node *Control[T]) string` (which doesn't compile — `err != nil` on a string) is rejected. Validation runs on a **snapshot** of the value taken under `RLock`, then errors are stored under `Lock` — validators never run while the node lock is held (fixes reentrancy deadlock).

### D6. `FormNode` interface & node shapes
```go
type FormNode interface {
    // Status & validation
    Status() Status                       // VALID | INVALID | PENDING | DISABLED
    Validate() bool
    Errors() map[string]string            // keyed by path for groups; controls join all errors via errors.Join

    // Lifecycle
    IsTouched() bool
    IsDirty() bool
    IsPristine() bool
    IsEnabled() bool

    // Actions
    MarkAsTouched()
    MarkAsUntouched()
    MarkAsPristine()
    SetDisabled(disabled bool)
    Reset()

    // Tree
    Parent() FormNode
    Children() []FormNode
    Path() string                         // dotted path from root
    Get(path string) (FormNode, bool)     // path resolver (D7)

    // Export
    RawValue() any

    // Internal hooks (unexported on impls; promoted via interface)
    setParent(p FormNode)
    onChildChanged()                      // propagation (D3)
}
```
- `Control[T]`: leaf. Holds `valueStore ValueStore[T]`, `initialValue`, validators, `errors []error`, `ownDisabled`, `touched`, `dirty`. `RawValue()` returns `nil` when effectively disabled.
- `Group`: `children map[string]FormNode`. `Status()`: `INVALID` if any child `INVALID`; `DISABLED` if *all* children disabled; `PENDING` if any child `PENDING`; else `VALID`. `RawValue()` returns `map[string]any` omitting disabled children.
- `Array`: `children []FormNode` + `Add/Insert/Remove/Length/At(index)`. Same status rules as Group. `RawValue()` returns `[]any` (ordered, omitting disabled). Removing/inserting shifts indices and fires `onChildChanged` up the chain. (`Get` is reserved on the `FormNode` interface for path resolution, so indexed access is `At(index)`.)

### D7. Path resolver (`path.go`)
Parse and resolve `a.b.c` and `items[0].name` (and mixed `group.items[2].field`):
- Split on `.`; segments are either map keys (`a`) or array indices (`items[0]`, parsed via `[`...`]`).
- `FormNode.Get(path)` walks segments; returns `(node, false)` on any miss or index-out-of-range.
- `Control.Get` returns `(nil, false)` for any non-empty path.
- Documented caveat: array indices are positional and shift on `Remove`/`Insert` — resolve early, hold the node, don't cache paths.

### D8. Change notification via `state.SubscriptionManager`
Per-node `OnValueChange/OnStatusChange/OnTouchChange` are backed by small typed wrappers over one `state.SubscriptionManager` (or three, one per stream kind). This reuses the production primitive (`Subscription` with idempotent `Unsubscribe`), so the engine has exactly one framework dependency, and `pkg/flow.MutableStateFlow` could replace it later without changing callers.

### D9. UI binding — version tick drives recomposition
```go
// form.go
func Form(state *FormState, content func(FormScope), options ...FormOption) api.Composable {
    // 1. remember a version MutableValue: version := state.MustRemember(c, ".../formVersion", ...)
    // 2. subscribe to the tree root's change stream → version.Set(version.Get()+1)
    // 3. render LazyColumn as today; FormScope.Field/Form resolve nodes by key via the path resolver
}
```
`RememberFormState(c)` builds the tree (a `Group`) and remembers it via `state.MustRemember`. Any engine mutation → tree notifier → `version.Set` → store notifier → `window.Invalidate()` → whole subtree recomposes → composables re-read `control.Value()/Errors()/IsTouched()`.

`FormScope` gains node binding: `Field(key any, ...)` and `Form(key any, header, children)` now look up the corresponding tree node by path and pass it into the content closure's field components.

### D10. `FormFieldState[T]` replaced by `FormFieldBinding[T]`
```go
type FormFieldBinding[T any] struct {
    control *Control[T]
}
```
`FormFieldBinding` wraps a `*Control[T]` only; the remembered `MutableValueTyped[T]` store lives *inside* the control, wired through the `MutableValueValueStore` adapter at remember time (`RememberFormFieldBinding` builds `NewControl(NewMutableValueValueStore(itemState), itemState.Get())`). The binding:
- reads `control.Value()/Errors()/IsTouched()` for the UI;
- `SetValue` → `control.Set` (through the `MutableValueValueStore` adapter, driving recomposition) + `MarkAsTouched()`.

`RememberFormFieldState(c, itemState)` is kept as a deprecated alias for `RememberFormFieldBinding(c, itemState)` so existing call sites keep working.

### D11. Touched-gated Material 3 errors
`components/textfield.go` switches to:
```go
textfield.WithError(control.HasErrors() && control.IsTouched())
textfield.WithSupportingText(allErrors) // errors.Join of every validator failure
```
and drops the hand-rolled red `BodySmall` error line (per user decision). `WithSupportingText` doubles as the error message (M3 convention, `outlined.go:98-101`).

### D12. Dead-code markers (no deletion)
Header `// DEAD CODE ...` comments, pointing at the replacement:
- `pkg/cforms/forms.go`, `pkg/cforms/errors.go` — superseded by the engine.
- `cmd/demo/form/form.go` — `Exampleorm()` never called; live demo is `ui.go`.
- `compose/foundation/form/types.go` — `Validatable`, `Orientation`.
- `compose/foundation/form/options.go` — `WithTextStyle`, `WithTextStyleOption`, `WithOrientation` (keep `WithModifier`).
- `components/textfield_options.go` — `TextFieldWithModifier/TextStyle/TextStyleOption/Inline`.

Also fix the stale demo title "Package CForms Demo" in `cmd/demo/form/main.go`.

## Risks / Trade-offs

- **Whole-tree recomposition on every keystroke** → Accepted: it is the framework's only mechanism; no per-field perf work until profiling proves a need.
- **Array index instability** (paths shift on add/remove) → Mitigation: documented in `path.go`; UI resolves paths early and holds node references.
- **Engine owns one framework dependency (`state.SubscriptionManager`)** → If the engine must be 100% dependency-free later, swap for a ~30-line internal subscription manager; the interface is stable.
- **`internal/` blocks external consumers** → Mitigation: `alias.go` re-exports on `compose/foundation/form`; engine stays private per decision.
- **Replace `FormFieldState` is a breaking API change for any external users** → Only known consumer is `cmd/demo/form/ui.go` (updated in this change); keep `RememberFormFieldState` as a compatibility alias if needed.
- **Async validators deferred** → `Pending` status reserved but never emitted in this change; document that async scheduling is follow-up work.

## Migration Plan

1. Land the engine (`internal/`) + its unit tests first — pure Go, no UI impact, demo still compiles.
2. Wire `FormState`/`FormScope`/`FormFieldBinding` to the engine; update `components/textfield.go`.
3. Update `cmd/demo/form/` to the new API; fix the stale title.
4. Add dead-code header comments.
5. `make test` green; run `go run ./cmd/demo/form/` manually.

Rollback: the engine is additive (`internal/`); only steps 2-3 touch the public surface, and the demo + textfield are the only consumers — revert those two commits to roll back.

## Open Questions

- Package name inside `internal/`: `package form` (mirrors parent, confusing in godoc) vs `package formengine` (clear). Default: `formengine` unless reviewer objects.
- Keep `RememberFormFieldState` as a deprecated alias for `RememberFormFieldBinding`, or rename outright in the demo? Default: keep alias, demo uses the new name.