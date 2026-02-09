package resource

var _ Resource[any] = (*resourceLoading[any])(nil)

type resourceLoading[T any] struct {
}

func (r *resourceLoading[T]) isLoading() bool {
	return true
}

func (r *resourceLoading[T]) getError() error {
	return nil
}
func (r *resourceLoading[T]) hasError() bool {
	return false
}

func (r *resourceLoading[T]) getData() T {
	panic("Data() called on loading resource")
}

func (r *resourceLoading[T]) hasData() bool {
	return false
}

func (r *resourceLoading[T]) isResource() {}
