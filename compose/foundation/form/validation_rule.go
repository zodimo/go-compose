package fform

type ValidationRule[T any] struct {
	id       string
	validate func(T) error
}

func (vr *ValidationRule[T]) ID() string {
	return vr.id
}
func (vr *ValidationRule[T]) Validate(v T) error {
	if vr.validate == nil {
		return nil
	}
	err := vr.validate(v)
	if err != nil {
		// include rule id ?
		return err
	}
	return nil
}

func NewValidationRule[T any](id string, validateFunc func(T) error) *ValidationRule[T] {
	return &ValidationRule[T]{
		id:       id,
		validate: validateFunc,
	}
}
