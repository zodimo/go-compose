package fform

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/lazy"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/pkg/api"
)

// FormScope binds named nodes of the form tree to their rendering content.
//
// Keys are resolved against the tree through the path resolver, relative to the
// scope's base path (the empty path at the root, a group's path inside a nested
// scope). The tree is the single source of truth: any key that does not resolve
// renders nothing rather than panicking, so content that lists a field the tree
// does not contain is silently skipped.
//
// Resolution helpers (ControlOf/Group/Array) are the type-safe escape hatch for
// callers that need the typed node directly; Field/Form/FormArray accept the raw
// FormNode for callers that do not.
type FormScope interface {
	// Has reports whether key resolves to a node under this scope.
	Has(key string) bool

	// Group resolves key and asserts it is a *Group.
	Group(key string) (*Group, bool)

	// Array resolves key and asserts it is an *Array.
	Array(key string) (*Array, bool)

	// Field resolves key against the tree and passes the resolved node into
	// content. An unresolvable key renders nothing.
	Field(key string, content func(node FormNode) api.Composable)

	// Form resolves key against the tree and, when found, renders a header row
	// then nests a child scope whose keys resolve relative to the subtree named
	// by key. An unresolvable key renders nothing.
	Form(key string, header api.Composable, children func(FormScope))

	// GroupScope is Form without a header row: it nests a child scope over the
	// group named by key while rendering no heading of its own.
	GroupScope(key string, children func(FormScope))

	// FormArray resolves key as an *Array and renders one list item per child,
	// in order. row receives the child's index and a scope rooted at that child
	// so nested keys resolve relative to the row, and returns the row's
	// composable. A group row typically calls
	// ControlFieldOf[string](rowScope, "number", ...).
	//
	// key is used as the list item key prefix; each row's item key is
	// "<base>key[i]". An unresolvable key renders nothing.
	FormArray(key string, row func(index int, rowScope FormScope) api.Composable)

	// ArrayLength returns the number of children of the array named by key, or
	// 0 when the key does not resolve to an array. Useful for rendering
	// add/remove affordances next to a FormArray.
	ArrayLength(key string) int

	// Resolve returns the raw node at key, or (nil, false) on a miss. It is the
	// untyped escape hatch for advanced wiring.
	Resolve(key string) (FormNode, bool)
}

var _ FormScope = (*formScopeImpl)(nil)

type formScopeImpl struct {
	listScope lazy.LazyListScope
	tree      *Group
	basePath  string
	// itemPrefix is prepended to item keys so rows from different arrays (or a
	// row and a scalar field) never collide in the lazy list key space.
	itemPrefix string
}

// Has reports whether key resolves under this scope.
func (s *formScopeImpl) Has(key string) bool {
	_, ok := s.Resolve(key)
	return ok
}

// Resolve returns the node at key relative to this scope's base path.
func (s *formScopeImpl) Resolve(key string) (FormNode, bool) {
	return ResolvePath(s.tree, s.joinPath(s.basePath, key))
}

// Resolver resolves a dotted/bracketed path to a node. A FormScope satisfies it
// directly; a bare tree node can be adapted with NodeResolver, so the resolution
// helpers work against either.
type Resolver interface {
	Resolve(key string) (FormNode, bool)
}

// nodeResolver adapts a FormNode's Get method to the Resolver vocabulary.
type nodeResolver struct {
	node FormNode
}

func (r nodeResolver) Resolve(key string) (FormNode, bool) {
	if r.node == nil {
		return nil, false
	}
	return r.node.Get(key)
}

// NodeResolver adapts a tree node (root Group, any FormNode) to a Resolver so it
// can be passed to ControlOf and ControlViewOf.
func NodeResolver(node FormNode) Resolver {
	return nodeResolver{node: node}
}

// ControlOf resolves key and asserts it is a *Control[T]. The bool result is
// false when the key does not resolve or resolves to a different node kind.
//
// ControlOf is a package-level generic helper, not a FormScope method: Go
// forbids type parameters on interface methods. It works on any Resolver: a
// FormScope, or a tree node wrapped with NodeResolver.
func ControlOf[T any](resolver Resolver, key string) (*Control[T], bool) {
	node, ok := resolver.Resolve(key)
	if !ok {
		return nil, false
	}
	control, ok := node.(*Control[T])
	if !ok {
		return nil, false
	}
	return control, true
}

// Group resolves key and asserts it is a *Group.
func (s *formScopeImpl) Group(key string) (*Group, bool) {
	node, ok := s.Resolve(key)
	if !ok {
		return nil, false
	}
	group, ok := node.(*Group)
	if !ok {
		return nil, false
	}
	return group, true
}

// Array resolves key and asserts it is an *Array.
func (s *formScopeImpl) Array(key string) (*Array, bool) {
	node, ok := s.Resolve(key)
	if !ok {
		return nil, false
	}
	array, ok := node.(*Array)
	if !ok {
		return nil, false
	}
	return array, true
}

