package resource

var _ Resource[any] = (*resourceSuccess[any])(nil)

type resourceSuccess[T any] struct {
	data T
}

func (r *resourceSuccess[T]) isLoading() bool {
	return false
}

func (r *resourceSuccess[T]) getError() error {
	return nil
}

func (r *resourceSuccess[T]) hasError() bool {
	return false
}

func (r *resourceSuccess[T]) getData() T {
	return r.data
}

func (r *resourceSuccess[T]) hasData() bool {
	return true
}

func (r *resourceSuccess[T]) isResource() {}
