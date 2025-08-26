package store

import (
	"gorm.io/gorm"
	"sync"
)

var (
	once sync.Once
	// S 全局变量，方便其它包直接调用已初始化好的 datastore 实例.
	S *dataStore
)

type IStore interface {
	User() UserStore
}

type dataStore struct {
	core *gorm.DB

	// 可以根据需要添加其他数据库实例
	// fake *gorm.DB
}

// 确保 datastore 实现了 IStore 接口.
var _ IStore = (*dataStore)(nil)

// NewStore 创建一个 IStore 类型的实例.
func NewStore(db *gorm.DB) *dataStore {
	// 确保 S 只被初始化一次
	once.Do(func() {
		S = &dataStore{db}
	})

	return S
}

// User 返回一个实现了 UserStore 接口的实例.
func (store *dataStore) User() UserStore {
	return newUserStore(store)
}
