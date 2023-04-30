package workflows

import (
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

// This is registration process where you register all your workflows
// and activity function handlers.
func init() {
	workflow.Register(LeadWorkflow)
}

func createLeadWorkflowState() WorkflowState {

	workflowState := WorkflowState{
		Current: WorkflowStep{
			Action: "select-degree",
			Index: 1,
			Status: "IN_PROGRESS",
			WorkflowID: nil,
		},
		Steps: []WorkflowStep{
			{
				Action: "select-degree",
				Index: 1,
				Status: "IN_PROGRESS",
				WorkflowID: nil,
			},
			{
				Action: "select-stream",
				Index: 2,
				Status: "NOT_STARTED",
				WorkflowID: nil,
			},
			{
				Action: "select-experience",
				Index: 3,
				Status: "NOT_STARTED",
				WorkflowID: nil,
			},
		},
	}

	return workflowState
}

func LeadWorkflow(ctx workflow.Context) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	logger := workflow.GetLogger(ctx)
	logger.Info("Teacher Setup workflow started")
	workflowState := createLeadWorkflowState()

	err := workflow.SetQueryHandler(ctx, "state", func(input []byte) (WorkflowState, error) {
		return workflowState, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed: " + err.Error())
	}

	// SELECT DEGREE
	var activityResult string
	err = workflow.ExecuteActivity(ctx, templateActivity, "Select Degree").Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Select Degree Activity failed.", zap.Error(err))
		return "", err
	}

	signalName := SignalName
  	selector := workflow.NewSelector(ctx)
 	var data Mystruct
	signalChan := workflow.GetSignalChannel(ctx, signalName)
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		workflowState.Steps[0].Status = "COMPLETED"
		workflowState.Steps[1].Status = "IN_PROGRESS"
		workflowState.Current = workflowState.Steps[1]
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))


	// SELECT STREAM
	err = workflow.ExecuteActivity(ctx, templateActivity, "Select Stream").Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Select Stream Activity failed.", zap.Error(err))
		return "", err
	}
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		workflowState.Steps[1].Status = "COMPLETED"
		workflowState.Steps[2].Status = "IN_PROGRESS"
		workflowState.Current = workflowState.Steps[2]
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))

	// SELECT EXPERIENCE
	err = workflow.ExecuteActivity(ctx, profileActivity, "Select Experience").Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Select Experience Activity failed.", zap.Error(err))
		return "", err
	}
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		workflowState.Steps[2].Status = "COMPLETED"
		workflowState.Current = workflowState.Steps[2]
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)

	// Wait for signal
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))

	return "Teacher Setup Completed", nil
}