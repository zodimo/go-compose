package fform

type ValidationRule struct {
	id       string
	validate func(string) error
}

func (vr *ValidationRule) Validate(v string) error {
	if vr.validate == nil {
		return nil
	}
	return vr.validate(v)
}
