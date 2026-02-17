package effect

import (
	"context"
	"fmt"

	"github.com/zodimo/go-compose/lifecycle"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
)

/**
 * Return a [CoroutineScope] bound to this point in the composition using the optional
 * [CoroutineContext] provided by [getContext]. [getContext] will only be called once and the same
 * [CoroutineScope] instance will be returned across recompositions.
 *
 * This scope will be [cancelled][CoroutineScope.cancel] when this call leaves the composition. The
 * [CoroutineContext] returned by [getContext] may not contain a [Job] as this scope is considered
 * to be a child of the composition.
 *
 * The default dispatcher of this scope if one is not provided by the context returned by
 * [getContext] will be the applying dispatcher of the composition's [Recomposer].
 *
 * Use this scope to launch jobs in response to callback events such as clicks or other user
 * interaction where the response to that event needs to unfold over time and be cancelled if the
 * composable managing that process leaves the composition. Jobs should never be launched into
 * **any** coroutine scope as a side effect of composition itself. For scoped ongoing jobs initiated
 * by composition, see [LaunchedEffect].
 *
 * This function will not throw if preconditions are not met, as composable functions do not yet
 * fully support exceptions. Instead the returned scope's [CoroutineScope.coroutineContext] will
 * contain a failed [Job] with the associated exception and will not be capable of launching child
 * jobs.
 */
// @Composable
// public inline fun rememberCoroutineScope(
//     crossinline getContext: @DisallowComposableCalls () -> CoroutineContext = {
//         EmptyCoroutineContext
//     }
// ): CoroutineScope {
//     val composer = currentComposer
//     return remember { createCompositionCoroutineScope(getContext(), composer) }
// }

type CoroutineScope interface {
	lifecycle.HasCoroutineScope
	Context() context.Context
	Launch(block func(ctx context.Context))
}

var _ CoroutineScope = (*coroutineScope)(nil)

type coroutineScope struct {
	ctx context.Context
}

func RememberCoroutineScope(c api.Composer) CoroutineScope {
	key := c.GenerateID()
	path := c.GetPath()

	coroutineScopePath := fmt.Sprintf("%d/%s/coroutineScope", key, path)

	currentCoroutineScope := state.MustRemember(c, coroutineScopePath, func() CoroutineScope {
		return &coroutineScope{}
	}).Get()

	return currentCoroutineScope
}

func (c *coroutineScope) SetCoroutineScope(ctx context.Context) {
	c.ctx = ctx
}

func (c *coroutineScope) Launch(block func(ctx context.Context)) {
	go func() {
		block(c.ctx)
	}()
}

func (c *coroutineScope) Context() context.Context {
	return c.ctx
}
