// Package client 基于 Kitex 框架提供 RPC 客户端实现。
//
// 该包封装了 Kitex 客户端，支持地址配置、超时设置和自定义客户端选项。
// Kitex 客户端通常需要通过生成代码创建后使用 SetClient 设置。
//
// 定义：
//
//   - Client: Kitex RPC 客户端
//   - Option: 客户端配置选项函数
//
// 快速开始:
//
//	// 创建客户端包装器
//	cli := client.New(
//	    client.WithAddress("localhost:8888"),
//	)
//	// 使用生成代码创建客户端后设置
//	// kcli := echoservice.NewClient("target=ip://:8888")
//	// cli.SetClient(kcli)
//	// cli.Connect()
package client

import (
	"fmt"
	"time"

	"github.com/cloudwego/kitex/client"
)

// Client 是 Kitex RPC 客户端，封装了 Kitex 客户端的连接和生命周期管理。
//
// 字段说明:
//   - client: Kitex 客户端实例
//   - address: 连接地址
//   - timeout: 连接超时时间
//   - opts: Kitex 客户端选项
type Client struct {
	client  client.Client
	address string
	timeout time.Duration
	opts    []client.Option
}

// Option 是 Kitex 客户端配置选项函数。
type Option func(*Client)

// WithAddress 设置客户端连接的目标地址
// 注意：kitex 客户端使用服务发现，地址格式通常为 "target=ip://127.0.0.1:8888"
// 或在创建客户端后通过 SetClient 设置由生成代码创建的客户端
func WithAddress(addr string) Option {
	return func(c *Client) {
		c.address = addr
	}
}

// WithClientOptions 设置 Kitex 客户端自定义选项。
//
// 参数:
//   - opts: Kitex 客户端选项列表
//
// 返回值:
//   - Option: 客户端配置选项函数
func WithClientOptions(opts ...client.Option) Option {
	return func(c *Client) {
		c.opts = append(c.opts, opts...)
	}
}

// WithTimeout 设置客户端超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// Timeout 返回客户端配置的超时时间
func (c *Client) Timeout() time.Duration {
	return c.timeout
}

// New 创建一个新的 Kitex 客户端包装器
// 注意：kitex 客户端需要服务名（serviceInfo）才能创建
// 通常流程：
// 1. 使用生成代码创建客户端：cli := echoservice.NewClient("target=ip://:8888")
// 2. 使用 SetClient 设置到底层：client.SetClient(cli)
func New(opts ...Option) *Client {
	c := &Client{
		address: "localhost:8888",
		timeout: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// SetClient 设置底层 kitex 客户端（通常由生成代码创建后设置）
func (c *Client) SetClient(cli client.Client) {
	c.client = cli
}

// Client 返回底层的 Kitex 客户端实例。
//
// 返回值:
//   - client.Client: Kitex 客户端实例
func (c *Client) Client() client.Client {
	return c.client
}

// Connect 连接到 Kitex 服务
// 注意：Kitex 客户端通常需要通过生成代码创建
// 推荐方式：使用生成代码创建客户端后调用 SetClient()
// 如果未通过 SetClient() 设置客户端，返回错误提示
func (c *Client) Connect() error {
	if c.client != nil {
		return nil // 如果客户端已设置（通过 SetClient），直接返回
	}
	// 提示用户需要通过生成代码或 SetClient 设置客户端
	return fmt.Errorf("kitex client not initialized; use SetClient() with generated client")
}

// Address 返回客户端配置的连接地址
func (c *Client) Address() string {
	return c.address
}

// Close 关闭客户端连接
// 注意：kitex 的 Client 接口没有 Close 方法，但实际实现有
// 这里使用类型断言来调用 Close（如果支持）
func (c *Client) Close() error {
	if c.client != nil {
		type closer interface {
			Close() error
		}
		if cl, ok := c.client.(closer); ok {
			return cl.Close()
		}
	}
	return nil
}
