//go:build wireinject
// +build wireinject

package apiserver

import (
	"github.com/google/wire"
	"github.com/pachirode/pkg/authz"

	"github.com/pachirode/go_blog_study/internal/apiserver/biz"
	"github.com/pachirode/go_blog_study/internal/apiserver/pkg/validation"
	"github.com/pachirode/go_blog_study/internal/apiserver/store"
	ginMw "github.com/pachirode/go_blog_study/internal/pkg/middleware/gin"
	"github.com/pachirode/go_blog_study/internal/pkg/server"
)

func InitializeWebServer(*Config) (server.Server, error) {
	wire.Build(
		wire.NewSet(NewWebServer, wire.FieldsOf(new(*Config), "ServerMode")),
		wire.Struct(new(ServerConfig), "*"), // * 表示注入全部字段
		wire.NewSet(store.ProviderSet, biz.ProviderSet),
		ProvideDB, // 提供数据库实例
		validation.ProviderSet,
		wire.NewSet(
			wire.Struct(new(UserRetriever), "*"),
			wire.Bind(new(ginMw.UserRetriever), new(*UserRetriever)),
		),
		authz.ProviderSet,
	)
	return nil, nil
}
