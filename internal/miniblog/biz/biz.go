package biz

import (
	"github.com/pachirode/go_blog_study/internal/miniblog/biz/user"
	"github.com/pachirode/go_blog_study/internal/miniblog/store"
)

type IBiz interface {
	Users() user.UserBiz
}

var _ IBiz = (*biz)(nil)

type biz struct {
	ds store.IStore
}

func NewBiz(ds store.IStore) *biz {
	return &biz{ds: ds}
}

func (b *biz) Users() user.UserBiz {
	return user.New(b.ds)
}
