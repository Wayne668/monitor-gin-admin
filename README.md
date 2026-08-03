# Monitor-Gin-Admin

> 基于 Golang + Gin + GORM 2.0 + Casbin 2.0 + Wire DI 后端，搭配 Ant Design React 前端的轻量级、灵活、优雅且功能齐全的 RBAC 脚手架。

## 项目结构

```text
├── server/    # 后端服务（Gin + GORM + Casbin + Wire）
└── web/       # 前端项目（Ant Design React + Vite）
```

## 功能特性

### 后端（server）

- :scroll: 优雅实现 `RESTful API`，采用接口化编程范式
- :house: 清晰简洁的模块化架构，代码结构一目了然
- :rocket: 基于 `GIN` 框架，集成丰富中间件（身份认证、跨域、日志、限流、链路追踪、权限控制、容错、压缩等）
- :closed_lock_with_key: 集成 `Casbin` 权限框架，灵活精准的 RBAC 权限控制
- :page_facing_up: 基于 `GORM 2.0` ORM 框架，优雅处理数据库操作
- :electric_plug: 采用 `WIRE` 依赖注入，简化模块依赖关系
- :memo: 基于 `Zap` 日志框架，配合 Context 链路追踪
- :key: 整合 `JWT` 认证机制，安全可靠
- :microscope: 自动集成 `Swagger` 接口文档
- :wrench: 完善的单元测试体系，基于 `testify` 框架
- :100: 无状态设计，支持水平扩展，搭配 Redis 实现动态权限管理

### 前端（web）

- :gem: **Neat Design**: 遵循 Ant Design 设计规范
- :triangular_ruler: **Common Templates**: 企业级应用典型模板
- :rocket: **State of The Art Development**: React/umi/dva/antd 最新开发栈
- :cn: **International**: 内置 i18n 国际化方案
- :closed_lock_with_key: **RBAC**: 支持 RBAC 权限管理

## 环境准备

### 后端

- [Go](https://golang.org/) 1.19+
- [Wire](github.com/google/wire) `go install github.com/google/wire/cmd/wire@latest`
- [Swag](github.com/swaggo/swag) `go install github.com/swaggo/swag/cmd/swag@latest`

### 前端

- Node.js v16.20.2（推荐使用 [nvm](https://github.com/nvm-sh/nvm) 管理 node 版本）

## 快速开始

### 后端

```bash
cd server

# 启动服务
make start
# 或
go run main.go start

# 编译服务
make build
# 或
go build -ldflags "-w -s -X main.VERSION=v1.0.0" -o server

# 生成 Swagger 文档
make swagger
# 或
swag init --parseDependency --generalInfo ./main.go --output ./internal/swagger

# 生成依赖注入代码
make wire
# 或
wire gen ./internal/wirex
```

> 通过更改 `configs/dev/server.toml` 配置文件中的 `MenuFile = "menu_cn.json"` 可以切换到中文菜单

### 前端

```bash
cd web

# 安装依赖
npm install
# 或
yarn

# 启动项目
vite --mode dev

# 构建项目
vite build --mode prod

# 代码风格检查
npm run lint
```

## License

Apache License 2.0
