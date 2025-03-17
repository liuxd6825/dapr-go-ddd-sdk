package schema

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/os/fs/localfs"
	"io/ioutil"
	"os"
	"testing"
)

func Test_NewSchemasWithJsFile(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(currentDir)
	}
	fs, err := localfs.NewFs(localfs.Config{Path: currentDir})
	if err != nil {
		t.Fatal(err)
		return
	}
	file, err := fs.Open("/testfile/scheam.js")
	if err != nil {
		t.Fatal(err)
		return
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
		return
	}
	library := NewSchemasLibrary()
	err = library.LoadJavaScript(bytes)
	if err != nil {
		t.Fatal(err)
		return
	}
	schemas, err := library.Get("domain")
	t.Log(schemas)
}
