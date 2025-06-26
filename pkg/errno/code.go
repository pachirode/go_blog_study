package errno

var (
	OK                  = &Errno{HTTP: 200, Code: "", Message: ""}
	InternalServerError = &Errno{HTTP: 500, Code: "InternalError", Message: "Internal server error."}
	ErrPageNotFound     = &Errno{HTTP: 404, Code: "ResourceNotFound.PageNotFound", Message: "Page not found."}
)
