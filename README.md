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

根据租户信息，查询出所属租户的资源数据，为了防止越权查询，租户信息需要从认证的 `Token` 中获取，如果通过认证说明信息是真实可信的

- 认证步骤
    - 用户登录
    - 签发 `Token`
    - `API` 请求
    - 解析 `Token` 获取 `UserID`
    - 查询 `UserID` 的数据

##### 请求参数默认值

借鉴了 `Kubernetes API` 请求参数默认值设置实现，基于 `API` 接口定义文件自动生成默认值

##### API 接口请求参数

对 `API` 请求参数进行校验是 `Web` 开发的核心功能，提高系统可靠性和安全性

- 系统稳定性
    - `API` 收到的用户请求可能是从客户端或者第三方应用发起的，参数可能是错误的；如果没有对参数进行校验，可能导致系统逻辑错误
        - 未校验分页参数
            - 数据库查询性能下降
- 数据的合法性和完整性
    - 可能提交不符合业务要求的数据
        - 必填字段缺失
        - 字符串格式错误
- 增强用户体验
    - 不进行参数校验，错误通常发生在逻辑处理，错误提示可能和用户实际遇到的问题无关
- 提高代码维护性
    - 只需要最外层参数校验
- 服务端可信
    - 不能相信客户端数据，服务器本身需要校验数据

###### 校验参数的方法

- 手动校验
    - 直接在代码里面判断
- 第三方库
    - 直接复用现成的校验逻辑，基于结构体标签来进行校验
    - `go-playground/validator`
    - `asaskevich/govalidator`
    - `ozzo-validation`
- 框架内置校验
    - `gin` 框架支持 `go-playground/validator`, 处理请求数据时，使用 `binding` 标签可以直接解析和校验
- 基于工具生成校验代码
    - `OpenAPI Generator`
    - `gqlgen（GraphQL 工具）`
- 中间件校验
- 一般校验要求
    - 支持自定义复杂校验逻辑
        - 通过创建统一的校验类型
    - 复用已有的参数校验逻辑
    - 灵活通用的校验方法
        - 所有请求接口的校验函数声明为统一的规范格式
    - 检验方式简单易维护

请求参数校验是几乎每个接口都要使用的功能，最理想的情况是通过 `Web` 中间件来校验请求参数

### 认证

应用层主要有一下三种手段来保障应用层安全

- 认证
    - 通过一定手段确认用户身份
    - 常见方式
        - 用户名密码
        - 数字证书
        - 令牌认证
        - 生物识别
- 授权
    - 身份认证之后，确定认证用户资源的访问权限以及操作范围
- 使用 `HTTPS` 协议通讯
    - 在 `HTTP` 基础上通过 `SSL/TLS` 加密通讯实现的安全通信协议，具有数据机密，完整性和身份验证，保障双方通信安全

##### 常用验证

- 基础认证
    - 用户名密码
- 令牌认证
    - 使用 `Token`
    - 优势
        - 无状态
        - 离线验证
        - 提高安全性

### 授权

### `HTTPS` 通讯

`HTTP` 以明文的方式传输数据；`HTTPS` 基于 `HTTP` 增加了 `SSL` 安全层，通过加密通道传输数据

- 数据传输加密
    - 通过 `HTTPS` 传输数据始终以加密的形式存在，确保数据安全
- 身份认证
    - 支持单向认证和双向认证，单向认证用于服务端真实性，双向认证服务端和客户端都需要

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