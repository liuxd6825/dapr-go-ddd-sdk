package restapp

import (
	"bufio"
	"bytes"
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/readexcel"
	"github.com/minio/minio-go/v7"
	"os"
	"testing"
)

func TestInitMinio(t *testing.T) {
	configs := map[string]*MinioConfig{}
	configs["default"] = &MinioConfig{
		Endpoint:        "122.143.11.104:9900",
		AccessKey:       "admin",
		SecretAccessKey: "admin9000",
		UseSSL:          false,
	}
	if err := InitMinio(configs); err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	client := GetMinioClient()

	bucketName := "case-document"
	if ok, err := client.BucketExists(ctx, bucketName); err != nil {
		t.Error(err)
		return
	} else if !ok {
		if err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			t.Error(err)
			return
		}
	}
	objectName := "record.xlsx"
	// 检查对象是否存在
	_, err := client.StatObject(context.Background(), bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			if err := uploadFile(ctx, client, bucketName, "../readexcel/test_files/record.xlsx", objectName); err != nil {
				t.Error(err)
				return
			}
		} else {
			t.Error(err)
		}
	}
	if butter, err := readObject(ctx, client, bucketName, objectName); err != nil {
		t.Error(err)
		return
	} else {
		if data, err := readexcel.ReadBytesToMap(butter.Bytes(), "", 100); err != nil {
			t.Error(err)
			return
		} else {
			t.Log("data.items.length:", len(data.Items))
		}
	}
}

func readObject(ctx context.Context, client *minio.Client, bucketName, objectName string) (*bytes.Buffer, error) {
	// 获取对象信息
	object, err := client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()
	// 读取文件流
	buf := make([]byte, 1024*1024)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		n, err := object.Read(buf)
		if n > 0 {
			buffer.Write(buf[:n])
		}
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
		}

	}
	return buffer, nil
}

// 下载文件
func downloadFile(ctx context.Context, minioClient *minio.Client, bucketName, localFilePath, objectName string) (err error) {
	// 创建一个新的可写流
	file, err := os.Create(localFilePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
	}(file)
	// 执行下载操作
	err = minioClient.FGetObject(ctx, bucketName, objectName, localFilePath, minio.GetObjectOptions{})
	return err
}

func uploadFile(ctx context.Context, minioClient *minio.Client, bucketName, localFilePath, objectName string) (err error) {
	// 打开本地文件
	file, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
	}(file)
	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	// 创建一个新的可读流
	reader := bufio.NewReader(file)
	// 设置上传选项
	opts := minio.PutObjectOptions{
		ContentType: "xlsx",
	}
	// 执行上传操作
	_, err = minioClient.PutObject(ctx, bucketName, objectName, reader, fileInfo.Size(), opts)
	return err
}
