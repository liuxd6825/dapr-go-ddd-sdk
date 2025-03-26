package schema

import "encoding/json"

type DBSortType string

const (
	DBSort_None DBSortType = ""
	DBSort_Asc  DBSortType = "asc"
	DBSort_Desc DBSortType = "desc"
)

type DBIndexType string

const (
	DBIndexType_None     DBIndexType = ""         // 无索引
	DBIndexType_Index    DBIndexType = "index"    // 普通索引
	DBIndexType_FullText DBIndexType = "fullText" // 全文索引
)

type DBField struct {
	NotField   bool        `json:"notField"`   // 不是数据库字段
	Name       string      `json:"name"`       // 字段名称
	Size       int64       `json:"size"`       // 字段大小
	PrimaryKey bool        `json:"primaryKey"` // 是主健
	NotNull    bool        `json:"notNull"`    // 不能为空
	Unique     bool        `json:"unique"`     // 唯一
	Creatable  bool        `json:"creatable"`  // 可创建
	Updatable  bool        `json:"updatable"`  // 可更新
	Readable   bool        `json:"readable"`   // 可读取
	SortType   DBSortType  `json:"sort"`       // 排序类型
	IndexType  DBIndexType `json:"indexType"`  // 索引类型
	IndexName  string      `json:"indexName"`  // 索引类型
}

type DBTable struct {
	Name string `json:"name"`
}

func (db *DBTable) init(values map[string]any) error {
	var err error
	for k, v := range values {
		switch k {
		case "name":
			db.Name = v.(string)
		}
	}
	return err
}

func (db *DBField) init(values map[string]any) error {
	var err error
	for key, value := range values {
		switch key {
		case "notField":
			if notField, ok := value.(bool); ok {
				db.NotField = notField
			}
		case "primaryKey":
			if val, ok := value.(bool); ok {
				db.PrimaryKey = val
			}
		case "creatable":
			if val, ok := value.(bool); ok {
				db.Creatable = val
			}
		case "size":
			if val, ok := value.(json.Number); ok {
				db.Size, err = val.Int64()
			}
		case "name":
			if val, ok := value.(string); ok {
				db.Name = val
			}
		case "notNull":
			if val, ok := value.(bool); ok {
				db.NotNull = val
			}
		case "unique":
			if val, ok := value.(bool); ok {
				db.Unique = val
			}
		case "updatable":
			if val, ok := value.(bool); ok {
				db.Updatable = val
			}
		case "readable":
			if val, ok := value.(bool); ok {
				db.Readable = val
			}
		case "sortType":
			if val, ok := value.(string); ok {
				db.SortType = DBSortType(val)
			}
		case "indexName":
			if val, ok := value.(string); ok {
				db.IndexName = val
			}
		case "indexType":
			if val, ok := value.(string); ok {
				db.IndexType = DBIndexType(val)
			}
		}
		if err != nil {
			return err
		}
	}
	return err
}
