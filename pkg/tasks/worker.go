package tasks

import (
	"context"
	"go.temporal.io/sdk/activity"
	tclient "go.temporal.io/sdk/client"
	tworker "go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

var client tclient.Client
var worker tworker.Worker

type ConnectConfig struct {
	HostPort  string
	Namespace string
	TaskQueue string
}

func Connect(cfg ConnectConfig) error {
	var err error
	client, err = tclient.NewLazyClient(tclient.Options{
		Namespace: cfg.Namespace,
		HostPort:  cfg.HostPort,
	})
	if err != nil {
		return err
	}
	worker = tworker.New(client, cfg.TaskQueue, tworker.Options{})
	return nil
}

func ExecuteWorkflow(ctx context.Context, options tclient.StartWorkflowOptions, workflow any, args ...any) (tclient.WorkflowRun, error) {
	return client.ExecuteWorkflow(ctx, options, workflow, args...)
}

func RegisterWorkflow(w any) {
	worker.RegisterWorkflow(w)
}

func RegisterWorkflowWithOptions(workflow any, options workflow.RegisterOptions) {
	worker.RegisterWorkflowWithOptions(workflow, options)
}

func RegisterActivity(activity any) {
	worker.RegisterActivity(activity)
}

func RegisterActivityWithOptions(activity any, options activity.RegisterOptions) {
	worker.RegisterActivityWithOptions(activity, options)
}

func Run() error {
	return worker.Run(tworker.InterruptCh())
}
