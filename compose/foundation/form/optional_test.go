package fform_test

import (
	"testing"

	fform "github.com/zodimo/go-compose/compose/foundation/form"
)

func TestOptional_SomeAndNone(t *testing.T) {
	if !fform.Some(7).IsSome() {
		t.Fatal("Some(7) should be some")
	}
	if got := fform.Some(7).UnwrapOr(0); got != 7 {
		t.Fatalf("expected 7, got %d", got)
	}
	if !fform.None[int]().IsNone() {
		t.Fatal("None[int]() should be none")
	}
	// The distinction that motivates Optional: None is not Some(0).
	if fform.None[int]() == fform.Some(0) {
		t.Fatal("None[int]() must differ from Some(0)")
	}
}

func TestOptionalOf(t *testing.T) {
	if fform.OptionalOf(5, false).IsSome() {
		t.Fatal("OptionalOf(v, false) should be none")
	}
	if got := fform.OptionalOf(5, true).UnwrapOr(0); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestRequiredOptional_Validator(t *testing.T) {
	validator := fform.RequiredOptional[int]()

	if err := validator(fform.None[int]()); err == nil {
		t.Fatal("expected RequiredOptional to fail on None")
	}
	if err := validator(fform.Some(0)); err != nil {
		t.Fatalf("expected RequiredOptional to accept Some(0), got %v", err)
	}
}

func TestOptionalControl_StartsEmptyAndValidates(t *testing.T) {
	// A control holding Optional[int] distinguishes "empty" from 0.
	control := fform.NewOptionalControl(fform.NewPlainValueStore(fform.None[int]()))

	if control.Value().IsSome() {
		t.Fatal("new optional control should start as None")
	}
	if control.HasErrors() {
		t.Fatal("no validators means no errors")
	}

	// RequiredOptional catches the empty state...
	required := fform.NewControl[fform.Optional[int]](
		fform.NewPlainValueStore(fform.None[int]()),
		fform.None[int](),
		fform.RequiredOptional[int](),
	)
	if required.Validate() {
		t.Fatal("RequiredOptional should reject None")
	}
	if !required.HasErrors() {
		t.Fatal("expected errors for None under RequiredOptional")
	}

	// ...and accepts a legitimate zero once entered.
	required.Set(fform.Some(0))
	if !required.Validate() {
		t.Fatalf("Some(0) should satisfy RequiredOptional, errors=%v", required.Errors())
	}

	// Clearing it makes it invalid again.
	required.Set(fform.None[int]())
	if required.Validate() {
		t.Fatal("clearing back to None should be invalid again")
	}
}

func TestOptionalControl_RequiredViaComparableSentinelAlsoWorks(t *testing.T) {
	// Required(None[T]()) relies on Optional being comparable; document that it
	// is an equivalent spelling of RequiredOptional for callers who prefer it.
	control := fform.NewControl[fform.Optional[int]](
		fform.NewPlainValueStore(fform.None[int]()),
		fform.None[int](),
		fform.Required(fform.None[int]()),
	)
	if control.Validate() {
		t.Fatal("Required(None[int]()) should reject None")
	}
	control.Set(fform.Some(3))
	if !control.Validate() {
		t.Fatal("Required(None[int]()) should accept Some(3)")
	}
}

func TestOptionalControl_ValueExportsDistinguishEmptyFromZero(t *testing.T) {
	root := fform.NewGroup(map[string]fform.FormNode{
		"empty": fform.NewControl(fform.NewPlainValueStore(fform.None[int]()), fform.None[int]()),
		"zero":  fform.NewControl(fform.NewPlainValueStore(fform.Some(0)), fform.Some(0)),
	})

	value := root.RawValue().(map[string]any)

	empty, ok := value["empty"].(fform.Optional[int])
	if !ok {
		t.Fatalf("expected Optional[int] in export, got %T", value["empty"])
	}
	zero, ok := value["zero"].(fform.Optional[int])
	if !ok {
		t.Fatalf("expected Optional[int] in export, got %T", value["zero"])
	}
	if empty.IsSome() {
		t.Fatal("empty field should export as None")
	}
	if !zero.IsSome() || zero.UnwrapOr(-1) != 0 {
		t.Fatal("zero field should export as Some(0)")
	}
}

func TestOptionalValue_UnwrapsControl(t *testing.T) {
	control := fform.NewOptionalControl(fform.NewPlainValueStore(fform.None[int]()))
	if _, ok := fform.OptionalValue(control); ok {
		t.Fatal("empty control should unwrap to not-ok")
	}
	control.Set(fform.Some(42))
	got, ok := fform.OptionalValue(control)
	if !ok || got != 42 {
		t.Fatalf("expected (42, true), got (%d, %t)", got, ok)
	}
}
