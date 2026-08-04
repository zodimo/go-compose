## ADDED Requirements

### Requirement: Seam packages are the only engine-type importers
The system SHALL confine direct imports of engine types (`gioui.org/layout`, `gioui.org/op`, `gioui.org/widget`, `gioui.org/text`, `gioui.org/font`) to the seam packages: `internal/layoutnode`, `internal/render`, and their subpackages. Non-seam packages SHALL NOT directly import these engine packages.

#### Scenario: Non-seam package imports engine types
- **WHEN** any package outside the seam whitelist imports `gioui.org/layout`, `gioui.org/op`, `gioui.org/widget`, `gioui.org/text`, or `gioui.org/font`
- **THEN** the API-purity check flags the package as a violation

#### Scenario: Seam package imports engine types
- **WHEN** a seam package (`internal/layoutnode`, `internal/render`, or a subpackage) imports engine types
- **THEN** the API-purity check permits it without violation

### Requirement: Framework context is engine-agnostic
The system SHALL represent the layout context used by components and modifiers as a framework-owned type. The framework-owned context SHALL NOT expose engine-specific fields or methods to consumers; any engine access SHALL occur through the seam.

#### Scenario: Component receives framework context
- **WHEN** a component's `Layout` or `Draw` method receives the layout context
- **THEN** the context is a framework-owned type with no engine-typed fields accessible to the component

#### Scenario: Engine access requires seam conversion
- **WHEN** framework code requires access to engine services (e.g. text shaping)
- **THEN** it does so through a conversion or accessor defined at the seam, and this conversion is not part of the public API

### Requirement: Conversion to engine types is confined to the seam
The system SHALL ensure that conversion between framework types and engine types is defined and invoked only within seam packages or engine-bound implementations. Non-engine-bound framework code SHALL NOT perform engine conversions.

#### Scenario: Framework component does not convert directly
- **WHEN** a pure-composition component renders
- **THEN** it produces no engine conversion calls; all engine interaction is handled by the seam or an engine-bound component

#### Scenario: Engine-bound component is explicitly classified
- **WHEN** a component's implementation requires direct engine interaction
- **THEN** it is classified as engine-bound in code, its public API remains gioui-free, and its engine usage is internal

### Requirement: Backends are registrable
The system SHALL support registering rendering backends by name, with at least two implementations: a gioui backend (production) and a software/recording backend (testing). Selection of the active backend SHALL NOT require changes to framework or component code.

#### Scenario: Register a backend
- **WHEN** a backend registers itself under a name
- **THEN** it can be selected as the active backend without modifying framework or component code

#### Scenario: Render composable tree on either backend
- **WHEN** the same composable tree is rendered with the gioui backend and then with the software backend
- **THEN** both runs complete without framework changes and the software backend produces a captured draw-call list

### Requirement: Runtime interface is engine-agnostic
The system SHALL define the runtime interface such that its methods reference only framework-owned types, including return values.

#### Scenario: Runtime returns opaque draw command
- **WHEN** a consumer calls the runtime to render a composable
- **THEN** the returned draw command is an opaque framework-owned value that the active backend can apply, and no `gioui.org` type is visible in the runtime interface signature
