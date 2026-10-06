## ADDED Requirements

### Requirement: Form composable remembers and binds the tree
The `Form()` composable SHALL remember the form tree via `RememberFormState(c)`, subscribe to the tree's change stream, and bump a remembered version `MutableValue` on any change so the existing store → `window.Invalidate()` loop triggers recomposition of the form subtree.

#### Scenario: Mutating a control recomposes the form
- **WHEN** a user types into a bound text field, mutating the underlying control
- **THEN** the version `MutableValue` is bumped and the form subtree is recomposed

#### Scenario: Form state survives recomposition
- **WHEN** the form subtree recomposes
- **THEN** the remembered tree retains its values, errors, and lifecycle state

### Requirement: FormScope binds named nodes
`FormScope.Field(key, content)` and `FormScope.Form(key, header, children)` SHALL resolve their key against the tree (through the path resolver) and pass the resolved node into the content so field components can read value/errors/touched state.

#### Scenario: Field resolves its node by path
- **WHEN** `fs.Field("billing.address.city", content)` is called inside a form
- **THEN** the content receives a binding whose control is the node at path `billing.address.city`

#### Scenario: Nested form groups resolve relative paths
- **WHEN** a nested `fs.Form("billing", header, children)` is used
- **THEN** children of that scope resolve against the `billing` subtree

### Requirement: FormFieldBinding replaces FormFieldState
The public API SHALL expose `FormFieldBinding[T]` (replacing `FormFieldState[T]`) wrapping a `*Control[T]` and a remembered `state.MutableValueTyped[T]` store adapter. It SHALL read `Value()`, `Errors()`, `IsTouched()`, `HasErrors()` from the control and write through `SetValue` (which touches the control). A deprecated `RememberFormFieldState` alias SHALL keep existing call sites working.

#### Scenario: Reading state reflects the control
- **WHEN** the control's value or errors change
- **THEN** the binding's `Value()`, `HasErrors()`, and `Error()` reflect the control state

#### Scenario: Writing through the binding drives recomposition
- **WHEN** `binding.SetValue(v)` is called
- **THEN** the control is set and touched, and the remembered store adapter triggers recomposition

### Requirement: Touched-gated error display
Field components SHALL display Material 3 errors only after the control is touched: `WithError(control.HasErrors() && control.IsTouched())` with `WithSupportingText(all errors joined via errors.Join)` conveying the messages. Untouched invalid controls SHALL NOT show error styling.

#### Scenario: Error hidden until touched
- **WHEN** a bound field has validation errors but has not been touched
- **THEN** the text field renders without error styling

#### Scenario: Error shown after touch
- **WHEN** the same field is then touched while still invalid
- **THEN** the text field renders with error styling and all error messages (joined) as supporting text

#### Scenario: Error clears on valid input
- **WHEN** a touched invalid field is corrected to a valid value
- **THEN** error styling and supporting error text are removed

### Requirement: TextFieldComponent binds to the engine
`components/textfield.go` SHALL consume a `FormFieldBinding[string]` and wire value, touch, validation, and error display through the control. The public signature and existing option names (`TextFieldWithLabel`, `TextFieldWithHintText`) SHALL remain usable by the demo. The inert `TextFieldWithValidators` option was removed: validation is wired exclusively through the engine's validators at tree-build time.

#### Scenario: Demo field uses the engine
- **WHEN** `cmd/demo/form/ui.go` renders a field via `TextFieldComponent(binding, ...)`
- **THEN** typing validates against the control's validators and touched-gated errors render per the M3 convention

### Requirement: Superseded residues removed
The prototype (coded errors + the `rules` catalog) superseded the engine's early scaffolding and the marked-but-kept dead code, so the residue SHALL be removed rather than marked: `pkg/cforms/` (`forms.go`, `errors.go`), `cmd/demo/form/form.go`, `compose/foundation/form/types.go` (`Validatable`, `Orientation`), the unwired `TextStyle`/`Orientation` option plumbing (`WithTextStyle`, `WithTextStyleOption`, `WithOrientation`), and the deprecated `RememberFormFieldState` alias SHALL NOT exist. Live options such as `WithModifier` SHALL remain.

