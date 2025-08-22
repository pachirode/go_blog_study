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

### 架构

- `Handler`
    - `API` 接口请求的参数解析，参数校验等功能
- `Biz`
    - 具体业务逻辑实现，根据不同的 `REST` 资源分为不同模块
- `Store`
    - 数据访问层，数据库进行交互