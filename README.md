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
 --grpc-gateway_out=allow_delete_body=true # 允许删除请求体
```