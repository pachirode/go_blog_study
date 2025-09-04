package rid

import "github.com/pachirode/pkg/id"

const defaultABC = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type ResourceID string

const (
	// UserID 定义用户资源标识符
	UserID ResourceID = "user"
	// PostID 定义博客资源标识符
	PostID ResourceID = "post"
)

// String 将资源标识符转换为字符串
func (rid ResourceID) String() string {
	return string(rid)
}

// NewResourceID 创建新的资源标识符
func (rid ResourceID) NewResourceID(counter uint64) string {
	// 使用自定义选项生成唯一标识符
	uniqueStr := id.NewCode(
		counter,
		id.WithCodeChars([]rune(defaultABC)),
		id.WithCodeL(6),
		id.WithCodeSalt(Salt()),
	)
	return rid.String() + "-" + uniqueStr
}
