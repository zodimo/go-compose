package main

import (
	"testing"

	"github.com/zodimo/go-compose/compose"
	fform "github.com/zodimo/go-compose/compose/foundation/form"
)

// TestBuildFormPathsResolve guards against drift between the tree built by
// buildForm and the paths formContent binds. Every key referenced by the demo UI
// must resolve to a node of the expected type, otherwise the field silently
// renders nothing.
func TestBuildFormPathsResolve(t *testing.T) {
	c := compose.NewComposer()
	root := buildForm(c)
	resolver := fform.NodeResolver(root)

	stringPaths := []string{"identity.name", "tier", "billing.street", "billing.city", "contact.email", "contact.phone"}
	for _, path := range stringPaths {
		if _, ok := fform.ControlOf[string](resolver, path); !ok {
			t.Errorf("expected %q to resolve to Control[string]", path)
		}
	}

	if _, ok := fform.ControlOf[fform.Optional[int]](resolver, "identity.age"); !ok {
		t.Error("expected identity.age to resolve to Control[Optional[int]]")
	}
	if _, ok := fform.ControlOf[bool](resolver, "sameAsShipping"); !ok {
		t.Error("expected sameAsShipping to resolve to Control[bool]")
	}

	phones, ok := root.Get("phones")
	if !ok {
		t.Fatal("expected phones to resolve")
	}
	arr, ok := phones.(*fform.Array)
	if !ok {
		t.Fatalf("expected phones to be *Array, got %T", phones)
	}
	if arr.Length() != 1 {
		t.Fatalf("expected 1 initial phone row, got %d", arr.Length())
	}
	if _, ok := fform.ControlOf[string](resolver, "phones[0].number"); !ok {
		t.Error("expected phones[0].number to resolve to Control[string]")
	}
}

// TestSameAsShippingDisablesBilling confirms the cross-field wiring added in
// buildForm disables the billing subtree when the checkbox flips.
func TestSameAsShippingDisablesBilling(t *testing.T) {
	c := compose.NewComposer()
	root := buildForm(c)
	resolver := fform.NodeResolver(root)

	same, ok := fform.ControlOf[bool](resolver, "sameAsShipping")
	if !ok {
		t.Fatal("expected sameAsShipping to resolve")
	}
	street, ok := fform.ControlOf[string](resolver, "billing.street")
	if !ok {
		t.Fatal("expected billing.street to resolve")
	}

	if !street.IsEnabled() {
		t.Fatal("expected billing.street enabled initially")
	}

	same.Set(true)
	if street.IsEnabled() {
		t.Fatal("expected billing.street disabled after sameAsShipping=true")
	}

	same.Set(false)
	if !street.IsEnabled() {
		t.Fatal("expected billing.street re-enabled after sameAsShipping=false")
	}
}

// TestSubmitCollectsValidValue exercises the full validate -> submit path with
// the demo's tree.
func TestSubmitCollectsValidValue(t *testing.T) {
	c := compose.NewComposer()
	formState := fform.RememberFormState(c, buildForm)

	set := func(path, value string) {
		ctl, ok := fform.ControlOf[string](fform.NodeResolver(formState.Root()), path)
		if !ok {
			t.Fatalf("expected %q to resolve", path)
		}
		ctl.Set(value)
	}
	set("identity.name", "Ada")
	set("contact.email", "ada@example.com")
	set("identity.password", "s3cret-pass")
	set("identity.confirm", "s3cret-pass")
	set("tier", "pro")
	set("billing.street", "42 Wallaby Way")
	set("billing.city", "Sydney")
	set("phones[0].number", "5550100123")

	age, ok := fform.ControlOf[fform.Optional[int]](fform.NodeResolver(formState.Root()), "identity.age")
	if !ok {
		t.Fatal("expected identity.age to resolve")
	}
	age.Set(fform.Some(36))

	var got map[string]any
	if !formState.Submit(func(value map[string]any) { got = value }) {
		t.Fatalf("expected submit to succeed, errors=%v", formState.Errors())
	}
	if got["tier"] != "pro" {
		t.Errorf("expected tier=pro, got %v", got["tier"])
	}
	identity, _ := got["identity"].(map[string]any)
	if identity["name"] != "Ada" {
		t.Errorf("expected identity.name=Ada, got %v", identity["name"])
	}
}

// TestContactGroupValidator exercises the demo's group-level validator: the
// contact group is invalid until at least one of email/phone is provided, and
// the failure is reported under the "contact" path.
func TestContactGroupValidator(t *testing.T) {
	c := compose.NewComposer()
	formState := fform.RememberFormState(c, buildForm)

	if !formState.Validate() {
		// Expected: other fields are empty too, but confirm the group error.
	}
	if _, ok := formState.Errors()["contact"]; !ok {
		t.Fatalf("expected a group error under 'contact', got %v", formState.Errors())
	}

	email, ok := fform.ControlOf[string](fform.NodeResolver(formState.Root()), "contact.email")
	if !ok {
		t.Fatal("expected contact.email to resolve")
	}
	email.Set("ada@example.com")

	// The group-level failure clears once a member is set.
	if _, ok := formState.Errors()["contact"]; ok {
		t.Fatalf("expected the contact group error to clear, got %v", formState.Errors())
	}
}

// TestContactGroupViaPhone shows the OR: setting phone instead of email also
// satisfies the group rule.
func TestContactGroupViaPhone(t *testing.T) {
	c := compose.NewComposer()
	formState := fform.RememberFormState(c, buildForm)

	phone, ok := fform.ControlOf[string](fform.NodeResolver(formState.Root()), "contact.phone")
	if !ok {
		t.Fatal("expected contact.phone to resolve")
	}
	phone.Set("5550100123")
	formState.Validate()

	if _, ok := formState.Errors()["contact"]; ok {
		t.Fatalf("expected the contact group error to clear via phone, got %v", formState.Errors())
	}
}
