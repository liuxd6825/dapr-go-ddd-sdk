package main

import (
	"context"
	"fmt"
	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/workflow"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	appcmd "github.com/liuxd6825/dapr-go-ddd-sdk/restapp/cmd"
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
			println("---- OnStartEvent ----")
			/*
				daprClient := server.DaprClient()
				w, err := workflow.NewWorker(workflow.WorkerWithDaprClient(daprClient))
				if err != nil {
					return err
				}
				newWorkflow(w, daprClient)
				fmt.Println("Worker initialized")
			*/
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

	defer daprClient.Close()
	ctx := context.Background()

	wfClient, err := workflow.NewClient(workflow.WithDaprClient(daprClient))
	if err != nil {
		log.Fatalf("failed to intialise client: %v", err)
	}
	// Start workflow test
	instanceID, err := wfClient.ScheduleNewWorkflow(ctx, "TestWorkflow", workflow.WithInstanceID("a7a4168d-3a1c-41da-8a4f-e7f6d9c718d9"), workflow.WithInput(1))
	if err != nil {
		log.Fatalf("failed to start workflow: %v", err)
	}
	fmt.Printf("workflow started with id: %v\n", instanceID)

	// Start workflow test
	respStart, err := daprClient.StartWorkflowBeta1(ctx, &client.StartWorkflowRequest{
		InstanceID:        "a7a4168d-3a1c-41da-8a4f-e7f6d9c718d9",
		WorkflowComponent: workflowComponent,
		WorkflowName:      "TestWorkflow",
		Options:           nil,
		Input:             1,
		SendRawInput:      false,
	})
	if err != nil {
		log.Fatalf("failed to start workflow: %v", err)
	}
	fmt.Printf("workflow started with id: %v\n", respStart.InstanceID)

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
