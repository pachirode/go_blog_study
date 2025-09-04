# Wire 依赖注入

项目中经常需要创建各种实例，这些实例可以依赖于其他实例；通常需要创建依赖实例，然后将他们作为参数传递给其他实例的构造函数
如果项目复杂，依赖实例的数量和依赖关系也会发生变化，还需要修改构造函数的参数列表

- `Provider`
- `ProviderSet`
- `Injector`

### 依赖注入

依赖注入是一种设计模式，用于管理对象之间的依赖关系，将对象的依赖关系从代码中分离出来。对象不再负责创建它所依赖的对象，由外部容器来负责创建和管理对象之间的依赖

- 实现
    - 构造函数注入
        - 对象初始化时，将所需的依赖作为参数传递给构造函数
        - 最常用，确保对象在创建时依赖关系被完整注入并初始化
    - 属性注入 / 字段注入
        - 通过将依赖声明为结构体的公开字段或属性，实例化后通过外部赋值方式注入依赖
        - 注入逻辑和对象创建逻辑解耦
    - 方法注入
        - 将依赖作为参数传递给对象，而不是长期持有
        - 特定场合短时间使用

### `Provider`

创建某种实例的构造函数，本质上时为依赖注入提供实际对象的工厂方法或者生成器，一个返回指定类型实例的函数
如果该函数本身存在依赖，`Wire` 会根据依赖关系自动解析并注入这些参数
`Provider` 函数返回的实例将被加入到依赖图

```go
package main

import "github.com/google/wire"

type Config struct {
	Name string
}

func NewConfig() *Config {
	return &Config{
		Name: "MyApp",
	}
}

type Service struct {
	Config *Config
}

func NewService(config *Config) *Service {
	return &Service{
		Config: config,
	}
}

var ProviderSet = wire.NewSet(NewConfig, NewService)
```

### `ProviderSet`

`ProviderSet` 是 `Provider` 的集合，用于声明应用程序所有的依赖是如何创建的

```go
var ProviderSet = wire.NewSet(
ProvideDB, // 提供数据库实例
ProvideUserRetriever, // 提供用户检索服务
// 将 UserRetriever 实现绑定到接口 mw.UserRetriever
wire.Bind(new(mw.UserRetriever), new(*UserRetriever)),
// 将 Config 的 MySQLOptions 字段注入到依赖树中
wire.FieldsOf(new(Config), "MySQLOptions"),
)
```

### `Injector`

依赖注入的入口函数，负责初始化所有对象及其所有依赖，并返回完整实例
由 `wire.Build` 的声明自动生成，通过调用相应的 `Provider`，按照依赖关系完成实例化和注入
只需要调用 `Injector`，无需管理依赖关系

```go
// wire.go  
// +build wireinject  

package main  

import "github.com/google/wire"  

// Wire 的声明入口：告诉 Wire 如何组织依赖  
func InitializeService() (*Service, error) {  
    panic(wire.Build(ProviderSet))  
}
```

### 生成依赖函数

`wire .` 该命令会生成 `wire_gen.go` 文件，其中依赖有顺序复杂，由 `Wire` 框架统一管理，生成代码自动注入以下行
```go
// 如果依赖更新，只需要更新 wire.go 文件，然后执行 go generate ./
//go:generate go run -mod=mod github.com/google/wire/cmd/wire
```
> `Go Workspace` 模式下，需要修改 `-mod=mod` 为 `-mod=readonly`
> 该模式下只能 `-mod=readonly/vendor` 否则错误