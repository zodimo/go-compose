package formengine

import "testing"

func TestArray_Structure(t *testing.T) {
	newItem := func(v string) FormNode {
		return NewControl(NewPlainValueStore(v), v)
	}

	t.Run("Add appends and wires parent", func(t *testing.T) {
		a := NewArray()
		item := newItem("a")
		a.Add(item)
		if a.Length() != 1 {
			t.Fatalf("Length() = %d, want 1", a.Length())
		}
		if item.Parent() != a {
			t.Fatal("added child Parent() != array")
		}
	})

	t.Run("Insert shifts indices", func(t *testing.T) {
		a := NewArray(newItem("a"), newItem("c"))
		a.Insert(1, newItem("b"))
		if a.Length() != 3 {
			t.Fatalf("Length() = %d, want 3", a.Length())
		}
		got, ok := a.At(1)
		if !ok {
			t.Fatal("At(1) = miss, want node")
		}
		if got.RawValue() != "b" {
			t.Fatalf("At(1) RawValue = %v, want %q", got.RawValue(), "b")
		}
	})

	t.Run("Remove shifts indices and detaches parent", func(t *testing.T) {
		a := NewArray(newItem("a"), newItem("b"), newItem("c"))
		removed, _ := a.At(1)
		a.Remove(1)
		if a.Length() != 2 {
			t.Fatalf("Length() = %d, want 2", a.Length())
		}
		got, _ := a.At(1)
		if got.RawValue() != "c" {
			t.Fatalf("At(1) RawValue = %v, want %q", got.RawValue(), "c")
		}
		if removed.Parent() != nil {
			t.Fatal("removed child still has a parent")
		}
	})

	t.Run("At out of range returns miss", func(t *testing.T) {
		a := NewArray(newItem("a"))
		if _, ok := a.At(5); ok {
			t.Fatal("At(5) = hit, want miss")
		}
		if _, ok := a.At(-1); ok {
			t.Fatal("At(-1) = hit, want miss")
		}
	})
}

func TestArray_RawValueExport(t *testing.T) {
	a := NewArray(
		NewControl(NewPlainValueStore(1), 1),
		NewControl(NewPlainValueStore(2), 2),
		NewControl(NewPlainValueStore(3), 3),
	)

	t.Run("ordered slice", func(t *testing.T) {
		raw := a.RawValue().([]any)
		if len(raw) != 3 {
			t.Fatalf("len(raw) = %d, want 3", len(raw))
		}
		for i, want := range []int{1, 2, 3} {
			if raw[i] != want {
				t.Fatalf("raw[%d] = %v, want %d", i, raw[i], want)
			}
		}
	})

	t.Run("disabled omitted", func(t *testing.T) {
		second, _ := a.At(1)
		second.SetDisabled(true)
		raw := a.RawValue().([]any)
		if len(raw) != 2 {
			t.Fatalf("len(raw) = %d after disabling one, want 2", len(raw))
		}
		if raw[0] != 1 || raw[1] != 3 {
			t.Fatalf("raw = %v, want [1 3]", raw)
		}
	})
}

func TestArray_StatusPropagation(t *testing.T) {
	t.Run("invalid child makes array invalid", func(t *testing.T) {
		child := NewControl(NewPlainValueStore(""), "", Required(""))
		a := NewArray(child)
		child.Set("")
		if a.Status() != StatusInvalid {
			t.Fatalf("Status() = %q, want %q", a.Status(), StatusInvalid)
		}
	})

	t.Run("Add fires status propagation up", func(t *testing.T) {
		child := NewControl(NewPlainValueStore(""), "", Required(""))
		child.Set("")
		inner := NewArray()
		root := NewGroup(map[string]FormNode{"items": inner})

		var got Status
		root.OnStatusChange(func(s Status) { got = s })

		inner.Add(child) // child is INVALID, so root must flip to INVALID
		if got != StatusInvalid {
			t.Fatalf("root status listener received %q after Add, want %q", got, StatusInvalid)
		}
	})
}
