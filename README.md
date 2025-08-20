# Go blog study

Go 语言学习项目

### 架构

- `Handler`
  - `API` 接口请求的参数解析，参数校验等功能
- `Biz`
  - 具体业务逻辑实现，根据不同的 `REST` 资源分为不同模块
- `Store`
  - 数据访问层，数据库进行交互