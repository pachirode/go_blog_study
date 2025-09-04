package store

import (
	"context"
	"sync"

	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/pachirode/pkg/store/where"
)

var ProviderSet = wire.NewSet(NewStore, wire.Bind(new(IStore), new(*dataStore)))

var (
	once sync.Once
	// S 全局变量，方便其它包直接调用已初始化好的 datastore 实例.
	S *dataStore
)

// transactionKey 用于在 context.Context 中存储事务上下文的键.
type transactionKey struct{}

// IStore 定义了 Store 层需要实现的方法.
type IStore interface {
	// DB 返回 Store 层的 *gorm.DB 实例，在少数场景下会被用到.
	DB(ctx context.Context, wheres ...where.Where) *gorm.DB
	TX(ctx context.Context, fn func(ctx context.Context) error) error

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

// DB 根据传入的条件（wheres）对数据库实例进行筛选.
// 如果未传入任何条件，则返回上下文中的数据库实例（事务实例或核心数据库实例）.
func (store *dataStore) DB(ctx context.Context, wheres ...where.Where) *gorm.DB {
	db := store.core
	// 从上下文中提取事务实例
	if tx, ok := ctx.Value(transactionKey{}).(*gorm.DB); ok {
		db = tx
	}

	// 遍历所有传入的条件并逐一叠加到数据库查询对象上
	for _, whr := range wheres {
		db = whr.Where(db)
	}
	return db
}

// TX 返回一个新的事务实例.
// nolint: fatcontext
func (store *dataStore) TX(ctx context.Context, fn func(ctx context.Context) error) error {
	return store.core.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			ctx = context.WithValue(ctx, transactionKey{}, tx)
			return fn(ctx)
		},
	)
}

// User 返回一个实现了 UserStore 接口的实例.
func (store *dataStore) User() UserStore {
	return newUserStore(store)
}
