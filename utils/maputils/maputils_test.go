package maputils

import (
	"testing"
	"time"
)

type userView struct {
	TenantId    string
	Remarks     string
	CreatedTime time.Time
	UpdatedTime *time.Time
}

func TestMapStructure_Decode(t *testing.T) {
	user := &userView{}
	props := map[string]interface{}{
		"graphId":     "graphId",
		"tenantId":    "tenantId",
		"remarks":     "remarks",
		"createdTime": "2022-08-31 13:47:36",
		"updatedTime": "2022-08-31T13:47:36.255551+08:00",
	}

	if err := Decode(props, user); err != nil {
		t.Error(err)
		return
	}
	t.Logf("tenantId = %v", user.TenantId)
	t.Logf("remarks = %v", user.Remarks)
	t.Logf("createdTime = %v", user.CreatedTime)
	t.Logf("updatedTime = %v", user.UpdatedTime)
}
