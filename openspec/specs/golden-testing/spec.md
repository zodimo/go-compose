# Golden Testing Specification

## Purpose

Define the behavioral contract for golden regression testing through the software/recording rendering backend: an ordered, deterministic capture of draw calls that verifies the engine-abstraction refactor preserves rendering behavior, with an explicit, extensible set of covered components. (Synced from change `abstract-render-backend`.)

## Requirements

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

### Requirement: Covered components are an explicit, extensible scope

The system SHALL define the set of components covered by golden baselines explicitly (a documented "covered set"), so the scope can be reviewed and extended incrementally without renegotiating the harness. Engine-bound (Door 2) components — slider, textfield, radiobutton, progress, badge, icon, tooltip — are EXCLUDED from the initial covered set; their draw calls remain gio-emitted until a later scope expansion routes them through the backend.

#### Scenario: Covered set is documented
- **WHEN** a developer reviews the golden-test harness
- **THEN** the set of components with committed baselines is enumerated in the harness or its documentation, and each component not covered is known to be excluded by the Door-2 carve-out

#### Scenario: Covered set can grow
- **WHEN** a Door-2 component is later routed through the render.Backend emission path
- **THEN** adding it to the covered set requires only emitting through the backend and committing its baseline — no harness redesign
