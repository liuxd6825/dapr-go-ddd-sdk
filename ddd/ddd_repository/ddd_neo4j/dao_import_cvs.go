package ddd_neo4j

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/gocsv"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"os"
	"reflect"
	"strings"
	"time"
)

// ImportSaveFileCallback
// fileName 文件名称
// data 保存的数据
// return string 存储文件的URI，可以是file:///或http://
type ImportSaveFileCallback func(ctx context.Context, tenantId, fileName string, data any) (string, ImportSaveCompleteCallback, error)
type ImportSaveCompleteCallback func() error
type importField struct {
	Name string
	Type importFieldType
}
type importFieldType int

const (
	importFieldTypeString importFieldType = iota
	importFieldTypeInt
	importFieldTypeFloat
)

func (d *Dao[T]) ImportNodeCsv(ctx context.Context, cmd ImportCsvCmd, opts ...ddd_repository.Options) (err error) {
	return d.importCsv(ctx, cmd, opts...)
}

func (d *Dao[T]) ImportRelationCsv(ctx context.Context, cmd ImportCsvCmd, opts ...ddd_repository.Options) (err error) {
	return d.importCsv(ctx, cmd, opts...)
}

func (d *Dao[T]) importCsv(ctx context.Context, cmd ImportCsvCmd, opts ...ddd_repository.Options) (err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()

	saveCallback := LocalFileImportSaveFileCallback
	if cmd.SaveFileCallback != nil {
		saveCallback = cmd.SaveFileCallback
	}
	fileUri := ""
	if saveCallback != nil {
		var complete ImportSaveCompleteCallback
		fileUri, complete, err = saveCallback(ctx, cmd.TenantId, cmd.ImportFile, cmd.Data.Data())
		if err != nil {
			return err
		}

		defer func() {
			if complete != nil {
				//三次重试
				for i := 1; i < 4; i++ {
					if err := complete(); err != nil {
						time.Sleep(time.Duration(3*i) * time.Second)
						continue
					}
					return
				}
			}
		}()
	}

	labels := ""
	for _, item := range cmd.Labels {
		labels += ":" + item
	}

	//创建索引
	_ = d.CreateIndex(ctx, "case_"+cmd.CaseId+"_id", "case_"+cmd.CaseId, "id")

	cypher := strings.Builder{}

	switch cmd.ImportType {
	case ImportTypeNode:
		fields, err := getStructFields(cmd.Data.Item(0))
		if err != nil {
			return err
		}
		_, setFields := getCreateUpdate(fields)

		cypher.WriteString(fmt.Sprintf(`
		LOAD CSV WITH HEADERS FROM '%s' AS a
		CALL {
			WITH a
			MERGE (n%s{id:a.id}) ON CREATE SET %v ON MATCH SET %v 
		}  IN TRANSACTIONS OF 100 ROWS;
        `, fileUri, labels, setFields.String(), setFields.String()))
		break

	case ImportTypeRelation:
		relTypes := make(map[string]int) // 存储去重后的关系
		cypher.WriteString(fmt.Sprintf(`
		LOAD CSV WITH HEADERS FROM '%s' AS a
		CALL {
			WITH a
			MATCH (s%s{id:a.startId}),(e%s{id:a.endId}) 
		`, fileUri, labels, labels))

		item := cmd.Data.Item(0)
		dataMap, ok := item.(map[string]any)
		if !ok {
			return errors.ErrorOf("map")
		}

		fields := getMapFields(dataMap)
		_, setFields := getCreateUpdate(fields)

		// 去重后的关系
		for i := 0; i < cmd.Data.Length(); i++ {
			item := cmd.Data.Item(i)
			if relMap, ok := item.(map[string]any); ok {
				relType := relMap["relType"].(string)
				if _, find := relTypes[relType]; !find {
					relTypes[relType] = 0
					// 按关系类型生成创建关系语句
					cypher.WriteString(fmt.Sprintf(`
					FOREACH (_ IN case when a.relType = '%s' then[1] else[] end|
						MERGE (s)-[n:%s{id:a.id}]->(e)
						ON CREATE SET %s
						ON MATCH  SET %s
					)`, relType, relType, setFields.String(), setFields.String()))
				}
			}
		}
		cypher.WriteString(`} IN TRANSACTIONS OF 100 ROWS;`)
		break
	}

	session := d.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		_ = session.Close(ctx)
	}()

	_, err = session.Run(ctx, cypher.String(), nil)
	if err != nil {
		logs.Error(ctx, cmd.TenantId, logs.Fields{"cypher": cypher.String()})
	}
	return err
}

func getStructFields(data any) ([]*importField, error) {
	refFields, err := reflectutils.GetFields(data)
	if err != nil {
		return nil, err
	}
	var list []*importField
	refFields.ForEach(func(index int, name string, sField *reflect.StructField) {
		field := &importField{
			Name: stringutils.FirstLower(name),
		}
		switch sField.Type.Kind() {
		case reflect.Int16:
		case reflect.Int64:
		case reflect.Int8:
		case reflect.Int:
			field.Type = importFieldTypeInt
			break
		case reflect.Float32:
		case reflect.Float64:
			field.Type = importFieldTypeFloat
			break
		default:
			field.Type = importFieldTypeString
		}
		list = append(list, field)
	})
	return list, nil

}

func getMapFields(data map[string]any) []*importField {
	var list []*importField
	for key, val := range data {
		field := &importField{
			Name: stringutils.FirstLower(key),
		}
		if val == nil {
			field.Type = importFieldTypeString
		} else {
			vt := reflect.TypeOf(val)
			switch vt.Kind() {
			case reflect.Int16:
			case reflect.Int64:
			case reflect.Int8:
			case reflect.Int:
				field.Type = importFieldTypeInt
				break
			case reflect.Float32:
			case reflect.Float64:
				field.Type = importFieldTypeFloat
				break
			default:
				field.Type = importFieldTypeString
			}
		}

		list = append(list, field)
	}
	return list

}

func getCreateUpdate(fields []*importField) (*strings.Builder, *strings.Builder) {
	fsLen := len(fields) - 1
	props := &strings.Builder{}
	setFields := &strings.Builder{}
	for i, field := range fields {
		name := stringutils.FirstLower(field.Name)
		props.WriteString(name + ":a." + name)
		switch field.Type {
		case importFieldTypeInt:
			setFields.WriteString(fmt.Sprintf("n.%s=toInteger(a.%s)", name, name))
			break
		case importFieldTypeFloat:
			setFields.WriteString(fmt.Sprintf("n.%s=toFloatOrNull(a.%s)", name, name))
			break
		case importFieldTypeString:
			setFields.WriteString("n." + name + "=a." + name)
			break
		}

		if i < fsLen {
			props.WriteString(",")
			setFields.WriteString(",")
		}
	}
	return props, setFields
}

func LocalFileImportSaveFileCallback(ctx context.Context, tenantId string, fileName string, data any) (string, ImportSaveCompleteCallback, error) {
	neo4jPath, err := restapp.GetConfigAppValue("neo4jPath")
	if err != nil {
		return "", nil, err
	}
	cvsFileName := fmt.Sprintf("%s/import/%s", neo4jPath, fileName)
	var csvFile *os.File
	csvFile, err = os.OpenFile(cvsFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return "", nil, err
	}

	if err = gocsv.MarshalHeadersToFirstLower(data, csvFile); err != nil {
		return "", nil, err
	}
	completeCallback := func() error {
		_ = csvFile.Close()
		return nil
		//return os.Remove(cvsFileName)
	}
	fileUri := fmt.Sprintf("file:///%s", fileName)
	return fileUri, completeCallback, nil
}
