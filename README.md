# RustFS Client

基于 MinIO Go SDK 封装的 RustFS 操作客户端，提供简单易用的文件存储操作接口。

## 功能特性

- ✅ 创建存储桶
- ✅ 设置桶权限（公共读取/私有访问）
- ✅ 文件上传（支持自动内容类型检测）
- ✅ 文件删除
- ✅ 文件列表
- ✅ 获取文件访问URL
- ✅ 桶管理（检查存在、删除桶）
- ✅ 预授权URL生成（支持临时上传/下载）
- ✅ 自动内容类型检测（支持60+种文件类型）
- ✅ 前端集成示例（NutUI Uploader）
- ✅ SSL 证书验证选项（支持自签名证书）

## 项目结构

```
rustfs-client/
├── README.md                    # 项目文档
├── example.go                   # 基础使用示例
├── go.mod                       # Go 模块文件
├── go.sum                       # 依赖校验文件
├── examples/                    # 示例代码目录
│   ├── presigned_upload_example.go  # 预授权上传示例
│   └── uploader.js              # 前端集成示例
└── rustfs/                      # 核心包目录
    ├── rustfs_client.go         # 客户端主文件
    ├── bucket.go                # 桶管理功能
    ├── handle_file.go           # 文件操作功能
    ├── pre_signed_url.go        # 预授权URL功能
    └── test_rust_fs.go          # 测试文件
```

## 安装依赖

```bash
go mod tidy
```

## 配置信息

- **RustFS 端口**: 9000
- **访问密钥**: rustfsadmin
- **秘密密钥**: rustfsadmin
- **SSL**: 默认关闭

## 使用示例

### 1. 创建客户端

#### 基础客户端创建

```go
client, err := NewRustFSClient("localhost:9000", "rustfsadmin", "rustfsadmin", false)
if err != nil {
    log.Fatalf("Failed to create RustFS client: %v", err)
}
```

#### 带 SSL 选项的客户端创建

```go
// 使用 HTTPS 且跳过 SSL 证书验证（适用于自签名证书）
client, err := NewRustFSClientWithSSLOptions("localhost:9000", "rustfsadmin", "rustfsadmin", true, true)
if err != nil {
    log.Fatalf("Failed to create RustFS client: %v", err)
}

// 使用 HTTPS 且验证 SSL 证书（适用于有效证书）
client, err := NewRustFSClientWithSSLOptions("localhost:9000", "rustfsadmin", "rustfsadmin", true, false)
if err != nil {
    log.Fatalf("Failed to create RustFS client: %v", err)
}
```

### 2. 创建存储桶

```go
ctx := context.Background()
err = client.CreateBucket(ctx, "my-bucket", "us-east-1")
if err != nil {
    log.Printf("Failed to create bucket: %v", err)
}
```

### 3. 设置桶权限

```go
// 设置为公共读取，私有写入
err = client.SetBucketPolicy(ctx, "my-bucket", rustfs.BucketPolicyPublicRead)

// 设置为完全私有
err = client.SetBucketPolicy(ctx, "my-bucket", rustfs.BucketPolicyPrivate)
```

### 4. 上传文件

```go
// 手动指定内容类型
fileContent := "Hello, RustFS!"
reader := strings.NewReader(fileContent)
err = client.UploadFile(ctx, "my-bucket", "hello.txt", reader, int64(len(fileContent)), "text/plain")

// 自动检测内容类型（推荐）
err = client.UploadFile(ctx, "my-bucket", "image.jpg", reader, int64(len(fileContent)), "")
// 系统会自动检测 .jpg 文件并设置为 "image/jpeg"
```

### 5. 删除文件

```go
err = client.DeleteFile(ctx, "my-bucket", "hello.txt")
```

### 6. 列出文件

```go
files, err := client.ListFiles(ctx, "my-bucket", "")
if err == nil {
    for _, file := range files {
        fmt.Println(file)
    }
}
```

### 7. 获取文件URL

```go
url, err := client.GetFileURL(ctx, "my-bucket", "hello.txt")
if err == nil {
    fmt.Printf("File URL: %s\n", url)
}
```

### 8. 预授权URL操作

```go
// 生成预授权下载URL（5分钟有效期）
downloadURL, err := client.GetPreSignedDownloadURL(ctx, "my-bucket", "hello.txt", 5*time.Minute)
if err == nil {
    fmt.Printf("Download URL: %s\n", downloadURL)
}

// 生成预授权上传URL（1小时有效期）
uploadURL, err := client.GetPreSignedUploadURL(ctx, "my-bucket", "new-file.txt", time.Hour)
if err == nil {
    fmt.Printf("Upload URL: %s\n", uploadURL)
    // 前端可以直接使用此URL进行PUT请求上传文件
}

// 通用预授权URL生成
genericURL, err := client.GetPreSignedURL(ctx, "my-bucket", "file.txt", "GET", 30*time.Minute)
```

