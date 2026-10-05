package fform_test

import (
	"testing"

	"github.com/zodimo/go-compose/compose"
	fform "github.com/zodimo/go-compose/compose/foundation/form"
	"github.com/zodimo/go-compose/pkg/api"
)

// buildTestForm builds a small tree: a required name with a min length, an
// optional nickname group that can be disabled, and a phones array.
func buildTestForm(c api.Composer) *fform.Group {
	phones := fform.NewArray(
		fform.NewGroup(map[string]fform.FormNode{
			"number": fform.NewControl(fform.NewPlainValueStore(""), "", fform.Required("")),
		}),
	)

	return fform.NewGroup(map[string]fform.FormNode{
		"name": fform.NewControl(fform.NewPlainValueStore(""), "", fform.Required(""), fform.MinLength(3)),
		"opt": fform.NewGroup(map[string]fform.FormNode{
			"nickname": fform.NewControl(fform.NewPlainValueStore(""), ""),
		}),
		"phones": phones,
	})
}

func rememberTestForm(t *testing.T) *fform.FormState {
	t.Helper()
	c := compose.NewComposer()
	fs := fform.RememberFormState(c, buildTestForm)
	if fs == nil {
		t.Fatal("RememberFormState returned nil")
	}
	return fs
}

// composer is a short alias for api.Composer used by test builders.
type composer = api.Composer

// rememberTestFormWith builds and remembers a form tree with the given builder.
func rememberTestFormWith(t *testing.T, build func(c api.Composer) *fform.Group) *fform.FormState {
	t.Helper()
	c := compose.NewComposer()
	fs := fform.RememberFormState(c, build)
	if fs == nil {
		t.Fatal("RememberFormState returned nil")
	}
	return fs
}

func TestFormState_ValidateAndErrors(t *testing.T) {
	fs := rememberTestForm(t)

	if fs.Validate() {
		t.Fatal("empty form with a required field should be invalid")
	}
	errs := fs.Errors()
	if _, ok := errs["name"]; !ok {
		t.Fatalf("expected error keyed by path 'name', got %v", errs)
	}
	if fs.Status() != fform.StatusInvalid {
		t.Fatalf("expected INVALID status, got %s", fs.Status())
	}

	name, ok := fform.ControlOf[string](fform.NodeResolver(fs.Root()), "name")
	if !ok {
		t.Fatal("could not resolve name control")
	}
	name.Set("Ada")

	// The default phone row is required too; fill it so the form is valid.
	phone, ok := fform.ControlOf[string](fform.NodeResolver(fs.Root()), "phones[0].number")
	if !ok {
		t.Fatal("could not resolve phone control")
	}
	phone.Set("555-0100")

	if !fs.Validate() {
		t.Fatalf("expected valid form, errors=%v", fs.Errors())
	}
	if fs.Status() != fform.StatusValid {
		t.Fatalf("expected VALID status, got %s", fs.Status())
	}
}

func TestFormState_ValueExportsTree(t *testing.T) {
	fs := rememberTestForm(t)

	name, _ := fform.ControlOf[string](fform.NodeResolver(fs.Root()), "name")
	nick, _ := fform.ControlOf[string](fform.NodeResolver(fs.Root()), "opt.nickname")
	name.Set("Ada")
	nick.Set("A")

	value := fs.Value()
	if got := value["name"]; got != "Ada" {
		t.Fatalf("expected name=Ada in value, got %v", got)
	}
	opt, ok := value["opt"].(map[string]any)
	if !ok {
		t.Fatalf("expected opt to export as map, got %T", value["opt"])
	}
	if got := opt["nickname"]; got != "A" {
		t.Fatalf("expected opt.nickname=A, got %v", got)
	}
	phones, ok := value["phones"].([]any)
	if !ok {
		t.Fatalf("expected phones to export as slice, got %T", value["phones"])
	}
	if len(phones) != 1 {
		t.Fatalf("expected 1 phone row, got %d", len(phones))
	}
}

func TestFormState_Reset(t *testing.T) {
	fs := rememberTestForm(t)

	name, _ := fform.ControlOf[string](fform.NodeResolver(fs.Root()), "name")
	name.Set("Ada")
	fs.MarkAllTouched()

	if !fs.IsDirty() {
		t.Fatal("expected form dirty after Set")
	}
	if !fs.IsTouched() {
		t.Fatal("expected form touched after MarkAllTouched")
	}

	fs.Reset()

	if fs.IsDirty() {
		t.Fatal("expected form pristine after Reset")
	}
	if fs.IsTouched() {
		t.Fatal("expected form untouched after Reset")
	}
	if got := name.Value(); got != "" {
		t.Fatalf("expected name reset to initial empty, got %q", got)
	}
	if len(fs.Errors()) != 0 {
		t.Fatalf("expected errors cleared after Reset, got %v", fs.Errors())
	}
}

