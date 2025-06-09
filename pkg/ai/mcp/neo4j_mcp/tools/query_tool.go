package tools

import (
	"context"
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/mcp/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/mcp/neo4j_mcp/neo4jdb"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/sirupsen/logrus"
)

type QueryTool struct {
	tool   mcp.Tool
	driver neo4j.Driver
}

type QueryInParameter struct {
	TenantId string `json:"tenantId"`
	CaseId   string `json:"caseId"`
	UserId   string `json:"userId"`
	Cypher   string `json:"cypher"`
}

func NewNeo4jQueryTool(driver neo4j.Driver) *QueryTool {
	tool := &QueryTool{driver: driver, tool: mcp.NewTool("neo4j-query",
		mcp.WithDescription("Neo4j Data Relationship Query"),
		mcp.WithString("tenantId",
			mcp.Required(),
			mcp.Description("Tenant ID in the system"),
		),
		mcp.WithString("caseId",
			mcp.Required(),
			mcp.Description("ID of the project case"),
		),
		mcp.WithString("userId",
			mcp.Required(),
			mcp.Description("User ID in the system"),
		),
		mcp.WithString("cypher",
			mcp.Required(),
			mcp.Description("cypher statements in the neo4j database"),
		),
	)}
	return tool
}

func (q *QueryTool) handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	inParam, err := q.getRequest(request)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	logrus.Infof("tenantId:%s; caseId:%s; userId:%s; cypher:%s", inParam.TenantId, inParam.CaseId, inParam.UserId, inParam.Cypher)
	// 在此解析请求内容并执行 Neo4j 查询
	session := q.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	result, err := session.Run(inParam.Cypher, nil)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	data := neo4jdb.GetData(ctx, result)
	dataJson, err := json.Marshal(data)
	if err != nil {
		logrus.Error(err.Error())
		return mcp.NewToolResultError(err.Error()), nil
	}
	resp := string(dataJson)
	logrus.Infof("resp:%s", resp)
	return mcp.NewToolResultText(resp), nil
}

func (q *QueryTool) getRequest(request mcp.CallToolRequest) (*QueryInParameter, error) {
	tenantId, err := request.RequireString("tenantId")
	if err != nil {
		return nil, err
	}

	caseId, err := request.RequireString("caseId")
	if err != nil {
		return nil, err
	}

	userId, err := request.RequireString("userId")
	if err != nil {
		return nil, err
	}

	cypher, err := request.RequireString("cypher")
	if err != nil {
		return nil, err
	}

	return &QueryInParameter{
		TenantId: tenantId,
		CaseId:   caseId,
		UserId:   userId,
		Cypher:   cypher,
	}, nil

}
func (q *QueryTool) Register(s common.IMCPServer) {
	s.AddTool(q.tool, q.handler)
}
