package xoss

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/golang-module/carbon/v2"
	"github.com/segmentio/ksuid"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Option struct {
	file          File
	bucketName    string // bucketName
	partSize      int64  // 分片大小
	checkpointDir string // 断点续传检查点目录
}

type Options func(*Option)

// WithFile WithFile
func WithFile(file File) Options {
	return func(c *Option) {
		c.file = file
	}
}

// WithBucketName bucketName
func WithBucketName(bucketName string) Options {
	return func(c *Option) {
		c.bucketName = bucketName
	}
}

// WithPartSize 设置分片大小
func WithPartSize(partSize int64) Options {
	return func(c *Option) {
		c.partSize = partSize
	}
}

func WithDefaultPartSize() Options {
	return func(c *Option) {
		c.partSize = DefaultPartSize
	}
}

// WithCheckpointDir 设置断点续传检查点目录
func WithCheckpointDir(checkpointDir string) Options {
	return func(c *Option) {
		c.checkpointDir = checkpointDir
	}
}

type File struct {
	Suffix   string // 文件后缀 .eg: png
	FileType int64  // 文件类型 1：文件类型 2：图片类型 3：音频类型 4：视频类型
}

type UploadResult struct {
	URL      string
	FileName string
}

const (
	DefaultPartSize = 10 * 1024 * 1024 // 默认分片大小10MB
)

const (
	FileTypeNil   = iota + 1 // 空类型
	FileTypeImage            // 图片类型
	FileTypeAudio            // 音频类型
	FileTypeVideo            // 视频类型
)

// Upload 从文件系统上传文件
func (a *Aliyun) Upload(filePath string, opts ...Options) (*UploadResult, error) {
	if filePath == "" {
		return nil, fmt.Errorf("缺少文件路径")
	}

	opt := a.uploadOption(opts...)

	fileOpen, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer fileOpen.Close()

	return a.uploadToOss(opt.bucketName, filePath, fileOpen, opt)
}

// UploadBinaryData 直接上传二进制数据
func (a *Aliyun) UploadBinaryData(data []byte, opts ...Options) (*UploadResult, error) {

	opt := a.uploadOption(opts...)
	return a.uploadToOss(opt.bucketName, "", bytes.NewReader(data), opt)
}

// uploadToOss 辅助函数，用于处理实际的上传过程
func (a *Aliyun) uploadToOss(bucketName string, filePath string, reader io.Reader, opt Option) (*UploadResult, error) {
	client, err := oss.New(a.conf.Endpoint, a.conf.AccessKey, a.conf.AccessSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create OSS client: %v", err)
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket: %v", err)
	}

	key := a.getFileName(filePath, opt.file)
	key = strings.ReplaceAll(key, "\\", "/")

	if opt.partSize > 0 && filePath != "" {
		// 分片上传
		err = a.multipartUpload(bucket, key, filePath, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to multipart upload data: %v", err)
		}
	} else {
		// 普通上传
		err = bucket.PutObject(key, reader)
		if err != nil {
			return nil, fmt.Errorf("failed to upload data: %v", err)
		}
	}

	url := fmt.Sprintf("https://%s.%s/%s", bucketName, a.conf.Endpoint, key)
	return &UploadResult{
		URL:      url,
		FileName: filepath.Base(key),
	}, nil
}

// multipartUpload 处理多部分上传
func (a *Aliyun) multipartUpload(bucket *oss.Bucket, key string, filePath string, opt Option) error {
	partSize := opt.partSize

	// 使用 UploadFile 实现断点续传和分片上传
	return bucket.UploadFile(key, filePath, partSize, oss.Routines(3), oss.CheckpointDir(true, opt.checkpointDir))
}

// getFileName 获取文件名
func (a *Aliyun) getFileName(filePath string, file File) string {

	if filePath != "" && file.Suffix == "" {
		file.Suffix = path.Ext(filePath)
	}

	// 确保 file.Suffix 以点开头
	if file.Suffix != "" && !strings.HasPrefix(file.Suffix, ".") {
		file.Suffix = "." + file.Suffix
	}

	if file.Suffix == "" {
		file.Suffix = ".png"
	}

	fileName := ksuid.New().String()

	dir := a.uploadFilePathPrefix(file.FileType, a.appName)
	date := carbon.Now().Format("Ymd")

	dir = filepath.Join(dir, date, fileName+file.Suffix)
	return dir
}

// uploadOption 应用选项并返回最终的Option结构
func (a *Aliyun) uploadOption(opts ...Options) Option {
	var o Option

	for _, opt := range opts {
		opt(&o)
	}
	if o.bucketName == "" {
		o.bucketName = a.conf.Bucket
	}
	if o.file.FileType == 0 {
		o.file.FileType = FileTypeImage
	}
	if o.file.Suffix == "" {
		switch o.file.FileType {
		case FileTypeImage:
			o.file.Suffix = "png"
			break
		case FileTypeAudio:
			o.file.Suffix = "mp3"
			break
		case FileTypeVideo:
			o.file.Suffix = "mp4"
			break
		default:
			o.file.Suffix = ""
			break
		}
	}
	return o
}

func (a *Aliyun) uploadFilePathPrefix(fileType int64, appName string) string {
	var filePath string
	switch fileType {
	case FileTypeImage:
		filePath = "images"
		break
	case FileTypeAudio:
		filePath = "audio"
		break
	case FileTypeVideo:
		filePath = "videos"
		break
	default:
		filePath = "files"
		break
	}
	return filepath.Join(appName, filePath)
}
