package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"go.temporal.io/sdk/workflow"
	"time"
)

func TaskAnalysisWorkflow(wctx workflow.Context, taskId string, workflowId string, ctxMap map[string]any) error {
	logs.Infofmt(nil, "TaskAnalysisWorkflow [%s]", taskId)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 24 * time.Hour,
	}
	actx := workflow.WithActivityOptions(wctx, ao)
	err := workflow.ExecuteActivity(actx, TaskAnalysisActivity, taskId, workflowId, ctxMap).Get(actx, nil)
	return err
}

func TaskAnalysisActivity(ctx context.Context, taskId, workflowId string, ctxMap map[string]any) (err error) {
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
