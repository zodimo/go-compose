package fform_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/zodimo/go-compose/compose"
	fform "github.com/zodimo/go-compose/compose/foundation/form"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
)

func TestValidationError_ErrorAndCode(t *testing.T) {
	err := fform.NewValidationError(fform.CodeRequired, "value is required")
	if err.Error() != "value is required" {
		t.Fatalf("unexpected message: %q", err.Error())
	}
	if code, ok := fform.CodeOf(err); !ok || code != fform.CodeRequired {
		t.Fatalf("expected CodeRequired, got %q ok=%t", code, ok)
	}
}

func TestCodeOf_ThroughJoinAndWrap(t *testing.T) {
	inner := fform.NewValidationError(fform.CodeMinLength, "too short")

	// Wrapped with fmt.Errorf %w.
	wrapped := fmt.Errorf("field x: %w", inner)
	if code, ok := fform.CodeOf(wrapped); !ok || code != fform.CodeMinLength {
		t.Fatalf("expected CodeMinLength through wrap, got %q ok=%t", code, ok)
	}

	// Joined with another error.
	joined := errors.Join(errors.New("plain"), inner)
	if code, ok := fform.CodeOf(joined); !ok || code != fform.CodeMinLength {
		t.Fatalf("expected CodeMinLength through join, got %q ok=%t", code, ok)
	}

	// No code present.
	if _, ok := fform.CodeOf(errors.New("plain")); ok {
		t.Fatal("plain error should carry no code")
	}
}

func TestCodesOf_DedupsAndSorts(t *testing.T) {
	joined := errors.Join(
		fform.NewValidationError(fform.CodeRequired, "a"),
		fform.NewValidationError(fform.CodeMinLength, "b"),
		fform.NewValidationError(fform.CodeRequired, "c"),
	)
	codes := fform.CodesOf(joined)
	if len(codes) != 2 {
		t.Fatalf("expected 2 distinct codes, got %v", codes)
	}
	if codes[0] != fform.CodeMinLength || codes[1] != fform.CodeRequired {
		t.Fatalf("expected sorted [min_length required], got %v", codes)
	}
}

func TestWithCode_WrapsExistingValidator(t *testing.T) {
	// A plain validator returning a bare error.
	plain := func(s string) error {
		if s == "" {
			return errors.New("nope")
		}
		return nil
	}

	coded := fform.WithCode("custom_rule", fform.ValidatorFunc[string](plain))
	err := coded("")
	if code, ok := fform.CodeOf(err); !ok || code != "custom_rule" {
		t.Fatalf("expected custom_rule, got %q ok=%t", code, ok)
	}
	if err.Error() != "nope" {
		t.Fatalf("message should be preserved, got %q", err.Error())
	}
	if err := coded("ok"); err != nil {
		t.Fatalf("expected pass, got %v", err)
	}
}

func TestWithCode_PreservesInnerCode(t *testing.T) {
	inner := func(s string) error {
		return fform.NewValidationError(fform.CodeMinLength, "too short")
	}
	// Wrapping a coded validator must not clobber its code.
	outer := fform.WithCode("should_not_win", fform.ValidatorFunc[string](inner))
	if code, _ := fform.CodeOf(outer("")); code != fform.CodeMinLength {
		t.Fatalf("expected inner code to win, got %q", code)
	}
}

func TestRequiredEmitsCode(t *testing.T) {
	control := fform.NewControl(fform.NewPlainValueStore(""), "", fform.Required(""))
	control.Validate()
	codes := controlCodes(control)
	if len(codes) != 1 || codes[0] != fform.CodeRequired {
		t.Fatalf("expected [required], got %v", codes)
	}
}

// controlCodes reads a control's distinct codes via the public helper surface.
func controlCodes(control *fform.Control[string]) []fform.ErrorCode {
	return fform.CodesOf(errors.Join(control.ValidationErrors()...))
}

