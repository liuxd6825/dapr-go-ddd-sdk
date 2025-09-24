package tasks

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"go.temporal.io/sdk/activity"
	tclient "go.temporal.io/sdk/client"
	tworker "go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

var _client tclient.Client
var _worker tworker.Worker
var _cfg Config

func GetTaskQueue() string {
	return _cfg.TaskQueue
}

func Connect(cfg Config) error {
	if cfg.HostPort == "" {
		return errors.New("temporal no hostport specified")
	}
	if cfg.Namespace == "" {
		return errors.New("temporal no namespace specified")
	}
	if cfg.TaskQueue == "" {
		return errors.New("temporal no task queue specified")
	}
	_cfg = cfg

	var err error
	_client, err = tclient.NewLazyClient(tclient.Options{
		Namespace: cfg.Namespace,
		HostPort:  cfg.HostPort,
	})
	if err != nil {
		return err
	}
	_worker = tworker.New(_client, _cfg.TaskQueue, tworker.Options{})
	return nil
}

func RunWorker() {
	if _client == nil {
		return
	}
	go func() {
		_worker.Run(tworker.InterruptCh())
	}()
}

func ExecuteWorkflow(ctx context.Context, options tclient.StartWorkflowOptions, workflow any, args ...any) (tclient.WorkflowRun, error) {
	return _client.ExecuteWorkflow(ctx, options, workflow, args...)
}

func RegisterWorkflow(ctx context.Context, w any) {
	_worker.RegisterWorkflow(w)
}

func RegisterWorkflowWithOptions(ctx context.Context, workflow any, options workflow.RegisterOptions) {
	_worker.RegisterWorkflowWithOptions(workflow, options)
}

func RegisterActivity(ctx context.Context, activity any) {
	_worker.RegisterActivity(activity)
}

func RegisterActivityWithOptions(ctx context.Context, activity any, options activity.RegisterOptions) {
	_worker.RegisterActivityWithOptions(activity, options)
}
