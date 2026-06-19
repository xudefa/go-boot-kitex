// Package server 基于 Kitex 框架提供 RPC 服务器实现。
//
// 该包封装了 Kitex 服务器，支持地址配置、自定义服务器选项和注册函数管理。
//
// 定义：
//
//   - Server: Kitex RPC 服务器
//   - Option: 服务器配置选项函数
//
// 快速开始:
//
//	// 创建 Kitex 服务器
//	srv := server.New(
//	    server.WithAddress(":8888"),
//	)
//	// 注册服务（通过生成代码创建的服务端）
//	srv.Register(func(svr server.Server) error {
//	    svr.RegisterService(yourServiceInfo, yourHandler)
//	    return nil
//	})
//	// 启动服务器
//	srv.Start()
package server

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/cloudwego/kitex/server"
)

// Server 是 Kitex RPC 服务器，封装了 Kitex 服务器的创建和生命周期管理。
//
// 字段说明:
//   - kitexServer: Kitex 服务器实例
//   - address: 监听地址
//   - opts: Kitex 服务器选项
//   - registerFn: 服务注册函数列表
type Server struct {
	kitexServer server.Server
	address     string
	opts        []server.Option
	registerFn  []func(server.Server) error
}

// Option 是 Kitex 服务器配置选项函数。
type Option func(*Server)

// WithAddress 设置服务器监听地址。
//
// 参数:
//   - addr: 监听地址，如 ":8888"
//
// 返回值:
//   - Option: 服务器配置选项函数
func WithAddress(addr string) Option {
	return func(s *Server) {
		s.address = addr
	}
}

// WithOptions 设置 Kitex 服务器自定义选项。
//
// 参数:
//   - opts: Kitex 服务器选项列表
//
// 返回值:
//   - Option: 服务器配置选项函数
func WithOptions(opts ...server.Option) Option {
	return func(s *Server) {
		s.opts = append(s.opts, opts...)
	}
}

// New 创建新的 Kitex 服务器实例。
//
// 参数:
//   - opts: 可选的配置选项
//
// 返回值:
//   - *Server: 服务器实例
func New(opts ...Option) *Server {
	s := &Server{
		address: ":8888",
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Register 注册一个服务到 Kitex 服务器（如通过生成代码创建的服务端实例）。
//
// 参数:
//   - fn: 注册函数，接受 Kitex server.Server 参数
//
// 返回值:
//   - *Server: 服务器实例，支持链式调用
func (s *Server) Register(fn func(server.Server) error) *Server {
	if s.kitexServer != nil {
		log.Println("warning: Register() called after server started, registration may not take effect")
	}
	s.registerFn = append(s.registerFn, fn)
	return s
}

// parseAddress 解析地址字符串为 TCPAddr
func parseAddress(addr string) (*net.TCPAddr, error) {
	// 处理空地址，使用默认值
	if addr == "" {
		addr = ":8888"
	}

	// 使用 net.SplitHostPort 解析地址（支持 "host:port" 和 ":port" 格式）
	host, port, err := net.SplitHostPort(addr)
	if err == nil {
		portNum, err := strconv.Atoi(port)
		if err != nil {
			return nil, err
		}
		ip := net.ParseIP(host)
		if ip == nil {
			// 如果不是有效IP地址，使用0.0.0.0
			ip = net.IPv4zero
		}
		return &net.TCPAddr{
			IP:   ip,
			Port: portNum,
		}, nil
	}

	// 如果 SplitHostPort 失败，尝试作为纯端口处理（以冒号开头）
	if len(addr) > 0 && addr[0] == ':' {
		port, err := strconv.Atoi(addr[1:])
		if err != nil {
			return nil, err
		}
		return &net.TCPAddr{
			IP:   net.IPv4zero,
			Port: port,
		}, nil
	}

	// 默认返回 :8888
	return &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: 8888,
	}, nil
}

// Start 启动 Kitex 服务器，开始监听并处理 RPC 请求。
//
// 在启动时会：
// 1. 解析监听地址
// 2. 创建 Kitex 服务器实例
// 3. 执行所有注册函数
// 4. 启动服务器（同步阻塞）
//
// 返回值:
//   - error: 启动错误
func (s *Server) Start() error {
	if s.kitexServer != nil {
		return fmt.Errorf("server already started")
	}

	// 解析地址
	addr, err := parseAddress(s.address)
	if err != nil {
		return err
	}

	opts := append(s.opts, server.WithServiceAddr(addr))

	svr := server.NewServer(opts...)

	// 执行注册函数
	for _, fn := range s.registerFn {
		if err := fn(svr); err != nil {
			return err // 注册失败，不设置 s.kitexServer
		}
	}

	// 注册成功后才设置
	s.kitexServer = svr

	// 同步启动服务器，错误会返回
	return svr.Run()
}

// Stop 优雅停止 Kitex 服务器。
//
// 返回值:
//   - error: 停止错误
func (s *Server) Stop() error {
	if s.kitexServer != nil {
		return s.kitexServer.Stop()
	}
	return nil
}

// KitexServer 返回底层的 Kitex 服务器实例，用于高级配置（如注册服务）。
//
// 返回值:
//   - server.Server: Kitex 服务器实例
func (s *Server) KitexServer() server.Server {
	return s.kitexServer
}

// Address 返回服务器配置的监听地址
func (s *Server) Address() string {
	return s.address
}
