package fform

type FieldOptions struct {
	Validators []Validator
}

type FieldOption func(o FieldOptions)
