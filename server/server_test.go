package server

import (
	"net"
	"testing"
	"time"

	"github.com/cloudwego/kitex/server"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	srv := New(WithAddress(":9999"))
	assert.NotNil(t, srv)
}

func TestNewWithDefaultAddress(t *testing.T) {
	srv := New()
	assert.NotNil(t, srv)
}

func TestWithOptions(t *testing.T) {
	srv := New(
		WithAddress(":1234"),
		WithOptions(),
	)
	assert.NotNil(t, srv)
}

func TestRegister_WithValidFunction_DoesNotPanic(t *testing.T) {
	t.Parallel()
	srv := New(WithAddress(":0"))

	// 测试 Register 不会 panic
	srv.Register(func(ks server.Server) error {
		return nil
	})

	assert.NotNil(t, srv)
}

func TestLifecycleRequiresService(t *testing.T) {
	srv := New(WithAddress(":0"))

	// 没有注册服务就启动会报错
	err := srv.Start()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no service")
}

func TestAddress_CustomAddress_ReturnsConfigured(t *testing.T) {
	t.Parallel()
	srv := New(WithAddress(":9999"))
	assert.Equal(t, ":9999", srv.Address())
}

func TestAddress_DefaultAddress_ReturnsDefault(t *testing.T) {
	t.Parallel()
	srv := New()
	assert.Equal(t, ":8888", srv.Address())
}

func TestStart_DoubleStart_ReturnsError(t *testing.T) {
	srv := New(WithAddress(":0"))

	// 先注册一个空服务（避免 no service 错误）
	srv.Register(func(ks server.Server) error {
		return nil
	})

	// 清理
	t.Cleanup(func() {
		_ = srv.Stop()
	})

	// 第一次启动（会在后台运行）
	go func() {
		_ = srv.Start()
	}()

	// 简单同步：等待服务器启动完成，避免竞态条件；生产测试建议使用 channel 同步
	time.Sleep(100 * time.Millisecond)

	// 第二次启动应该返回错误
	err := srv.Start()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already started")
}

func TestRegister_AfterStart_LogsWarning(t *testing.T) {
	srv := New(WithAddress(":0"))

	// 先注册一个服务，避免 no service 错误
	srv.Register(func(ks server.Server) error {
		return nil
	})

	// 启动服务器（在后台）
	go func() {
		_ = srv.Start()
	}()

	// 简单同步：等待服务器启动完成，避免竞态条件；生产测试建议使用 channel 同步
	time.Sleep(100 * time.Millisecond)

	// 启动后注册应该记录警告（这里主要测试不 panic）
	srv.Register(func(ks server.Server) error {
		return nil
	})

	// 停止服务器
	_ = srv.Stop()
}

func TestParseAddress(t *testing.T) {
	tests := []struct {
		name     string
		addr     string
		wantIP   net.IP
		wantPort int
		wantErr  bool
	}{
		{"empty address", "", net.IPv4zero, 8888, false},
		{"port only", ":1234", net.IPv4zero, 1234, false},
		{"ip:port", "127.0.0.1:5678", net.IPv4(127, 0, 0, 1), 5678, false},
		{"invalid port", ":abc", nil, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAddress(tt.addr)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantIP, got.IP)
			assert.Equal(t, tt.wantPort, got.Port)
		})
	}
}

func TestServerWithAdditionalOptions(t *testing.T) {
	opts := []server.Option{
		// 这里可以添加实际可用的server选项
	}
	srv := New(WithAddress(":0"), WithOptions(opts...))
	assert.NotNil(t, srv)
}

func TestServerStopWithoutStart(t *testing.T) {
	srv := New(WithAddress(":0"))
	err := srv.Stop()
	assert.NoError(t, err)
}
