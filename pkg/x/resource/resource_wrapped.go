package resource

import (
	"encoding/json"
	"fmt"
	"sync"
)

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
			return resource.isLoading()
		},
		errorFunc: func() error {
			return opts.mapErrorFunc(resource.getError())
		},
		dataFunc: func() U {
			return mapDataFunc(resource.getData())
		},
		hasDataFunc: func() bool {
			return resource.hasData()
		},
	}
}

func (w *wrappedResource[T]) MarshalJSON() ([]byte, error) {
	resourceType := fmt.Sprintf("%T", w)
	if w.isLoading() {
		return json.Marshal(map[string]any{
			"resource_type": resourceType,
			"loading":       true,
		})
	}
	if w.hasError() {
		return json.Marshal(map[string]any{
			"resource_type": resourceType,
			"error":         w.getError().Error(),
		})
	}
	return json.Marshal(map[string]any{
		"resource_type": resourceType,
		"data":          w.getData(),
	})
}

func (w *wrappedResource[T]) isLoading() bool {
	return w.isLoadingFunc()
}

func (w *wrappedResource[T]) getError() error {
	return w.errorFunc()
}

func (w *wrappedResource[T]) hasError() bool {
	return w.errorFunc() != nil
}

func (w *wrappedResource[T]) getData() T {
	w.once.Do(func() {
		val := w.dataFunc()
		w.cachedData = &val
	})
	return *w.cachedData
}

func (w *wrappedResource[T]) hasData() bool {
	return w.hasDataFunc()
}

func (w *wrappedResource[T]) isResource() {}
