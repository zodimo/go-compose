## ADDED Requirements

### Requirement: Software backend captures draw calls
The system SHALL provide a software/recording rendering backend that captures every draw operation emitted by a composable tree as an ordered, inspectable list of draw calls, independent of any real window or GPU.

#### Scenario: Render produces draw-call list
- **WHEN** a composable tree is rendered through the software backend
- **THEN** the backend exposes the ordered list of draw calls emitted during the render

#### Scenario: Draw call list is inspectable
- **WHEN** a test inspects the captured draw-call list
- **THEN** each call identifies its operation and payload (e.g. rectangle, color, position) without requiring engine knowledge

### Requirement: Software backend output is deterministic
The system SHALL ensure the software backend produces identical draw-call lists for identical input (same composable tree, same constraints, same state). Iteration order SHALL be stable and independent of map ordering, wall-clock time, or memory addresses.

#### Scenario: Identical input yields identical output
- **WHEN** the same composable tree is rendered twice through the software backend with identical input
- **THEN** both render passes produce identical draw-call lists

### Requirement: Discarded recordings emit no draw calls
The system SHALL ensure that recordings which are recorded but discarded (e.g. an input-only pass that records and drops its output) produce no draw calls in the software backend's captured list.

#### Scenario: Pointer-phase recording is discarded
- **WHEN** a frame includes an input-only pass that records ops and discards them
- **THEN** the software backend's final draw-call list contains no calls from the discarded pass

### Requirement: Refactor preserves rendering behavior
The system SHALL ensure that the engine-abstraction refactor does not change observable rendering behavior: a composable rendered through the gioui backend before the refactor and the same composable rendered through the software backend after the refactor SHALL produce equivalent output as verified by golden baselines.

#### Scenario: Existing component renders identically
- **WHEN** an existing component (e.g. a Material 3 component) is rendered through the software backend after the refactor
- **THEN** its captured draw calls match the committed golden baseline for that component

#### Scenario: Golden baseline mismatch fails
- **WHEN** a component's captured draw calls differ from its golden baseline
- **THEN** the test run fails, identifying the component and the divergence

### Requirement: Golden baselines are committed
The system SHALL commit golden baselines for covered components so regressions are detectable in CI and locally via `make test`.

#### Scenario: Baseline available in repo
- **WHEN** a developer runs the test suite
- **THEN** golden baselines for covered components exist in the repository and are compared automatically
