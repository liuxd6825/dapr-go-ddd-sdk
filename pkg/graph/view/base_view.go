package view

import (
	"time"
)

// BaseView
// @Description: 视图基类
type BaseView struct {
	Id       string `json:"id,omitempty"  bson:"_id" gorm:"primaryKey" `                                 // 主键
	CaseId   string `json:"caseId,omitempty"  bson:"case_id"  index:""  gorm:"index:idx_case_id" `       // 案件ID
	TenantId string `json:"tenantId,omitempty"  bson:"tenant_id"  index:""  gorm:"index:idx_tenant_id" ` // 租户ID

	CreatedTime *time.Time `json:"createdTime,omitempty"  index:""  bson:"created_time" gorm:"" `                          // 创建时间
	CreatorId   string     `json:"creatorId,omitempty"  index:""   bson:"creator_id" gorm:"index:idx_creator_id" `         // 创建人ID
	CreatorName string     `json:"creatorName,omitempty"  index:""   bson:"creator_name"  gorm:"index:idx_creator_name"  ` // 创建人名称

	DeletedTime *time.Time `json:"deletedTime,omitempty"  index:""  bson:"deleted_time" gorm:"" `                       // 删除时间
	DeleterId   string     `json:"deleterId,omitempty"  index:""  bson:"deleter_id" gorm:"index:idx_deleter_id" `       // 删除人ID
	DeleterName string     `json:"deleterName,omitempty"  index:""  bson:"deleter_name" gorm:"index:idx_deleter_name" ` // 删除人名称

	UpdatedTime *time.Time `json:"updatedTime,omitempty"  index:""  bson:"updated_time" gorm:"" `                       // 修改时间
	UpdaterId   string     `json:"updaterId,omitempty"  index:""  bson:"updater_id" gorm:"index:idx_updater_id" `       // 修改人ID
	UpdaterName string     `json:"updaterName,omitempty"  index:""  bson:"updater_name" gorm:"index:idx_updater_name" ` // 修改人名称

	IsDeleted bool   `json:"isDeleted,omitempty" index:""   bson:"is_deleted" gorm:""  gorm:"index:idx_is_del"` // 是否删除
	Remarks   string `json:"remarks,omitempty"  bson:"remarks" gorm:"" `                                        // 备注

}

// BaseNodeView
// @Description: 节点视图基类
type BaseNodeView struct {
	BaseView `bson:",inline"`
	Nid      int64    `json:"-" bson:"nid"`
	GraphId  string   `json:"graphId,omitempty"  bson:"graph_id" validate:"gt=0" gorm:"index:idx_graph_id"  desc:"图形Id"` // neo4j主键id
	Labels   []string `json:"labels" bson:"labels" gorm:"-"`                                                             // neo4j标签
	NodeType string   `json:"nodeType" bson:"nodeType" gorm:"index:idx_node_type"`                                       // 节点类型
}

// BaseRelationView
// @Description: 关系视图基类
type BaseRelationView struct {
	BaseView   `bson:",inline"`
	Nid        int64                  `json:"-" bson:"nid" gorm:"-"`                                       // neo4j主键id
	Sid        int64                  `json:"-" bson:"sid" gorm:"-"`                                       // neo4j关系开始id
	Eid        int64                  `json:"-" bson:"eid" gorm:"-"`                                       // neo4j关系结束id
	GraphId    string                 `json:"graphId,omitempty" bson:"graph_id" gorm:"index:idx_graph_id"` // 关系图id
	RelType    string                 `json:"relType,omitempty" bson:"rel_type" gorm:"index:idx_rel_type"` // 关系类型
	StartId    string                 `json:"startId,omitempty" bson:"start_id" gorm:"index:idx_start_id"` // 关系开始id
	EndId      string                 `json:"endId,omitempty" bson:"end_id" gorm:"index:idx_end_id"`       // 关系结束id
	Properties map[string]interface{} `json:"properties,omitempty" bson:"properties" `                     // 扩展属性
}

// GetCaseId
// @Description: 获取 案件ID
func (v *BaseView) GetCaseId() string {
	return v.CaseId
}

// SetCaseId
// @Description: 设置 案件ID
func (v *BaseView) SetCaseId(value string) {
	v.CaseId = value
}

// GetCreatedTime
// @Description: 获取 创建时间
func (v *BaseView) GetCreatedTime() *time.Time {
	return v.CreatedTime
}

// SetCreatedTime
// @Description: 设置 创建时间
func (v *BaseView) SetCreatedTime(value *time.Time) {
	v.CreatedTime = value
}

// GetCreatorId
// @Description: 获取 创建人ID
func (v *BaseView) GetCreatorId() string {
	return v.CreatorId
}

// SetCreatorId
// @Description: 设置 创建人ID
func (v *BaseView) SetCreatorId(value string) {
	v.CreatorId = value
}

// GetCreatorName
// @Description: 获取 创建人名称
func (v *BaseView) GetCreatorName() string {
	return v.CreatorName
}

// SetCreatorName
// @Description: 设置 创建人名称
func (v *BaseView) SetCreatorName(value string) {
	v.CreatorName = value
}

