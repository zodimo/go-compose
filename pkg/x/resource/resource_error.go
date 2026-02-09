package resource

var _ Resource[any] = (*resourceError[any])(nil)

type resourceError[T any] struct {
	isLoading bool
	error     error
}

func (r *resourceError[T]) IsLoading() bool {
	return false
}

func (r *resourceError[T]) Error() error {
	return r.error
}

func (r *resourceError[T]) Data() T {
	panic("Data() called on error resource")
}

func (r *resourceError[T]) HasData() bool {
	return false
}

func (r *resourceError[T]) isResource() {}
