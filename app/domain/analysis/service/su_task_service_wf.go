package service

import (
	"context"
	"go.temporal.io/sdk/workflow"
	"time"
)

func TaskAnalysisWorkflow(ctx workflow.Context, taskId string, workflowId string) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 24 * time.Hour,
	}
	aCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(aCtx, taskId, workflowId).Get(aCtx, nil)
	return err
}

func TaskAnalysisActivity(ctx context.Context, taskId, workflowId string) (err error) {
	taskService := NewSuTaskService()
	return taskService.analysis(ctx, taskId, workflowId)
}
