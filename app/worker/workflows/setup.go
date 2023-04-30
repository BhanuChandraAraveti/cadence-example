package workflows

import (
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

// This is registration process where you register all your workflows
// and activity function handlers.
func init() {
	workflow.Register(SetupWorkflow)
}

func createWorkflowState() WorkflowState {

	workflowState := WorkflowState{
		Current: WorkflowStep{
			Action: "basic-details",
			Index: 1,
			Status: "IN_PROGRESS",
			WorkflowID: nil,
		},
		Steps: []WorkflowStep{
			{
				Action: "basic-details",
				Index: 1,
				Status: "IN_PROGRESS",
				WorkflowID: nil,
			},
			{
				Action: "agreement",
				Index: 2,
				Status: "NOT_STARTED",
				WorkflowID: nil,
			},
			{
				Action: "profile",
				Index: 3,
				Status: "NOT_STARTED",
				WorkflowID: nil,
			},
			{
				Action: "availability",
				Index: 4,
				Status: "NOT_STARTED",
				WorkflowID: nil,
			},
		},
	}

	return workflowState
}

func SetupWorkflow(ctx workflow.Context, applicantID string) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	logger := workflow.GetLogger(ctx)
	logger.Info("Teacher Setup workflow started")
	logger.Info("Applicant ID: " + applicantID)

	info := workflow.GetInfo(ctx)
  	workflowID := info.WorkflowExecution.ID
	runID := info.WorkflowExecution.RunID
	
	workflowState := createWorkflowState()

	err := workflow.SetQueryHandler(ctx, "state", func(input []byte) (WorkflowState, error) {
		return workflowState, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed: " + err.Error())
	}
	// BASIC DETAILS
	var activityResult string
	err = workflow.ExecuteActivity(ctx, basicDetailsActivity, applicantID, workflowID, runID).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Basic Details Activity failed.", zap.Error(err))
		return "", err
	}

	signalName := SignalName
  	selector := workflow.NewSelector(ctx)
 	var data Mystruct
	signalChan := workflow.GetSignalChannel(ctx, signalName)
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		index := workflowState.Current.Index
		workflow.GetLogger(ctx).Info("&&&&&&&&&& index", zap.Any("index", workflowState))
		workflowState.Steps[index-1].Status = "COMPLETED"
		workflowState.Steps[index].Status = "IN_PROGRESS"
		workflowState.Current = workflowState.Steps[index]
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

	// AGREEMENT
	err = workflow.ExecuteActivity(ctx, agreementActivity, applicantID, workflowID, runID).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Agreement Activity failed.", zap.Error(err))
		return "", err
	}
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		index := workflowState.Current.Index
		workflow.GetLogger(ctx).Info("&&&&&&&&&& index", zap.Int("index", index))
		workflowState.Steps[index-1].Status = "COMPLETED"
		workflowState.Steps[index].Status = "IN_PROGRESS"
		workflowState.Current = workflowState.Steps[index]
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)

	// Wait for signal
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))

	// call BE API
	msg, err = call(data)
	logger.Info(msg)

	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		index := workflowState.Current.Index
		workflow.GetLogger(ctx).Info("&&&&&&&&&& index", zap.Any("index", workflowState))
		workflowState.Steps[index-1].Status = "COMPLETED"
		workflowState.Steps[index].Status = "IN_PROGRESS"
		workflowState.Current = workflowState.Steps[index]
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)

	// Wait for signal
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))

	// call BE API
	msg, err = call(data)
	logger.Info(msg)

	// PROFILE
	err = workflow.ExecuteActivity(ctx, profileActivity, applicantID, workflowID, runID).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Profile Activity failed.", zap.Error(err))
		return "", err
	}
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		index := workflowState.Current.Index
		workflow.GetLogger(ctx).Info("&&&&&&&&&& index", zap.Any("index", workflowState))
		workflowState.Steps[index-1].Status = "COMPLETED"
		workflowState.Steps[index].Status = "IN_PROGRESS"
		workflowState.Current = workflowState.Steps[index]
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)

	// Wait for signal
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))

	// call BE API
	msg, err = call(data)
	logger.Info(msg)


	// Availability
	err = workflow.ExecuteActivity(ctx, availabilityActivity, applicantID, workflowID, runID).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Profile Activity failed.", zap.Error(err))
		return "", err
	}
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		workflowState.Current.Status = "COMPLETED"
		index := workflowState.Current.Index
		workflow.GetLogger(ctx).Info("&&&&&&&&&& index", zap.Any("index", workflowState))
		workflowState.Steps[index-1].Status = "COMPLETED"
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)

	// Wait for signal
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))

	// call BE API
	msg, err = call(data)
	logger.Info(msg)

	return "Teacher Setup Completed", nil
}