## 运行示例

```bash
go run .
```

## API 参考

### RustFSClient 结构体

#### NewRustFSClient(endpoint, accessKey, secretKey string, useSSL bool) (*RustFSClient, error)
创建新的 RustFS 客户端实例（默认不跳过 SSL 证书验证）。

#### NewRustFSClientWithSSLOptions(endpoint, accessKey, secretKey string, useSSL, skipSSLVerify bool) (*RustFSClient, error)
创建带 SSL 选项的 RustFS 客户端实例。
- `useSSL`: 是否使用 HTTPS 连接
- `skipSSLVerify`: 是否跳过 SSL 证书验证（适用于自签名证书场景）

#### CreateBucket(ctx context.Context, bucketName, region string) error
创建新的存储桶。

#### SetBucketPolicy(ctx context.Context, bucketName, policyType string) error
设置桶的访问策略。
- `policyType`: 使用常量 `rustfs.BucketPolicyPublicRead` (公共读取) 或 `rustfs.BucketPolicyPrivate` (私有)

#### UploadFile(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) error
上传文件到指定桶。
- `contentType`: 文件MIME类型，如果为空字符串，系统会根据文件扩展名自动检测

#### DeleteFile(ctx context.Context, bucketName, objectName string) error
删除指定桶中的文件。

#### ListFiles(ctx context.Context, bucketName, prefix string) ([]string, error)
列出桶中的文件。

#### GetFileURL(ctx context.Context, bucketName, objectName string) (string, error)
获取文件的访问URL。

#### BucketExists(ctx context.Context, bucketName string) (bool, error)
检查桶是否存在。

#### DeleteBucket(ctx context.Context, bucketName string) error
删除空桶。

#### GetPreSignedURL(ctx context.Context, bucketName, objectName, method string, expiry time.Duration) (string, error)
生成预授权URL，用于临时访问私有文件。
- `method`: HTTP方法，支持 "GET"（下载）、"PUT"（上传）
- `expiry`: URL有效期，最长7天

#### GetPreSignedDownloadURL(ctx context.Context, bucketName, objectName string, expiry time.Duration) (string, error)
生成预授权下载URL的便捷方法。

#### GetPreSignedUploadURL(ctx context.Context, bucketName, objectName string, expiry time.Duration) (string, error)
生成预授权上传URL的便捷方法。

#### GetContentTypeByExtension(filename string) string
根据文件扩展名自动检测内容类型，支持60+种常见文件格式。

## 注意事项

1. 确保 RustFS 服务正在运行并监听 9000 端口
2. 删除桶之前需要先删除桶中的所有文件
3. 公共读取策略允许任何人访问桶中的文件，请谨慎使用
4. 文件上传时可以留空 Content-Type，系统会自动检测
5. 预授权URL有时效性，最长不超过7天
6. 预授权上传URL仅支持PUT方法，不支持POST方法
7. 预授权DELETE操作不被支持，请使用认证方式删除文件
8. **SSL 证书验证**：
   - 生产环境建议使用有效的 SSL 证书并启用证书验证
   - 开发环境或使用自签名证书时，可以设置 `skipSSLVerify=true` 跳过证书验证
   - 跳过 SSL 证书验证会降低安全性，仅在必要时使用
   - 自签名证书场景下，建议在内网环境中使用

## 前端集成示例

项目提供了与 NutUI Uploader 组件的集成示例，支持使用预授权URL进行文件上传：

```javascript
// 参考 examples/uploader.js
const upload = async (file) => {
  const uploadURL = await getPresignedUploadURL(file.name);
  const response = await fetch(uploadURL, {
    method: 'PUT',
    body: file,
    headers: { 'Content-Type': file.type || 'application/octet-stream' }
  });
  return response.ok ? { url: uploadURL.split('?')[0] } : null;
};
```

## 支持的文件类型

系统支持自动检测60+种文件类型，包括：
- **图片**: jpg, png, gif, bmp, webp, svg, ico, tiff
- **文档**: pdf, doc, docx, xls, xlsx, ppt, pptx, txt, rtf
- **音频**: mp3, wav, flac, aac, ogg, wma, m4a
- **视频**: mp4, avi, mov, wmv, flv, webm, mkv, 3gp
- **压缩**: zip, rar, 7z, tar, gz, bz2, xz
- **代码**: html, css, js, json, xml, go, py, java, c, cpp
- **其他**: bin, exe, dmg, iso, deb, rpm, apk

## 错误处理

所有方法都返回 error，建议在生产环境中进行适当的错误处理和日志记录。