package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chenleijava/rustfs-client/rustfs"
)

func main() {
	// 1. 创建RustFS客户端
	endpoint := "cdn.juchuangjiapin.dpdns.org" // 这里的URL来自于NewRustFSClient的endpoint参数
	accessKey := "Xia0xSEuyD5t8bQYwMP6"
	secretKey := "qrQ0bgB36x81DIK75WCe4fXZTOcNhdFslVkyYG2H"

	client, err := rustfs.NewRustFSClient(endpoint, accessKey, secretKey)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	ctx := context.Background()
	bucketName := "test-bucket"
	fileName := fmt.Sprintf("upload-via-presigned-%s.txt", time.Now().Format("20060102150405"))

	// 2. 确保桶存在
	err = client.CreateBucket(ctx, bucketName, "us-east-1")
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		log.Fatalf("创建桶失败: %v", err)
	}
	_ = client.SetBucketPolicy(ctx, bucketName, rustfs.BucketPolicyPublicRead)

	// 3. 生成预授权上传URL
	fmt.Println("正在生成预授权上传URL...")
	uploadURL, err := client.GetPreSignedUploadURL(ctx, bucketName, fileName, 5*time.Minute)
	if err != nil {
		log.Fatalf("生成预授权上传URL失败: %v", err)
	}

	fmt.Printf("预授权上传URL: %s\n", uploadURL)
	fmt.Printf("URL中的域名部分来自于: endpoint参数 (%s)\n", endpoint)

	// 4. 使用预授权URL上传文件
	fileContent := "这是通过预授权URL上传的文件内容！\n上传时间: " + time.Now().Format("2006-01-02 15:04:05")
	err = uploadFileWithPresignedURL(uploadURL, fileContent)
	if err != nil {
		log.Fatalf("使用预授权URL上传文件失败: %v", err)
	}

	fmt.Println("✓ 文件上传成功！")

	// 5. 验证文件是否上传成功
	files, err := client.ListFiles(ctx, bucketName, "")
	if err != nil {
		log.Printf("列出文件失败: %v", err)
	} else {
		fmt.Printf("桶中的文件: %v\n", files)
	}

	// 6. 生成预授权下载URL来验证内容
	downloadURL, err := client.GetPreSignedDownloadURL(ctx, bucketName, fileName, 5*time.Minute)
	if err != nil {
		log.Printf("生成下载URL失败: %v", err)
	} else {
		fmt.Printf("预授权下载URL: %s\n", downloadURL)

		// 下载并显示内容
		content, err := downloadFileWithPresignedURL(downloadURL)
		if err != nil {
			log.Printf("下载文件失败: %v", err)
		} else {
			fmt.Printf("文件内容:\n%s\n", content)
		}
	}
}

// uploadFileWithPresignedURL 使用预授权URL上传文件
func uploadFileWithPresignedURL(presignedURL, content string) error {
	// 创建HTTP PUT请求
	req, err := http.NewRequest("PUT", presignedURL, strings.NewReader(content))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置Content-Type（可选）
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("上传失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	return nil
}

// downloadFileWithPresignedURL 使用预授权URL下载文件
func downloadFileWithPresignedURL(presignedURL string) (string, error) {
	// 创建HTTP GET请求
	resp, err := http.Get(presignedURL)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	// 读取内容
	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取内容失败: %v", err)
	}

	return buf.String(), nil
}
