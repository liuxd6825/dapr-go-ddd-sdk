package merge_html

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_Merge(t *testing.T) {
	fs := afero.NewOsFs()
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前目录失败:", err)
		return
	}
	err = MergeFile(fs, dir+"/test", ".html", "./test/parts.js")
	assert.NoError(t, err)
}