// GetDeletedTime
// @Description: 获取 删除时间
func (v *BaseView) GetDeletedTime() *time.Time {
	return v.DeletedTime
}

// SetDeletedTime
// @Description: 设置 删除时间
func (v *BaseView) SetDeletedTime(value *time.Time) {
	v.DeletedTime = value
}

// GetDeleterId
// @Description: 获取 删除人ID
func (v *BaseView) GetDeleterId() string {
	return v.DeleterId
}

// SetDeleterId
// @Description: 设置 删除人ID
func (v *BaseView) SetDeleterId(value string) {
	v.DeleterId = value
}

// GetDeleterName
// @Description: 获取 删除人名称
func (v *BaseView) GetDeleterName() string {
	return v.DeleterName
}

// SetDeleterName
// @Description: 设置 删除人名称
func (v *BaseView) SetDeleterName(value string) {
	v.DeleterName = value
}

// GetId
// @Description: 获取 主键
func (v *BaseView) GetId() string {
	return v.Id
}

// SetId
// @Description: 设置 主键
func (v *BaseView) SetId(value string) {
	v.Id = value
}

// GetIsDeleted
// @Description: 获取 是否删除
func (v *BaseView) GetIsDeleted() bool {
	return v.IsDeleted
}

// SetIsDeleted
// @Description: 设置 是否删除
func (v *BaseView) SetIsDeleted(value bool) {
	v.IsDeleted = value
}

// GetRemarks
// @Description: 获取 备注
func (v *BaseView) GetRemarks() string {
	return v.Remarks
}

// SetRemarks
// @Description: 设置 备注
func (v *BaseView) SetRemarks(value string) {
	v.Remarks = value
}

// GetTenantId
// @Description: 获取 租户ID
func (v *BaseView) GetTenantId() string {
	return v.TenantId
}

// SetTenantId
// @Description: 设置 租户ID
func (v *BaseView) SetTenantId(value string) {
	v.TenantId = value
}

// GetUpdatedTime
// @Description: 获取 修改时间
func (v *BaseView) GetUpdatedTime() *time.Time {
	return v.UpdatedTime
}

// SetUpdatedTime
// @Description: 设置 修改时间
func (v *BaseView) SetUpdatedTime(value *time.Time) {
	v.UpdatedTime = value
}

// GetUpdaterId
// @Description: 获取 修改人ID
func (v *BaseView) GetUpdaterId() string {
	return v.UpdaterId
}

// SetUpdaterId
// @Description: 设置 修改人ID
func (v *BaseView) SetUpdaterId(value string) {
	v.UpdaterId = value
}

// GetUpdaterName
// @Description: 获取 修改人名称
func (v *BaseView) GetUpdaterName() string {
	return v.UpdaterName
}

// SetUpdaterName
// @Description: 设置 修改人名称
func (v *BaseView) SetUpdaterName(value string) {
	v.UpdaterName = value
}

// BaseNodeView

func (b *BaseNodeView) GetNid() int64 {
	return b.Nid
}

func (b *BaseNodeView) SetNid(int2 int64) {
	b.Nid = int2
}

func (b *BaseNodeView) GetLabels() []string {
	return b.Labels
}

func (b *BaseNodeView) SetLabels(v []string) {
	b.Labels = v
}

func (b *BaseNodeView) GetId() string {
	return b.Id
}

func (b *BaseNodeView) SetId(v string) {
	b.Id = v
}

func (b *BaseNodeView) SetCaseId(v string) {
	b.GraphId = v
}

func (b *BaseNodeView) GetGraphId() string {
	return b.GraphId
}

// BaseRelationView

func (b *BaseRelationView) GetId() string {
	return b.Id
}

func (b *BaseRelationView) SetId(v string) {
	b.Id = v
}

func (b *BaseRelationView) SetNid(s int64) {
	b.Nid = s
}

func (b *BaseRelationView) GetNid() int64 {
	return b.Nid
}

func (b *BaseRelationView) SetRelType(s string) {
	b.RelType = s
}

func (b *BaseRelationView) GetRelType() string {
	return b.RelType
}

func (b *BaseRelationView) SetSid(s int64) {
	b.Sid = s
}

func (b *BaseRelationView) GetSid() int64 {
	return b.Sid
}

func (b *BaseRelationView) SetEid(s int64) {
	b.Eid = s
}

func (b *BaseRelationView) GetEid() int64 {
	return b.Eid
}

func (b *BaseRelationView) SetStartId(s string) {
	b.StartId = s
}

func (b *BaseRelationView) GetStartId() string {
	return b.StartId
}

func (b *BaseRelationView) SetEndId(v string) {
	b.EndId = v
}

func (b *BaseRelationView) GetEndId() string {
	return b.EndId
}

func (b *BaseRelationView) SetGraphId(v string) {
	b.GraphId = v
}

func (b *BaseRelationView) GetGraphId() string {
	return b.GraphId
}

func (b *BaseRelationView) SetProperties(v map[string]interface{}) {
	b.Properties = v
}

func (b *BaseRelationView) GetProperties() map[string]interface{} {
	return b.Properties
}
