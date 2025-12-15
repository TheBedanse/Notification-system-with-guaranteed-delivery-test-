package workflow

import (
	"fmt"
	"math/rand"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func GreetingWorkflow(ctx workflow.Context, name string) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 10,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var result string
	err := workflow.ExecuteActivity(ctx, ComposeGreeting, name).Get(ctx, &result)
	if err != nil {
		return "", err
	}

	workflow.Sleep(ctx, time.Second*1)

	return result, nil
}

func ComposeGreeting(name string) (string, error) {
	//Test Retry Policy
	if rand.Intn(100) < 40 { //40% err
		return "", fmt.Errorf("random error in Activity")
	}
	greeting := fmt.Sprintf("Hello %s! Application processed via Temporal!", name)
	time.Sleep(500 * time.Millisecond)
	return greeting, nil
}
