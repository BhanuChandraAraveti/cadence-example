// app/worker/workflows/helloworldworkflows.go
package workflows

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

/**
 * This is the hello world workflow sample.
 */

// ApplicationName is the task list for this sample
const TaskListName = "helloWorldGroup"
const SignalName = "helloWorldSignal"

// This is registration process where you register all your workflows
// and activity function handlers.
func init() {
	workflow.Register(Workflow)
	activity.Register(overviewActivity)
	activity.Register(evalSOPActivity)
	activity.Register(evalCETActivity)
	activity.Register(degreeDetailsActivity)
	activity.Register(watchVideoActivity)
}

var activityOptions = workflow.ActivityOptions{
	ScheduleToStartTimeout: time.Minute,
	StartToCloseTimeout:    time.Minute,
	HeartbeatTimeout:       time.Second * 20,
	// RetryPolicy: &cadence.RetryPolicy{
	// 	InitialInterval:          time.Second,
	// 	BackoffCoefficient:       2.0,
	// 	MaximumInterval:          time.Minute,
	// 	ExpirationInterval:       time.Minute * 5,
	// 	MaximumAttempts:          5,
	// 	NonRetriableErrorReasons: []string{"bad-error"},
	// },
}


func call_api() {
	resp, err := http.Get("https://64397c471b9a7dd5c968fa7d.mockapi.io/tasks/3")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(body))
}

func evalSOPActivity(ctx context.Context, name string) (int, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("SOP evaluation activity started")
	call_api()
	return 70, nil
}


func evalCETActivity(ctx context.Context, name string) (int, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Cat evaluation activity started")
	call_api()
	return 80, nil
}

func overviewActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Overview activity started")
	call_api()
	return "Hello " + name + "! How old are you!", nil
}

// func Workflow(ctx workflow.Context, name string) (string, error) {
// 	ctx = workflow.WithActivityOptions(ctx, activityOptions)

// 	logger := workflow.GetLogger(ctx)
// 	logger.Info("Teacher workflow started")
// 	var activityResult string
// 	err := workflow.ExecuteActivity(ctx, overviewActivity, name).Get(ctx, &activityResult)
// 	if err != nil {
// 		logger.Error("Overview Activity failed.", zap.Error(err))
// 		return "", err
// 	}

// 	// After saying hello, the workflow will wait for you to inform it of your age!
// 	signalName := SignalName
// 	selector := workflow.NewSelector(ctx)
// 	var ageResult int

// 	for {
// 		signalChan := workflow.GetSignalChannel(ctx, signalName)
// 		selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
// 			c.Receive(ctx, &ageResult)
// 			workflow.GetLogger(ctx).Info("Received age results from signal!", zap.String("signal", signalName), zap.Int("value", ageResult))
// 		})
// 		workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
// 		// Wait for signal
// 		selector.Select(ctx)

// 		// We can check the age and return an appropriate response
// 		if ageResult > 50{
// 			logger.Info("Workflow completed.", zap.String("NameResult", activityResult), zap.Int("AgeResult", ageResult))

// 			return fmt.Sprintf("Hello "+name+"! Let's make teaching fun!", ageResult), nil
// 		}
			
// 		var futures []workflow.Future
// 		// starts activities in parallel
// 		ao := workflow.ActivityOptions{
// 			ScheduleToStartTimeout: time.Minute,
// 			StartToCloseTimeout:    time.Minute,
// 			HeartbeatTimeout:       time.Second * 20,
// 		}
// 		ctx = workflow.WithActivityOptions(ctx, ao)

// 		totalBranches := 2
// 		//evaluationActivities := [evalCETActivity, evalSOPActivity]
// 		for i := 1; i <= totalBranches; i++ {
// 			activityInput := fmt.Sprintf("branch %d of %d.", i, totalBranches)
// 			future := workflow.ExecuteActivity(ctx, evalCETActivity, activityInput)
// 			futures = append(futures, future)
// 		}

// 		// wait until all futures are done
// 		var sum int
// 		for _, future := range futures {
// 			var result int
// 			if err := future.Get(ctx, &result); err != nil {
// 				return "", err
// 			}
// 			sum += result
// 		}

// 		if sum < 150 {
// 			workflow.NewContinueAsNewError(ctx, Workflow, name)
// 		}

// 		workflow.GetLogger(ctx).Info("Workflow completed.")
// 	}
// }


func degreeDetailsActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("degree details activity started")
	// Ask frontend to show the degreeDetails Screen
	call_api()
	logger.Info("degree details activity ended")
	return "", nil
}

func watchVideoActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("watch video activity started")
	// Ask frontend to show the watchVideo Screen
	call_api()
	logger.Info("watch video activity ended")
	return "", nil
}

type State struct{

}


func Workflow(ctx workflow.Context, name string) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	logger := workflow.GetLogger(ctx)
	logger.Info("Teacher signup workflow started")
	var activityResult string
	err := workflow.ExecuteActivity(ctx, degreeDetailsActivity, name).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Overview Activity failed.", zap.Error(err))
		return "", err
	}
	signalChan := workflow.GetSignalChannel(ctx, signalName)
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &ageResult)
		workflow.GetLogger(ctx).Info("Received age results from signal!", zap.String("signal", signalName), zap.Int("value", ageResult))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
	// Wait for signal
	selector.Select(ctx)

	err = workflow.ExecuteActivity(ctx, watchVideoActivity, name).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Overview Activity failed.", zap.Error(err))
		return "", err
	}

	logger.Info("Workflow completed.")
	return "Workflow completed.", nil
}

