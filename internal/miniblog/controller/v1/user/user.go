package user

import (
	"github.com/pachirode/go_blog_study/internal/miniblog/biz"
	"github.com/pachirode/go_blog_study/internal/miniblog/store"
	"github.com/pachirode/go_blog_study/pkg/auth"
	pb "github.com/pachirode/go_blog_study/pkg/proto/miniblog/v1"
)

type UserController struct {
	a *auth.Authz
	b biz.IBiz
	pb.UnimplementedMiniBlogServer
}

func New(ds store.IStore, a *auth.Authz) *UserController {
	return &UserController{a: a, b: biz.NewBiz(ds)}
}
