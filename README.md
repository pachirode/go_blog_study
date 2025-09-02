# Go blog study

Go 语言学习项目

### 热加载 Go 应用

自动编译并重启程序

```bash
go install github.com/air-verse/air@latest

# 配置 air .air.toml

air
```

### 添加版权

```bash
go install github.com/nishanths/license/v5@1atest
license -n 'pachirode' -o 'LICENSE' mit
```

### 应用配置

用于对应用程序进行配置

- 不修改代码的情况下切换不同的环境
- 提高安全，用来保留不能硬编码的敏感信息
- 提高发布效率和应用程序的稳定性

配置来源

- 命令行选项
- 命令行参数
- 配置文件
- 分布式配置存储服务

### 实现服务

外部一般通过 `HTTP` 接口访问服务，而内部系统则基于性能和便捷，使用 `gRPC` 接口通讯
一般情况下，服务只需要对外提供一种类型的通信协议，但是有些服务没有 `API` 网关，内置一个反向代理服务器，将请求转化为 `gRPC`

- `Gin` 框架实现 `HTTP` 服务
- `gRPC` 实现 `RPC` 服务
- `grpc-gateway` 实现 `HTTP` 反向代理，将 `HTTP` 转换为 `gRPC`
- 支持开启 `TLS` 认证

##### `gRPC`

使用 `reflection.Register()` 注册反射服务，从而使得服务支持反射功能。反射功能允许客户端动态查询 `gRPC`
服务器上的服务信息且无需拥有 `Protobuf` 文件

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
grpcurl -plaintext localhost:6666 list # 需要先启动 miniblog gRPC服务，可以稍后测试
grpc.reflection.v1.ServerReflection
grpc.reflection.v1alpha.ServerReflection
v1.MiniBlog
```

##### `grpc-gateway`

传统 `gRPC` 应用中，通常会创建一个 `gRPC` 客户端和 `gRPC` 服务器进行交互；`RESTful`
服务，为每一个远程方法暴露 `RESTful API`，接收到请求之后将其转换为 `gRPC` 请求，并调用 `gRPC` 服务

`protoc` 的一个插件，读取 `gRPC` 服务定义，并生成反向代理服务器
反向代理服务器根据服务定义中的 `google.api.http` 注释生成，能给将 `RESTful` 请求映射为 `gRPC`

插件

- `protoc-gen-grpc-gateway`
    - 生成 `HTTP/REST API` 反向代理代码
- `protoc-gen-openapiv2`
    - 生成 `OpenAPI v2(Swagger)` 定义文件

```bash
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.24.0
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.24.0
 --grpc-gateway_out=allow_delete_body=true # 允许删除请求体
```

实现反向代理服务器步骤

- 给 `gRPC` 服务添加 `HTTP` 映射规则
- 生成反向代理代码
- 实现反向代理服务器
- `HTTP` 请求测试

##### `gin`

实现复合 `REST` 规范的 `HTTP` 服务器

### 中间件

位于应用程序请求-响应处理循环中的特殊函数，可以在请求到达业务立即之前修改、处理请求，或者在响应返回客户端之前修改、处理响应。中间件可以分为客户端中间件和服务端中间件
核心作用是对请求和响应进行预处理，后处理或者监控

- 认证和授权
- 日志记录
- 错误处理
- 性能监视
- 提高代码复用
- 提高代码模块化
- 方便处理全局功能
- 支持中间件链式调用

##### `gPRC` 添加拦截器

是一个 `Web` 中间件拦截器
`gRPC` 通信模式分为 `Unary` 和 `Streaming` 两种模式，因此拦截器也分为两种，还可以分为服务端和客户端

- `UnaryInterceptor`
    - 一元拦截器
- `StreamInterceptor`
    - 流式拦截器

- 请求 `ID` 拦截器
    - 每次请求注入一个唯一的 `RequestID` 或者 `TraceID`，开发人员即可快速定位与该请求相关的所有日志记录
        - 请求注入 `ID`
        - 日志中打印 `ID`
    - 将 `ID` 保存到 `gRPC` 元数据中，元数据可能很少被处理因此还需要使用日志打印

### 业务层

使用三层架构，`Handler` 依赖 `Biz` 层，`Biz` 层依赖 `Model` 层, `Model` 层依赖数据库

##### `Store`

`Store` 层依赖一些数据类型，这些数据类型实际上就是 `GORM Model`

###### 根据数据库表生成 `Model` 文件

官方提供了 `Gen` 工具以及 `gorm.io/gen`，用于读取数据库表结构并自动生成 `GORM Model`

- 自定义结构体名
- 自定义结构体字段名
- 自定义 `GORM` 标签
- 自定义结构体字段类型
- 自定义结构体字段注释
- 自定义代码生成路径

###### `SQLite`

`In-Memory Mode` 内存模式使用 `SQLite` 内存数据库
`Sqlite` 有些字段不支持，因此自动生成代码使用 `FieldWithTypeTag = false`

###### `GORM` 钩子

添加了数据库表 `xxxID` 字段用来自动生成钩子，用来生成并保留记录的唯一标识符
> 用户名本质上是用户提供的信息，如果依赖用户名作为标识，修改用户名时会造成影响

###### 唯一标识

开发过程中需要为每一条 `REST` 资源生成唯一标识，用来唯一定位该资源

- 冲突问题
    - 业务范围内不冲突，尤其是分布式
- 性能瓶颈
    - 高并发系统中，大量请求同时生成 `ID`，需要选择性能较高的算法（雪花），避免数据库自增主键
- 信息泄漏
    - `ID` 不能透露敏感信息，避免使用主键
- 唯一性范围
    - 如果使用时间戳生成，需要判断唯一性范围，全局唯一，表唯一
- 可扩展
- 数据库主键
    - 不会使用数据库自增主键作为 `ID`，会暴露系统数据规模，容易模拟
- `UUID`
    - 通用唯一标识符，长度较长不易记住
- 雪花算法
    - 分布式 `ID` 生成算法，通过时间戳，机器编号和自增序列号组成唯一标识
    - `github.com/sony/sonyflake`
- 数据库自增 `ID` 配合随机化
    - 自增主键加一个随机后缀
- 基于时间戳自定义生成
    - 可预测，易冲突
- 分布式 `ID`

##### `Biz`

- 标准资源 `CURD` 接口
- 扩展接口
    - 用户登录
    - `Token` 刷新
    - 密码修改

##### `Handler`

通过调用 `Biz` 层完成业务逻辑处理

### 请求处理

常使用的是对请求进行鉴权，请求参数设置默认值和参数合法性校验

##### 添加 `Bypass` 认证中间件

# 约定的命名方式

- `XXXOr`
    - `Or` 表示或者含义，暗示这个函数会有两种或者多种选择中一种可能性
- `MustXXX`
    - 函数必须成功，如果失败直接 `panic`; 表示不可恢复的错误或者强制性的操作
- `XXXOrDie`
    - 明确表示失败会导致程序退出，通常调用 `os.Exit`
        - `InitOrDie`
        - `RunOrDie`
- `XXXOrPanic`
    - 操作失败直接 `panic`
- `TryXXX`
    - 尝试操作，如果失败则返回错误，通常不会直接退出，由业务层决定
    - `TryParse`
    - `TryConnect`
- `EnsureXXX`
    - 常用于确保某个操作成功，未成功处理失败的情况，常伴随日志记录或者其他逻辑处理