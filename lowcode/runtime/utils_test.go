package runtime

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/localfs"
	"github.com/spf13/afero"
	"io"
	"os"
	"testing"
)

func Test_TransformTSCodeToJS(t *testing.T) {
	tsCode := `
		(function (self: Server) {
			import { createParams as params } from "./human-service.d.ts";
			let data = params.cmd.data;
			data.tenantId = tenantId;
			let err = service.schema.validate(data);
			if (err != null){
				return err;
			}
			service.dao.create(ctx, data);
			return data
		})()
	`

	jsCode, err := TransformTSCodeToJS(tsCode)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(jsCode)
	}
}

func Test_Require(t *testing.T) {
	code := `
	(function () {
		let my = require("./js/math.js");
		let sumCount = my.sum(3, 1);
		return sumCount 
	})()`

	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fs, err := localfs.NewFs(localfs.Config{
		Name: "file",
		Path: exePath,
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	vm := NewRuntime(&Reader{fs: fs})
	val, err := vm.RunString(code)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(val.Export())
}

type Reader struct {
	fs afero.Fs
}

func (r *Reader) ReadFile(filename string, opts ...*fsopts.Options) ([]byte, error) {
	file, err := r.fs.Open(filename)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(file)
	return data, err
}