// Field resolves key relative to the scope's base path and, when found, renders
// content with the resolved node. An unresolvable key renders nothing.
func (s *formScopeImpl) Field(key string, content func(node FormNode) api.Composable) {
	node, ok := s.Resolve(key)
	if !ok {
		return
	}

	s.listScope.Item(s.itemKey(key), func(c api.Composer) api.Composer {
		return content(node)(c)
	})
}

// ControlFieldOf renders content with a binding over the *Control[T] at key as
// its own list item, rendering nothing on a miss or kind mismatch. Use it for
// fields that occupy a full row of the form.
//
// ControlFieldOf must be called during content construction (it appends to the
// scope's list), so it cannot be a FormScope method for the same
// no-type-parameters-on-methods reason as ControlOf. Use ControlViewOf instead
// when the binding must be embedded inside other content (for example an array
// row that also renders a remove button).
func ControlFieldOf[T any](scope FormScope, key string, content func(binding *FormFieldBinding[T]) api.Composable) {
	impl, ok := scope.(*formScopeImpl)
	if !ok {
		return
	}
	control, ok := ControlOf[T](scope, key)
	if !ok {
		return
	}

	binding := NewFormFieldBinding(control)
	impl.listScope.Item(impl.itemKey(key), func(c api.Composer) api.Composer {
		return content(binding)(c)
	})
}

// ControlViewOf returns content bound to the *Control[T] at key for embedding
// inline (rather than as its own list item). It returns a no-op composable on a
// miss or kind mismatch, so it is safe to place directly in a row or column.
func ControlViewOf[T any](resolver Resolver, key string, content func(binding *FormFieldBinding[T]) api.Composable) api.Composable {
	control, ok := ControlOf[T](resolver, key)
	if !ok {
		return func(c api.Composer) api.Composer { return c }
	}
	return content(NewFormFieldBinding(control))
}

// Form resolves key relative to the scope's base path and, when found, renders
// the header row then nests a child scope resolving relative paths against the
// subtree named by key. An unresolvable key renders nothing (the tree is the
// single source of truth, matching Field).
func (s *formScopeImpl) Form(key string, header api.Composable, children func(FormScope)) {
	if !s.Has(key) {
		return
	}

	s.listScope.Item(s.itemKey(key), func(c api.Composer) api.Composer {
		return row.Row(
			c.Sequence(
				row.Row(
					c.Sequence(
						header,
					),
					row.WithAlignment(row.Middle),
					row.WithModifier(
						weight.Weight(1).
							Then(padding.All(4)),
					),
				),
			),
			row.WithAlignment(row.Middle),
			row.WithModifier(
				size.FillMaxWidth(),
			),
		)(c)
	})

	children(s.childScope(key))
}

// GroupScope nests a child scope over the group named by key, rendering no
// header of its own.
func (s *formScopeImpl) GroupScope(key string, children func(FormScope)) {
	if !s.Has(key) {
		return
	}
	children(s.childScope(key))
}

// FormArray resolves key as an *Array and renders one list item per child in
// order, each row receiving its index and a scope rooted at that child.
func (s *formScopeImpl) FormArray(key string, row func(index int, rowScope FormScope) api.Composable) {
	array, ok := s.Array(key)
	if !ok {
		return
	}

	length := array.Length()
	base := s.joinPath(s.basePath, key)
	for i := 0; i < length; i++ {
		index := i
		rowScope := &formScopeImpl{
			listScope:  s.listScope,
			tree:       s.tree,
			basePath:   fmt.Sprintf("%s[%d]", base, index),
			itemPrefix: s.itemPrefix,
		}
		content := row(index, rowScope)
		s.listScope.Item(s.arrayItemKey(key, index), content)
	}
}

// ArrayLength returns the number of children of the array named by key, or 0
// when the key does not resolve to an array.
func (s *formScopeImpl) ArrayLength(key string) int {
	array, ok := s.Array(key)
	if !ok {
		return 0
	}
	return array.Length()
}

// childScope returns a scope whose base path is key resolved against this
// scope's base path.
func (s *formScopeImpl) childScope(key string) *formScopeImpl {
	return &formScopeImpl{
		listScope:  s.listScope,
		tree:       s.tree,
		basePath:   s.joinPath(s.basePath, key),
		itemPrefix: s.itemPrefix,
	}
}

// itemKey namespaces a field/form item key so it cannot collide with other
// scope items.
func (s *formScopeImpl) itemKey(key string) string {
	if s.basePath == "" {
		return s.itemPrefix + key
	}
	return s.itemPrefix + s.basePath + "/" + key
}

// arrayItemKey namespaces a row item key by array path and index. The index is
// positional: adding or removing an earlier row shifts later rows' keys, which
// matches the Array node's positional index model.
func (s *formScopeImpl) arrayItemKey(key string, index int) string {
	return fmt.Sprintf("%s%s/%s[%d]", s.itemPrefix, s.basePath, key, index)
}

// joinPath joins a base path and key with a dot separator, returning key alone
// when base is empty.
func (s *formScopeImpl) joinPath(base, key string) string {
	if base == "" {
		return key
	}
	return base + "." + key
}
