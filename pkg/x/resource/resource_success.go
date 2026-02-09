package resource

var _ Resource[any] = (*resourceSuccess[any])(nil)

type resourceSuccess[T any] struct {
	data T
}

func (r *resourceSuccess[T]) IsLoading() bool {
	return false
}

func (r *resourceSuccess[T]) Error() error {
	return nil
}

func (r *resourceSuccess[T]) Data() T {
	return r.data
}

func (r *resourceSuccess[T]) HasData() bool {
	return true
}

func (r *resourceSuccess[T]) isResource() {}
