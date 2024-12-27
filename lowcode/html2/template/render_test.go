package template

import (
	"bytes"
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/localfs"
	"os"
	"testing"
)

func Test_Render(t *testing.T) {
	ctx := context.Background()
	w := &bytes.Buffer{}
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
		return
	}
	fs, err := localfs.NewFs(&localfs.Config{Path: path})
	if err != nil {
		t.Fatal(err)
		return
	}
	if err := Render(ctx, w, fs, "/testfile/form.html"); err != nil {
		t.Fatal(err)
	} else {
		t.Log(w.String())
	}
}
