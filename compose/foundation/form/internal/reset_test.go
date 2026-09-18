package formengine

import "testing"

func TestReset_RestoresGroup(t *testing.T) {
	name := NewControl(NewPlainValueStore("Jane"), "Jane")
	city := NewControl(NewPlainValueStore("Sydney"), "Sydney")
	address := NewGroup(map[string]FormNode{"city": city})
	item0 := NewGroup(map[string]FormNode{
		"name":  NewControl(NewPlainValueStore("A"), "A"),
		"price": NewControl(NewPlainValueStore(10), 10),
	})
	items := NewArray(item0)
	root := NewGroup(map[string]FormNode{"name": name, "address": address, "items": items})

	// Mutate everything.
	name.Set("John")
	name.MarkAsTouched()
	city.Set("Melbourne")
	city.MarkAsTouched()
	item0Name, _ := item0.Get("name")
	item0Name.(*Control[string]).Set("Z")

	root.Reset()

	t.Run("values restored", func(t *testing.T) {
		if name.Value() != "Jane" {
			t.Fatalf("name = %q, want %q", name.Value(), "Jane")
		}
		if city.Value() != "Sydney" {
			t.Fatalf("city = %q, want %q", city.Value(), "Sydney")
		}
		if v := item0Name.(*Control[string]).Value(); v != "A" {
			t.Fatalf("item0 name = %q, want %q", v, "A")
		}
	})

	t.Run("lifecycle restored", func(t *testing.T) {
		if root.IsDirty() {
			t.Fatal("root IsDirty() = true after Reset, want false")
		}
		if root.IsTouched() {
			t.Fatal("root IsTouched() = true after Reset, want false")
		}
		if !root.IsPristine() {
			t.Fatal("root IsPristine() = false after Reset, want true")
		}
		if name.IsDirty() || name.IsTouched() {
			t.Fatal("leaf not reset to pristine/untouched")
		}
	})
}
