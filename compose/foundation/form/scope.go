package fform

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/lazy"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/pkg/api"
)

// FormScope binds named nodes of the form tree to their rendering content.
type FormScope interface {
	// Field resolves key against the tree and passes the resolved node into content.
	Field(key string, content func(node FormNode) api.Composable)
	// Form resolves key against the tree and, when found, renders a header row
	// then nests a child scope whose keys resolve relative to the subtree named
	// by key. An unresolvable key renders nothing (the tree is the single source
	// of truth).
	Form(key string, header api.Composable, children func(FormScope))
}

var _ FormScope = (*formScopeImpl)(nil)

type formScopeImpl struct {
	listScope lazy.LazyListScope
	tree      *Group
	basePath  string
}

// Field resolves key relative to the scope's base path and, when found, renders
// content with the resolved node. An unresolvable key renders nothing.
func (s *formScopeImpl) Field(key string, content func(node FormNode) api.Composable) {
	path := s.joinPath(s.basePath, key)
	node, ok := ResolvePath(s.tree, path)
	if !ok {
		return
	}

	s.listScope.Item(key, func(c api.Composer) api.Composer {
		return content(node)(c)
	})
}

// Form resolves key relative to the scope's base path and, when found, renders
// the header row then nests a child scope resolving relative paths against the
// subtree named by key. An unresolvable key renders nothing (the tree is the
// single source of truth, matching Field).
func (s *formScopeImpl) Form(key string, header api.Composable, children func(FormScope)) {
	path := s.joinPath(s.basePath, key)
	if _, ok := ResolvePath(s.tree, path); !ok {
		return
	}

	s.listScope.Item(key, func(c api.Composer) api.Composer {
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

	childScope := &formScopeImpl{
		listScope: s.listScope,
		tree:      s.tree,
		basePath:  path,
	}
	children(childScope)
}

// joinPath joins a base path and key with a dot separator, returning key alone
// when base is empty.
func (s *formScopeImpl) joinPath(base, key string) string {
	if base == "" {
		return key
	}
	return base + "." + key
}
