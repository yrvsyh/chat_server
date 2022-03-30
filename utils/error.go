package utils

type Error struct {
	code    int
	message string
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.message
}

func (e *Error) Error() string {
	return e.message
}
