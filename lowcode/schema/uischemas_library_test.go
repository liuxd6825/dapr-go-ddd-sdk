package schema

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/localfs"
	"io/ioutil"
	"os"
	"testing"
)

func Test_LoadUISchemasLibrary(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(currentDir)
	}
	fs, err := localfs.NewFs(&localfs.Config{Id: "test", Path: currentDir})
	if err != nil {
		t.Fatal(err)
		return
	}
	file, err := fs.Open("/testfile/uischema.js")
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
	library := NewUiSchemasLibrary()
	err = library.LoadJavaScript(bytes)
	if err != nil {
		t.Fatal(err)
		return
	}
	schemas, err := library.Get("default")
	t.Log(schemas)
}
