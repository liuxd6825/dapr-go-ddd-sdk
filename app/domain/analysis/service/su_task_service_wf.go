package service

import (
	"context"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
	"go.temporal.io/sdk/workflow"
)

// Analysis_SuTaskAnalysisWorkflow
// @Description: 可疑分析工作流
// @param wctx
// @param taskId
// @param workflowId
// @param ctxMap
// @return error
func Analysis_SuTaskAnalysisWorkflow(wctx workflow.Context, taskId string, workflowId string, ctxMap map[string]any) error {
	logs.Infofmt(nil, "SuTaskAnalysisWorkflow [%s]", taskId)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 24 * time.Hour,
	}
	actx := workflow.WithActivityOptions(wctx, ao)
	err := workflow.ExecuteActivity(actx, Analysis_SuTaskAnalysisActivity, taskId, workflowId, ctxMap).Get(actx, nil)
	return err
}

// Analysis_SuTaskAnalysisActivity
// @Description: 执行可疑分析活动
// @param ctx
// @param taskId
// @param workflowId
// @param ctxMap
// @return err
func Analysis_SuTaskAnalysisActivity(ctx context.Context, taskId, workflowId string, ctxMap map[string]any) (err error) {
	gp.Try(func() error {
		ctx, err := appctx.NewContextWithMap(ctx, ctxMap)
		if err != nil {
			return err
		}
		taskService := NewSuTaskService()
		return taskService.analysisTask(ctx, taskId, workflowId)
	}).Catch(func(e error) {
		err = e
	})
	return err
}
