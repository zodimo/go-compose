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
Field components SHALL display Material 3 errors only after the control is touched: `WithError(control.HasErrors() && control.IsTouched())` with `WithSupportingText(firstError)` conveying the message. Untouched invalid controls SHALL NOT show error styling.

#### Scenario: Error hidden until touched
- **WHEN** a bound field has validation errors but has not been touched
- **THEN** the text field renders without error styling

#### Scenario: Error shown after touch
- **WHEN** the same field is then touched while still invalid
- **THEN** the text field renders with error styling and the first error message as supporting text

#### Scenario: Error clears on valid input
- **WHEN** a touched invalid field is corrected to a valid value
- **THEN** error styling and supporting error text are removed

### Requirement: TextFieldComponent binds to the engine
`components/textfield.go` SHALL consume a `FormFieldBinding[string]` and wire value, touch, validation, and error display through the control. The public signature and existing option names (`TextFieldWithLabel`, `TextFieldWithHintText`, `TextFieldWithValidators`) SHALL remain usable by the demo.

#### Scenario: Demo field uses the engine
- **WHEN** `cmd/demo/form/ui.go` renders a field via `TextFieldComponent(binding, ...)`
- **THEN** typing validates against the control's validators and touched-gated errors render per the M3 convention

### Requirement: Dead code markers
The fully-dead files `pkg/cforms/forms.go`, `pkg/cforms/errors.go`, and `cmd/demo/form/form.go` SHALL carry a header `// DEAD CODE` comment referencing the replacement. The dead declarations `Validatable`, `Orientation`, `WithTextStyle`, `WithTextStyleOption`, `WithOrientation`, and the dead `TextFieldWith*` options SHALL each carry a `// DEAD` marker. No dead code SHALL be deleted by this change.

#### Scenario: Dead files are marked
- **WHEN** a developer opens any fully-dead file listed above
- **THEN** a header comment explains the file is dead and what replaces it

#### Scenario: Live options remain unmarked
- **WHEN** a developer opens `compose/foundation/form/options.go`
- **THEN** `WithModifier` carries no dead marker while the dead options carry `// DEAD` markers

### Requirement: Demo title fix
`cmd/demo/form/main.go` SHALL use a window title matching the actual demo content instead of the stale "Package CForms Demo" title.

#### Scenario: Demo has accurate title
- **WHEN** the form demo launches
- **THEN** the window title reflects the new form engine demo