package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	cli := New(WithAddress("localhost:9999"))
	assert.NotNil(t, cli)
	// Client() 初始为 nil，需要 SetClient() 后才设置
	assert.Nil(t, cli.Client())
}

func TestNewClientWithDefaultAddress(t *testing.T) {
	cli := New()
	assert.NotNil(t, cli)
}

func TestWithAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
	}{
		{"localhost", "localhost:8888"},
		{"ip address", "127.0.0.1:9999"},
		{"only port", ":1234"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := New(WithAddress(tt.address))
			assert.NotNil(t, cli)
		})
	}
}

func TestWithClientOptions(t *testing.T) {
	cli := New(
		WithAddress("localhost:8888"),
		WithClientOptions(),
	)
	assert.NotNil(t, cli)
}

func TestSetClientAndClose_WithoutClient_HandlesGracefully(t *testing.T) {
	t.Parallel()
	cli := New(WithAddress("localhost:0"))

	// 测试 SetClient（实际使用中由生成代码创建后设置）
	// 这里无法创建真实的 Kitex client，所以测试 SetClient 不 panic
	cli.SetClient(nil) // 应该不 panic

	// 测试 Close（没有客户端时应该不报错）
	err := cli.Close()
	assert.NoError(t, err)

	// 测试 Connect（没有客户端时应该返回错误）
	err = cli.Connect()
	assert.Error(t, err)
	assert.Equal(t, "kitex client not initialized; use SetClient() with generated client", err.Error())
}

func TestClientMethodChain(t *testing.T) {
	// 测试链式调用
	cli := New(
		WithAddress(":9999"),
		WithClientOptions(),
	)
	assert.NotNil(t, cli)

	// 测试多次 WithClientOptions
	cli2 := New(
		WithAddress(":8888"),
		WithClientOptions(),
		WithClientOptions(),
	)
	assert.NotNil(t, cli2)
}

func TestAddress_CustomAddress_ReturnsConfigured(t *testing.T) {
	t.Parallel()
	cli := New(WithAddress("localhost:9999"))
	assert.Equal(t, "localhost:9999", cli.Address())
}

func TestAddress_DefaultAddress_ReturnsDefault(t *testing.T) {
	t.Parallel()
	cli := New()
	assert.Equal(t, "localhost:8888", cli.Address())
}

func TestWithTimeout_SetsTimeout_ReturnsCorrectValue(t *testing.T) {
	t.Parallel()
	cli := New(WithTimeout(10 * time.Second))
	assert.NotNil(t, cli)
	assert.Equal(t, 10*time.Second, cli.Timeout())
}

func TestConnect_WithoutClient_ReturnsError(t *testing.T) {
	t.Parallel()
	cli := New(WithAddress("localhost:9999"))
	err := cli.Connect()
	assert.Error(t, err)
	assert.Equal(t, "kitex client not initialized; use SetClient() with generated client", err.Error())
}

func TestConnect_Idempotent_ReturnsSameError(t *testing.T) {
	t.Parallel()
	cli := New(WithAddress("localhost:9999"))
	// 第一次调用
	err1 := cli.Connect()
	assert.Error(t, err1)

	// 第二次调用应该返回相同的错误
	err2 := cli.Connect()
	assert.Error(t, err2)
	assert.Equal(t, err1.Error(), err2.Error())
}

func TestClientWithTimeout(t *testing.T) {
	timeout := 15 * time.Second
	cli := New(
		WithAddress("localhost:9999"),
		WithTimeout(timeout),
	)

	assert.Equal(t, "localhost:9999", cli.Address())
	assert.Equal(t, timeout, cli.Timeout())
}

func TestClientWithClientOptions(t *testing.T) {
	cli := New(
		WithAddress("localhost:9999"),
		WithClientOptions(), // Empty options, but should work
	)

	assert.NotNil(t, cli)
	assert.Equal(t, "localhost:9999", cli.Address())
}
