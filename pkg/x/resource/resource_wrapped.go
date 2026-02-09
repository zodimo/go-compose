package resource

import "sync"

var _ Resource[any] = (*wrappedResource[any])(nil)

type wrappedResource[T any] struct {
	isLoadingFunc func() bool
	errorFunc     func() error
	dataFunc      func() T
	hasDataFunc   func() bool
	cachedData    *T
	once          sync.Once
}

type wrappedResourceOptions[T any] struct {
	mapErrorFunc func(error) error
}

func withErrorMapper[T any](fn func(error) error) wrappedResourceOption[T] {
	return func(opts *wrappedResourceOptions[T]) {
		opts.mapErrorFunc = fn
	}
}

type wrappedResourceOption[T any] func(*wrappedResourceOptions[T])

func defaultWrappedResourceOptions[T any]() wrappedResourceOptions[T] {
	return wrappedResourceOptions[T]{
		mapErrorFunc: func(err error) error {
			return err
		},
	}
}

func newWrappedResource[T, U any](resource Resource[T], mapDataFunc func(T) U, options ...wrappedResourceOption[U]) *wrappedResource[U] {
	opts := defaultWrappedResourceOptions[U]()
	for _, option := range options {
		option(&opts)
	}
	return &wrappedResource[U]{
		isLoadingFunc: func() bool {
			return resource.IsLoading()
		},
		errorFunc: func() error {
			return opts.mapErrorFunc(resource.Error())
		},
		dataFunc: func() U {
			return mapDataFunc(resource.Data())
		},
		hasDataFunc: func() bool {
			return resource.HasData()
		},
	}
}

func (w *wrappedResource[T]) IsLoading() bool {
	return w.isLoadingFunc()
}

func (w *wrappedResource[T]) Error() error {
	return w.errorFunc()
}

func (w *wrappedResource[T]) Data() T {
	w.once.Do(func() {
		val := w.dataFunc()
		w.cachedData = &val
	})
	return *w.cachedData
}

func (w *wrappedResource[T]) HasData() bool {
	return w.hasDataFunc()
}

func (w *wrappedResource[T]) isResource() {}
