package rustfs

import (
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
)

// UploadFile 上传文件到指定桶
// 如果 contentType 为空，将根据文件扩展名自动检测
func (r *Client) UploadFile(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) error {
	// 如果没有指定 contentType，根据文件扩展名自动检测
	if contentType == "" {
		contentType = GetContentTypeByExtension(objectName)
	}

	// 上传文件
	info, err := r.client.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	log.Printf("Successfully uploaded file: %s, size: %d bytes, etag: %s, content-type: %s", objectName, info.Size, info.ETag, contentType)
	return nil
}

// DeleteFile 删除指定桶中的文件
func (r *Client) DeleteFile(ctx context.Context, bucketName, objectName string) error {
	err := r.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	log.Printf("Successfully deleted file: %s from bucket: %s", objectName, bucketName)
	return nil
}

// ListFiles 列出桶中的文件
func (r *Client) ListFiles(ctx context.Context, bucketName, prefix string) ([]string, error) {
	var files []string

	// 列出对象
	objectCh := r.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("failed to list objects: %v", object.Err)
		}
		files = append(files, object.Key)
	}

	return files, nil
}

// GetFileURL 获取文件的访问URL（用于公共读取的文件）
func (r *Client) GetFileURL(bucketName, objectName string) (string, error) {
	// 对于公共读取的桶，可以直接构造URL
	endpoint := r.client.EndpointURL().String()
	if strings.HasSuffix(endpoint, "/") {
		endpoint = strings.TrimSuffix(endpoint, "/")
	}
	return fmt.Sprintf("%s/%s/%s", endpoint, bucketName, objectName), nil
}

// 常见文件类型映射
var contentTypeMap = map[string]string{
	// 图片类型
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".bmp":  "image/bmp",
	".webp": "image/webp",
	".svg":  "image/svg+xml",
	".ico":  "image/x-icon",
	".tiff": "image/tiff",
	".tif":  "image/tiff",

	// 文档类型
	".pdf":  "application/pdf",
	".doc":  "application/msword",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".xls":  "application/vnd.ms-excel",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".ppt":  "application/vnd.ms-powerpoint",
	".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	".txt":  "text/plain",
	".rtf":  "application/rtf",
	".odt":  "application/vnd.oasis.opendocument.text",
	".ods":  "application/vnd.oasis.opendocument.spreadsheet",
	".odp":  "application/vnd.oasis.opendocument.presentation",

	// 音频类型
	".mp3":  "audio/mpeg",
	".wav":  "audio/wav",
	".flac": "audio/flac",
	".aac":  "audio/aac",
	".ogg":  "audio/ogg",
	".wma":  "audio/x-ms-wma",
	".m4a":  "audio/mp4",

	// 视频类型
	".mp4":  "video/mp4",
	".avi":  "video/x-msvideo",
	".mov":  "video/quicktime",
	".wmv":  "video/x-ms-wmv",
	".flv":  "video/x-flv",
	".webm": "video/webm",
	".mkv":  "video/x-matroska",
	".3gp":  "video/3gpp",
	".m4v":  "video/x-m4v",

	// 压缩文件
	".zip": "application/zip",
	".rar": "application/vnd.rar",
	".7z":  "application/x-7z-compressed",
	".tar": "application/x-tar",
	".gz":  "application/gzip",
	".bz2": "application/x-bzip2",
	".xz":  "application/x-xz",

	// 代码文件
	".html": "text/html",
	".htm":  "text/html",
	".css":  "text/css",
	".js":   "application/javascript",
	".json": "application/json",
	".xml":  "application/xml",
	".yaml": "application/x-yaml",
	".yml":  "application/x-yaml",
	".go":   "text/x-go",
	".py":   "text/x-python",
	".java": "text/x-java-source",
	".c":    "text/x-c",
	".cpp":  "text/x-c++",
	".h":    "text/x-c",
	".php":  "application/x-httpd-php",
	".rb":   "text/x-ruby",
	".sh":   "application/x-sh",
	".sql":  "application/sql",

	// 其他常见类型
	".bin": "application/octet-stream",
	".exe": "application/octet-stream",
	".dmg": "application/x-apple-diskimage",
	".iso": "application/x-iso9660-image",
	".deb": "application/vnd.debian.binary-package",
	".rpm": "application/x-rpm",
	".apk": "application/vnd.android.package-archive",
	".ipa": "application/octet-stream",
}

// GetContentTypeByExtension 根据文件扩展名获取 Content-Type
func GetContentTypeByExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if contentType, exists := contentTypeMap[ext]; exists {
		return contentType
	}
	// 默认返回二进制流类型
	return "application/octet-stream"
}
