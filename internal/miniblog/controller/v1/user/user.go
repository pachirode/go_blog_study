package user

import (
	"github.com/marmotedu/Miniblog/internal/miniblog/biz"
	"github.com/marmotedu/Miniblog/internal/miniblog/store"
)

type UserController struct {
	b biz.IBiz
}

func New(ds store.IStore) *UserController {
	return &UserController{b: biz.NewBiz(ds)}
}

