package light_rag

import (
	"context"
	"github.com/spf13/afero"
	"log"
	"testing"
)

func Test_Upload(t *testing.T) {
	ctx := context.Background()
	fs := afero.NewOsFs()

	// 需要读取的文件路径
	filePath := "./test_file/xyj-22-1.txt"

	// 打开文件
	file, err := fs.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	err = client().Document().Upload(ctx, filePath, fs)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_DeleteByFileName(t *testing.T) {
	ctx := context.Background()
	// 需要读取的文件路径
	filePath := "xyj-22-1.txt"

	res, err := client().Document().DeleteByFileName(ctx, filePath)
	if err != nil {
		t.Fatal(err)
		return
	}
	log.Println(res.Message)
}
