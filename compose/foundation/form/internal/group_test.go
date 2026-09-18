package formengine

import "testing"

// pendingNode is a test-only FormNode that always reports StatusPending. It
// exists to exercise the Pending branch of Group/Array status aggregation, since
// no real engine node emits StatusPending in this change.
type pendingNode struct{}

func (pendingNode) Status() Status              { return StatusPending }
func (pendingNode) Validate() bool              { return true }
func (pendingNode) Errors() map[string]string   { return nil }
func (pendingNode) IsTouched() bool             { return false }
func (pendingNode) IsDirty() bool               { return false }
func (pendingNode) IsPristine() bool            { return true }
func (pendingNode) IsEnabled() bool             { return true }
func (pendingNode) MarkAsTouched()              {}
func (pendingNode) MarkAsUntouched()            {}
func (pendingNode) MarkAsPristine()             {}
func (pendingNode) SetDisabled(bool)            {}
func (pendingNode) Reset()                      {}
func (pendingNode) Parent() FormNode            { return nil }
func (pendingNode) Children() []FormNode        { return nil }
func (pendingNode) Path() string                { return "" }
func (pendingNode) Get(string) (FormNode, bool) { return nil, false }
func (pendingNode) RawValue() any               { return nil }
func (pendingNode) setParent(FormNode)          {}
func (pendingNode) onChildChanged()             {}
func (pendingNode) setInheritedDisabled(bool)   {}

func TestGroup_StatusPropagation(t *testing.T) {
	t.Run("invalid child makes group invalid", func(t *testing.T) {
		child := NewControl(NewPlainValueStore(""), "", Required(""))
		g := NewGroup(map[string]FormNode{"c": child})
		child.Set("")
		if g.Status() != StatusInvalid {
			t.Fatalf("Status() = %q, want %q", g.Status(), StatusInvalid)
		}
	})

	t.Run("all children disabled makes group disabled", func(t *testing.T) {
		c1 := NewControl(NewPlainValueStore("a"), "a")
		c2 := NewControl(NewPlainValueStore("b"), "b")
		g := NewGroup(map[string]FormNode{"a": c1, "b": c2})
		c1.SetDisabled(true)
		c2.SetDisabled(true)
		if g.Status() != StatusDisabled {
			t.Fatalf("Status() = %q, want %q", g.Status(), StatusDisabled)
		}
	})

	t.Run("pending propagates", func(t *testing.T) {
		g := NewGroup(map[string]FormNode{"p": pendingNode{}})
		if g.Status() != StatusPending {
			t.Fatalf("Status() = %q, want %q", g.Status(), StatusPending)
		}
	})

	t.Run("invalid takes precedence over pending", func(t *testing.T) {
		child := NewControl(NewPlainValueStore(""), "", Required(""))
		child.Set("")
		g := NewGroup(map[string]FormNode{"c": child, "p": pendingNode{}})
		if g.Status() != StatusInvalid {
			t.Fatalf("Status() = %q, want %q", g.Status(), StatusInvalid)
		}
	})
}

func TestGroup_DisabledOwnVsInherited(t *testing.T) {
	t.Run("re-enabling parent restores enabled child", func(t *testing.T) {
		child := NewControl(NewPlainValueStore("x"), "x")
		g := NewGroup(map[string]FormNode{"c": child})

		g.SetDisabled(true)
		if child.IsEnabled() {
			t.Fatal("child IsEnabled() = true while parent disabled, want false")
		}
		if child.Status() != StatusDisabled {
			t.Fatalf("child Status() = %q, want %q", child.Status(), StatusDisabled)
		}

		g.SetDisabled(false)
		if !child.IsEnabled() {
			t.Fatal("child IsEnabled() = false after parent re-enabled, want true")
		}
	})

	t.Run("re-enabling parent preserves child own-disabled state", func(t *testing.T) {
		child := NewControl(NewPlainValueStore("y"), "y")
		child.SetDisabled(true)
		g := NewGroup(map[string]FormNode{"c": child})

		g.SetDisabled(true)
		g.SetDisabled(false)

		if child.IsEnabled() {
			t.Fatal("own-disabled child IsEnabled() = true after round-trip, want false")
		}
	})

	t.Run("child reports own and inherited independently", func(t *testing.T) {
		child := NewControl(NewPlainValueStore("z"), "z")
		g := NewGroup(map[string]FormNode{"c": child})
		if child.InheritedDisabled() {
			t.Fatal("fresh child InheritedDisabled() = true, want false")
		}
		g.SetDisabled(true)
		if !child.InheritedDisabled() {
			t.Fatal("child InheritedDisabled() = false after parent disable, want true")
		}
	})
}

func TestGroup_RawValueExport(t *testing.T) {
	name := NewControl(NewPlainValueStore("Jane"), "Jane")
	city := NewControl(NewPlainValueStore("Sydney"), "Sydney")
	address := NewGroup(map[string]FormNode{"city": city})
	root := NewGroup(map[string]FormNode{"name": name, "address": address})

	t.Run("nested map", func(t *testing.T) {
		raw := root.RawValue().(map[string]any)
		if raw["name"] != "Jane" {
			t.Fatalf("raw[name] = %v, want %q", raw["name"], "Jane")
		}
		addr, ok := raw["address"].(map[string]any)
		if !ok {
			t.Fatalf("raw[address] is %T, want map[string]any", raw["address"])
		}
		if addr["city"] != "Sydney" {
			t.Fatalf("addr[city] = %v, want %q", addr["city"], "Sydney")
		}
	})

	t.Run("disabled child omitted", func(t *testing.T) {
		city.SetDisabled(true)
		raw := root.RawValue().(map[string]any)
		addr, ok := raw["address"].(map[string]any)
		if !ok {
			t.Fatalf("raw[address] is %T, want map[string]any", raw["address"])
		}
		if _, present := addr["city"]; present {
			t.Fatal("disabled city should be omitted from export")
		}
	})
}

func TestGroup_StatusListenerFiresOnChildInvalid(t *testing.T) {
	child := NewControl(NewPlainValueStore(""), "", Required(""))
	g := NewGroup(map[string]FormNode{"c": child})

	var got Status
	g.OnStatusChange(func(s Status) { got = s })

	child.Set("")

	if got != StatusInvalid {
		t.Fatalf("group status listener received %q, want %q", got, StatusInvalid)
	}
}

func TestGroup_MarkAsTouchedPropagates(t *testing.T) {
	child := NewControl(NewPlainValueStore("a"), "a")
	inner := NewGroup(map[string]FormNode{"c": child})
	root := NewGroup(map[string]FormNode{"inner": inner})

	child.MarkAsTouched()

	if !root.IsTouched() {
		t.Fatal("root IsTouched() = false after child touched, want true")
	}
	if !inner.IsTouched() {
		t.Fatal("inner IsTouched() = false after child touched, want true")
	}
}
