## ADDED Requirements

### Requirement: Public API signatures are gioui-free
The system SHALL ensure that no exported symbol in the public package trees (`compose/`, `modifiers/`, `theme/`, `runtime/`, `pkg/`) references a type from any `gioui.org` package in its signature. A signature includes function parameters, return values, struct fields, embedded types, type aliases, and underlying types of defined types.

#### Scenario: Exported function has gioui parameter
- **WHEN** an exported function in a public package declares a parameter whose type is from a `gioui.org` package
- **THEN** the API-purity check flags the function as a violation

#### Scenario: Exported function returns gioui type
- **WHEN** an exported function in a public package returns a `gioui.org` type (e.g. `op.CallOp`)
- **THEN** the API-purity check flags the function as a violation

#### Scenario: Exported struct embeds gioui type
- **WHEN** an exported struct in a public package embeds a `gioui.org` type (e.g. `widget.Editor`, `widget.List`)
- **THEN** the API-purity check flags the struct as a violation

#### Scenario: Exported type aliases gioui type
- **WHEN** an exported type declaration in a public package aliases a `gioui.org` type (e.g. `type BasicTheme = material.Theme`)
- **THEN** the API-purity check flags the declaration as a violation

#### Scenario: Exported type uses gioui underlying type
- **WHEN** an exported defined type in a public package has a `gioui.org` underlying type (e.g. `type TextAlign gioText.Alignment`)
- **THEN** the API-purity check flags the declaration as a violation

#### Scenario: Runtime interface is gioui-free
- **WHEN** a consumer inspects the `runtime.Runtime` interface
- **THEN** no method signature references a `gioui.org` type, including return values

### Requirement: Public API types are go-compose-owned
The system SHALL expose go-compose-defined types where it previously exposed gioui aliases. These types SHALL be distinct types (not aliases) with go-compose-owned constants, and SHALL NOT be assignable from gioui types without an explicit conversion. This applies to `theme.BasicTheme`, `box.Direction`, `box.Stack`, `column.Spacing`, `column.Alignment`, `intl.Locale`, `style.TextAlign`, `style.LineBreak`, and `padding.RTL`.

#### Scenario: Direction constants are go-compose-owned
- **WHEN** a consumer uses `box.Direction` and its constants
- **THEN** the type and constants are defined by go-compose and do not require importing `gioui.org/layout`

#### Scenario: Locale is a distinct type
- **WHEN** a consumer uses `intl.Locale`
- **THEN** it is a go-compose-defined type, not an alias for `gioui.org/io/system.Locale`

#### Scenario: Text enum types are distinct
- **WHEN** a consumer uses `style.TextAlign` or `style.LineBreak`
- **THEN** these are go-compose-defined enum types with go-compose-owned constants, not gioui enum types

### Requirement: Engine conversion helpers are not public
The system SHALL NOT export functions whose purpose is converting between go-compose types and gioui types from the public packages. Conversion functions SHALL live outside the public package trees.

#### Scenario: Unit conversion helpers removed from public packages
- **WHEN** a consumer searches public packages (`compose/`, `modifiers/`, `theme/`, `runtime/`, `pkg/`) for conversion helpers such as `DpToGioUnit`, `TextUnitToGioSp`, or `AsGioSp`
- **THEN** no such exported function exists

#### Scenario: Font conversion helpers removed from public packages
- **WHEN** a consumer searches public packages for conversion helpers such as `ToGioFont`, `FromGioFont`, `ToGioWeight`, or `ToGioStyle`
- **THEN** no such exported function exists

#### Scenario: Shape interface is gioui-free
- **WHEN** a consumer implements or calls `shape.Shape` / `shape.Outline`
- **THEN** no method parameter or return value references a `gioui.org` type (e.g. `clip.Stack`, `clip.Op`, `clip.PathSpec`)
