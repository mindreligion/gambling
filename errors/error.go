package errors

type Code interface {
	Code() int
}

type gamblingError struct {
	code int
	text string
}
func (e *gamblingError) Error() string {
	return e.text
}
func (e *gamblingError) Code() int {
	return e.code
}
func New(code int, text string) error{
	return &gamblingError{code, text}
}
