package view

import (
	"reflect"
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/timeutils"
)

// GraphView
// @Description: 案件分析图（人员/公司/合同/产品/账户）
type GraphView struct {
	BaseView `bson:",inline"`
	Name     string `json:"name,omitempty" bson:"name"  gorm:""` // 分析图名称
}

// NewGraphView
// @Description: 案件分析图（人员/公司/合同/产品/账户）
func NewGraphView() *GraphView {
	return &GraphView{}
}

// Equal
// @Description:  对比案件分析图（人员/公司/合同/产品/账户）
func (v *GraphView) Equal(v2 *GraphView) error {
	var msg []string
	if !reflect.DeepEqual(v.CaseId, v2.CaseId) {
		msg = append(msg, "CaseId")
	}
	if !timeutils.Equal(v.CreatedTime, v2.CreatedTime) {
		msg = append(msg, "CreatedTime")
	}
	if !reflect.DeepEqual(v.CreatorId, v2.CreatorId) {
		msg = append(msg, "CreatorId")
	}
	if !reflect.DeepEqual(v.CreatorName, v2.CreatorName) {
		msg = append(msg, "CreatorName")
	}
	if !timeutils.Equal(v.DeletedTime, v2.DeletedTime) {
		msg = append(msg, "DeletedTime")
	}
	if !reflect.DeepEqual(v.DeleterId, v2.DeleterId) {
		msg = append(msg, "DeleterId")
	}
	if !reflect.DeepEqual(v.DeleterName, v2.DeleterName) {
		msg = append(msg, "DeleterName")
	}
	if !reflect.DeepEqual(v.Id, v2.Id) {
		msg = append(msg, "Identity")
	}
	if !reflect.DeepEqual(v.IsDeleted, v2.IsDeleted) {
		msg = append(msg, "IsDeleted")
	}
	if !reflect.DeepEqual(v.Name, v2.Name) {
		msg = append(msg, "AName")
	}
	if !reflect.DeepEqual(v.Remarks, v2.Remarks) {
		msg = append(msg, "Remarks")
	}
	if !reflect.DeepEqual(v.TenantId, v2.TenantId) {
		msg = append(msg, "TenantId")
	}
	if !timeutils.Equal(v.UpdatedTime, v2.UpdatedTime) {
		msg = append(msg, "UpdatedTime")
	}
	if !reflect.DeepEqual(v.UpdaterId, v2.UpdaterId) {
		msg = append(msg, "UpdaterId")
	}
	if !reflect.DeepEqual(v.UpdaterName, v2.UpdaterName) {
		msg = append(msg, "UpdaterName")
	}
	var err error
	if len(msg) > 0 {
		err = errors.New(strings.Join(msg, ","))
	}
	return err
}