func TestFormState_CodedErrorsAndFirstInvalidPath(t *testing.T) {
	build := func(c api.Composer) *fform.Group {
		return fform.NewGroup(map[string]fform.FormNode{
			"identity": fform.NewGroup(map[string]fform.FormNode{
				"name": fform.NewControl(fform.NewPlainValueStore(""), "", fform.Required(""), fform.MinLength(3)),
			}),
			"email": fform.NewControl(fform.NewPlainValueStore(""), "", fform.Required("")),
		})
	}
	c := compose.NewComposer()
	fs := fform.RememberFormState(c, build)

	// Errors (and therefore coded errors) are populated by validation, matching
	// FormState.Errors semantics.
	if fs.Validate() {
		t.Fatal("expected the empty form to be invalid")
	}

	coded := fs.CodedErrors()
	if len(coded) == 0 {
		t.Fatal("expected coded errors for an empty invalid form")
	}
	// Every entry carries a path and a code.
	for _, ve := range coded {
		if ve.Path == "" {
			t.Fatalf("expected a path on every coded error, got %+v", ve)
		}
		if ve.Code == "" {
			t.Fatalf("expected a code on every coded error, got %+v", ve)
		}
	}

	// Deterministic: the lexically smallest failing path. With "email" and
	// "identity.name" both invalid, "email" sorts first.
	if path, ok := fs.FirstInvalidPath(); !ok || path != "email" {
		t.Fatalf("expected first invalid path email, got %q ok=%t", path, ok)
	}

	// Codes are deduped and sorted across the tree.
	codes := fs.Codes()
	if len(codes) < 1 || codes[0] == "" {
		t.Fatalf("expected non-empty codes, got %v", codes)
	}

	// Fixing everything clears codes and the first-invalid path.
	name, _ := fform.ControlOf[string](fform.NodeResolver(fs.Root()), "identity.name")
	name.Set("Ada")
	email, _ := fform.ControlOf[string](fform.NodeResolver(fs.Root()), "email")
	email.Set("ada@example.com")

	if !fs.Validate() {
		t.Fatalf("expected valid form, errors=%v", fs.Errors())
	}
	if len(fs.CodedErrors()) != 0 {
		t.Fatalf("expected no coded errors, got %v", fs.CodedErrors())
	}
	if _, ok := fs.FirstInvalidPath(); ok {
		t.Fatal("expected no first-invalid path on a valid form")
	}
}

func TestFormFieldBinding_Codes(t *testing.T) {
	control := fform.NewControl(fform.NewPlainValueStore(""), "", fform.Required(""), fform.MinLength(3))
	binding := fform.NewFormFieldBinding(control)

	if len(binding.Codes()) != 0 {
		t.Fatalf("expected no codes before validation, got %v", binding.Codes())
	}
	control.Validate()
	codes := binding.Codes()
	if len(codes) != 2 {
		t.Fatalf("expected 2 codes, got %v", codes)
	}
	// ValidationErrors are stamped with the control's path.
	for _, err := range binding.ValidationErrors() {
		code, ok := fform.CodeOf(err)
		if !ok {
			t.Fatalf("expected coded error, got %v", err)
		}
		if code == "" {
			t.Fatal("expected non-empty code")
		}
	}
}

// TestRememberFormFieldBinding_Standalone verifies the tree-free standalone path:
// the binding owns a control whose value is driven by the remembered
// MutableValueTyped, so writes flow through and validators run.
func TestRememberFormFieldBinding_Standalone(t *testing.T) {
	c := compose.NewComposer()
	itemState := state.MustRemember[string](c, "standalone-name", func() string {
		return ""
	})

	binding := fform.RememberFormFieldBinding(c, itemState)
	if binding == nil {
		t.Fatal("expected a binding")
	}
	if got := binding.Value(); got != "" {
		t.Fatalf("expected empty initial value, got %q", got)
	}

	binding.SetValue("Ada")
	if got := itemState.Get(); got != "Ada" {
		t.Fatalf("expected SetValue to write through to itemState, got %q", got)
	}
	if !binding.IsTouched() {
		t.Fatal("expected SetValue to mark the control touched")
	}
	if !binding.IsDirty() {
		t.Fatal("expected SetValue to mark the control dirty")
	}
}
