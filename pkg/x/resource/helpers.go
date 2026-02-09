package resource

func MapData[T, U any](resource Resource[T], fn func(T) U) Resource[U] {
	return newWrappedResource(resource, fn)
}

func MapError[T any](resource Resource[T], fn func(error) error) Resource[T] {
	return newWrappedResource(
		resource,
		func(t T) T {
			return t
		},
		withErrorMapper[T](fn),
	)
}

func FlatMapData[T, U any](resource Resource[T], fn func(T) Resource[U]) Resource[U] {
	return &wrappedResource[U]{
		isLoadingFunc: func() bool {
			if resource.isLoading() {
				return true
			}
			if resource.hasData() {
				return fn(resource.getData()).isLoading()
			}
			return false
		},
		errorFunc: func() error {
			if resource.hasError() {
				return resource.getError()
			}
			if resource.hasData() {
				return fn(resource.getData()).getError()
			}
			return nil
		},
		dataFunc: func() U {
			return fn(resource.getData()).getData()
		},
		hasDataFunc: func() bool {
			if !resource.hasData() {
				return false
			}
			return fn(resource.getData()).hasData()
		},
	}
}
