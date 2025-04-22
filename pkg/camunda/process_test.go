package camunda

import (
	_ "embed"
	"encoding/json"
	camundaclientgo "github.com/citilinkru/camunda-client-go/v3"
	"log"
	"testing"
)

//go:embed bpmn.xml
var bpmnXml string

func Test_Start(t *testing.T) {
	client := getClient()
	procDef := camundaclientgo.QueryProcessDefinitionBy{
		Key:      getPString("Process_0hy60rq"),
		TenantId: getPString("lxd"),
	}
	req := camundaclientgo.ReqStartInstance{
		Variables:   &map[string]camundaclientgo.Variable{},
		BusinessKey: getPString("busKey"),
	}
	res, err := client.ProcessDefinition.StartInstance(procDef, req)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(res)

}

func Test_ProcessInstanceGetList(t *testing.T) {
	client := getClient()
	pList, err := client.ProcessInstance.GetList(map[string]string{"tenantIdIn": "lxd"})
	if err != nil {
		t.Fatal(err)
		return
	}
	for _, p := range pList {
		log.Println(p)
	}
}

func Test_ProcessInstanceGet(t *testing.T) {
	client := getClient()
	ins, err := client.ProcessInstance.Get("69918aaa-1863-11f0-b3fc-0242ac110002")
	if err != nil {
		t.Fatal(err)
		return
	}
	log.Println(ins)
	qry := camundaclientgo.QueryProcessDefinitionBy{
		Id: getPString(ins.DefinitionId),
	}
	def, err := client.ProcessDefinition.Get(qry)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(def.DeploymentId)
	resp, err := client.ProcessDefinition.GetXML(qry)
	if err != nil {
		t.Fatal(err)
		return
	}
	log.Println(resp.Bpmn20Xml)

}

func Test_GetNextNodes(t *testing.T) {
	bpmn := NewProcessWidthXml(bpmnXml)
	nodes := bpmn.GetNextNodes("Activity_1sc4zxa")
	t.Log("Activity_1sc4zxa", "初审")
	for _, n := range nodes {
		printNode(t, n)
	}

	nodes = bpmn.GetNextNodes("Activity_00h3ft5")
	t.Log("Activity_00h3ft5", "测试经理一")
	for _, n := range nodes {
		printNode(t, n)
	}

	nodes = bpmn.GetNextNodes("Activity_0b1stkn")
	t.Log("Activity_0b1stkn", "并行任务2")
	for _, n := range nodes {
		printNode(t, n)
	}
}

func printNode(t *testing.T, node *Node) {
	data, err := json.Marshal(node)
	if err != nil {
		t.Fatal(err)
	}
	jsonNode := string(data)
	log.Println(jsonNode)
}
