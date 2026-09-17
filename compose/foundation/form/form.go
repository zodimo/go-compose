package fform

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/lazy"
	"github.com/zodimo/go-compose/pkg/api"
)

func Form(
	state *FormState,
	content func(FormScope),
	options ...FormOption,
) api.Composable {
	opts := DefaultFormOptions()
	for _, opt := range options {
		if opt != nil {
			opt(&opts)
		}
	}

	lazyOpts := make([]lazy.LazyListOption, 0)
	if opts.Modifier != nil {
		lazyOpts = append(lazyOpts, lazy.WithModifier(opts.Modifier))
	}

	return compose.CompositionLocalProvider(
		[]api.ProvidedValue{compose.LocalTextStyle.Provides(opts.TextStyle)},
		lazy.LazyColumn(
			func(scope lazy.LazyListScope) {
				tScope := &formScopeImpl{
					listScope: scope,
					state:     state,

					options: &opts,
				}
				content(tScope)
			},
			lazyOpts...,
		),
	)
}
