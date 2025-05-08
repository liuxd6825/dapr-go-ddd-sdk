package minio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/localcache"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/minio/minio-go/v7"
	"io"
	"strings"
	"sync"
)

type Client interface {
	ReadObject(ctx context.Context, bucketName, objectName string) (*bytes.Buffer, error)
	ReadFile(ctx context.Context, tenantId, fileId, fileName string) (*bytes.Buffer, error)
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts PutObjectOptions) (info minio.UploadInfo, err error)
	RemoveObject(ctx context.Context, bucketName, objectName string, opts RemoveObjectOptions) error
}

type clientProxy struct {
	client *minio.Client
}

var _clientOnce sync.Once
var _proxy Client

type PutObjectOptions = minio.PutObjectOptions
type UploadInfo = minio.UploadInfo
type RemoveObjectOptions = minio.RemoveObjectOptions

func GetClient() Client {
	_clientOnce.Do(func() {
		_proxy = &clientProxy{
			client: restapp.GetMinioClient(),
		}
	})
	return _proxy
}

func (c *clientProxy) ReadFile(ctx context.Context, tenantId, fileId, fileName string) (*bytes.Buffer, error) {
	var buffer *bytes.Buffer
	objectName, err := GetObjectName(fileId, fileName)
	if err != nil {
		return nil, err
	}
	cacheKey := fmt.Sprintf("field-%v-%v", tenantId, objectName)
	if data, err := localcache.GetCache(cacheKey); err == nil {
		buffer = bytes.NewBuffer(data)
	} else {
		buffer, err = GetClient().ReadObject(ctx, tenantId, objectName)
		if err != nil {
			return nil, err
		}
		if err = localcache.SetCache(cacheKey, buffer.Bytes()); err != nil {
			return nil, err
		}
	}
	return buffer, nil
}

func (c *clientProxy) WriteFile(ctx context.Context, tenantId, fileId, fileName string, data *bytes.Buffer) (*bytes.Buffer, error) {
	var buffer *bytes.Buffer
	objectName, err := GetObjectName(fileId, fileName)
	if err != nil {
		return nil, err
	}
	cacheKey := fmt.Sprintf("field-%v-%v", tenantId, objectName)
	if data, err := localcache.GetCache(cacheKey); err == nil {
		buffer = bytes.NewBuffer(data)
	} else {
		buffer, err = GetClient().ReadObject(ctx, tenantId, objectName)
		if err != nil {
			return nil, err
		}
		if err = localcache.SetCache(cacheKey, buffer.Bytes()); err != nil {
			return nil, err
		}
	}
	return buffer, nil
}

func GetObjectName(fileId, fileName string) (string, error) {
	var extName string
	if i := strings.LastIndex(fileName, "."); i > 0 {
		extName = fileName[i:]
	} else {
		return "", errors.New("文件名【%s】, 没有扩展名.", fileName)
	}
	return fileId + extName, nil
}

func (c *clientProxy) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts PutObjectOptions) (info UploadInfo, err error) {
	return c.client.PutObject(ctx, bucketName, objectName, reader, objectSize, opts)
}

func (c *clientProxy) RemoveObject(ctx context.Context, bucketName, objectName string, opts RemoveObjectOptions) error {
	return c.client.RemoveObject(ctx, bucketName, objectName, opts)
}

func (c *clientProxy) ReadObject(ctx context.Context, bucketName, objectName string) (buffer *bytes.Buffer, err error) {
	_, err = c.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	// 获取对象信息
	object, err := c.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer func(object *minio.Object) {
		err = object.Close()
	}(object)

	// 读取文件流
	buf := make([]byte, 1024*1024)
	buffer = bytes.NewBuffer(make([]byte, 0))
	for {
		n, e := object.Read(buf)
		if n > 0 {
			buffer.Write(buf[:n])
		}
		if e != nil {
			if e.Error() == "EOF" {
				e = nil
			}
			break
		}
	}
	return buffer, err
}