func TestFormState_SubmitOnlyProceedsWhenValid(t *testing.T) {
	fs := rememberTestForm(t)

	called := 0
	if fs.Submit(func(map[string]any) { called++ }) {
		t.Fatal("Submit should report false for an invalid form")
	}
	if called != 0 {
		t.Fatal("onValid must not be called for an invalid form")
	}
	// Submit must have marked controls touched so errors render.
	if !fs.IsTouched() {
		t.Fatal("Submit should mark all controls touched even when invalid")
	}

	name, _ := fform.ControlOf[string](fform.NodeResolver(fs.Root()), "name")
	name.Set("Ada")
	phone, _ := fform.ControlOf[string](fform.NodeResolver(fs.Root()), "phones[0].number")
	phone.Set("555-0100")

	var got map[string]any
	if !fs.Submit(func(value map[string]any) { got = value }) {
		t.Fatalf("Submit should succeed, errors=%v", fs.Errors())
	}
	if called != 0 {
		t.Fatal("count incremented unexpectedly")
	}
	if got == nil || got["name"] != "Ada" {
		t.Fatalf("onValid should receive the serialized value, got %v", got)
	}
}

func TestFormState_SetDisabledOmitsFromValue(t *testing.T) {
	fs := rememberTestForm(t)

	nick, _ := fform.ControlOf[string](fform.NodeResolver(fs.Root()), "opt.nickname")
	nick.Set("A")
	if _, ok := fs.Value()["opt"]; !ok {
		t.Fatal("expected opt present while enabled")
	}

	fs.SetDisabled(true)
	if _, ok := fs.Value()["opt"]; ok {
		t.Fatal("expected disabled subtree omitted from value")
	}
	if fs.Status() != fform.StatusDisabled {
		t.Fatalf("expected DISABLED status, got %s", fs.Status())
	}

	fs.SetDisabled(false)
	if _, ok := fs.Value()["opt"]; !ok {
		t.Fatal("expected opt restored after re-enable")
	}
}

func TestFormScope_ArrayLengthAndResolution(t *testing.T) {
	fs := rememberTestForm(t)

	// A minimal scope stand-in would require a lazy list; the tree itself is the
	// resolution surface the helpers build on, so exercise ResolvePath directly.
	if _, ok := fform.ResolvePath(fs.Root(), "phones[0].number"); !ok {
		t.Fatal("expected phones[0].number to resolve")
	}
	if _, ok := fform.ResolvePath(fs.Root(), "phones[1].number"); ok {
		t.Fatal("expected phones[1].number to miss")
	}

	array, ok := fform.ResolvePath(fs.Root(), "phones")
	if !ok {
		t.Fatal("expected phones to resolve")
	}
	arr, ok := array.(*fform.Array)
	if !ok {
		t.Fatalf("expected *Array, got %T", array)
	}
	if arr.Length() != 1 {
		t.Fatalf("expected length 1, got %d", arr.Length())
	}

	arr.Add(fform.NewGroup(map[string]fform.FormNode{
		"number": fform.NewControl(fform.NewPlainValueStore(""), "", fform.Required("")),
	}))
	if arr.Length() != 2 {
		t.Fatalf("expected length 2 after Add, got %d", arr.Length())
	}
	if _, ok := fform.ResolvePath(fs.Root(), "phones[1].number"); !ok {
		t.Fatal("expected phones[1].number to resolve after Add")
	}

	arr.Remove(0)
	if arr.Length() != 1 {
		t.Fatalf("expected length 1 after Remove, got %d", arr.Length())
	}
	if _, ok := fform.ResolvePath(fs.Root(), "phones[1].number"); ok {
		t.Fatal("expected phones[1].number to miss after Remove")
	}
}

func TestWalkVisitsEveryNode(t *testing.T) {
	fs := rememberTestForm(t)

	count := 0
	fform.Walk(fs.Root(), func(fform.FormNode) { count++ })

	// root + name + opt + nickname + phones + phones[0] + phones[0].number
	if count != 7 {
		t.Fatalf("expected 7 nodes, got %d", count)
	}
}
