package fform

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/lazy"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/pkg/api"
)

type FormScope interface {
	Field(key any, content api.Composable)
	Form(key any, header api.Composable, children func(FormScope))
}

var _ FormScope = (*formScopeImpl)(nil)

type formScopeImpl struct {
	listScope lazy.LazyListScope
	state     *FormState
	options   *FormOptions
}

func (s *formScopeImpl) Field(key any, content api.Composable) {
	// Capture state and options for closure
	// state := s.state
	// opts := s.options

	s.listScope.Item(key, func(c api.Composer) api.Composer {
		return content(c)
	})
}
func (s *formScopeImpl) Form(key any, header api.Composable, children func(FormScope)) {

	// Capture state and options for closure
	// state := s.state
	// opts := s.options

	// Branch Header
	s.listScope.Item(key, func(c api.Composer) api.Composer {
		return row.Row(
			c.Sequence(
				row.Row(
					c.Sequence(
						// Header Content
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

	// Children
	childScope := &formScopeImpl{
		listScope: s.listScope,
		state:     s.state,
		options:   s.options,
	}
	children(childScope)
}
