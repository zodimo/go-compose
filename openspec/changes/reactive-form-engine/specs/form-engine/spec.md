## ADDED Requirements

### Requirement: Form node tree structure
The form engine SHALL model a form as a tree of nodes implementing the `FormNode` interface: leaf `Control[T]` nodes, `Group` nodes with named children, and `Array` nodes with indexed children. Every node SHALL maintain a parent reference, expose its children, and support depth-first traversal.

#### Scenario: Build a nested form tree
- **WHEN** a `Group` is created with child nodes including another `Group` and an `Array`
- **THEN** every child node reports its `Parent()`, and the root `Group` reports all descendants through `Children()`

#### Scenario: Array nodes support indexed children
- **WHEN** an `Array` node has `Add`, `Insert`, `Remove`, and `Length` operations performed on it
- **THEN** children are stored in order, indices shift correctly on insert/remove, and `At(index)` returns the node at that position

#### Scenario: Parent pointer integrity
- **WHEN** a node is added to a `Group` or `Array`
- **THEN** the node's `Parent()` returns that container, and the container's `Children()` includes the node

### Requirement: Control value management
A `Control[T]` SHALL read and write its typed value through a `ValueStore[T]` seam. The engine SHALL provide a plain mutex-guarded store for non-UI use and accept an adapter over `state.MutableValueTyped[T]` for UI use.

#### Scenario: Set and read a control value
- **WHEN** `control.Set(v)` is called on a `Control[string]`
- **THEN** `control.Value()` returns `v`, the control is marked dirty, and validators are re-run

#### Scenario: Value store seam is swappable
- **WHEN** a `Control[T]` is constructed with a plain store
- **THEN** it can be constructed and mutated in pure Go without any UI or composition dependency

### Requirement: Status propagation
Node status SHALL be one of `VALID`, `INVALID`, `PENDING`, or `DISABLED`, computed recursively: a `Group` or `Array` is `INVALID` if any child is `INVALID`, `PENDING` if any child is `PENDING`, `DISABLED` if all children are disabled, otherwise `VALID`. A disabled node SHALL report `DISABLED` regardless of children.

#### Scenario: Invalid child makes group invalid
- **WHEN** a `Group` contains a `Control` that fails validation
- **THEN** the `Group.Status()` returns `INVALID`

#### Scenario: All children disabled makes group disabled
- **WHEN** every child of a `Group` is effectively disabled
- **THEN** the `Group.Status()` returns `DISABLED`

#### Scenario: Pending status is reserved
- **WHEN** a node has pending asynchronous validation work
- **THEN** its `Status()` returns `PENDING`

### Requirement: Disabled state distinguishes own and inherited
A node SHALL track its own disabled flag separately from disabled state inherited from ancestors. Disabling a parent SHALL NOT destroy child-owned disabled settings.

#### Scenario: Re-enabling parent restores child state
- **WHEN** a `Group` is `SetDisabled(true)` and then `SetDisabled(false)`
- **THEN** each child returns to its own previously-set enabled/disabled state

#### Scenario: Effective disabled is own or inherited
- **WHEN** a child control is enabled but its parent group is disabled
- **THEN** the child's `IsEnabled()` returns false and its status is `DISABLED`

### Requirement: Lifecycle state
Nodes SHALL track `touched`, `dirty`, and `pristine` state. `IsPristine()` SHALL be the negation of `IsDirty()`. `MarkAsTouched`, `MarkAsUntouched`, and `MarkAsPristine` SHALL propagate to children of groups and arrays.

#### Scenario: Touching a control marks the group touched
- **WHEN** a leaf `Control` is `MarkAsTouched()`
- **THEN** the control and its ancestor `Group`s report `IsTouched() == true`

#### Scenario: Setting a value makes control dirty
- **WHEN** `control.Set(v)` is called on a control holding a different value
- **THEN** `IsDirty()` returns true and `IsPristine()` returns false until `MarkAsPristine()` is called

