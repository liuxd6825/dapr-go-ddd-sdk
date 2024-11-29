package schema_pkg

import (
	"bytes"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
)

type SchemaPkg struct {
	server pkg.Server
}

func New(server pkg.Server) *SchemaPkg {
	return &SchemaPkg{
		server: server,
	}
}

func (s *SchemaPkg) LoadFile(fileUrl string, workPath string) *schema.Schema {
	data := s.server.GetFsm().ReadFile(fileUrl, &fsopts.Options{WorkPath: workPath})
	if len(data) == 0 {
		return nil
	}
	reader := bytes.NewReader(data)
	res, err := schema.NewSchema(reader)
	if err != nil {
		panic(err)
	}
	return res
}
