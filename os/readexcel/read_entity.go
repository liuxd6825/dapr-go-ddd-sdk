package readexcel

import (
	"bytes"
	"context"
)

type Options struct {
	TenantId  string `json:"tenantId,omitempty" bson:"tenant_id" desc:"租户标识"`
	CaseId    string `json:"caseId,omitempty" bson:"case_id"  desc:"案件ID"`
	DocId     string `json:"docId,omitempty" bson:"doc_id"  desc:"文档id"`
	FileId    string `json:"fileId,omitempty" bson:"file_id"  desc:"文件id"`
	BatchId   string `json:"batchId,omitempty" bson:"batch_id"  desc:"批量号"`
	BatchSize int64  `json:"batchSize" bson:"batch_size"  desc:"批量大小"`
}

type Batching struct {
	BatchIndex       int `json:"batchIndex"` // 当前批号
	BatchSize        int `json:"batchSize"`  // 批数据量
	BatchCount       int `json:"batchCount"` // 总批数
	RowTotal         int `json:"rowTotal"`   // 总数据量
	BatchRecordCount int `json:"batchRecordCount"`
}

func ReadFileToEntity[T any](ctx context.Context, fileName string, sheetName string, temp *Template, isView bool,
	newItem func(ctx context.Context, row *DataRow, temp *Template) (T, error),
	batchFunc func(ctx context.Context, list []T, paging Batching) error, opts ...*Options) (*DataTable, error) {
	bs, err := readFile(fileName)
	if err != nil {
		return nil, err
	}
	return ReadByteToEntity(ctx, bytes.NewBuffer(bs), sheetName, temp, isView, newItem, batchFunc, opts...)
}

func ReadByteToEntity[T any](ctx context.Context, buffer *bytes.Buffer, sheetName string, temp *Template, isView bool,
	newItem func(ctx context.Context, row *DataRow, temp *Template) (T, error),
	batchFunc func(ctx context.Context, list []T, paging Batching) error, opts ...*Options) (*DataTable, error) {

	opt := NewOptions(opts...)

	table, err := ReadBytes(ctx, buffer, sheetName, temp)
	if err != nil {
		return nil, err
	}

	list := make([]T, 0)
	rowTotal := int64(len(table.Rows))
	batchCount := rowTotal / opt.BatchSize

	if rowTotal%opt.BatchSize > 0 {
		batchCount++
	}

	paging := Batching{
		BatchIndex:       1,
		BatchSize:        int(opt.BatchSize),
		BatchCount:       int(batchCount),
		RowTotal:         int(rowTotal),
		BatchRecordCount: 0,
	}

	for i, dataRow := range table.Rows {
		item, err := newItem(ctx, dataRow, temp)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
		paging.BatchRecordCount = i + 1

		// 当前数据量是否可以进行批处理
		if int64(len(list)) == opt.BatchSize || int64(i) == rowTotal-1 {
			if err := batchFunc(ctx, list, paging); err != nil {
				return nil, err
			}
			// 是数据预览
			if isView {
				break
			}
			list = make([]T, 0)
			paging.BatchIndex++
		}
	}
	return table, err
}

func ReadMapToEntity[T any](ctx context.Context, mapList []map[string]any, temp *Template,
	newItem func(ctx context.Context, row *DataRow, temp *Template) (T, error),
	opts ...*Options) (*DataTable, []T, error) {

	table, err := ReadByMap(ctx, mapList, temp)
	if err != nil {
		return nil, nil, err
	}

	list := make([]T, 0)

	for _, dataRow := range table.Rows {
		item, err := newItem(ctx, dataRow, temp)
		if err != nil {
			return nil, nil, err
		}
		list = append(list, item)
	}
	return table, list, nil
}

func NewOptions(opts ...*Options) *Options {
	o := &Options{BatchSize: 10000}
	for _, i := range opts {
		if len(i.DocId) != 0 {
			o.DocId = i.DocId
		}
		if len(i.BatchId) != 0 {
			o.BatchId = i.BatchId
		}
		if len(i.TenantId) != 0 {
			o.TenantId = i.TenantId
		}
		if len(i.FileId) != 0 {
			o.FileId = i.FileId
		}
		if i.BatchSize != 0 {
			o.BatchSize = i.BatchSize
		}
	}
	return o
}
