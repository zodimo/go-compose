# Spec Delta

## Purpose

Defines how Material Symbol icons (`SymbolName` sources) render: which variable-font axes and
style a caller can select, how optical size adapts to icon size, how fill and stroke are painted,
and how converted geometry is cached across frames.

## ADDED Requirements

### Requirement: Icons honor variable-font axes

A symbol icon SHALL render at the variable-font instance selected by the caller's axis options
rather than always at the font's default instance. The selectable axes are fill (`0`–`1`), grade
(`-50`–`200`), optical size (`20`–`48`) and weight (`100`–`700`).

#### Scenario: Filled icon differs from outlined
- **WHEN** the same symbol is rendered once with fill `0` and once with fill `1`
- **THEN** the two renderings use different geometry, and both are non-empty

#### Scenario: Weight changes stroke thickness
- **WHEN** the same symbol is rendered at weight `100` and at weight `700`
- **THEN** the rendering at `700` has strictly greater outline bounds than the one at `100`

#### Scenario: Unset axes use the font default
- **WHEN** an icon is rendered with no axis options
- **THEN** it renders at the font's default instance for every axis

### Requirement: Out-of-range axis values are clamped

Each axis option SHALL reject or clamp values outside its supported range so that rendering never
silently produces a fallback instance. A clamped value SHALL still render a valid icon.

#### Scenario: Grade above maximum
- **WHEN** grade `500` is requested, above the `200` maximum
- **THEN** the icon renders at grade `200`

#### Scenario: Negative fill
- **WHEN** fill `-1` is requested
- **THEN** the icon renders at fill `0`

### Requirement: Icons select a symbol style

A symbol icon SHALL render in one of three styles — outlined, rounded or sharp — with outlined as
the default. Selecting a style SHALL change the rendered geometry while preserving color, size,
axis and fill behavior.

#### Scenario: Style selection changes geometry
- **WHEN** the same symbol is rendered as outlined and as rounded at identical axes
- **THEN** the two renderings differ in geometry but match in size and color

#### Scenario: An unavailable style fails visibly
- **WHEN** a style is requested for which no font asset is available
- **THEN** rendering reports a diagnosable error or falls back to outlined, and does not render a
  blank or missing glyph silently

### Requirement: Optical size adapts to rendered icon size

When optical size is not set explicitly, the renderer SHALL derive it from the icon's resolved
size, clamped to the axis range, so that small icons get the small-size design and large icons get
the large-size design.

#### Scenario: Small icon selects small optical size
- **WHEN** an icon renders at 20dp with no explicit optical size
- **THEN** the effective optical size is at the axis minimum

#### Scenario: Large icon selects large optical size
- **WHEN** an icon renders at 48dp or larger with no explicit optical size
- **THEN** the effective optical size is at the axis maximum

#### Scenario: Explicit optical size wins
- **WHEN** an icon renders at 48dp with optical size explicitly set to `24`
- **THEN** the effective optical size is `24`

### Requirement: Icons fill with color, brush or stroke

A symbol icon SHALL be paintable with a solid color, a gradient brush, or a stroked outline, and
these paints SHALL combine with axis and style selection rather than overriding them.

#### Scenario: Gradient fill
- **WHEN** an icon is rendered with a linear-gradient brush spanning its bounds
- **THEN** the icon's geometry is filled by that gradient instead of a solid color

#### Scenario: Stroked icon
- **WHEN** an icon is rendered with a stroke width greater than zero
- **THEN** only the outline is painted at that width, and the icon's reported size is unchanged

#### Scenario: Default paint is the inherited content color
- **WHEN** an icon is rendered with no color or brush specified
- **THEN** it is painted with the ambient content color

### Requirement: Icon geometry is cached across frames

The renderer SHALL cache converted glyph geometry keyed by symbol, style, effective axes and
optically-scaled size, so that repeated layouts of the same icon reuse geometry instead of
re-extracting outlines on every frame. The cache SHALL be bounded and SHALL NOT affect rendered
output.

#### Scenario: Repeated layout reuses geometry
- **WHEN** the same icon is laid out in identical configuration on consecutive frames
- **THEN** outline extraction and path conversion occur once, not once per frame

#### Scenario: Distinct configurations are cached separately
- **WHEN** the same symbol is laid out at two different weights
- **THEN** each configuration renders its own geometry rather than reusing the other's

#### Scenario: Cache growth is bounded
- **WHEN** more distinct configurations are rendered than the cache holds
- **THEN** least-recently-used entries are evicted and rendering remains correct
