package restapi

import (
	"context"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/events"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/jsonutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/maputils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/stringutils"
)

type CDCRecord struct {
	DB           string             `json:"db"`     // 数据库
	Table        string             `json:"table"`  // 数据表
	Before       map[string]any     `json:"before"` // 之前数据
	After        map[string]any     `json:"after"`  // 之后数据
	OpType       OpType             `json:"opType"` // 操作状态 "r" for read/backfill, "c" for create, "u" for update, "d" for delete
	CdcTimestamp time.Time          `json:"cdcTimestamp"`
	DBSchema     *dbschema.DBSchema `json:"-"`
}

type OpType string

const (
	OpTypeRead   OpType = "r"
	OpTypeCreate OpType = "c"
	OpTypeUpdate OpType = "u"
	OpTypeDelete OpType = "d"
)

func GetCDCParams(ictx iris.Context) (res any, ctx context.Context, err error) {
	request := ictx.Request()
	if request.Method != iris.MethodPost {
		return nil, nil, errors.New("request body is null")
	}

	if request.ContentLength == 0 {
		return nil, nil, errors.New("request body is null")
	}
	cdcRecord, err := NewCDCRecordWithIris(ictx)
	if err != nil {
		return nil, nil, err
	}
	ctx = context.Background()
	return cdcRecord, ctx, err
}

func NewCDCRecordWithIris(ictx iris.Context) (*CDCRecord, error) {
	request := ictx.Request()
	if request.Method != iris.MethodPost {
		return nil, errors.New("request body is null")
	}

	if request.ContentLength == 0 {
		return nil, errors.New("request body is null")
	}
	var cloudEvent events.CloudEvent
	err := ictx.ReadJSON(&cloudEvent)
	if err != nil {
		return nil, err
	}

	dataJson, err := cloudEvent.GetData()
	if err != nil {
		return nil, err
	}
	var cdcRecord CDCRecord
	data := []byte(dataJson)
	if err = jsonutils.Unmarshal(data, &cdcRecord); err != nil {
		return nil, err
	}
	return &cdcRecord, nil
}

// IsMaster 是主数据
func (r *CDCRecord) IsMaster() bool {
	return !strings.Contains(r.Table, "_")
}

// IsRelation 是关系数据
func (r *CDCRecord) IsRelation() bool {
	return strings.Contains(r.Table, "_")
}

// IsRename 是否数据更新
func (r *CDCRecord) IsRename() bool {
	if r.OpType == "u" {
		newName, _ := maputils.GetString(r.After, "name", "")
		oldName, _ := maputils.GetString(r.Before, "name", "")
		if newName != oldName {
			return true
		}
	}
	return false
}

// IsChangedRelType 是否数据更新
func (r *CDCRecord) IsChangedRelType() bool {
	if r.OpType == "u" {
		newType, _ := maputils.GetString(r.After, "relation_type", "")
		oldType, _ := maputils.GetString(r.Before, "relation_type", "")
		if newType != oldType {
			return true
		}
	}
	return false
}

// IsChangedBusFields 是否更新业务字段
func (r *CDCRecord) IsChangedBusFields() bool {
	if r.OpType == "u" {
		for k, _ := range r.After {
			if k != "relation_type" && k != "name" {
				return true
			}
		}
	}
	return false
}

func (r *CDCRecord) AfterMap() map[string]any {
	return r.After
}

func (r *CDCRecord) BeforeMap() map[string]any {
	return r.Before
}

func (r *CDCRecord) AfterName() string {
	val, _ := maputils.GetString(r.After, "name", "")
	return val
}

func (r *CDCRecord) BeforeName() string {
	val, _ := maputils.GetString(r.Before, "name", "")
	return val
}

func (r *CDCRecord) newMap(vals map[string]any) map[string]any {
	mapData := make(map[string]any)
	for key, value := range vals {
		propName := stringutils.FirstLowerCamelString(key)
		mapData[propName] = value
	}
	return mapData
}

func (r *CDCRecord) SetAfterValue(key string, val any) *CDCRecord {
	r.After[key] = val
	return r
}

func (r *CDCRecord) SetBeforeValue(key string, val any) *CDCRecord {
	r.Before[key] = val
	return r
}
