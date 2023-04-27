package workflows

import (
	"context"

	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

// This is registration process where you register all your workflows
// and activity function handlers.
func init() {
	workflow.Register(teacherSignupWorkflow)
}

func createTeacherSignupState() WorkflowState {

	workflowState := WorkflowState{
		Current: WorkflowStep{
			Action: "signup",
			Index: 1,
			Status: "IN_PROGRESS",
			WorkflowID: nil,
			RunID: nil,
		},
		Steps: []WorkflowStep{
			{
				Action: "signup",
				Index: 1,
				Status: "IN_PROGRESS",
				WorkflowID: nil,
				RunID: nil,
			},
			{
				Action: "degree",
				Index: 2,
				Status: "NOT_STARTED",
				WorkflowID: nil,
				RunID: nil,
			},
			{
				Action: "stream",
				Index: 3,
				Status: "NOT_STARTED",
				WorkflowID: nil,
				RunID: nil,
			},
			{
				Action: "experience",
				Index: 4,
				Status: "NOT_STARTED",
				WorkflowID: nil,
				RunID: nil,
			},
		},
	}

	return workflowState
}

func signupActivity(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Agreement activity started")
	call_api()
	logger.Info("Agreement activity ended")
	return "Agreement activity ended", nil
}

func teacherSignupWorkflow(ctx workflow.Context, applicantID string) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	logger := workflow.GetLogger(ctx)
	logger.Info("Teacher Signup workflow started")
	logger.Info("Applicant ID: " + applicantID)

	info := workflow.GetInfo(ctx)
  	workflowID := info.WorkflowExecution.ID
	runID := info.WorkflowExecution.RunID

	workflowState := createTeacherSignupState()

	err := workflow.SetQueryHandler(ctx, "state", func(input []byte) (WorkflowState, error) {
		return workflowState, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed: " + err.Error())
	}

	var activityResult string
	err = workflow.ExecuteActivity(ctx, orientationActivity, applicantID, workflowID, runID).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Degree Details Activity failed.", zap.Error(err))
		return "", err
	}

	signalName := SignalName
  	selector := workflow.NewSelector(ctx)
 	var data Mystruct
	signalChan := workflow.GetSignalChannel(ctx, signalName)
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		var msg1 string
		msg1, err = call(data)
		workflowState.Current.Status = msg1
		//workflowState.Current.Status = "COMPLETED"
		workflowState.Steps[0].Status = "COMPLETED"
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)

	// Wait for signal
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))

	// call BE API
	var msg string
	msg, err = call(data)
	logger.Info(msg)

	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)

	// Wait for signal
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))

	return "Teacher Orientation Completed", nil
}