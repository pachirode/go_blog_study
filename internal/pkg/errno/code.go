package errno

import (
	"github.com/pachirode/pkg/errorsx"
	"net/http"
)

var (
	// OK 请求成功
	OK = &errorsx.ErrorX{Code: http.StatusOK, Reason: "", Message: ""}
	// ErrBind 绑定参数错误
	ErrBind = errorsx.ErrBind
	// ErrInternal 服务器内部错误
	ErrInternal = errorsx.ErrInternal
	// ErrNotFound 资源不存在
	ErrNotFound = errorsx.ErrNotFound

	// ErrInvalidArgument 表示参数验证失败
	ErrInvalidArgument = errorsx.ErrInvalidArgument
	// ErrUnauthenticated 表示认证失败.
	ErrUnauthenticated = errorsx.ErrUnauthenticated
	// ErrPermissionDenied 表示请求没有权限.
	ErrPermissionDenied = errorsx.ErrPermissionDenied
	// ErrOperationFailed 表示操作失败.
	ErrOperationFailed = errorsx.ErrOperationFailed
)

var (
	// ErrPageNotFound 为找到页面
	ErrPageNotFound = &errorsx.ErrorX{Code: http.StatusNotFound, Reason: "NotFound.PageNotFound", Message: "Page not found."}
	// ErrSignToken 签发令牌 JWT Token 错误
	ErrSignToken = &errorsx.ErrorX{Code: http.StatusUnauthorized, Reason: "Unauthenticated.SignTokenError", Message: "Error occurred while signing the JSON web token."}
	// ErrTokenInvalid JWT Token 格式错误
	ErrTokenInvalid = &errorsx.ErrorX{Code: http.StatusUnauthorized, Reason: "Unauthenticated.TokenInvalid", Message: "Token is invalid."}
	ErrDBRead       = &errorsx.ErrorX{Code: http.StatusInternalServerError, Reason: "InternalError.DBRead", Message: "Error occurred while reading from database."}
	ErrDBWrite      = &errorsx.ErrorX{Code: http.StatusInternalServerError, Reason: "InternalError.DBWrite", Message: "Error occurred while writing to database."}
	ErrAddRole      = &errorsx.ErrorX{Code: http.StatusInternalServerError, Reason: "InternalError.AddRole", Message: "Error occurred while adding role."}
	ErrRemoveRole   = &errorsx.ErrorX{Code: http.StatusInternalServerError, Reason: "InternalError.RemoveRole", Message: "Error occurred while removing role."}
)
