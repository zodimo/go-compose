package resource

type Resource[T any] interface {
	isLoading() bool
	getError() error
	hasError() bool
	getData() T
	hasData() bool
	isResource()
}

func ResourceError[T any](err error) Resource[T] {
	return &resourceError[T]{
		error: err,
	}
}

func ResourceLoading[T any]() Resource[T] {
	return &resourceLoading[T]{}
}

func ResourceSuccess[T any](data T) Resource[T] {
	return &resourceSuccess[T]{
		data: data,
	}
}
