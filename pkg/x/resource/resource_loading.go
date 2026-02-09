package resource

var _ Resource[any] = (*resourceLoading[any])(nil)

type resourceLoading[T any] struct {
	isLoading bool
}

func (r *resourceLoading[T]) IsLoading() bool {
	return true
}

func (r *resourceLoading[T]) Error() error {
	return nil
}

func (r *resourceLoading[T]) Data() T {
	panic("Data() called on loading resource")
}

func (r *resourceLoading[T]) HasData() bool {
	return false
}

func (r *resourceLoading[T]) isResource() {}
