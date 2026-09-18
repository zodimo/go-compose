package formengine

import "testing"

func newTestTree() *Group {
	city := NewControl(NewPlainValueStore("Sydney"), "Sydney")
	zip := NewControl(NewPlainValueStore("2000"), "2000")
	address := NewGroup(map[string]FormNode{"city": city, "zip": zip})
	billing := NewGroup(map[string]FormNode{"address": address})

	item0 := NewGroup(map[string]FormNode{
		"name":  NewControl(NewPlainValueStore("A"), "A"),
		"price": NewControl(NewPlainValueStore(10), 10),
	})
	item1 := NewGroup(map[string]FormNode{
		"name":  NewControl(NewPlainValueStore("B"), "B"),
		"price": NewControl(NewPlainValueStore(20), 20),
	})
	item2 := NewGroup(map[string]FormNode{
		"name":  NewControl(NewPlainValueStore("C"), "C"),
		"price": NewControl(NewPlainValueStore(30), 30),
	})
	items := NewArray(item0, item1, item2)

	return NewGroup(map[string]FormNode{"billing": billing, "items": items})
}

func TestResolvePath(t *testing.T) {
	root := newTestTree()

	tests := []struct {
		name   string
		path   string
		want   string // expected RawValue of the resolved control, "" for group/array nodes
		wantOK bool
	}{
		{name: "dotted path", path: "billing.address.city", want: "Sydney", wantOK: true},
		{name: "dotted path deep", path: "billing.address.zip", want: "2000", wantOK: true},
		{name: "bracketed path", path: "items[0].name", want: "A", wantOK: true},
		{name: "bracketed index 2", path: "items[2].name", want: "C", wantOK: true},
		{name: "mixed path", path: "items[1].price", want: "20", wantOK: true},
		{name: "group node resolution", path: "billing", want: "", wantOK: true},
		{name: "unknown key miss", path: "billing.address.missing", want: "", wantOK: false},
		{name: "unknown top key miss", path: "nope.address.city", want: "", wantOK: false},
		{name: "out of range index miss", path: "items[5].name", want: "", wantOK: false},
		{name: "index into non-array miss", path: "billing[0].city", want: "", wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, ok := root.Get(tt.path)
			if ok != tt.wantOK {
				t.Fatalf("Get(%q) ok = %v, want %v", tt.path, ok, tt.wantOK)
			}
			if !tt.wantOK {
				if node != nil {
					t.Fatalf("Get(%q) returned non-nil node %v on miss", tt.path, node)
				}
				return
			}
			if node == nil {
				t.Fatalf("Get(%q) returned nil node", tt.path)
			}
		})
	}

	t.Run("resolved node identity", func(t *testing.T) {
		city, ok := root.Get("billing.address.city")
		if !ok {
			t.Fatal("failed to resolve billing.address.city")
		}
		if city.RawValue() != "Sydney" {
			t.Fatalf("resolved RawValue = %v, want %q", city.RawValue(), "Sydney")
		}
	})

	t.Run("empty path resolves to root", func(t *testing.T) {
		node, ok := root.Get("")
		if !ok || node != root {
			t.Fatalf("Get(\"\") = (%v, %v), want (root, true)", node, ok)
		}
	})
}

func TestPathString(t *testing.T) {
	root := newTestTree()

	tests := []struct {
		resolve string
		want    string
	}{
		{"billing.address.city", "billing.address.city"},
		{"items[0].name", "items[0].name"},
		{"items[2].price", "items[2].price"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			node, ok := root.Get(tt.resolve)
			if !ok {
				t.Fatalf("failed to resolve %q", tt.resolve)
			}
			if node.Path() != tt.want {
				t.Fatalf("Path() = %q, want %q", node.Path(), tt.want)
			}
		})
	}

	t.Run("root path is empty", func(t *testing.T) {
		if root.Path() != "" {
			t.Fatalf("root Path() = %q, want empty", root.Path())
		}
	})
}
