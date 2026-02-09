package resource

type Resource[T any] interface {
	IsLoading() bool
	Error() error
	Data() T
	HasData() bool
	isResource()
}

func ResourceError[T any](err error) Resource[T] {
	return &resourceError[T]{
		error: err,
	}
}

func ResourceLoading[T any]() Resource[T] {
	return &resourceLoading[T]{
		isLoading: true,
	}
}

func ResourceSuccess[T any](data T) Resource[T] {
	return &resourceSuccess[T]{
		data: data,
	}
}
