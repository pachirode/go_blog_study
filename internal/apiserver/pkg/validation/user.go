package validation

import (
	"context"

	genericValidation "github.com/pachirode/pkg/validation"

	"github.com/pachirode/go_blog_study/internal/pkg/contextx"
	"github.com/pachirode/go_blog_study/internal/pkg/errno"
	apiV1 "github.com/pachirode/go_blog_study/pkg/api/apiserver/v1/proto"
)

func (v *Validator) ValidateUserRules() genericValidation.Rules {
	// 通用的密码校验函数
	validatePassword := func() genericValidation.ValidatorFunc {
		return func(value any) error {
			return isValidPassword(value.(string))
		}
	}

	// 定义各字段的校验逻辑，通过一个 map 实现模块化和简化
	return genericValidation.Rules{
		"Password":    validatePassword(),
		"OldPassword": validatePassword(),
		"NewPassword": validatePassword(),
		"UserID": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("userID cannot be empty")
			}
			return nil
		},
		"Username": func(value any) error {
			if !isValidUsername(value.(string)) {
				return errno.ErrUsernameInvalid
			}
			return nil
		},
		"Nickname": func(value any) error {
			if len(value.(string)) >= 30 {
				return errno.ErrInvalidArgument.WithMessage("nickname must be less than 30 characters")
			}
			return nil
		},
		"Email": func(value any) error {
			return isValidEmail(value.(string))
		},
		"Phone": func(value any) error {
			return isValidPhone(value.(string))
		},
		"Limit": func(value any) error {
			if value.(int64) <= 0 {
				return errno.ErrInvalidArgument.WithMessage("limit must be greater than 0")
			}
			return nil
		},
		"Offset": func(value any) error {
			return nil
		},
	}
}

// ValidateLoginRequest 校验修改密码请求.
func (v *Validator) ValidateLoginRequest(ctx context.Context, request *apiV1.LoginRequest) error {
	return genericValidation.ValidateAllFields(request, v.ValidateUserRules())
}

// ValidateChangePasswordRequest 校验 ChangePasswordRequest 结构体的有效性.
func (v *Validator) ValidateChangePasswordRequest(ctx context.Context, request *apiV1.ChangePasswordRequest) error {
	if request.GetUserID() != contextx.UserID(ctx) {
		return errno.ErrPermissionDenied.WithMessage("The logged-in user `%s` does not match request user `%s`", contextx.UserID(ctx), request.GetUserID())
	}
	return genericValidation.ValidateAllFields(request, v.ValidateUserRules())
}

// ValidateCreateUserRequest 校验 CreateUserRequest 结构体的有效性.
func (v *Validator) ValidateCreateUserRequest(ctx context.Context, request *apiV1.CreateUserRequest) error {
	return genericValidation.ValidateAllFields(request, v.ValidateUserRules())
}

// ValidateUpdateUserRequest 校验更新用户请求.
func (v *Validator) ValidateUpdateUserRequest(ctx context.Context, request *apiV1.UpdateUserRequest) error {
	if request.GetUserID() != contextx.UserID(ctx) {
		return errno.ErrPermissionDenied.WithMessage("The logged-in user `%s` does not match request user `%s`", contextx.UserID(ctx), request.GetUserID())
	}
	return genericValidation.ValidateSelectedFields(request, v.ValidateUserRules(), "UserID")
}

// ValidateDeleteUserRequest 校验 DeleteUserRequest 结构体的有效性.
func (v *Validator) ValidateDeleteUserRequest(ctx context.Context, request *apiV1.DeleteUserRequest) error {
	return genericValidation.ValidateAllFields(request, v.ValidateUserRules())
}

// ValidateGetUserRequest 校验 GetUserRequest 结构体的有效性.
func (v *Validator) ValidateGetUserRequest(ctx context.Context, request *apiV1.GetUserRequest) error {
	if request.GetUserID() != contextx.UserID(ctx) {
		return errno.ErrPermissionDenied.WithMessage("The logged-in user `%s` does not match request user `%s`", contextx.UserID(ctx), request.GetUserID())
	}
	return genericValidation.ValidateAllFields(request, v.ValidateUserRules())
}

// ValidateListUserRequest 校验 ListUserRequest 结构体的有效性.
func (v *Validator) ValidateListUserRequest(ctx context.Context, request *apiV1.ListUserRequest) error {
	return genericValidation.ValidateAllFields(request, v.ValidateUserRules())
}
