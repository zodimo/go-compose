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