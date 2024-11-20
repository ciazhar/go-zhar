package temporal

import (
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

func GetDefaultActivityOptions() workflow.ActivityOptions {
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    100 * time.Second,
		MaximumAttempts:    500,
	}

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy:         retryPolicy,
	}

	return options
}

//
//func ExecuteActivity[I any, R any](ctx workflow.Context, activityFunc interface{}, input I) (R, error) {
//	var res R
//	err := workflow.ExecuteActivity(ctx, activityFunc, input).Get(ctx, &res)
//	return res, err
//}
