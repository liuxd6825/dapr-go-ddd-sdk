package server

import (
	"bytes"
	"context"

	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/jsonschema/v6"
)

type Proxy struct {
	server *Server
	vm     *goja.Runtime
	funcs  map[string]goja.Value
}

func NewProxy(server *Server, vm *goja.Runtime) *Proxy {
	p := &Proxy{
		vm:     vm,
		server: server,
	}
	p.init()
	return p
}

func (s *Proxy) init() {
	s.funcs = map[string]goja.Value{}
	s.funcs["init"] = s.vm.ToValue(s.server.Init)
	s.funcs["loadPkg"] = s.vm.ToValue(s.server.LoadPkg)
	s.funcs["workPath"] = s.vm.ToValue(s.server.FsOpts().WorkPath)
	s.funcs["logs"] = s.vm.ToValue(s.server.Logs)
	s.funcs["app"] = s.vm.ToValue(s.server.App)
	s.funcs["autoMigrateAll"] = s.vm.ToValue(s.AutoMigrateAll)
	s.funcs["autoMigrateTable"] = s.vm.ToValue(s.AutoMigrateTable)
}

// Get 方法：获取键对应的值
func (s *Proxy) Get(name string) goja.Value {
	var res = goja.Undefined()
	if v, ok := s.funcs[name]; ok {
		res = v
	} else if value, exists := s.server.RunValues().Get(name); exists {
		res = s.vm.ToValue(value)
	}
	return res

}

// Set 方法：设置键值
func (s *Proxy) Set(name string, val goja.Value) bool {
	s.server.RunValues().Set(name, val.Export())
	return true
}

// Has 方法：检查键是否存在
func (s *Proxy) Has(name string) bool {
	exists := s.server.RunValues().Has(name)
	return exists
}

// Delete 方法：删除键
func (s *Proxy) Delete(name string) bool {
	s.server.RunValues().Remove(name)
	return true
}

// Keys 方法：获取所有键
func (s *Proxy) Keys() []string {
	return s.server.RunValues().Keys()
}

func (s *Proxy) InitVM(vm *goja.Runtime) error {
	return s.server.InitVM(vm)
}

func (s *Proxy) AutoMigrateAll(path string) error {
	server := s.server
	ctx := context.Background()
	files := server.fsPkg.ReadAllPath(path)
	for _, file := range files {
		if file.IsDir {
			err := s.AutoMigrateAll(file.Path + "/" + file.Name)
			if err != nil {
				return err
			}
		} else {
			s.AutoMigrateTable(ctx, file.Path+"/"+file.Name)
		}
	}
	return nil
}

func (s *Proxy) AutoMigrateTable(ctx context.Context, schFile string) {
	sch := s.LoadSchemaFile(schFile, "")
	aggField, aggType, tableName, _, dbKey := s.getSchemaDBInfos(sch)
	dbSch := dbschema.NewDBSchemaWithJsonSchema(sch)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		AggField:  aggField,
		AggType:   aggType,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	vDao := dao.NewDao[map[string]any](newCfg)
	vDao.Table().AutoMigrate(ctx)
}

func (s *Proxy) LoadSchemaFile(fileUrl string, workPath string) *jsonschema.Schema {
	server := s.server
	data := server.FsPkg().ReadFile(fileUrl, &fsopts.Options{WorkPath: workPath})
	if len(data) == 0 {
		return nil
	}

	reader, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	schemaFile := fileUrl
	compiler := schema.NewCompiler()
	compiler.UseLoader(server.SchemaLoader())

	if err := compiler.AddResource(schemaFile, reader); err != nil {
		panic(err)
	}
	sch, err := compiler.Compile(schemaFile)
	if err != nil {
		panic(err)
	}
	return sch
}

func (s *Proxy) getSchemaDBInfos(sch *jsonschema.Schema) (aggField string, aggType string, tableName string, isPubEvent bool, dbKey string) {
	meta := schema.GetMetaExtension(sch)
	tableName = sch.Name()
	isPubEvent = false
	if meta != nil && meta.DBTable != nil {
		tableName = meta.DBTable.Name
		dbKey = meta.DBTable.DBKey
	}
	if meta != nil && meta.DDD != nil {
		aggField = meta.DDD.AggField
		aggType = meta.DDD.AggType
		isPubEvent = meta.DDD.IsPubEvent
	}
	return
}
