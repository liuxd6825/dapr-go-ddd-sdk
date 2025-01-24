package main

import (
	"context"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/workflow"
	iris_context "github.com/kataras/iris/v12/context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	appcmd "github.com/liuxd6825/dapr-go-ddd-sdk/restapp/cmd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/samber/lo"
	"log"
	"time"
)

var (
	Version   = "1.0.0"
	BuildTime = ""
	GitHead   = ""
)

func main() {
	appcmd.StartApp(&appcmd.AppStartOptions{
		AppTitle:  "Web服务器",
		Version:   Version,
		BuildTime: BuildTime,
		GitHead:   GitHead,
		Actors:    nil,
		OnInitEvent: func(server *restapp.HttpServer) error {
			println("---- OnInitEvent ----")
			return nil
		},
		OnStartEvent: func(server *restapp.HttpServer) error {
			return nil
			println("---- OnStartEvent ----")

			daprClient := server.DaprClient()
			w, err := workflow.NewWorker(workflow.WorkerWithDaprClient(daprClient))
			if err != nil {
				return err
			}

			go func() {
				newWorkflow(w, daprClient)
				fmt.Println("Worker initialized")
			}()

			server.App().Get("/api/v1.0/workflow", func(ictx *iris_context.Context) {
				iid := idutils.NewId()
				lo.TryCatchWithErrorValue(func() error {
					ctx := context.Background()
					// Start workflow test
					respStart, e := daprClient.StartWorkflowBeta1(ctx, &client.StartWorkflowRequest{
						InstanceID:        iid,
						WorkflowComponent: workflowComponent,
						WorkflowName:      "TestWorkflow",
						Options:           nil,
						Input:             1,
						SendRawInput:      false,
					})
					if e == nil && respStart != nil {
						fmt.Printf("workflow started with id: %v\n", respStart.InstanceID)
					}
					return err
				}, func(e any) {
					if v, ok := e.(error); ok {
						ictx.SetErr(v)
					}
					ictx.StatusCode(500)
				})
				println("/api/v1.0/workflow")
			})
			return nil
		},
	})
}

var stage = 0

const (
	workflowComponent = "dapr"
)

func newWorkflow(w *workflow.WorkflowWorker, daprClient dapr.Client) {

	w, err := workflow.NewWorker(workflow.WorkerWithDaprClient(daprClient))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Worker initialized")

	if err := w.RegisterWorkflow(TestWorkflow); err != nil {
		log.Fatal(err)
	}
	fmt.Println("TestWorkflow registered")

	if err := w.RegisterActivity(TestActivity); err != nil {
		log.Fatal(err)
	}
	fmt.Println("TestActivity registered")

	// Start workflow runner
	if err := w.Start(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("runner started")

}

func TestWorkflow(ctx *workflow.WorkflowContext) (any, error) {
	var input int
	if err := ctx.GetInput(&input); err != nil {
		return nil, err
	}
	var output string
	if err := ctx.CallActivity(TestActivity, workflow.ActivityInput(input)).Await(&output); err != nil {
		return nil, err
	}

	err := ctx.WaitForExternalEvent("testEvent", time.Second*60).Await(&output)
	if err != nil {
		return nil, err
	}

	if err := ctx.CallActivity(TestActivity, workflow.ActivityInput(input)).Await(&output); err != nil {
		return nil, err
	}

	return output, nil
}

func TestActivity(ctx workflow.ActivityContext) (any, error) {
	var input int
	if err := ctx.GetInput(&input); err != nil {
		return "", err
	}

	stage += input

	return fmt.Sprintf("Stage: %d", stage), nil
}
