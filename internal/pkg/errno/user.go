package errno

var (
	ErrUserAlreadyExist  = &Errno{HTTP: 400, Code: "FailedOperation.UserAlreadyExist", Message: "User already exist."}
	ErrUserNotFound      = &Errno{HTTP: 400, Code: "FailedOperation.UserNotFound", Message: "User not found."}
	ErrPasswordIncorrect = &Errno{HTTP: 400, Code: "FailedOperation.PasswordIncorrect", Message: "Password incorrect."}
)