### Requirement: Validation
`Control[T]` SHALL accept one or more validators of signature `func(T) error`. Validation SHALL run on a snapshot of the value, aggregate failures, and expose them via `Errors()`. Validators MUST NOT run while the node's internal lock is held.

#### Scenario: Multiple validators aggregate errors
- **WHEN** a control has a required validator and a length validator, and both fail
- **THEN** `Errors()` exposes all error messages (joined via `errors.Join`) and `Validate()` returns false

#### Scenario: Validation does not deadlock on reentrancy
- **WHEN** a validator calls back into the control's read methods
- **THEN** validation completes without deadlock

### Requirement: Raw value export
`RawValue()` SHALL export the tree as plain Go values: `Control` exports its value, `Group` exports `map[string]any` keyed by child name, `Array` exports `[]any` in child order. Disabled nodes SHALL be omitted from exports.

#### Scenario: Export nested form value
- **WHEN** `RawValue()` is called on a root group containing a control, a nested group, and an array
- **THEN** the result is a `map[string]any` with the control's value, the nested group's map, and the array's slice in order

#### Scenario: Disabled nodes omitted from export
- **WHEN** a control is disabled and `RawValue()` is called on its parent
- **THEN** the disabled control's key is absent from the exported map

### Requirement: Path resolver
Nodes SHALL resolve paths of the form `a.b.c` (map children) and `items[0].name` (array indices), including mixed nesting. Resolution SHALL return the node or a miss indicator.

#### Scenario: Resolve dotted path through groups
- **WHEN** `Get("billing.address.city")` is called on a root group with nested groups
- **THEN** the matching control is returned

#### Scenario: Resolve bracketed path through arrays
- **WHEN** `Get("items[2].price")` is called on a root group containing an array
- **THEN** the price control of the third item is returned

#### Scenario: Path miss returns miss indicator
- **WHEN** `Get("a.b.c")` is called and `c` does not exist under `a.b`
- **THEN** the resolver returns the miss indicator

### Requirement: Change notification
Nodes SHALL expose change subscriptions (`OnValueChange`, `OnStatusChange`, `OnTouchChange`) backed by `state.SubscriptionManager`. Mutations SHALL notify the node's own listeners and propagate to ancestors so group-level listeners fire when children change.

#### Scenario: Group status listener fires on child change
- **WHEN** a `Control` inside a `Group` is set to an invalid value and the group has a status listener
- **THEN** the group's status listener is invoked with the group's new status

#### Scenario: Value listener fires on set
- **WHEN** `control.Set(v)` is called and a value listener is registered
- **THEN** the listener is invoked with `v`

### Requirement: Reset
`Reset()` SHALL restore a node's value to its initial value, clear errors, mark it pristine and untouched, and propagate to descendants.

#### Scenario: Reset restores initial state
- **WHEN** `Reset()` is called on a group whose children were modified
- **THEN** all values return to initial values, `IsDirty()` is false, and `IsTouched()` is false
### Requirement: Form-level lifecycle surface

`FormState` SHALL expose form-level operations over the tree root: `Validate() bool`, `Status()`, `Errors() map[string]string`, `Value() map[string]any`, `RawValue() any`, `IsTouched()`, `IsDirty()`, `IsPristine()`, `IsEnabled()`, `MarkAllTouched()`, `MarkAllUntouched()`, `MarkAllPristine()`, `Reset()`, `SetDisabled(bool)`, `Root()`, and `Submit(onValid func(map[string]any)) bool`.

#### Scenario: Validate reports aggregate validity

- **WHEN** `FormState.Validate()` is called on a tree containing an invalid control
- **THEN** it returns false and `Errors()` returns an entry keyed by the control's dotted path

#### Scenario: Value exports the whole tree

