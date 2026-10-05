package fform_test

import (
	"errors"
	"testing"

	fform "github.com/zodimo/go-compose/compose/foundation/form"
)

// buildGroupForm builds a form whose "contact" group carries a group-level
// validator requiring at least one of email/phone.
func buildGroupForm(t *testing.T, validator fform.GroupValidatorFunc) *fform.FormState {
	t.Helper()
	return rememberTestFormWith(t, func(c composer) *fform.Group {
		return fform.NewGroup(map[string]fform.FormNode{
			"contact": fform.NewGroup(map[string]fform.FormNode{
				"email": fform.NewControl(fform.NewPlainValueStore(""), ""),
				"phone": fform.NewControl(fform.NewPlainValueStore(""), ""),
			}, fform.WithGroupValidator(validator)),
		})
	})
}

func TestGroupValidator_RunsOnRead(t *testing.T) {
	runs := 0
	validator := func(g *fform.Group) error {
		runs++
		return nil
	}
	fs := buildGroupForm(t, validator)

	// A clean form validates; the validator must have run.
	if !fs.Validate() {
		t.Fatalf("expected valid, errors=%v", fs.Errors())
	}
	if runs == 0 {
		t.Fatal("expected the group validator to run during Validate")
	}
}

func TestGroupValidator_FailsAndReportsGroupPath(t *testing.T) {
	atLeastOne := func(g *fform.Group) error {
		resolver := fform.NodeResolver(g)
		email, _ := fform.ControlOf[string](resolver, "email")
		phone, _ := fform.ControlOf[string](resolver, "phone")
		if email.Value() == "" && phone.Value() == "" {
			return fform.NewValidationError("at_least_one", "provide an email or phone")
		}
		return nil
	}
	fs := buildGroupForm(t, atLeastOne)

	if fs.Validate() {
		t.Fatal("expected the group to be invalid with both children empty")
	}
	// The failure is keyed by the group's own path, not a child's.
	if _, ok := fs.Errors()["contact"]; !ok {
		t.Fatalf("expected error keyed by group path 'contact', got %v", fs.Errors())
	}
	if fs.Status() != fform.StatusInvalid {
		t.Fatalf("expected INVALID status, got %s", fs.Status())
	}
	// CodedErrors stamps the group path onto the coded failure.
	coded := fs.CodedErrors()
	var found bool
	for _, ve := range coded {
		if ve.Path == "contact" && ve.Code == "at_least_one" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a coded error at path contact, got %+v", coded)
	}

	// Setting one child clears the group-level failure.
	email, _ := fform.ControlOf[string](fform.NodeResolver(fs.Root()), "contact.email")
	email.Set("ada@example.com")
	if !fs.Validate() {
		t.Fatalf("expected valid after setting email, errors=%v", fs.Errors())
	}
}

func TestGroupValidator_SkippedWhenDisabled(t *testing.T) {
	runs := 0
	validator := func(g *fform.Group) error {
		runs++
		return fform.NewValidationError("always", "should not run while disabled")
	}
	fs := buildGroupForm(t, validator)

	// Disabling the whole form disables the group; the validator must not run
	// and the form is valid.
	fs.SetDisabled(true)
	runs = 0
	if !fs.Validate() {
		t.Fatal("expected a fully disabled form to validate")
	}
	if runs != 0 {
		t.Fatalf("expected group validator not to run while disabled, ran %d times", runs)
	}
	if fs.Status() != fform.StatusDisabled {
		t.Fatalf("expected DISABLED status, got %s", fs.Status())
	}
}

func TestGroupValidator_MultipleAggregate(t *testing.T) {
	fs := rememberTestFormWith(t, func(c composer) *fform.Group {
		return fform.NewGroup(map[string]fform.FormNode{
			"pair": fform.NewGroup(map[string]fform.FormNode{
				"a": fform.NewControl(fform.NewPlainValueStore(""), ""),
			},
				fform.WithGroupValidator(func(*fform.Group) error {
					return fform.NewValidationError("rule_a", "a failed")
				}),
				fform.WithGroupValidator(func(*fform.Group) error {
					return fform.NewValidationError("rule_b", "b failed")
				}),
			),
		})
	})

	if fs.Validate() {
		t.Fatal("expected invalid")
	}
	codes := fs.Codes()
	if len(codes) != 2 {
		t.Fatalf("expected 2 distinct codes, got %v", codes)
	}
	// Both messages land on the group path, joined.
	if msg := fs.Errors()["pair"]; msg == "" {
		t.Fatalf("expected a joined message under 'pair', got %v", fs.Errors())
	}
}

func TestGroupValidator_UncodedErrorStillReported(t *testing.T) {
	fs := rememberTestFormWith(t, func(c composer) *fform.Group {
		return fform.NewGroup(map[string]fform.FormNode{
			"g": fform.NewGroup(map[string]fform.FormNode{}, fform.WithGroupValidator(func(*fform.Group) error {
				return errors.New("plain failure")
			})),
		})
	})

	if fs.Validate() {
		t.Fatal("expected invalid")
	}
	// A non-coded group failure still surfaces in CodedErrors with an empty code.
	var found bool
	for _, ve := range fs.CodedErrors() {
		if ve.Path == "g" && ve.Code == "" && ve.Message == "plain failure" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected an uncoded group error at path g, got %+v", fs.CodedErrors())
	}
}
