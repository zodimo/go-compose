# API Purity Analyzer Specification

## Purpose

Define the behavioral contract for the API-purity analyzer: a tool that inspects exported symbols of the public package trees (`compose/`, `modifiers/`, `theme/`, `runtime/`, `pkg/`) and detects references to `gioui.org` types in their signatures, so the framework's public API stays engine-agnostic. (Synced from change `abstract-render-backend`.)

## Requirements

### Requirement: Analyzer flags leaked signatures

The system SHALL provide an API-purity analyzer that inspects exported symbols of the public package trees (`compose/`, `modifiers/`, `theme/`, `runtime/`, `pkg/`) and detects references to `gioui.org` types in their signatures. A signature includes parameters, returns, fields, embedded types, aliases, and underlying types.

#### Scenario: Analyzer detects a leaked parameter type
- **WHEN** the analyzer runs against a package containing an exported function with a `gioui.org` parameter type
- **THEN** it reports a violation for that function

#### Scenario: Analyzer detects a leaked return type
- **WHEN** the analyzer runs against a package containing an exported function returning a `gioui.org` type
- **THEN** it reports a violation for that function

#### Scenario: Analyzer detects leaked struct fields and embeddings
- **WHEN** the analyzer runs against a package containing an exported struct with a `gioui.org` field or embedded type
- **THEN** it reports a violation for that struct

#### Scenario: Analyzer detects gioui underlying types
- **WHEN** the analyzer runs against a package containing an exported defined type whose underlying type is from `gioui.org`
- **THEN** it reports a violation for that type

### Requirement: Analyzer reports offending symbols

The system SHALL ensure the analyzer reports each violation with the package, the exported symbol name, and the offending `gioui.org` type, so an implementer can locate and fix it without further investigation.

#### Scenario: Violation report is actionable
- **WHEN** the analyzer reports a violation
- **THEN** the report names the package, the exported symbol, and the gioui type that leaks

### Requirement: Analyzer fails the build on violations

The system SHALL exit non-zero when any violation is detected, and exit zero when the public API is clean.

#### Scenario: Clean API passes
- **WHEN** the analyzer runs against a codebase with no leaked signatures
- **THEN** it exits with status zero

#### Scenario: Leaked API fails
- **WHEN** the analyzer runs against a codebase containing at least one leaked signature
- **THEN** it exits with a non-zero status

### Requirement: Seam packages are whitelisted

The system SHALL whitelist the seam packages (`internal/layoutnode`, `internal/render`, and subpackages) from the analyzer's import restriction, while still applying the exported-signature rule to all other packages.

#### Scenario: Whitelisted package imports engine types
- **WHEN** a whitelisted seam package imports engine types
- **THEN** the analyzer does not flag it

#### Scenario: Non-whitelisted internal package is checked
- **WHEN** a non-whitelisted package (including other `internal/` packages) exposes a `gioui.org` type in an exported signature
- **THEN** the analyzer flags it

### Requirement: Analyzer runs via make check-api

The system SHALL expose the analyzer through a `make check-api` target so it can be run locally and in CI without special setup.

#### Scenario: Invoke via make target
- **WHEN** a developer runs `make check-api`
- **THEN** the analyzer executes against the public package trees and reports any violations with the documented exit behavior

### Requirement: Known leaks are regression fixtures

The system SHALL treat the known gioui leaks present before this change as test fixtures: the analyzer SHALL flag every one of them before the cleanup, and SHALL flag none of them after the cleanup.

#### Scenario: Fixture coverage before cleanup
- **WHEN** the analyzer runs against the pre-cleanup codebase
- **THEN** it flags each of the known leak sites captured as fixtures

#### Scenario: Fixture coverage after cleanup
- **WHEN** the analyzer runs against the post-cleanup codebase
- **THEN** none of the former leak sites are flagged
