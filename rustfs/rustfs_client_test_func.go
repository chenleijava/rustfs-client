package rustfs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// TestRustFSClient 测试 RustFS 客户端功能
func TestRustFSClient() {

	// 连接配置
	endpoint := "cdn.juchuangjiapin.dpdns.org" // 这里的URL来自于NewRustFSClient的endpoint参数
	// default :rustfsadmin
	accessKey := "Xia0xSEuyD5t8bQYwMP6"
	secretKey := "qrQ0bgB36x81DIK75WCe4fXZTOcNhdFslVkyYG2H"

	// 创建客户端
	client, err := NewRustFSClient(endpoint, accessKey, secretKey)
	if err != nil {
		log.Printf("创建客户端失败: %v", err)
		return
	}

	ctx := context.Background()
	testBucket := "test-rustfs-bucket"
	testFile := "test-file.txt"
	testContent := fmt.Sprintf("这是一个测试文件，用于验证 RustFS 客户端功能。\n测试时间: %s", time.Now().Format("2006-01-02 15:04:05"))

	fmt.Println("开始测试 RustFS 客户端...")

	// 1. 测试创建桶
	fmt.Println("\n1. 测试创建桶...")
	err = client.CreateBucket(ctx, testBucket, "us-east-1")
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			fmt.Printf("✓ 桶 %s 已存在\n", testBucket)
		} else {
			fmt.Printf("✗ 创建桶失败: %v\n", err)
			return
		}
	} else {
		fmt.Printf("✓ 成功创建桶: %s\n", testBucket)
	}

	// 2. 测试设置桶权限
	fmt.Println("\n2. 测试设置桶权限...")
	err = client.SetBucketPolicy(ctx, testBucket, BucketPolicyPublicRead)
	if err != nil {
		fmt.Printf("✗ 设置桶权限失败: %v\n", err)
	} else {
		fmt.Printf("✓ 成功设置桶为公共读取权限\n")
	}

	// 3. 测试上传文件
	fmt.Println("\n3. 测试上传文件...")
	reader := strings.NewReader(testContent)
	err = client.UploadFile(ctx, testBucket, testFile, reader, int64(len(testContent)), "")
	if err != nil {
		fmt.Printf("✗ 上传文件失败: %v\n", err)
	} else {
		fmt.Printf("✓ 成功上传文件: %s\n", testFile)
	}

	// 4. 测试列出文件
	fmt.Println("\n4. 测试列出文件...")
	files, err := client.ListFiles(ctx, testBucket, "")
	if err != nil {
		fmt.Printf("✗ 列出文件失败: %v\n", err)
	} else {
		fmt.Printf("✓ 桶中的文件列表:\n")
		for i, file := range files {
			fmt.Printf("  %d. %s\n", i+1, file)
		}
	}

	// 5. 测试获取文件URL
	fmt.Println("\n5. 测试获取文件URL...")
	fileURL, err := client.GetFileURL(testBucket, testFile)
	if err != nil {
		fmt.Printf("✗ 获取文件URL失败: %v\n", err)
	} else {
		fmt.Printf("✓ 文件访问URL: %s\n", fileURL)
	}

	// 6. 测试检查桶是否存在
	fmt.Println("\n6. 测试检查桶是否存在...")
	exists, err := client.BucketExists(ctx, testBucket)
	if err != nil {
		fmt.Printf("✗ 检查桶是否存在失败: %v\n", err)
	} else {
		fmt.Printf("✓ 桶 %s 存在状态: %t\n", testBucket, exists)
	}

	// 7. 测试自动检测Content-Type功能
	fmt.Println("\n7. 测试自动检测Content-Type功能...")
	testFiles := map[string]string{
		"test.jpg":     "image/jpeg",
		"test.png":     "image/png",
		"test.pdf":     "application/pdf",
		"test.mp4":     "video/mp4",
		"test.json":    "application/json",
		"test.unknown": "application/octet-stream",
	}

	for filename, expectedType := range testFiles {
		actualType := GetContentTypeByExtension(filename)
		if actualType == expectedType {
			fmt.Printf("✓ %s -> %s\n", filename, actualType)
		} else {
			fmt.Printf("✗ %s -> 期望: %s, 实际: %s\n", filename, expectedType, actualType)
		}
	}

	// 8. 测试预授权URL生成
	fmt.Println("\n8. 测试预授权URL生成...")

	// 测试生成下载预授权URL
	downloadURL, err := client.GetPreSignedDownloadURL(ctx, testBucket, testFile, 1*time.Hour)
	if err != nil {
		fmt.Printf("✗ 生成下载预授权URL失败: %v\n", err)
	} else {
		fmt.Printf("✓ 下载预授权URL (1小时有效): %s\n", downloadURL)
	}

	// 测试生成上传预授权URL
	uploadURL, err := client.GetPreSignedUploadURL(ctx, testBucket, "upload-test.txt", 30*time.Minute)
	if err != nil {
		fmt.Printf("✗ 生成上传预授权URL失败: %v\n", err)
	} else {
		fmt.Printf("✓ 上传预授权URL (30分钟有效): %s\n", uploadURL)
	}

	// 测试不支持的方法
	_, err = client.GetPreSignedURL(ctx, testBucket, testFile, "DELETE", 1*time.Hour)
	if err != nil {
		fmt.Printf("✓ DELETE方法正确返回错误: %v\n", err)
	} else {
		fmt.Printf("✗ DELETE方法应该返回错误\n")
	}

	// 测试无效的过期时间
	_, err = client.GetPreSignedURL(ctx, testBucket, testFile, "GET", 8*24*time.Hour) // 超过7天
	if err != nil {
		fmt.Printf("✓ 过期时间验证正确: %v\n", err)
	} else {
		fmt.Printf("✗ 过期时间验证失败\n")
	}

	// 9. 测试上传文件时自动检测Content-Type
	fmt.Println("\n9. 测试上传文件时自动检测Content-Type...")
	testImageContent := "fake image content"
	testImageReader := strings.NewReader(testImageContent)
	err = client.UploadFile(ctx, testBucket, "auto-detect.jpg", testImageReader, int64(len(testImageContent)), "") // 空的contentType
	if err != nil {
		fmt.Printf("✗ 上传文件失败: %v\n", err)
	} else {
		fmt.Printf("✓ 成功上传文件并自动检测Content-Type\n")
	}

	// 10. 测试删除文件
	// 删除自动检测的测试文件
	err = client.DeleteFile(ctx, testBucket, "auto-detect.jpg")
	if err != nil {
		fmt.Printf("✗ 删除测试文件失败: %v\n", err)
	} else {
		fmt.Printf("✓ 成功删除测试文件: auto-detect.jpg\n")
	}

	// 11. 验证文件已删除
	fmt.Println("\n11. 验证文件已删除...")
	files, err = client.ListFiles(ctx, testBucket, "")
	if err != nil {
		fmt.Printf("✗ 列出文件失败: %v\n", err)
	} else {
		if len(files) == 0 {
			fmt.Printf("✓ 桶中已无文件，删除成功\n")
		} else {
			fmt.Printf("! 桶中仍有 %d 个文件\n", len(files))
		}
	}

	fmt.Println("\n=== 测试完成 ===")
	fmt.Println("\n注意: 如果要删除测试桶，请确保桶为空后运行:")
	fmt.Printf("client.DeleteBucket(ctx, \"%s\")\n", testBucket)
}
