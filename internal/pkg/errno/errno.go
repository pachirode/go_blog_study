package errno

import "fmt"

type Errno struct {
	HTTP    int
	Code    string
	Message string
}

func (err *Errno) String() string {
	return fmt.Sprintf("Http: %d; Code: %s; Message: %s", err.HTTP, err.Code, err.Message)
}

func (err *Errno) Error() string {
	return err.Message
}

func (err *Errno) SetMessage(msgTemplate string, args ...interface{}) *Errno {
	err.Message = fmt.Sprintf(msgTemplate, args...)
	return err
}

func Derrcode(err error) (int, string, string) {
	if err == nil {
		return OK.HTTP, OK.Code, OK.Message
	}

	switch typed := err.(type) {
	case *Errno:
		return typed.HTTP, typed.Code, typed.Message
	default:
		return InternalServerError.HTTP, InternalServerError.Code, err.Error()
	}
}
