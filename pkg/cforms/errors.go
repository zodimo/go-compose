package cforms

type CustomError struct {
	Message string
}

func (c CustomError) Error() string {
	return c.Message
}
