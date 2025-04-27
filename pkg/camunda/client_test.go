package camunda

import (
	camunda_client_go "github.com/citilinkru/camunda-client-go/v3"
	"testing"
)

var client *camunda_client_go.Client
var tenantId string = "lxd"

func init() {
	cfg := &NewConfig{
		Url:         NewDaprUrl("localhost", 48081, "camunda-service"),
		ApiUser:     "demo",
		ApiPassword: "demo",
	}
	cfg = &NewConfig{
		Url:         "http://localhost:48080/engine-rest",
		ApiUser:     "demo",
		ApiPassword: "demo",
	}
	client = NewClient(cfg)
}

func TestClient_StartInstance(t *testing.T) {
	depList, err := client.Deployment.GetList(nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("depList", len(depList))

	procDef := camunda_client_go.QueryProcessDefinitionBy{Key: String("Process_0hy60rq"), TenantId: String(tenantId)}
	reqStartInst := camunda_client_go.ReqStartInstance{BusinessKey: String("001")}
	variables := &map[string]camunda_client_go.Variable{
		"userId":   {Value: "001", Type: "string"},
		"userName": {Value: "张三", Type: "string"},
	}
	reqStartInst.Variables = variables
	startInstResp, err := client.ProcessDefinition.StartInstance(procDef, reqStartInst)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(startInstResp)
}

func TestClient_TaskList(t *testing.T) {
	list, err := client.UserTask.GetList(&camunda_client_go.UserTaskGetListQuery{
		TenantIdIn: []string{tenantId},
		Assignee:   "001",
	})
	variables := map[string]camunda_client_go.Variable{
		"userId":   {Value: "002", Type: "string"},
		"userName": {Value: "李四", Type: "string"},
	}
	for _, task := range list {
		task.Complete(camunda_client_go.QueryUserTaskComplete{
			Variables: variables,
		})
	}
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(list)
}
