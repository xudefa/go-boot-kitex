# go-boot-kitex

[![Go Version](https://img.shields.io/github/go-mod/go-version/xudefa/go-boot-kitex)](https://go.dev/) [![License](https://img.shields.io/github/license/xudefa/go-boot-kitex)](./LICENSE) [![Build Status](https://img.shields.io/github/actions/workflow/status/xudefa/go-boot-kitex/test.yml?branch=master)](https://github.com/xudefa/go-boot-kitex/actions) [![Go Reference](https://pkg.go.dev/badge/github.com/xudefa/go-boot-kitex.svg)](https://pkg.go.dev/github.com/xudefa/go-boot-kitex) [![Go Report Card](https://goreportcard.com/badge/github.com/xudefa/go-boot-kitex)](https://goreportcard.com/report/github.com/xudefa/go-boot-kitex)

基于 [go-boot](https://github.com/xudefa/go-boot) 的 Kitex RPC 框架集成模块。将 Kitex 无缝集成到 go-boot 的 IoC 容器和自动配置体系中,提供声明式的 RPC 服务器和客户端能力。

> 设计理念:遵循 go-boot 的开发规范,将 Kitex Server 作为 RPC 服务实现,通过自动配置实现零代码启动 RPC 服务。

## 整体架构

```
┌───────────────────────────────────────────────────────────────────────┐
│                    go-boot ApplicationContext                         │
│  ┌───────────┐ ┌──────────────┐ ┌───────────┐ ┌───────────┐           │
│  │ Container │ │  Environment │ │ Lifecycle │ │ EventBus  │           │
│  └───────────┘ └──────────────┘ └───────────┘ └───────────┘           │
│                       ┌─────────────────────┐                         │
│                       │ AutoConfig Registry │                         │
│                       └─────────────────────┘                         │
└───────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
                    ┌───────────────────────────────┐
                    │     go-boot-kitex Starter     │
                    │  ┌─────────────────────────┐  │
                    │  │ Kitex Server Bean       │  │
                    │  │ Kitex Client Bean       │  │
                    │  │ Service Registration    │  │
                    │  │ Middleware Chain        │  │
                    │  └─────────────────────────┘  │
                    └───────────────────────────────┘
```

## 目录

- [快速开始](#快速开始)
- [功能特性](#功能特性)
- [Kitex 服务器](#kitex-服务器)
- [Kitex 客户端](#kitex-客户端)
- [配置选项](#配置选项)
- [项目结构](#项目结构)
- [开发指南](#开发指南)
- [贡献](#贡献)
- [许可证](#许可证)

## 快速开始

### 安装

```bash
# 安装核心框架
go get github.com/xudefa/go-boot

# 安装 Kitex 集成模块
go get github.com/xudefa/go-boot-kitex
```

### 最小示例

```go
package main

import (
    "context"

    "github.com/xudefa/go-boot/boot"
    kitexserver "github.com/xudefa/go-boot-kitex/server"
    pb "your/kitex/package"
)

func main() {
    app, err := boot.NewApplication(
        boot.WithAppName("my-rpc-app"),
        boot.WithVersion("1.0.0"),
        boot.WithProperty("kitex.server.enabled", "true"),
        boot.WithProperty("kitex.server.address", ":8888"),
    )
    if err != nil {
        panic(err)
    }
    defer app.Stop()

    // 注册 Kitex 服务
    server := app.Container().Get("kitexServer").(*kitexserver.Server)
    server.RegisterService(pb.NewEchoServiceHandler(&echoServiceImpl{}))

    // 启动应用（自动启动 Kitex 服务器）
    app.Start()

    // 等待终止信号
    app.WaitForSignal()
}

// Kitex 服务实现
type echoServiceImpl struct{}

func (s *echoServiceImpl) Echo(ctx context.Context, req *pb.Request) (resp *pb.Response, err error) {
    return &pb.Response{Message: req.Message}, nil
}
```

## 功能特性

| 特性 | 说明 |
|------|------|
| Kitex 集成 | 将 Kitex Server/Client 注册为 Bean,支持依赖注入 |
| 自动配置 | 通过 `kitex.server.enabled=true` 自动启动 RPC 服务 |
| 优雅启停 | 支持优雅关闭和生命周期管理 |
| 声明式服务 | 支持通过 IDL 生成服务并注册 |
| 中间件支持 | 支持全局和路由级中间件配置 |
| 配置驱动 | 地址、超时、重试等均可通过配置控制 |

## Kitex 服务器

### 基本服务器

```go
import kitexserver "github.com/xudefa/go-boot-kitex/server"

// 创建 Kitex 服务器
server := kitexserver.New(
    kitexserver.WithAddress(":8888"),
)

// 注册服务
server.RegisterService(pb.NewEchoServiceHandler(&echoServiceImpl{}))
```

### 服务实现示例

```go
type UserService struct{}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    return &pb.User{
        Id:    req.Id,
        Name:  "John Doe",
        Email: "john@example.com",
    }, nil
}

func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
    return &pb.ListUsersResponse{
        Users: []*pb.User{
            {Id: 1, Name: "Alice", Email: "alice@example.com"},
            {Id: 2, Name: "Bob", Email: "bob@example.com"},
        },
    }, nil
}
```

## Kitex 客户端

### 基本客户端

```go
import kitexclient "github.com/xudefa/go-boot-kitex/client"

// 创建 Kitex 客户端
client := kitexclient.New(
    kitexclient.WithAddress("localhost:8888"),
    kitexclient.WithTimeout(5*time.Second),
)

// 创建服务客户端
echoClient := pb.NewEchoClient(client.Client())
resp, err := echoClient.Echo(context.Background(), &pb.Request{
    Message: "Hello Kitex",
})
```

### 带重试的客户端

```go
import kitexclient "github.com/xudefa/go-boot-kitex/client"

client := kitexclient.New(
    kitexclient.WithAddress("localhost:8888"),
    kitexclient.WithTimeout(5*time.Second),
    kitexclient.WithRetry(3),
)
```

## 配置选项

通过 `boot.WithProperty()` 或配置文件设置:

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `kitex.server.enabled` | `false` | 是否启用 Kitex 服务器 |
| `kitex.server.address` | `:8888` | 服务器监听地址 |
| `kitex.client.address` | `localhost:8888` | 客户端连接地址 |
| `kitex.client.timeout` | `5` | 客户端超时(秒) |

### 示例配置

```yaml
# application.yml
kitex:
  server:
    enabled: true
    address: ":8888"
  client:
    address: "localhost:8888"
    timeout: 5
```

## 项目结构

```
go-boot-kitex/
├── autoconfig.go           # 自动配置注册
├── client/                 # Kitex 客户端
│   ├── client.go           # Kitex 客户端实现
│   └── options.go          # 客户端选项配置
├── server/                 # Kitex 服务器
│   ├── server.go           # Kitex 服务器实现
│   ├── options.go          # 服务器选项配置
│   └── middleware/         # 中间件
│       ├── logging.go      # 日志中间件
│       └── recovery.go     # 恢复中间件
├── README.md
├── LICENSE
└── go.mod
```

## 开发指南

### 构建

```bash
go build ./...
```

### 测试

```bash
go test ./...
go test -cover ./...       # 带覆盖率
go test -race ./...        # 数据竞争检测
```

### 代码规范

```bash
go fmt ./...
golangci-lint run
```

## 贡献

欢迎提交 Issue 和 Pull Request!详细贡献指南请参阅 [CONTRIBUTING.md](./CONTRIBUTING.md)。

## 许可证

本项目采用 MIT 许可证 — 详情请参阅 [LICENSE](./LICENSE) 文件。