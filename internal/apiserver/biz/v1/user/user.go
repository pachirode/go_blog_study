package user

import (
	"context"
	"github.com/pachirode/go_blog_study/internal/apiserver/pkg/conversion"
	"sync"

	"github.com/jinzhu/copier"
	"github.com/pachirode/pkg/authn"
	"github.com/pachirode/pkg/authz"
	"github.com/pachirode/pkg/store/where"
	"github.com/pachirode/pkg/token"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/pachirode/go_blog_study/internal/apiserver/model"
	"github.com/pachirode/go_blog_study/internal/apiserver/store"
	"github.com/pachirode/go_blog_study/internal/pkg/contextx"
	"github.com/pachirode/go_blog_study/internal/pkg/errno"
	"github.com/pachirode/go_blog_study/internal/pkg/known"
	"github.com/pachirode/go_blog_study/internal/pkg/log"
	apiV1 "github.com/pachirode/go_blog_study/pkg/api/apiserver/v1/proto"
)

// UserBiz 定义处理用户请求所需的方法.
type UserBiz interface {
	Create(ctx context.Context, request *apiV1.CreateUserRequest) (*apiV1.CreateUserResponse, error)
	Update(ctx context.Context, request *apiV1.UpdateUserRequest) (*apiV1.UpdateUserResponse, error)
	Delete(ctx context.Context, request *apiV1.DeleteUserRequest) (*apiV1.DeleteUserResponse, error)
	Get(ctx context.Context, request *apiV1.GetUserRequest) (*apiV1.GetUserResponse, error)
	List(ctx context.Context, request *apiV1.ListUserRequest) (*apiV1.ListUserResponse, error)

	UserExpansion
}

// UserExpansion 定义用户操作的扩展方法.
type UserExpansion interface {
	Login(ctx context.Context, request *apiV1.LoginRequest) (*apiV1.LoginResponse, error)
	RefreshToken(ctx context.Context, request *apiV1.RefreshTokenRequest) (*apiV1.RefreshTokenResponse, error)
	ChangePassword(ctx context.Context, request *apiV1.ChangePasswordRequest) (*apiV1.ChangePasswordResponse, error)
	ListWithBadPerformance(ctx context.Context, request *apiV1.ListUserRequest) (*apiV1.ListUserResponse, error)
}

// userBiz 是 UserBiz 接口的实现.
type userBiz struct {
	store store.IStore
	authz *authz.Authz
}

// 确保 userBiz 实现了 UserBiz 接口.
var _ UserBiz = (*userBiz)(nil)

func New(store store.IStore, authz *authz.Authz) *userBiz {
	return &userBiz{store: store, authz: authz}
}

// Login 实现 UserBiz 接口中的 Login 方法.
func (b *userBiz) Login(ctx context.Context, request *apiV1.LoginRequest) (*apiV1.LoginResponse, error) {
	// 获取登录用户的所有信息
	whr := where.F("username", request.GetUsername())
	userM, err := b.store.User().Get(ctx, whr)
	if err != nil {
		return nil, errno.ErrUserNotFound
	}

	// 对比传入的明文密码和数据库中已加密过的密码是否匹配
	if err := authn.Compare(userM.Password, request.GetPassword()); err != nil {
		log.C(ctx).Errorw("Failed to compare password", "err", err)
		return nil, errno.ErrPasswordInvalid
	}

	// 如果匹配成功，说明登录成功，签发 token 并返回
	tokenStr, expireAt, err := token.Sign(userM.UserID)
	if err != nil {
		log.C(ctx).Errorw("Failed to sign token", "err", err)
		return nil, errno.ErrSignToken
	}

	return &apiV1.LoginResponse{Token: tokenStr, ExpireAt: timestamppb.New(expireAt)}, nil
}

// RefreshToken 用于刷新用户的身份验证令牌.
// 当用户的令牌即将过期时，可以调用此方法生成一个新的令牌.
func (b *userBiz) RefreshToken(ctx context.Context, request *apiV1.RefreshTokenRequest) (*apiV1.RefreshTokenResponse, error) {
	tokenStr, expireAt, err := token.Sign(contextx.UserID(ctx))
	if err != nil {
		log.C(ctx).Errorw("Failed to sign token", "err", err)
		return nil, errno.ErrSignToken
	}

	return &apiV1.RefreshTokenResponse{Token: tokenStr, ExpireAt: timestamppb.New(expireAt)}, nil
}

// ChangePassword 实现 UserBiz 接口中的 ChangePassword 方法.
func (b *userBiz) ChangePassword(ctx context.Context, request *apiV1.ChangePasswordRequest) (*apiV1.ChangePasswordResponse, error) {
	userM, err := b.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err
	}

	if err := authn.Compare(userM.Password, request.GetOldPassword()); err != nil {
		log.C(ctx).Errorw("Failed to compare password", "err", err)
		return nil, errno.ErrPasswordInvalid
	}

	userM.Password, _ = authn.Encrypt(request.GetNewPassword())
	if err := b.store.User().Update(ctx, userM); err != nil {
		return nil, err
	}

	return &apiV1.ChangePasswordResponse{}, nil
}

