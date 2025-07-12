package user

import (
	"github.com/marmotedu/Miniblog/internal/miniblog/biz"
	"github.com/marmotedu/Miniblog/internal/miniblog/store"
	"github.com/marmotedu/Miniblog/pkg/auth"
	pb "github.com/marmotedu/Miniblog/pkg/proto/miniblog/v1"
)

type UserController struct {
	a *auth.Authz
	b biz.IBiz
	pb.UnimplementedMiniBlogServer
}

func New(ds store.IStore, a *auth.Authz) *UserController {
	return &UserController{a: a, b: biz.NewBiz(ds)}
}
