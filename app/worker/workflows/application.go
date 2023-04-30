package workflows

import (
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)


func init() {
	workflow.Register(ApplicationWorkflow)
}

func createApplicationWorkflowState() WorkflowState {

	workflowState := WorkflowState{
		Current: WorkflowStep{
			Action: "watch-video",
			Index: 1,
			Status: "IN_PROGRESS",
			WorkflowID: nil,
		},
		Steps: []WorkflowStep{
			{
				Action: "watch-video",
				Index: 1,
				Status: "IN_PROGRESS",
				WorkflowID: nil,
			},
			{
				Action: "select-grade",
				Index: 2,
				Status: "NOT_STARTED",
				WorkflowID: nil,
			},
			{
				Action: "screening",
				Index: 3,
				Status: "NOT_STARTED",
				WorkflowID: nil,
			},
		},
	}

	return workflowState
}

func ApplicationWorkflow(ctx workflow.Context) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	logger := workflow.GetLogger(ctx)
	logger.Info("Teacher Application workflow started")
	
	workflowState := createLeadWorkflowState()

	err := workflow.SetQueryHandler(ctx, "state", func(input []byte) (WorkflowState, error) {
		return workflowState, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed: " + err.Error())
	}

	// WATCH VIDEO
	var activityResult string
	err = workflow.ExecuteActivity(ctx, templateActivity, "Watch Video").Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Watch Video Activity failed.", zap.Error(err))
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


	// SELECT GRADE
	err = workflow.ExecuteActivity(ctx, templateActivity, "Select Grade").Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Select Stream Grade failed.", zap.Error(err))
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

	// SCREENING
	err = workflow.ExecuteActivity(ctx, profileActivity, "Screening").Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Screening Activity failed.", zap.Error(err))
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