// Create 实现 UserBiz 接口中的 Create 方法.
func (b *userBiz) Create(ctx context.Context, request *apiV1.CreateUserRequest) (*apiV1.CreateUserResponse, error) {
	var userM model.UserM
	_ = copier.Copy(&userM, request)

	if err := b.store.User().Create(ctx, &userM); err != nil {
		return nil, err
	}

	if _, err := b.authz.AddGroupingPolicy(userM.UserID, known.RoleUser); err != nil {
		log.C(ctx).Errorw("Failed to add grouping policy for user", "user", userM.UserID, "role", known.RoleUser)
		return nil, errno.ErrAddRole.WithMessage(err.Error())
	}

	return &apiV1.CreateUserResponse{UserID: userM.UserID}, nil
}

// Update 实现 UserBiz 接口中的 Update 方法.
func (b *userBiz) Update(ctx context.Context, request *apiV1.UpdateUserRequest) (*apiV1.UpdateUserResponse, error) {
	userM, err := b.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err
	}

	if request.Username != nil {
		userM.Username = request.GetUsername()
	}
	if request.Email != nil {
		userM.Email = request.GetEmail()
	}
	if request.Nickname != nil {
		userM.Nickname = request.GetNickname()
	}
	if request.Phone != nil {
		userM.Phone = request.GetPhone()
	}

	if err := b.store.User().Update(ctx, userM); err != nil {
		return nil, err
	}

	return &apiV1.UpdateUserResponse{}, nil
}

// Delete 实现 UserBiz 接口中的 Delete 方法.
func (b *userBiz) Delete(ctx context.Context, request *apiV1.DeleteUserRequest) (*apiV1.DeleteUserResponse, error) {
	// 只有 `root` 用户可以删除用户，并且可以删除其他用户
	// 所以这里不用 where.T()，因为 where.T() 会查询 `root` 用户自己
	if err := b.store.User().Delete(ctx, where.F("userID", request.GetUserID())); err != nil {
		return nil, err
	}

	if _, err := b.authz.RemoveGroupingPolicy(request.GetUserID(), known.RoleUser); err != nil {
		log.C(ctx).Errorw("Failed to remove grouping policy for user", "user", request.GetUserID(), "role", known.RoleUser)
		return nil, errno.ErrRemoveRole.WithMessage(err.Error())
	}

	return &apiV1.DeleteUserResponse{}, nil
}

// Get 实现 UserBiz 接口中的 Get 方法.
func (b *userBiz) Get(ctx context.Context, request *apiV1.GetUserRequest) (*apiV1.GetUserResponse, error) {
	userM, err := b.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err
	}

	return &apiV1.GetUserResponse{User: conversion.UserModelToUserV1(userM)}, nil
}

// List 实现 UserBiz 接口中的 List 方法.
func (b *userBiz) List(ctx context.Context, request *apiV1.ListUserRequest) (*apiV1.ListUserResponse, error) {
	whr := where.P(int(request.GetOffset()), int(request.GetLimit()))
	if contextx.Username(ctx) != known.AdminUsername {
		whr.T(ctx)
	}

	count, userList, err := b.store.User().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	var m sync.Map
	eg, ctx := errgroup.WithContext(ctx)

	// 设置最大并发数量为常量 MaxConcurrency
	eg.SetLimit(known.MaxErrGroupConcurrency)

	// 使用 goroutine 提高接口性能
	for _, user := range userList {
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				converted := conversion.UserModelToUserV1(user)
				m.Store(user.ID, converted)

				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		log.C(ctx).Errorw("Failed to wait all function calls returned", "err", err)
		return nil, err
	}

	users := make([]*apiV1.User, 0, len(userList))
	for _, item := range userList {
		user, _ := m.Load(item.ID)
		users = append(users, user.(*apiV1.User))
	}

	log.C(ctx).Debugw("Get users from backend storage", "count", len(users))

	return &apiV1.ListUserResponse{TotalCount: count, Users: users}, nil
}

// ListWithBadPerformance 是性能较差的实现方式（已废弃）.
func (b *userBiz) ListWithBadPerformance(ctx context.Context, request *apiV1.ListUserRequest) (*apiV1.ListUserResponse, error) {
	whr := where.P(int(request.GetOffset()), int(request.GetLimit()))
	if contextx.Username(ctx) != known.AdminUsername {
		whr.T(ctx)
	}

	count, userList, err := b.store.User().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	users := make([]*apiV1.User, 0, len(userList))
	for _, user := range userList {
		converted := conversion.UserModelToUserV1(user)
		users = append(users, converted)
	}

	log.C(ctx).Debugw("Get users from backend storage", "count", len(users))

	return &apiV1.ListUserResponse{TotalCount: count, Users: users}, nil
}
