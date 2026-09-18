package fform

import (
	"fmt"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/lazy"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
)

// Form renders the form tree as a LazyColumn wrapped in a CompositionLocalProvider.
//
// It subscribes once to the tree root's change stream and bumps a remembered
// version on every mutation, driving the store -> window.Invalidate()
// recomposition loop (design D9).
func Form(
	formState *FormState,
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

	return func(c api.Composer) api.Composer {
		versionKey := fmt.Sprintf("%d/%s/formVersion", c.GenerateID(), c.GetPath())
		version := state.MustRemember(c, versionKey, func() int { return 0 })

		// Subscribe ONCE to the tree's change stream. Group.OnValueChange fires on
		// any descendant mutation (propagation), so a single root subscription
		// covers the whole tree. Capturing the remembered version value in the
		// closure is safe across recompositions.
		subKey := fmt.Sprintf("%d/%s/formSub", c.GenerateID(), c.GetPath())
		state.MustRemember(c, subKey, func() state.Subscription {
			return formState.root.OnValueChange(func(_ map[string]any) {
				version.Set(version.Get() + 1)
			})
		})

		return compose.CompositionLocalProvider(
			[]api.ProvidedValue{compose.LocalTextStyle.Provides(opts.TextStyle)},
			lazy.LazyColumn(
				func(scope lazy.LazyListScope) {
					tScope := &formScopeImpl{
						listScope: scope,
						tree:      formState.root,
					}
					content(tScope)
				},
				lazyOpts...,
			),
		)(c)
	}
}
