package store

import (
	"context"

	genericStore "github.com/pachirode/pkg/store"
	"github.com/pachirode/pkg/store/where"

	"github.com/pachirode/go_blog_study/internal/apiserver/model"
)

type UserStore interface {
	Create(ctx context.Context, obj *model.UserM) error
	Update(ctx context.Context, obj *model.UserM) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.UserM, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.UserM, error)

	UserExpansion
}

// UserExpansion 定义了用户操作的附加方法.
// nolint: iface
type UserExpansion interface{}

// userStore 是 UserStore 接口的实现.
type userStore struct {
	*genericStore.Store[model.UserM]
}

// 确保 userStore 实现了 UserStore 接口.
var _ UserStore = (*userStore)(nil)

// newUserStore 创建 userStore 的实例.
func newUserStore(store *dataStore) *userStore {
	return &userStore{
		Store: genericStore.NewStore[model.UserM](store, NewLogger()),
	}
}
