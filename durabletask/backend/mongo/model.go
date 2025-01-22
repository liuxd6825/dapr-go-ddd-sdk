package mongo

import "time"

type Instances struct {
	InstanceID       string     `bson:"instance_id"`
	ExecutionID      string     `bson:"execution_id"`
	Name             string     `bson:"name"`
	Version          string     `bson:"version"`
	RuntimeStatus    string     `bson:"runtime_status"`
	CreatedTime      time.Time  `bson:"created_time"`
	LastUpdatedTime  *time.Time `bson:"lastUpdated_time"`
	CompletedTime    *time.Time `bson:"completed_time"`
	LockedBy         string     `bson:"locked_by"`
	LockExpiration   *time.Time `bson:"lock_expiration"`
	Input            string     `bson:"input"`
	Output           string     `bson:"output"`
	CustomStatus     string     `bson:"custom_status"`
	FailureDetails   []byte     `bson:"failure_details"`
	ParentInstanceID string     `bson:"parent_instance_id"`
	TenantID         string     `bson:"tenant_id"`
}

func (i *Instances) GetTenantId() string {
	return i.TenantID
}

func (i *Instances) SetTenantId(v string) {
	i.TenantID = v
}

func (i *Instances) GetId() string {
	return i.InstanceID
}

func (i *Instances) SetId(v string) {
	i.InstanceID = v
}

type History struct {
	InstanceID     string `bson:"instance_id"`
	SequenceNumber int64  `bson:"sequence_number"`
	EventPayload   []byte `bson:"event_payload"`
	TenantID       string `bson:"tenantID"`
}

func (i *History) GetTenantId() string {
	return i.TenantID
}

func (i *History) SetTenantId(v string) {
	i.TenantID = v
}

func (i *History) GetId() string {
	return i.InstanceID
}

func (i *History) SetId(v string) {
	i.InstanceID = v
}

type NewEvents struct {
	SequenceNumber int64      `bson:"seq"`
	InstanceID     string     `bson:"instance_id"`
	ExecutionID    string     `bson:"execution_id"`
	Timestamp      *time.Time `bson:"timestamp"`
	VisibleTime    *time.Time `bson:"visible_time"`
	DequeueCount   int64      `bson:"dequeue_count"`
	LockedBy       string     `bson:"locked_by"`
	EventPayload   []byte     `bson:"event_payload"`
	TenantID       string     `bson:"tenant_id"`
}

func (i *NewEvents) GetTenantId() string {
	return i.TenantID
}

func (i *NewEvents) SetTenantId(v string) {
	i.TenantID = v
}

func (i *NewEvents) GetId() string {
	return i.InstanceID
}

func (i *NewEvents) SetId(v string) {
	i.InstanceID = v
}

type NewTasks struct {
	SequenceNumber int64      `bson:"seq"`
	InstanceID     string     `bson:"instance_id"`
	ExecutionID    string     `bson:"execution_id"`
	Timestamp      *time.Time `bson:"timestamp"`
	DequeueCount   *time.Time `bson:"dequeue_count"`
	LockedBy       string     `bson:"locked_by"`
	LockExpiration string     `bson:"lock_expiration"`
	EventPayload   []byte     `bson:"event_payload"`
	TenantID       string     `bson:"tenant_id"`
}

func (i *NewTasks) GetTenantId() string {
	return i.TenantID
}

func (i *NewTasks) SetTenantId(v string) {
	i.TenantID = v
}

func (i *NewTasks) GetId() string {
	return i.InstanceID
}

func (i *NewTasks) SetId(v string) {
	i.InstanceID = v
}