- **WHEN** `FormState.Value()` is called
- **THEN** it returns the root group's `RawValue()` as `map[string]any`, with nested groups as maps, arrays as ordered slices, and effectively-disabled nodes omitted

#### Scenario: Submit validates, touches, and proceeds only when valid

- **WHEN** `FormState.Submit(onValid)` is called on an invalid form
- **THEN** it marks every control touched, does not invoke `onValid`, and returns false
- **WHEN** the form is valid
- **THEN** it invokes `onValid` with `Value()` and returns true

#### Scenario: Reset restores the whole form

- **WHEN** `FormState.Reset()` is called after mutation
- **THEN** all values return to initial, errors clear, and the form is pristine and untouched

### Requirement: Typed node resolution

The public package SHALL expose generic resolution helpers so callers do not perform raw type assertions: `ControlOf[T](Resolver, key) (*Control[T], bool)`, `ControlFieldOf[T](FormScope, key, content)`, and `ControlViewOf[T](Resolver, key, content)`. A `Resolver` SHALL be satisfied by a `FormScope` and by any tree node wrapped with `NodeResolver(node)`.

#### Scenario: ControlOf resolves a typed control

- **WHEN** `ControlOf[string](resolver, "name")` is called and the node at `name` is a `Control[string]`
- **THEN** the control and true are returned
- **WHEN** the node is a different kind or does not resolve
- **THEN** nil and false are returned without panicking

#### Scenario: NodeResolver adapts a bare node

- **WHEN** `ControlOf[string](NodeResolver(root), "a.b")` is called on a tree root
- **THEN** resolution proceeds through the node's `Get`

### Requirement: Array rendering and traversal helpers

`FormScope` SHALL expose `FormArray(key, row)` to render one item per array child in order (passing the index and a row-relative scope), `ArrayLength(key)`, `Array(key)`, `Group(key)`, and `Resolve(key)`. Rows SHALL receive a scope whose base path is the array's path plus `[index]`, so nested keys resolve relative to the row.

#### Scenario: FormArray iterates every child with relative resolution

- **WHEN** `FormArray("phones", row)` is called on an array with two children
- **THEN** `row` is invoked with indices 0 and 1, and each row scope resolves `"number"` to that row's control

#### Scenario: ArrayLength reports child count

- **WHEN** `ArrayLength(key)` is called for an array, a non-array node, and a missing key
- **THEN** it returns the child count, 0, and 0 respectively

### Requirement: Optional values for empty-vs-zero

The engine SHALL support absent values through `Optional[T]`, an alias for `github.com/zodimo/go-maybe.Maybe`, so a typed control can distinguish "empty" from a legitimate zero value. The public package SHALL expose `Some`, `None`, `OptionalOf`, `RequiredOptional[T]()`, `NewOptionalControl[T]`, and `OptionalValue[T]`.

#### Scenario: Optional distinguishes empty from zero

- **WHEN** a `Control[Optional[int]]` holds `None[int]()`
- **THEN** its exported value is `None`, distinct from `Some(0)`

#### Scenario: RequiredOptional rejects only the empty state

- **WHEN** `RequiredOptional[int]()` validates `None[int]()`
- **THEN** it returns an error
- **WHEN** it validates `Some(0)`
- **THEN** it returns nil

#### Scenario: Optional is comparable for Required

- **WHEN** `Required(None[int]())` is used as a validator on a `Control[Optional[int]]`
- **THEN** it behaves equivalently to `RequiredOptional`, failing on `None` and accepting `Some(T)`

### Requirement: Coded validation errors

Validators SHALL be able to attach a machine-readable `ErrorCode` to failures via `*ValidationError`, which SHALL satisfy the `error` interface, carry `Code`, `Message`, and `Path`, and remain recoverable through `errors.Join` and `%w` wrapping. The package SHALL expose `NewValidationError`, `WithCode`, `CodeOf`, and `CodesOf`.

#### Scenario: Code survives join and wrap