#### Scenario: Superseded files and declarations are gone
- **WHEN** a developer looks for `pkg/cforms/`, `cmd/demo/form/form.go`, or `compose/foundation/form/types.go`
- **THEN** they do not exist, having been replaced by the form engine and its `rules` catalog

#### Scenario: Live options remain
- **WHEN** a developer opens `compose/foundation/form/options.go`
- **THEN** `WithModifier` is present and no superseded `TextStyle`/`Orientation` options remain

### Requirement: Demo title fix
`cmd/demo/form/main.go` SHALL use a window title matching the actual demo content instead of the stale "Package CForms Demo" title.

#### Scenario: Demo has accurate title
- **WHEN** the form demo launches
- **THEN** the window title reflects the new form engine demo
### Requirement: Bound field components

`compose/foundation/form/components` SHALL provide bound Material 3 field components that take a `*FormFieldBinding[T]`, read value/errors/touched state from the control, write through `SetValue`, and honor the control's enabled state: `TextFieldComponent` (`string`), `NumberFieldComponent` (`int`), `CheckboxComponent` and `SwitchComponent` (`bool`), and `SelectComponent` (`string` with `[]SelectOption`). Errors SHALL be touched-gated and rendered as supporting text.

#### Scenario: Disabled control renders a disabled widget

- **WHEN** a bound field's control is effectively disabled (own or inherited)
- **THEN** the field component renders with error/enabled styling reflecting the disabled state and its writes are ignored

#### Scenario: Number field tolerates partial input

- **WHEN** the user types a value that does not parse as an int into a `NumberFieldComponent`
- **THEN** the raw text is retained for editing, the control keeps its last parsed value, and a touched-gated error is shown

#### Scenario: Select writes the option value

- **WHEN** the user picks an option in a `SelectComponent`
- **THEN** the option's `Value` is written through the binding and the menu closes

### Requirement: FormFieldBinding read surface

`FormFieldBinding[T]` SHALL expose `Value()`, `SetValue(v)`, `Errors()`, `HasErrors()`, `IsTouched()`, `IsEnabled()`, `IsDirty()`, `Status()`, `ErrorMessage()`, and `Control()`.

#### Scenario: Binding reflects control state

- **WHEN** the control's value, errors, or enabled state change
- **THEN** the binding's corresponding accessors reflect the new state

### Requirement: Optional-capable number field

`compose/foundation/form/components` SHALL provide `OptionalNumberFieldComponent(*FormFieldBinding[Optional[int]], ...)` alongside the int `NumberFieldComponent`. The raw text buffer SHALL be the authority for emptiness: an empty field writes `None`, a parseable entry writes `Some(n)`, and a non-empty unparsable entry is a touched-gated supporting error that leaves the control's last value in place. An external control change (Reset, programmatic Set) SHALL re-seed the displayed text rather than being shadowed by a stale buffer.

#### Scenario: Clearing an optional number field yields None

- **WHEN** the user clears an `OptionalNumberFieldComponent`
- **THEN** the control's value becomes `None` and the field renders empty

#### Scenario: Zero is preserved as Some(0)

- **WHEN** the user types "0"
- **THEN** the control's value becomes `Some(0)` and the field renders "0"

#### Scenario: External reset re-seeds the field

- **WHEN** the bound control is reset to `None` while the field previously showed a number
- **THEN** the field renders empty again rather than the stale text

### Requirement: Group-level rule catalog

The `rules` package SHALL provide group-level rules that read named children and raise coded errors: `GroupAtLeastOneSet`, `GroupMutuallyExclusive`, `GroupRequiredTogether`, `GroupOrdered`, and `GroupCustom`. These SHALL be `fform.GroupValidatorFunc` values attachable with `fform.WithGroupValidator`. A length rule SHALL exist for optional string fields that passes on empty input (`MinLengthIfNotEmpty`), and `Length` SHALL pass on empty input.

#### Scenario: At-least-one reads the group's children

- **WHEN** `GroupAtLeastOneSet[string]("email", "phone")` is attached to a group and both children are empty
- **THEN** the group is invalid with a `cross_field` code; setting either child clears it

#### Scenario: Disabled children are ignored by group rules

- **WHEN** a group rule names a child that is disabled
- **THEN** that child does not count toward the rule's outcome
