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
			if resource.IsLoading() {
				return true
			}
			if resource.HasData() {
				return fn(resource.Data()).IsLoading()
			}
			return false
		},
		errorFunc: func() error {
			if resource.Error() != nil {
				return resource.Error()
			}
			if resource.HasData() {
				return fn(resource.Data()).Error()
			}
			return nil
		},
		dataFunc: func() U {
			return fn(resource.Data()).Data()
		},
		hasDataFunc: func() bool {
			if !resource.HasData() {
				return false
			}
			return fn(resource.Data()).HasData()
		},
	}
}
