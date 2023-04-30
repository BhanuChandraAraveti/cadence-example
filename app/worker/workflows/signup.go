package workflows

import (
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

// This is registration process where you register all your workflows
// and activity function handlers.
func init() {
	workflow.Register(SignupWorkflow)
}

func SignupWorkflow(ctx workflow.Context) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	logger := workflow.GetLogger(ctx)
	logger.Info("Teacher Signup workflow started")

	info := workflow.GetInfo(ctx)
  	workflowID := info.WorkflowExecution.ID
	runID := info.WorkflowExecution.RunID

	//parentState.Current.WorkflowID = &workflowID
	//parentState.Current.Status = "IT WORKED"

	workflowStep := WorkflowStep{
		Action: "signup",
		Index: 1,
		Status: "IN_PROGRESS",
		WorkflowID: &workflowID,
		RunID: &runID,
	}
	
	workflowState := WorkflowState{
		Current: workflowStep,
		Steps: []WorkflowStep{workflowStep},
	}

	err := workflow.SetQueryHandler(ctx, "state", func(input []byte) (WorkflowState, error) {
		return workflowState, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed: " + err.Error())
	}
	
	var activityResult string
	err = workflow.ExecuteActivity(ctx, templateActivity, "Signup").Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Signup Activity failed.", zap.Error(err))
		return "", err
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

	return "Teacher Signup Completed", nil
}