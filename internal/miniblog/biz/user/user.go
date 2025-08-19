package user

import (
	"context"
	"errors"
	"regexp"

	"github.com/jinzhu/copier"
	"github.com/pachirode/go_blog_study/internal/miniblog/store"
	"github.com/pachirode/go_blog_study/internal/pkg/errno"
	"github.com/pachirode/go_blog_study/internal/pkg/log"
	"github.com/pachirode/go_blog_study/internal/pkg/model"
	v1 "github.com/pachirode/go_blog_study/pkg/api/miniblog/v1"
	"github.com/pachirode/go_blog_study/pkg/auth"
	"github.com/pachirode/go_blog_study/pkg/token"
	"gorm.io/gorm"
)

type UserBiz interface {
	ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error
	Create(ctx context.Context, r *v1.CreateUserRequest) error
	Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error)
	Get(ctx context.Context, username string) (*v1.GetUserResponse, error)
	List(ctx context.Context, offset, limit int) (*v1.ListUserResponse, error)
}

type userBiz struct {
	ds store.IStore
}

var _ UserBiz = (*userBiz)(nil)

func New(ds store.IStore) *userBiz {
	return &userBiz{ds: ds}
}

func (b *userBiz) ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error {
	userM, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		return nil
	}

	if err := auth.Compare(userM.Password, r.OldPassword); err != nil {
		return errno.ErrPasswordIncorrect
	}

	userM.Password, _ = auth.Encrypt(r.NewPassword)
	if err := b.ds.Users().Update(ctx, userM); err != nil {
		return err
	}

	return nil
}

func (b *userBiz) Create(ctx context.Context, r *v1.CreateUserRequest) error {
	var userM model.UserM
	_ = copier.Copy(&userM, r)

	if err := b.ds.Users().Create(ctx, &userM); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'username' ", err.Error()); match {
			return errno.ErrUserAlreadyExist
		}

		return err
	}

	return nil
}

func (b *userBiz) Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error) {
	userM, err := b.ds.Users().Get(ctx, r.Username)
	if err != nil {
		return nil, errno.ErrUserNotFound
	}

	if err := auth.Compare(userM.Password, r.Password); err != nil {
		return nil, errno.ErrPasswordIncorrect
	}

	t, err := token.Sign(r.Username)
	if err != nil {
		return nil, errno.ErrSignToken
	}

	return &v1.LoginResponse{Token: t}, nil
}

func (b *userBiz) Get(ctx context.Context, username string) (*v1.GetUserResponse, error) {
	userM, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrUserNotFound
		}

		return nil, err
	}

	var resp v1.GetUserResponse
	_ = copier.Copy(&resp, userM)

	resp.CreatedAt = userM.CreatedAt.Format("2001-01-01 01:00:00")
	resp.UpdatedAt = userM.UpdatedAt.Format("2001-01-01 01:00:00")

	return &resp, nil
}

func (b *userBiz) List(ctx context.Context, offset, limit int) (*v1.ListUserResponse, error) {
	count, list, err := b.ds.Users().List(ctx, offset, limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list users from storage", "err", err)
		return nil, err
	}

	users := make([]*v1.UserInfo, 0, len(list))
	for _, user := range list {
		users = append(users, &v1.UserInfo{
			Username:  user.Username,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt.Format("2001-01-01 01:00:00"),
			UpdatedAt: user.UpdatedAt.Format("2001-01-01 01:00:00"),
		})
	}

	log.C(ctx).Debugw("Get users from storage", "count", len(users))

	return &v1.ListUserResponse{TotalCount: count, Users: users}, nil
}
