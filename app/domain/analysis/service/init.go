package service

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/tasks"

// RegisterWorkflow
// @Description: 注册工作流
func RegisterWorkflow() {
	// 分析工作流
	tasks.RegisterWorkflow(TaskAnalysisWorkflow)
	tasks.RegisterActivity(TaskAnalysisActivity)
}
