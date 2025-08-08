package rustfs

import (
	"crypto/tls"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/http"
	"time"
)

// Client  封装了 MinIO 客户端，用于操作 RustFS
type Client struct {
	client *minio.Client
}

// 用于处理自签名证书的场景
// getInsecureTransport 返回一个不验证 SSL 证书的 HTTP 传输对象
func getInsecureTransport() *http.Transport {
	return &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		DisableCompression:    false,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// NewRustFSClient 创建新的 RustFS 客户端（向后兼容）, 使用 SSL 证书，默认不跳过SSL证书验证
// endpoint: RustFS 服务端点
// accessKey: 访问密钥
// secretKey: 秘密密钥
func NewRustFSClient(endpoint, accessKey, secretKey string) (*Client, error) {
	// 默认不跳过 SSL 证书验证
	return NewRustFSClientWithSSLOptions(endpoint, accessKey, secretKey, true, false)
}

// NewRustFSClientWithSSLOptions 创建新的 RustFS 客户端，支持 SSL 证书验证选项
// endpoint: RustFS 服务端点
// accessKey: 访问密钥
// secretKey: 秘密密钥
// useSSL: 是否使用 SSL 连接
// skipSSLVerify: 是否跳过 SSL 证书验证（用于自签名证书）
func NewRustFSClientWithSSLOptions(endpoint, accessKey, secretKey string, useSSL bool, skipSSLVerify bool) (*Client, error) {
	// 初始化 MinIO 客户端，配置跳过 SSL 证书验证
	options := &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	}

	// 如果使用 SSL 且需要跳过证书验证（自签名证书场景）
	if useSSL && skipSSLVerify {
		options.Transport = getInsecureTransport()
	}

	minioClient, err := minio.New(endpoint, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %v", err)
	}

	return &Client{
		client: minioClient,
	}, nil
}
