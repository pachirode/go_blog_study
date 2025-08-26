package conversion

import (
	"github.com/pachirode/pkg/core"

	"github.com/pachirode/go_blog_study/internal/apiserver/model"
	apiV1 "github.com/pachirode/go_blog_study/pkg/api/apiserver/v1/proto"
)

// UserModelToUserV1 将模型层的 UserM（用户模型对象）转换为 Protobuf 层的 User（v1 用户对象）.
func UserModelToUserV1(userModel *model.UserM) *apiV1.User {
	var protoUser apiV1.User
	_ = core.CopyWithConverters(&protoUser, userModel)
	return &protoUser
}

// UserV1ToUserModel 将 Protobuf 层的 User（v1 用户对象）转换为模型层的 UserM（用户模型对象）.
func UserV1ToUserModel(protoUser *apiV1.User) *model.UserM {
	var userModel model.UserM
	_ = core.CopyWithConverters(&userModel, protoUser)
	return &userModel
}
