package resource

import "github.com/zodimo/go-maybe"

func Match[T, U any](resource Resource[T], onLoading func() U, onError func(error) U, onData func(T) U) U {
	if resource.isLoading() {
		return onLoading()
	}
	if resource.hasData() {
		return onData(resource.getData())
	}
	return onError(resource.getError())
}

type MatchOptions[T, U any] struct {
	onLoading maybe.Maybe[func() maybe.Maybe[U]]
	onError   maybe.Maybe[func(error) maybe.Maybe[U]]
	onData    maybe.Maybe[func(T) maybe.Maybe[U]]
	onDefault func() maybe.Maybe[U]
}

type MatchOption[T, U any] func(*MatchOptions[T, U])

func WithOnLoading[T, U any](onLoading func() U) MatchOption[T, U] {
	return func(opts *MatchOptions[T, U]) {
		opts.onLoading = maybe.Some(func() maybe.Maybe[U] {
			return maybe.Some(onLoading())
		})
	}
}

func WithOnError[T, U any](onError func(error) U) MatchOption[T, U] {
	return func(opts *MatchOptions[T, U]) {
		opts.onError = maybe.Some(func(err error) maybe.Maybe[U] {
			return maybe.Some(onError(err))
		})
	}
}

func WithOnData[T, U any](onData func(T) U) MatchOption[T, U] {
	return func(opts *MatchOptions[T, U]) {
		opts.onData = maybe.Some(func(data T) maybe.Maybe[U] {
			return maybe.Some(onData(data))
		})
	}
}

func WithDefault[T, U any](onDefault func() U) MatchOption[T, U] {
	return func(opts *MatchOptions[T, U]) {
		opts.onDefault = func() maybe.Maybe[U] {
			return maybe.Some(onDefault())
		}
	}
}

func MatchPartial[T, U any](resource Resource[T], options ...MatchOption[T, U]) maybe.Maybe[U] {

	opts := MatchOptions[T, U]{
		onLoading: maybe.None[func() maybe.Maybe[U]](),
		onError:   maybe.None[func(error) maybe.Maybe[U]](),
		onData:    maybe.None[func(T) maybe.Maybe[U]](),
		onDefault: func() maybe.Maybe[U] {
			return maybe.None[U]()
		},
	}
	for _, opt := range options {
		opt(&opts)
	}
	defaultFunc := opts.onDefault
	if resource.isLoading() {
		return opts.onLoading.OrElse(defaultFunc)()
	}
	if resource.hasData() {
		return opts.onData.OrElse(func(T) maybe.Maybe[U] {
			return defaultFunc()
		})(resource.getData())
	}
	return opts.onError.OrElse(func(error) maybe.Maybe[U] {
		return defaultFunc()
	})(resource.getError())
}

func MatchOnLoading[T, U any](resource Resource[T], onLoading func() U, defaultValue U) U {
	if resource.isLoading() {
		return onLoading()
	}
	return defaultValue
}

func MatchOnLoadingLazy[T, U any](resource Resource[T], onLoading func() U, defaultValue func() U) U {
	if resource.isLoading() {
		return onLoading()
	}
	return defaultValue()
}

func MatchOnError[T, U any](resource Resource[T], onError func(error) U, defaultValue U) U {
	if resource.hasError() {
		return onError(resource.getError())
	}
	return defaultValue
}

func MatchOnErrorLazy[T, U any](resource Resource[T], onError func(error) U, defaultValue func() U) U {
	if resource.hasError() {
		return onError(resource.getError())
	}
	return defaultValue()
}

func MatchOnData[T, U any](resource Resource[T], onData func(T) U, defaultValue U) U {
	if resource.hasData() {
		return onData(resource.getData())
	}
	return defaultValue
}

func MatchOnDataLazy[T, U any](resource Resource[T], onData func(T) U, defaultValue func() U) U {
	if resource.hasData() {
		return onData(resource.getData())
	}
	return defaultValue()
}
