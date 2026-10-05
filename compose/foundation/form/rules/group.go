package rules

import (
	"cmp"
	"fmt"

	fform "github.com/zodimo/go-compose/compose/foundation/form"
)

// Group-level rules. Unlike the control rules, these validate a Group as a whole
// by reading named children, so the failure is reported against the group's own
// path rather than an arbitrary member field. Attach them with
// fform.WithGroupValidator:
//
//	fform.NewGroup(map[string]fform.FormNode{
//	    "email": email,
//	    "phone": phone,
//	}, fform.WithGroupValidator(rules.GroupAtLeastOneSet[string]("email", "phone")))

// GroupAtLeastOneSet returns a group validator that fails when every named child
// control holds its zero value. It raises CodeCrossField.
//
// Keys are resolved relative to the group (dotted paths are allowed for nested
// children). Keys that do not resolve, or resolve to a non-control, are skipped;
// if no key resolves to a control the rule passes.
func GroupAtLeastOneSet[T comparable](keys ...string) fform.GroupValidatorFunc {
	return func(g *fform.Group) error {
		resolver := fform.NodeResolver(g)
		var zero T
		checked := 0
		for _, key := range keys {
			control, ok := fform.ControlOf[T](resolver, key)
			if !ok || !control.IsEnabled() {
				continue
			}
			checked++
			if control.Value() != zero {
				return nil
			}
		}
		if checked == 0 {
			return nil
		}
		return fail(CodeCrossField, "at least one of %v must be set", keys)
	}
}

// GroupMutuallyExclusive returns a group validator that fails when more than one
// named child control holds a non-zero value. It raises CodeCrossField.
//
// Disabled children are ignored (they cannot participate), and keys that do not
// resolve are skipped.
func GroupMutuallyExclusive[T comparable](keys ...string) fform.GroupValidatorFunc {
	return func(g *fform.Group) error {
		resolver := fform.NodeResolver(g)
		var zero T
		set := 0
		for _, key := range keys {
			control, ok := fform.ControlOf[T](resolver, key)
			if !ok || !control.IsEnabled() {
				continue
			}
			if control.Value() != zero {
				set++
			}
		}
		if set > 1 {
			return fail(CodeCrossField, "only one of %v may be set", keys)
		}
		return nil
	}
}

// GroupRequiredTogether returns a group validator that fails when some but not
// all of the named child controls are set (all-or-nothing). It raises
// CodeCrossField. Disabled children are ignored.
func GroupRequiredTogether[T comparable](keys ...string) fform.GroupValidatorFunc {
	return func(g *fform.Group) error {
		resolver := fform.NodeResolver(g)
		var zero T
		set, checked := 0, 0
		for _, key := range keys {
			control, ok := fform.ControlOf[T](resolver, key)
			if !ok || !control.IsEnabled() {
				continue
			}
			checked++
			if control.Value() != zero {
				set++
			}
		}
		if set == 0 || set == checked {
			return nil
		}
		return fail(CodeCrossField, "either all or none of %v must be set", keys)
	}
}

// GroupOrdered returns a group validator that fails when the value of the control
// named by key is not strictly greater than that of the control named by after.
// It raises CodeCompareFields. Resolves both keys relative to the group;
// unresolvable or disabled controls make the rule pass.
func GroupOrdered[T cmp.Ordered](key, after string) fform.GroupValidatorFunc {
	return func(g *fform.Group) error {
		resolver := fform.NodeResolver(g)
		a, ok := fform.ControlOf[T](resolver, key)
		if !ok || !a.IsEnabled() {
			return nil
		}
		b, ok := fform.ControlOf[T](resolver, after)
		if !ok || !b.IsEnabled() {
			return nil
		}
		if a.Value() <= b.Value() {
			return fail(CodeCompareFields, "%q must be greater than %q", key, after)
		}
		return nil
	}
}

// GroupCustom builds a group validator from an arbitrary predicate over the
// group, attaching code and a message. It is the escape hatch for one-off group
// logic that the named rules do not cover.
func GroupCustom(code Code, ok func(g *fform.Group) bool, format string, args ...any) fform.GroupValidatorFunc {
	message := fmt.Sprintf(format, args...)
	return func(g *fform.Group) error {
		if !ok(g) {
			return fail(code, "%s", message)
		}
		return nil
	}
}
