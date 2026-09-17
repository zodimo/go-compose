## 1. Engine scaffolding (compose/foundation/form/internal)

- [ ] 1.1 Create `internal/` package skeleton with `package formengine` (per design D1, D12 open question — confirm name before starting)
- [ ] 1.2 Implement `valuestore.go`: `ValueStore[T]` interface, `PlainValueStore[T]` (mutex-guarded), `NewPlainValueStore[T]`, and the `MutableValueValueStore[T]` adapter over `state.MutableValueTyped[T]` (design D2)
- [ ] 1.3 Implement `notify.go`: typed change streams (`OnValueChange`, `OnStatusChange`, `OnTouchChange`) backed by `state.SubscriptionManager`; `notify` + `onChildChanged` ancestor propagation (design D3, D8)
- [ ] 1.4 Implement `status.go`: `Status` type with `StatusValid`/`StatusInvalid`/`StatusPending`/`StatusDisabled` constants and `StatusPending` reserved path

## 2. Engine nodes

- [ ] 2.1 Implement `formnode.go`: `FormNode` interface (status, validation, lifecycle, actions, tree hooks, export — design D6)
- [ ] 2.2 Implement `control.go`: `Control[T]` with `ValueStore[T]`, initial value, validators, `Set`/`Get`/`Reset`/`MarkAsTouched`/`MarkAsUntouched`/`MarkAsPristine`/`SetDisabled`/`Validate`/`Errors`/`RawValue` (design D5: validators run on a value snapshot, never under lock)
- [ ] 2.3 Implement own-vs-inherited disabled: `ownDisabled` + `inheritedDisabled`, effective = own || inherited; `SetDisabled` touches only own (design D4)
- [ ] 2.4 Implement `group.go`: `Group` with `map[string]FormNode` children, status aggregation (INVALID > PENDING > DISABLED-when-all > VALID), lifecycle propagation, `RawValue()` as `map[string]any` omitting disabled, `Get(path)` delegation to resolver
- [ ] 2.5 Implement `array.go`: `Array` with `[]FormNode` children, `Add`/`Insert`/`Remove`/`Length`/`Get(index)`, same status rules as Group, `RawValue()` as ordered `[]any`, index-shift notifications
- [ ] 2.6 Implement `validator.go`: `ValidatorFunc[T] func(T) error`, built-in `Required`/`MinLength` (matching existing `ValidationRule[T]` convention), `errors.Join` aggregation

## 3. Traversal and path resolution

- [ ] 3.1 Implement `walk.go`: DFS/BFS visitor traversal over `FormNode` (root-to-leaf and leaf-to-root for propagation helpers)
- [ ] 3.2 Implement `path.go`: dotted + bracketed path parser/resolver — `a.b.c`, `items[0].name`, mixed nesting; miss indicator on unknown key or out-of-range index; document array index-instability caveat (design D7)
- [ ] 3.3 Implement `Node.Path()` (dotted path from root) used by `Errors()` keying

## 4. Engine tests (pure Go, no UI)

- [ ] 4.1 Unit tests: control value/lifecycle/validation (`Set`, dirty/pristine, touched, validators aggregate, reentrancy no-deadlock)
- [ ] 4.2 Unit tests: status propagation (invalid child, all-disabled group, pending), own-vs-inherited disabled round-trip
- [ ] 4.3 Unit tests: group/array raw value export (map, ordered slice, disabled omitted)
- [ ] 4.4 Unit tests: path resolver (dotted, bracketed, mixed, miss)
- [ ] 4.5 Unit tests: change notification (group status listener fires on child change; value listener fires on set; unsubscribe works)
- [ ] 4.6 Unit tests: reset restores initial state
- [ ] 4.7 Run `go test ./compose/foundation/form/internal/...` green

## 5. Public API binding (compose/foundation/form)

- [ ] 5.1 Rework `form_state.go`: `RememberFormState(c)` builds and remembers the tree root (`Group`); `FormFieldBinding[T]` replaces `FormFieldState[T]`; keep `RememberFormFieldState` as a deprecated alias for `RememberFormFieldBinding` (design D9, D10)
- [ ] 5.2 Rework `form.go`: `Form()` remembers a version `MutableValue`, subscribes to the tree change stream, bumps version on any change → recomposition (design D9)
- [ ] 5.3 Rework `scope.go`: `FormScope.Field`/`Form` resolve keys against the tree via the path resolver and pass node bindings into content closures; nested scopes resolve relative paths
- [ ] 5.4 Add `alias.go` re-exporting engine types (`FormNode`, `Control`, `Group`, `Array`, `Status`, validators) per repo alias convention
- [ ] 5.5 Remove `FormFieldState[T]` and its dead methods (`ID`, `ClearError`, `Touched`, `Disable`, `Enable`, `Update`, `Clone`)
- [ ] 5.6 Update `components/textfield.go` to consume `FormFieldBinding[string]` with touched-gated errors: `WithError(control.HasErrors() && control.IsTouched())` + `WithSupportingText(firstError)`, dropping the hand-rolled red `BodySmall` (design D11)

## 6. Demo update

- [ ] 6.1 Update `cmd/demo/form/ui.go` to the new API — nested groups, array usage, cross-field disable (sameAsShipping pattern), touched-gated errors
- [ ] 6.2 Fix stale window title in `cmd/demo/form/main.go` ("Package CForms Demo" → accurate title)
- [ ] 6.3 Run `go run ./cmd/demo/form/` and verify the form renders and validates

## 7. Dead-code markers (no deletion)

- [ ] 7.1 Add header `// DEAD CODE` comments to `pkg/cforms/forms.go` and `pkg/cforms/errors.go`
- [ ] 7.2 Add header `// DEAD CODE` comment to `cmd/demo/form/form.go`
- [ ] 7.3 Add `// DEAD` markers to `Validatable` and `Orientation` in `types.go`
- [ ] 7.4 Add `// DEAD` markers to `WithTextStyle`, `WithTextStyleOption`, `WithOrientation` in `options.go` (leave `WithModifier` live)
- [ ] 7.5 Add `// DEAD` markers to `TextFieldWithModifier`, `TextFieldWithTextStyle`, `TextFieldWithTextStyleOption`, `TextFieldWithInline` in `components/textfield_options.go`

## 8. Verification

- [ ] 8.1 `make test` (go test ./...) passes
- [ ] 8.2 `lsp_diagnostics` clean on all changed files
- [ ] 8.3 Confirm dead-code markers present and no dead code deleted
- [ ] 8.4 Manual smoke test of form demo (web + desktop if feasible)