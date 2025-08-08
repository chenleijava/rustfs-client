package rustfs

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"log"
)

// 桶权限类型常量
const (
	// BucketPolicyPublicRead 公共读取权限 - 允许任何人读取桶中的对象，私有写入
	BucketPolicyPublicRead = "public-read"
	// BucketPolicyPrivate 私有权限 - 只有桶所有者可以访问
	BucketPolicyPrivate = "private"
)

// CreateBucket 创建存储桶
func (r *Client) CreateBucket(ctx context.Context, bucketName, region string) error {
	// 检查桶是否已存在
	exists, err := r.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %v", err)
	}

	if exists {
		return fmt.Errorf("bucket %s already exists", bucketName)
	}

	// 创建桶
	err = r.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
		Region: region,
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %v", err)
	}

	log.Printf("Successfully created bucket: %s", bucketName)
	return nil
}

// BucketExists 检查桶是否存在
func (r *Client) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	exists, err := r.client.BucketExists(ctx, bucketName)
	if err != nil {
		return false, fmt.Errorf("failed to check bucket existence: %v", err)
	}
	return exists, nil
}

// DeleteBucket 删除空桶
func (r *Client) DeleteBucket(ctx context.Context, bucketName string) error {
	err := r.client.RemoveBucket(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to delete bucket: %v", err)
	}

	log.Printf("Successfully deleted bucket: %s", bucketName)
	return nil
}

// SetBucketPolicy 设置桶的访问策略
// policyType: 使用 BucketPolicyPublicRead 表示公共读取，私有写入；BucketPolicyPrivate 表示完全私有
func (r *Client) SetBucketPolicy(ctx context.Context, bucketName, policyType string) error {
	var policy string

	switch policyType {
	case BucketPolicyPublicRead:
		// 公共读取，私有写入策略
		policy = fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {
						"AWS": ["*"]
					},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, bucketName)
	case BucketPolicyPrivate:
		// 完全私有策略（默认）
		policy = ""
	default:
		return fmt.Errorf("unsupported policy type: %s. Use '%s' or '%s'", policyType, BucketPolicyPublicRead, BucketPolicyPrivate)
	}

	err := r.client.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		return fmt.Errorf("failed to set bucket policy: %v", err)
	}
	//log.Printf("Successfully set bucket policy for %s to %s", bucketName, policyType)
	return nil

}
