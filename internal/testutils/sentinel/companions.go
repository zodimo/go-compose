package sentinel

import (
	"testing"

	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// === Shared conformance assertions for sentinel companions ===
//
// Every sentinel-carrying type in compose/ui/** exports exactly one value — its
// Companion — and wires the ten canonical operations through
// go-sentinel-helper/sentinel/companions. The library fixes the argument order:
//
//	NewScalarCompanion(unspecified, string, isSpecified, isUnspecified,
//	                   takeOrElse, merge, same, semanticEqual, equal, copy)                       // 10 args
//	NewStructCompanion(unspecified, string, isSpecified, isUnspecified,
//	                   takeOrElse, merge, same, semanticEqual, equal,
//	                   coalesce, copy)                                                           // 11 args, coalesce before copy
//
// `isUnspecified` is the fourth argument — not the third. Transposing
// isSpecified/isUnspecified compiles cleanly and inverts the semantics silently,
// which is why assertUnspecifiedIsNotSpecified exists below as the canonical guard.
//
// These helpers are intentionally generic over the companion interface so per-type
// tests can reuse them without importing a shared test-only package.

// assertSentinelIsNotZero records that the reserved "unspecified" value is
// distinguishable from the type's zero value. A companion whose sentinel equals
// the zero value cannot tell "unset" from "deliberately zero".
//
// zero must be a genuine zero value of the type whose zero is a real value.
// For types where the zero value IS the sentinel (e.g. packed TextUnit), pass
// wantDistinct=false and the check only verifies the companion reports it as
// unspecified.
func AssertSentinelIsNotZero[T any](t *testing.T, c companions.Companion[T], zero T, wantDistinct bool) {
	t.Helper()

	unspecified := c.Unspecified()
	if !c.IsUnspecified(unspecified) {
		t.Errorf("IsUnspecified(Unspecified()) = false, want true — isSpecified/isUnspecified are transposed")
	}
	if c.IsSpecified(unspecified) {
		t.Errorf("IsSpecified(Unspecified()) = true, want false — isSpecified/isUnspecified are transposed")
	}

	if !wantDistinct {
		return
	}
	if c.Equal(unspecified, zero) {
		t.Errorf("Equal(Unspecified(), zero) = true, want false — the sentinel must not be the zero value")
	}
}

// assertMergeDirection records both directions of the Merge contract:
//
//	Merge(base, Unspecified()) == base     (the sentinel contributes nothing)
//	Merge(base, override)      == override (specified values win)
//
// override must be a fully specified value distinct from base.
func AssertMergeDirection[T any](t *testing.T, c companions.Companion[T], base, override T) {
	t.Helper()

	if got := c.Merge(base, c.Unspecified()); !c.SemanticEqual(got, base) {
		t.Errorf("Merge(base, Unspecified()) != base — the sentinel must contribute nothing")
	}
	if got := c.Merge(c.Unspecified(), override); !c.SemanticEqual(got, override) {
		t.Errorf("Merge(Unspecified(), override) != override — a specified value must win")
	}
	if got := c.Merge(base, override); !c.SemanticEqual(got, override) {
		t.Errorf("Merge(base, override) != override — the override direction is inverted")
	}
}

// assertUnspecifiedSemanticEqual records that two independent unspecified values
// compare equal. NaN-based sentinels (Dp, Offset, ...) fail this if `same` or
// `semanticEqual` is written with `==` instead of delegating to floatutils.
func AssertUnspecifiedSemanticEqual[T any](t *testing.T, c companions.Companion[T]) {
	t.Helper()

	u1 := c.Unspecified()
	u2 := c.Unspecified()

	if !c.SemanticEqual(u1, u2) {
		t.Errorf("SemanticEqual(Unspecified(), Unspecified()) = false, want true")
	}
	if !c.Equal(u1, u2) {
		t.Errorf("Equal(Unspecified(), Unspecified()) = false, want true")
	}
	if !c.Same(u1, u2) {
		t.Errorf("Same(Unspecified(), Unspecified()) = false, want true")
	}
}

// assertCopyLeavesSource records that Copy returns a value reflecting the mutator
// while leaving the source untouched. mutate must change at least one observable
// fact; observable reports that fact for a value, so the assertion does not need
// to know the type's internals.
//
// For scalar companions pass the same signature; the variadic mutator list is
// ignored by scalar Copy.
func AssertCopyLeavesSource[T any](
	t *testing.T,
	c companions.StructCompanion[T],
	source T,
	mutate func(T),
	observable func(T) any,
) {
	t.Helper()

	before := observable(source)

	derived := c.Copy(source, mutate)
	if got := observable(derived); got == before {
		t.Errorf("Copy(source, mutate) did not apply the mutator")
	}
	if got := observable(source); got != before {
		t.Errorf("Copy mutated its source: observable changed from %v to %v", before, got)
	}
}

// assertCoalesceFoldsNil records that a struct companion maps nil onto the default.
func AssertCoalesceFoldsNil[T any](t *testing.T, c companions.StructCompanion[T], def T) {
	t.Helper()

	var nilValue T
	if got := c.Coalesce(nilValue, def); !c.Same(got, def) {
		t.Errorf("Coalesce(nil, def) != def — nil must fold onto the default")
	}
}
