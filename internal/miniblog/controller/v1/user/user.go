package user

import (
	"github.com/marmotedu/Miniblog/internal/miniblog/biz"
	"github.com/marmotedu/Miniblog/internal/miniblog/store"
	"github.com/marmotedu/Miniblog/pkg/auth"
)

type UserController struct {
	a *auth.Authz
	b biz.IBiz
}

func New(ds store.IStore, a *auth.Authz) *UserController {
	return &UserController{a: a, b: biz.NewBiz(ds)}
}
