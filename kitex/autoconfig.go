// Package kitex 提供 Kitex 服务端和客户端的自动配置。
//
// 当 kitex.server.enabled=true 时自动启用，从 Environment 中读取 kitex.server.address、kitex.client.address、
// kitex.client.timeout 等配置项，
// 创建并注册 Kitex Server Bean（Bean ID: kitexServer）和 Kitex Client Bean（Bean ID: kitexClient）到 IoC 容器中。
package kitex

import (
	"time"

	"github.com/xudefa/go-boot-kitex/client"
	"github.com/xudefa/go-boot-kitex/server"

	"github.com/xudefa/go-boot/boot"
	"github.com/xudefa/go-boot/condition"
	"github.com/xudefa/go-boot/constants"
	"github.com/xudefa/go-boot/core"
)

// KitexAutoConfiguration Kitex 服务端和客户端的自动配置
//
// 从 Environment 中读取 kitex.server.address、kitex.client.address、kitex.client.timeout 等配置项，
// 创建 Kitex Server 和 Kitex Client 实例并注册到 IoC 容器中。
// 启用条件：kitex.server.enabled=true
type KitexAutoConfiguration struct{}

// init 注册 Kitex 自动配置，由 kitex.server.enabled=true 条件控制
func init() {
	boot.RegisterAutoConfig(&KitexAutoConfiguration{},
		condition.OnProperty(constants.KitexServerEnabled, constants.ConditionTrue),
	)
}

// Configure 执行自动配置逻辑，创建 Kitex Server 和 Client 并注册为 Bean
func (k *KitexAutoConfiguration) Configure(ctx boot.ApplicationContext) error {
	env := ctx.Environment()

	// 配置 Kitex Server
	if addr := env.GetString(constants.KitexServerAddress, ""); addr != "" {
		srv := server.New(server.WithAddress(addr))
		if err := ctx.Register(constants.KitexServerBeanID,
			core.Bean(srv),
			core.Singleton(),
		); err != nil {
			return err
		}
	}

	// 配置 Kitex Client
	if addr := env.GetString(constants.KitexClientAddress, ""); addr != "" {
		timeout := env.GetInt(constants.KitexClientTimeout, constants.DefaultKitexClientTimeout)
		cli := client.New(
			client.WithAddress(addr),
			client.WithTimeout(time.Duration(timeout)*time.Second),
		)
		if err := ctx.Register(constants.KitexClientBeanID,
			core.Bean(cli),
			core.Singleton(),
		); err != nil {
			return err
		}
	}

	return nil
}
