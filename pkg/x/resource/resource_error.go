package resource

var _ Resource[any] = (*resourceError[any])(nil)

type resourceError[T any] struct {
	error error
}

func (r *resourceError[T]) isLoading() bool {
	return false
}

func (r *resourceError[T]) getError() error {
	return r.error
}

func (r *resourceError[T]) hasError() bool {
	return true
}

func (r *resourceError[T]) getData() T {
	panic("Data() called on error resource")
}

func (r *resourceError[T]) hasData() bool {
	return false
}

func (r *resourceError[T]) isResource() {}
