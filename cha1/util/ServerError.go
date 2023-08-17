package util

type ServerError struct {
	Code int
	Msg  string
}

// Error implements error.
func (se *ServerError) Error() string {
	return se.Msg
}

func NewServerError(code int, msg string) error {
	return &ServerError{
		Code: code,
		Msg:  msg,
	}
}
