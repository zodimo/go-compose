## ADDED Requirements

### Requirement: Sentinel-carrying types export exactly one value

Every sentinel-carrying type in the compose core surface SHALL export exactly one public symbol
for its sentinel family: the Companion. The reserved "unspecified" value MUST be package-private
and reachable by callers only through `XCompanion.Unspecified()`.

#### Scenario: No exported sentinel remains

- **WHEN** `go doc -all` is run for an in-scope package (e.g. `compose/ui/geometry`)
- **THEN** no exported symbol whose name ends in `Unspecified` appears in the output
- **AND** the only exported sentinel-family symbol is `XCompanion`

#### Scenario: Sentinel is reachable only through the companion

- **WHEN** a caller outside the package needs the unspecified value
- **THEN** the caller obtains it via `XCompanion.Unspecified()`
- **AND** there is no exported `XUnspecified` variable or constant to compare against

### Requirement: Each sentinel type wires the ten canonical operations

Each in-scope type SHALL implement the ten canonical operations as package-private functions and
register them via `companions.NewScalarCompanion` (value types) or
`companions.NewStructCompanion` (pointer types), matching the argument order of
`go-sentinel-helper` v0.2.0 — including `isUnspecified` as the fourth argument.

#### Scenario: Scalar type registers a scalar companion

- **WHEN** the type is a value type (e.g. `Offset`, `Dp`, `TextUnit`, `Color`, `FontWeight`, `FontStyle`)
- **THEN** it registers `XCompanion = companions.NewScalarCompanion(unspecified, string, isSpecified, isUnspecified, takeOrElse, merge, same, semanticEqual, equal, copy)` with all ten arguments in that order

#### Scenario: Pointer type registers a struct companion

- **WHEN** the type is a pointer or interface type (e.g. `Shadow`, `Brush`)
- **THEN** it registers `XCompanion = companions.NewStructCompanion(...)` with the ten arguments in the same order plus `coalesce` before `copy`
- **AND** `copy` has signature `func(T, ...func(T)) T`

#### Scenario: isSpecified and isUnspecified are not transposed

- **WHEN** the companion is constructed
- **THEN** `isSpecified(Unspecified())` returns false and `isUnspecified(Unspecified())` returns true
- **AND** `isSpecified(zeroValue)` returns true for a type whose zero value is a real value

### Requirement: Companion operations obey the specified semantics

The operations registered by each companion SHALL follow the semantics defined by the
`go-sentinel-helper` contract: `Merge` prefers the incoming specified value, `SemanticEqual`
treats two unspecified values as equal, and `Coalesce` maps nil to the default.

#### Scenario: Merge prefers the override

- **WHEN** `Merge(base, sentinel)` is called
- **THEN** the result equals `base`
- **AND WHEN** `Merge(base, override)` is called with a specified override
- **THEN** the result equals `override`

#### Scenario: Two unspecified values are semantically equal

- **WHEN** `SemanticEqual(unspecified, unspecified)` is called
- **THEN** the result is true

#### Scenario: Coalesce folds nil onto the default

- **WHEN** `Coalesce(nil, def)` is called on a struct companion
- **THEN** the result is `def`

#### Scenario: Copy does not modify the source

- **WHEN** `Copy` is called with a mutator on a specified value
- **THEN** the returned value reflects the mutator
- **AND** the original source value is unchanged

### Requirement: pkg/sentinel is removed in favour of go-sentinel-helper

The `pkg/sentinel` package SHALL be deleted, and every importer SHALL use
`github.com/zodimo/go-sentinel-helper` directly. The two implementations SHALL NOT coexist.

#### Scenario: Package no longer exists

- **WHEN** the repository is searched for imports of `github.com/zodimo/go-compose/pkg/sentinel`
- **THEN** no file references it
- **AND** the `pkg/sentinel/` directory does not exist

#### Scenario: String sentinel comes from the library companion

- **WHEN** a former `pkg/sentinel` consumer needs the unspecified string
- **THEN** it uses `stringutils.StringsCompanion.Unspecified()` instead of `sentinel.StringValueUnspecified` or the deprecated `sentinel.StringUnspecified`

#### Scenario: String fallback comes from the library companion

- **WHEN** a former `pkg/sentinel` consumer needs a fallback (`sentinel.TakeOrElseString(a, "")`)
- **THEN** it uses `stringutils.StringsCompanion.TakeOrElse(a, "")`

#### Scenario: Generic TakeOrElse is replaced by per-field companions

- **WHEN** `compose/ui/graphics/shadow.go` needs to merge its `Color`, `Offset` and `BlurRadius` fields
- **THEN** each field goes through its own type's companion `TakeOrElse` (e.g. `ColorCompanion.TakeOrElse(a, b)`) instead of the three-argument `sentinel.TakeOrElse`
- **AND** no in-scope type's raw sentinel is passed as an argument

#### Scenario: The library's own exported sentinel is not used as a comparison

- **WHEN** code needs to test whether a string value is unspecified
- **THEN** it calls `stringutils.StringsCompanion.IsSpecified` / `IsUnspecified`
- **AND** it does not compare against the library's exported `StringValueUnspecified` constant, because the compose contract requires the companion to be the single door

### Requirement: The module depends directly on go-sentinel-helper

The module SHALL list `github.com/zodimo/go-sentinel-helper` as a direct requirement, since
`compose/ui/**` and `pkg/sentinel` now import `sentinel/companions` in non-test code.

#### Scenario: Dependency is direct

- **WHEN** `go.mod` is inspected
- **THEN** `github.com/zodimo/go-sentinel-helper` appears in the primary require block without an `// indirect` comment
