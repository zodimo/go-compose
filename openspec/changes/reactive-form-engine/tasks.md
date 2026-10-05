## 1. Engine scaffolding (compose/foundation/form/internal)

- [x] 1.1 Create `internal/` package skeleton with `package formengine` (per design D1, D12 open question — confirm name before starting)
- [x] 1.2 Implement `valuestore.go`: `ValueStore[T]` interface, `PlainValueStore[T]` (mutex-guarded), `NewPlainValueStore[T]`, and the `MutableValueValueStore[T]` adapter over `state.MutableValueTyped[T]` (design D2)
- [x] 1.3 Implement `notify.go`: typed change streams (`OnValueChange`, `OnStatusChange`, `OnTouchChange`) backed by `state.SubscriptionManager`; `notify` + `onChildChanged` ancestor propagation (design D3, D8)
- [x] 1.4 Implement `status.go`: `Status` type with `StatusValid`/`StatusInvalid`/`StatusPending`/`StatusDisabled` constants and `StatusPending` reserved path

## 2. Engine nodes

- [x] 2.1 Implement `formnode.go`: `FormNode` interface (status, validation, lifecycle, actions, tree hooks, export — design D6)
- [x] 2.2 Implement `control.go`: `Control[T]` with `ValueStore[T]`, initial value, validators, `Set`/`Get`/`Reset`/`MarkAsTouched`/`MarkAsUntouched`/`MarkAsPristine`/`SetDisabled`/`Validate`/`Errors`/`RawValue` (design D5: validators run on a value snapshot, never under lock)
- [x] 2.3 Implement own-vs-inherited disabled: `ownDisabled` + `inheritedDisabled`, effective = own || inherited; `SetDisabled` touches only own (design D4)
- [x] 2.4 Implement `group.go`: `Group` with `map[string]FormNode` children, status aggregation (INVALID > PENDING > DISABLED-when-all > VALID), lifecycle propagation, `RawValue()` as `map[string]any` omitting disabled, `Get(path)` delegation to resolver
- [x] 2.5 Implement `array.go`: `Array` with `[]FormNode` children, `Add`/`Insert`/`Remove`/`Length`/`At(index)`, same status rules as Group, `RawValue()` as ordered `[]any`, index-shift notifications
- [x] 2.6 Implement `validator.go`: `ValidatorFunc[T] func(T) error`, built-in `Required`/`MinLength` (matching existing `ValidationRule[T]` convention), `errors.Join` aggregation

## 3. Traversal and path resolution

- [x] 3.1 Implement `walk.go`: DFS/BFS visitor traversal over `FormNode` (root-to-leaf and leaf-to-root for propagation helpers)
- [x] 3.2 Implement `path.go`: dotted + bracketed path parser/resolver — `a.b.c`, `items[0].name`, mixed nesting; miss indicator on unknown key or out-of-range index; document array index-instability caveat (design D7)
- [x] 3.3 Implement `Node.Path()` (dotted path from root) used by `Errors()` keying

## 4. Engine tests (pure Go, no UI)

- [x] 4.1 Unit tests: control value/lifecycle/validation (`Set`, dirty/pristine, touched, validators aggregate, reentrancy no-deadlock)
- [x] 4.2 Unit tests: status propagation (invalid child, all-disabled group, pending), own-vs-inherited disabled round-trip
- [x] 4.3 Unit tests: group/array raw value export (map, ordered slice, disabled omitted)
- [x] 4.4 Unit tests: path resolver (dotted, bracketed, mixed, miss)
- [x] 4.5 Unit tests: change notification (group status listener fires on child change; value listener fires on set; unsubscribe works)
- [x] 4.6 Unit tests: reset restores initial state
- [x] 4.7 Run `go test ./compose/foundation/form/internal/...` green

## 5. Public API binding (compose/foundation/form)

- [x] 5.1 Rework `form_state.go`: `RememberFormState(c)` builds and remembers the tree root (`Group`); `FormFieldBinding[T]` replaces `FormFieldState[T]`; keep `RememberFormFieldState` as a deprecated alias for `RememberFormFieldBinding` (design D9, D10)
- [x] 5.2 Rework `form.go`: `Form()` remembers a version `MutableValue`, subscribes to the tree change stream, bumps version on any change → recomposition (design D9)
- [x] 5.3 Rework `scope.go`: `FormScope.Field`/`Form` resolve keys against the tree via the path resolver and pass node bindings into content closures; nested scopes resolve relative paths
- [x] 5.4 Add `alias.go` re-exporting engine types (`FormNode`, `Control`, `Group`, `Array`, `Status`, validators) per repo alias convention
- [x] 5.5 Remove `FormFieldState[T]` and its dead methods (`ID`, `ClearError`, `Touched`, `Disable`, `Enable`, `Update`, `Clone`)
- [x] 5.6 Update `components/textfield.go` to consume `FormFieldBinding[string]` with touched-gated errors: `WithError(control.HasErrors() && control.IsTouched())` + `WithSupportingText(all errors joined via errors.Join)`, dropping the hand-rolled red `BodySmall` (design D11)

