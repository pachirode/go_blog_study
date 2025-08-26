/*
Package contextx 扩展 Go 中 context 功能，允许存储提取用户相关的信息

实例：

	ctx := context.Background()
	ctx = contextx.WithUserID(ctx, "user-id")
	ctx = contextx.WithUserName(ctx, "user-name")

	userID = contextx.UserID(ctx)
*/
package contextx
