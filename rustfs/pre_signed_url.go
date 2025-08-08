package rustfs

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// GetPreSignedDownloadURL 生成预授权下载URL的便捷方法
func (r *Client) GetPreSignedDownloadURL(ctx context.Context, bucketName, objectName string, expiry time.Duration) (string, error) {
	return r.GetPreSignedURL(ctx, bucketName, objectName, "GET", expiry)
}

// GetPreSignedUploadURL 生成预授权上传URL的便捷方法
func (r *Client) GetPreSignedUploadURL(ctx context.Context, bucketName, objectName string, expiry time.Duration) (string, error) {
	return r.GetPreSignedURL(ctx, bucketName, objectName, "PUT", expiry)
}

// GetPreSignedURL 生成预授权URL，用于临时访问私有文件
// expiry: URL有效期，建议不超过7天
// method: HTTP方法，如 "GET"、"PUT"、"DELETE" 等
func (r *Client) GetPreSignedURL(ctx context.Context, bucketName, objectName, method string, expiry time.Duration) (string, error) {
	// 验证过期时间
	if expiry <= 0 {
		return "", fmt.Errorf("expiry must be positive")
	}
	if expiry > 7*24*time.Hour {
		return "", fmt.Errorf("expiry cannot exceed 7 days")
	}

	// 根据方法类型生成不同的预授权URL
	var preSignedURL *url.URL
	var err error

	switch strings.ToUpper(method) {
	case "GET":
		// 生成下载用的预授权URL
		preSignedURL, err = r.client.PresignedGetObject(ctx, bucketName, objectName, expiry, nil)
	case "PUT":
		// 生成上传用的预授权URL
		preSignedURL, err = r.client.PresignedPutObject(ctx, bucketName, objectName, expiry)
	case "DELETE":
		// DELETE操作不支持预授权URL
		return "", fmt.Errorf("presigned DELETE not supported, use authenticated DELETE instead")
	default:
		return "", fmt.Errorf("unsupported method: %s. Use GET, PUT, or DELETE", method)
	}

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	//log.Printf("Generated presigned %s URL for %s/%s, expires in %v", method, bucketName, objectName, expiry)

	return preSignedURL.String(), nil
}