- **WHEN** a `*ValidationError` with code `min_length` is joined with other errors or wrapped with `fmt.Errorf("%w")`
- **THEN** `CodeOf` returns `min_length` and true

#### Scenario: CodesOf dedups and sorts

- **WHEN** `CodesOf` is given a joined error containing duplicate and distinct codes
- **THEN** it returns each distinct code once, sorted

#### Scenario: WithCode wraps an uncoded validator

- **WHEN** `WithCode(code, v)` wraps a validator returning a bare error
- **THEN** failures carry `code` with the original message preserved, and a validator that already carries a code keeps its own

### Requirement: Structured form error access

`FormState` SHALL expose `CodedErrors() []*ValidationError` (each stamped with the failing control's path), `Codes() []ErrorCode`, and `FirstInvalidPath() (string, bool)`. `FormFieldBinding[T]` SHALL expose `ValidationErrors() []error` and `Codes() []ErrorCode`. These reflect the last validation run, matching `Errors()`.

#### Scenario: CodedErrors reports path and code per failure

- **WHEN** the form is validated with invalid controls
- **THEN** `CodedErrors()` returns one entry per validator failure, each with a non-empty Path and Code

#### Scenario: FirstInvalidPath is deterministic

- **WHEN** multiple controls fail
- **THEN** `FirstInvalidPath()` returns the lexically smallest failing dotted path — group children are stored in a map, so there is no stable insertion order to walk and "first" means deterministic lexical order (positions within a path are preserved, e.g. `phones[0].number` sorts before `phones[1].number`) — or ("", false) when valid

### Requirement: Rule catalog

`compose/foundation/form/rules` SHALL provide a catalog of ready-made validators, each a plain `fform.ValidatorFunc[T]` raising a documented `ErrorCode`: required/not-empty, length (runes and items), string shape (regex, contains/excludes, starts/ends-with, alpha/numeric/alphanumeric), formats (email, url, uuid, ip, dns, base64, hex, semver), comparables (one-of, eq/neq, gt/gte/lt/lte, in-range), collections (unique, unique-by, slice/map contains), and cross-field (equal-to, greater/less-than-field, required-when, forbidden-when, at-least-one-set, mutually-exclusive). The implementations SHALL be original to this repository and SHALL NOT copy third-party source.

#### Scenario: Rules compose with the engine

- **WHEN** multiple rules are passed to `NewControl`
- **THEN** each failure carries its own code and all codes survive the engine's errors.Join aggregation

#### Scenario: Cross-field rule reads a sibling at validation time

- **WHEN** `EqualTo(func() string { return password.Value() })` validates the confirmation field
- **THEN** it compares against the password control's current value and raises `compare_fields` on mismatch

### Requirement: Group-level validation

A `Group` SHALL accept group-level validators (`GroupValidatorFunc`, `func(*Group) error`) via `NewGroup(children, WithValidator(...))`. Group validators SHALL run on read — during `Validate()`, `Errors()`, and `Status()` — never cached, and SHALL NOT run for an effectively disabled group. Their failures SHALL be reported against the group's own path and SHALL be exposed via `Group.ValidationErrors()` so structured (coded) errors flow through `FormState.CodedErrors`.

#### Scenario: Group validator failure invalidates the group

- **WHEN** a group with a failing group validator is validated
- **THEN** `Validate()` returns false, `Status()` returns INVALID, and `Errors()` contains an entry keyed by the group's own path

#### Scenario: Group validators are skipped when disabled

- **WHEN** a group (or an ancestor) is disabled
- **THEN** its group validators do not run and the group reports DISABLED

#### Scenario: Multiple group validators aggregate

- **WHEN** a group has two failing validators
- **THEN** both failures are joined under the group's path and both codes appear in `Codes()`

#### Scenario: Group validators run on read

- **WHEN** `Status()` or `Validate()` is called
- **THEN** the group's validators execute, and no cached result is reused across calls