## 6. Demo update

- [x] 6.1 Update `cmd/demo/form/ui.go` to the new API — nested groups, array usage, cross-field disable (sameAsShipping pattern), touched-gated errors
- [x] 6.2 Fix stale window title in `cmd/demo/form/main.go` ("Package CForms Demo" → accurate title)
- [x] 6.3 Run `go run ./cmd/demo/form/` and verify the form renders and validates

## 7. Dead-code markers (no deletion)

- [x] 7.1 Add header `// DEAD CODE` comments to `pkg/cforms/forms.go` and `pkg/cforms/errors.go`
- [x] 7.2 Add header `// DEAD CODE` comment to `cmd/demo/form/form.go`
- [x] 7.3 Add `// DEAD` markers to `Validatable` and `Orientation` in `types.go`
- [x] 7.4 Add `// DEAD` markers to `WithTextStyle`, `WithTextStyleOption`, `WithOrientation` in `options.go` (leave `WithModifier` live)
- [x] 7.5 Add `// DEAD` markers to `TextFieldWithModifier`, `TextFieldWithTextStyle`, `TextFieldWithTextStyleOption`, `TextFieldWithInline` in `components/textfield_options.go`

## 8. Verification

- [x] 8.1 `make test` (go test ./...) passes
- [x] 8.2 `lsp_diagnostics` clean on all changed files
- [x] 8.3 Confirm dead-code markers present and no dead code deleted
- [x] 8.4 Manual smoke test of form demo (web + desktop if feasible)
## 9. Follow-up: form lifecycle, typed binding, array UI, extra components

- [x] 9.1 `FormState` lifecycle surface: `Validate`/`Status`/`Errors`/`Value`/`RawValue`/`IsTouched`/`IsDirty`/`IsPristine`/`IsEnabled`/`MarkAllTouched`/`MarkAllUntouched`/`MarkAllPristine`/`Reset`/`SetDisabled`/`Root`/`Submit`
- [x] 9.2 `Resolver` interface + `NodeResolver` adapter; `ControlOf[T]`, `ControlFieldOf[T]`, `ControlViewOf[T]` generic helpers removing manual node type assertions
- [x] 9.3 `FormScope` additions: `FormArray` (row iteration), `GroupScope`, `ArrayLength`, `Array`, `Group`, `Has`, `Resolve`; namespaced lazy item keys
- [x] 9.4 `FormFieldBinding` read surface: `IsEnabled`/`IsDirty`/`Status`/`ErrorMessage`
- [x] 9.5 Components: `NumberFieldComponent`, `CheckboxComponent`, `SwitchComponent`, `SelectComponent`; shared `field.go` touched-gating helpers; disabled-aware `TextFieldComponent`
- [x] 9.6 Components option files for checkbox, switch, select, number; un-marked live `TextFieldWith*` options and wired `Inline`/`TextStyle`
- [x] 9.7 Demo: real Submit/Reset, dynamic phone add/remove via `FormArray`, select and number fields, live status line, submission summary
- [x] 9.8 Engine exports: `Walk`/`WalkUp` in `alias.go`
- [x] 9.9 Tests: `FormState` lifecycle, scope/array resolution (fake lazy scope), demo path/type guards and cross-field disable
- [x] 9.10 `compose/foundation/form/README.md` documenting the public API; openspec spec deltas for the new capabilities
- [x] 9.11 `go build ./...`, `go vet ./...`, `go test ./...` all green

## 10. Prototype: rules catalog + coded errors (inspired by nobl9/govy)

- [x] 10.1 `error.go`: `ErrorCode`, `*ValidationError{Code,Message,Path}`, `NewValidationError`, `CodeOf`, `CodesOf`, `WithCode`; survives `errors.Join`/`%w`
- [x] 10.2 `Control.ValidationErrors() []error` returning raw failures for structured inspection
- [x] 10.3 `FormFieldBinding.ValidationErrors`/`Codes`; `FormState.CodedErrors`/`Codes`/`FirstInvalidPath`
- [x] 10.4 Base `Required`/`MinLength` now emit `CodeRequired`/`CodeMinLength`
- [x] 10.5 `rules` package: length, string shape, formats, comparables, collections, cross-field; all original implementations (no govy source)
- [x] 10.6 Tests: rules catalog, code survival through join/wrap, WithCode preservation, FormState coded access
- [x] 10.7 Demo: email + password/confirm (EqualTo) + phone rules + code-aware status line
- [x] 10.8 README + spec deltas; provenance note (MPL-2.0 inspiration, original code)
