package fform

import (
	"testing"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/lazy"
	"github.com/zodimo/go-compose/pkg/api"
)

// fakeListScope records the keys of items appended to it so tests can assert
// what a FormScope resolved and emitted without a real composer.
type fakeListScope struct {
	keys    []string
	items   map[string]api.Composable
	ordered []string
}

func newFakeListScope() *fakeListScope {
	return &fakeListScope{items: map[string]api.Composable{}}
}

func (s *fakeListScope) Item(key any, content api.Composable) {
	k, _ := key.(string)
	s.keys = append(s.keys, k)
	s.ordered = append(s.ordered, k)
	s.items[k] = content
}

func (s *fakeListScope) Items(count int, key func(index int) any, itemContent func(index int) api.Composable) {
	for i := 0; i < count; i++ {
		var k any
		if key != nil {
			k = key(i)
		}
		s.Item(k, itemContent(i))
	}
}

func (s *fakeListScope) StickyHeader(key any, content api.Composable) {
	s.Item(key, content)
}

var _ lazy.LazyListScope = (*fakeListScope)(nil)

func newTestScope(t *testing.T, root *Group) (*formScopeImpl, *fakeListScope) {
	t.Helper()
	list := newFakeListScope()
	return &formScopeImpl{
		listScope: list,
		tree:      root,
		basePath:  "",
	}, list
}

func TestFormScope_FieldSkipsUnresolvableKey(t *testing.T) {
	root := NewGroup(map[string]FormNode{
		"name": NewControl(NewPlainValueStore(""), ""),
	})
	scope, list := newTestScope(t, root)

	scope.Field("missing", func(FormNode) api.Composable { return noop })
	scope.Field("name", func(FormNode) api.Composable { return noop })

	if len(list.keys) != 1 {
		t.Fatalf("expected exactly 1 item for the resolvable key, got %v", list.keys)
	}
}

func TestFormScope_ControlOfResolvesTypedControl(t *testing.T) {
	root := NewGroup(map[string]FormNode{
		"name": NewControl(NewPlainValueStore(""), "", Required("")),
		"age":  NewControl(NewPlainValueStore(0), 0),
	})
	scope, _ := newTestScope(t, root)

	if _, ok := ControlOf[string](scope, "name"); !ok {
		t.Fatal("expected name to resolve as Control[string]")
	}
	// Wrong type must miss rather than panic.
	if _, ok := ControlOf[int](scope, "name"); ok {
		t.Fatal("expected Control[int] lookup on a string control to miss")
	}
	if _, ok := ControlOf[int](scope, "age"); !ok {
		t.Fatal("expected age to resolve as Control[int]")
	}
	if _, ok := ControlOf[string](scope, "missing"); ok {
		t.Fatal("expected missing key to miss")
	}
}

func TestFormScope_ControlFieldOfEmitsItem(t *testing.T) {
	root := NewGroup(map[string]FormNode{
		"name": NewControl(NewPlainValueStore(""), ""),
	})
	scope, list := newTestScope(t, root)

	ControlFieldOf[string](scope, "name", func(b *FormFieldBinding[string]) api.Composable {
		if b == nil {
			t.Fatal("expected content to receive a non-nil binding")
		}
		return noop
	})

	if len(list.keys) != 1 {
		t.Fatalf("expected 1 item, got %v", list.keys)
	}
	// ControlFieldOf defers content until the list item composes; execute it to
	// confirm the binding flows through.
	if item, ok := list.items[list.keys[0]]; !ok {
		t.Fatalf("expected an item stored under key %q", list.keys[0])
	} else {
		item(compose.NewComposer())
	}
}

func TestFormScope_FormArrayIteratesRowsWithRelativePaths(t *testing.T) {
	root := NewGroup(map[string]FormNode{
		"phones": NewArray(
			phoneGroup("111"),
			phoneGroup("222"),
		),
	})
	scope, list := newTestScope(t, root)

	var indices []int
	var resolved []string
	scope.FormArray("phones", func(index int, rowScope FormScope) api.Composable {
		indices = append(indices, index)
		control, ok := ControlOf[string](rowScope, "number")
		if ok {
			resolved = append(resolved, control.Value())
		}
		return noop
	})

	if len(indices) != 2 || indices[0] != 0 || indices[1] != 1 {
		t.Fatalf("expected indices [0 1], got %v", indices)
	}
	if len(resolved) != 2 || resolved[0] != "111" || resolved[1] != "222" {
		t.Fatalf("expected row-relative resolution to find 111/222, got %v", resolved)
	}
	if len(list.keys) != 2 {
		t.Fatalf("expected 2 emitted rows, got %v", list.keys)
	}
	if list.keys[0] == list.keys[1] {
		t.Fatalf("row item keys must be distinct, got %v", list.keys)
	}
}

func TestFormScope_GroupScopeNestsRelativePaths(t *testing.T) {
	root := NewGroup(map[string]FormNode{
		"billing": NewGroup(map[string]FormNode{
			"city": NewControl(NewPlainValueStore(""), ""),
		}),
	})
	scope, _ := newTestScope(t, root)

	var got *Control[string]
	scope.GroupScope("billing", func(child FormScope) {
		got, _ = ControlOf[string](child, "city")
	})

	if got == nil {
		t.Fatal("expected nested scope to resolve billing.city")
	}
}

func TestFormScope_ArrayLength(t *testing.T) {
	root := NewGroup(map[string]FormNode{
		"phones": NewArray(phoneGroup(""), phoneGroup(""), phoneGroup("")),
		"notarr": NewControl(NewPlainValueStore(""), ""),
	})
	scope, _ := newTestScope(t, root)

	if n := scope.ArrayLength("phones"); n != 3 {
		t.Fatalf("expected 3, got %d", n)
	}
	if n := scope.ArrayLength("notarr"); n != 0 {
		t.Fatalf("expected 0 for a non-array, got %d", n)
	}
	if n := scope.ArrayLength("missing"); n != 0 {
		t.Fatalf("expected 0 for a miss, got %d", n)
	}
}

func phoneGroup(number string) *Group {
	return NewGroup(map[string]FormNode{
		"number": NewControl(NewPlainValueStore(number), number),
	})
}

func noop(c api.Composer) api.Composer { return c }
