package user

import (
	"context"
	"regexp"

	"github.com/jinzhu/copier"
	"github.com/marmotedu/Miniblog/internal/miniblog/store"
	"github.com/marmotedu/Miniblog/internal/pkg/errno"
	"github.com/marmotedu/Miniblog/internal/pkg/model"
	v1 "github.com/marmotedu/Miniblog/pkg/api/miniblog/v1"
)

type UserBiz interface {
	Create(ctx context.Context, r *v1.CreateUserRequest) error
}

type userBiz struct {
	ds store.IStore
}

var _ UserBiz = (*userBiz)(nil)

func New(ds store.IStore) *userBiz {
	return &userBiz{ds: ds}
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
