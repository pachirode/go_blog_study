package errno

import (
	"github.com/pachirode/pkg/errorsx"
	"net/http"
)

var (
	ErrUsernameInvalid = &errorsx.ErrorX{
		Code:    http.StatusBadRequest,
		Reason:  "InvalidArgument.UsernameInvalid",
		Message: "Invalid username: Username must consists of letters, digits and _, and its length must between 3 and 20",
	}

	ErrPasswordInvalid = &errorsx.ErrorX{
		Code:    http.StatusBadRequest,
		Reason:  "InvalidArgument.PasswordInvalid",
		Message: "Password is incorrect",
	}

	ErrUserAlreadyExist = &errorsx.ErrorX{Code: http.StatusBadRequest, Reason: "AlreadyExist.UserAlreadyExist", Message: "User already exist."}
	ErrUserNotFound     = &errorsx.ErrorX{Code: http.StatusNotFound, Reason: "NotFound.UserNotFound", Message: "User not found."}
)
