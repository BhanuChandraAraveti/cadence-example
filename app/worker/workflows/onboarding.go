// app/worker/workflows/cueteacherOnboarding.go
package workflows

import (
	"fmt"
	"time"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

// This is registration process where you register all your workflows
// and activity function handlers.
func init() {
	workflow.Register(OnboardingWorkflow)
}

func createOnboardingWorkflowState() WorkflowState {

	workflowState := WorkflowState{
		Current: WorkflowStep{
			Action: "orientation",
			Index: 1,
			Status: "IN_PROGRESS",
			WorkflowID: nil,
		},
		Steps: []WorkflowStep{
			{
				Action: "orientation",
				Index: 1,
				Status: "IN_PROGRESS",
				WorkflowID: nil,
			},
			{
				Action: "setup",
				Index: 2,
				Status: "NOT_STARTED",
				WorkflowID: nil,
			},
		},
	}

	return workflowState
}

func OnboardingWorkflow(ctx workflow.Context, applicantID string) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	logger := workflow.GetLogger(ctx)
	logger.Info("Teacher Onboarding workflow started")
	logger.Info("Applicant ID: " + applicantID)
	
	workflowState := createOnboardingWorkflowState()

	err := workflow.SetQueryHandler(ctx, "state", func(input []byte) (WorkflowState, error) {
		return workflowState, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed: " + err.Error())
	}

	signalName := SignalName
  	selector := workflow.NewSelector(ctx)
 	var data Mystruct
	signalChan := workflow.GetSignalChannel(ctx, signalName)
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)

	// Wait for signal
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))
	
	// Orientation Workflow
	execution := workflow.GetInfo(ctx).WorkflowExecution
	// Parent workflow can choose to specify it's own ID for child execution.  Make sure they are unique for each execution.
	childID := fmt.Sprintf("orientation:%v", execution.RunID)
	cwo := workflow.ChildWorkflowOptions{
		// Do not specify WorkflowID if you want cadence to generate a unique ID for child execution
		WorkflowID:                   childID,
		ExecutionStartToCloseTimeout: time.Hour*24*7*1000,
	}
	ctx = workflow.WithChildOptions(ctx, cwo)
	var result string
	err = workflow.ExecuteChildWorkflow(ctx, OrientationWorkflow, applicantID).Get(ctx, &result)
	if err != nil {
		logger.Error("Parent execution received child execution failure.", zap.Error(err))
		return "", err
	}

	index := workflowState.Current.Index
	workflowState.Steps[index-1].Status = "COMPLETED"
	workflowState.Steps[index].Status = "IN_PROGRESS"
	workflowState.Current = workflowState.Steps[index]


	// Setup Workflow
	childID = fmt.Sprintf("setup:%v", execution.RunID)
	cwo = workflow.ChildWorkflowOptions{
		WorkflowID:                   childID,
		ExecutionStartToCloseTimeout: time.Hour,
	}
	ctx = workflow.WithChildOptions(ctx, cwo)
	err = workflow.ExecuteChildWorkflow(ctx, SetupWorkflow, applicantID).Get(ctx, &result)
	if err != nil {
		logger.Error("Parent execution received child execution failure.", zap.Error(err))
		return "", err
	}

	index = workflowState.Current.Index
	workflowState.Steps[index-1].Status = "COMPLETED"
	workflowState.Current.Status = "COMPLETED"


	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)

	// Wait for signal
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))

	return "Teacher Onboarding Completed", nil
}