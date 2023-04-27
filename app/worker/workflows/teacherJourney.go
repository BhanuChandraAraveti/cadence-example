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
	workflow.Register(TeacherJourneyWorkflow)
}

func createTeacherJourneyState() WorkflowState {

	workflowState := WorkflowState{
		Current: WorkflowStep{
			Action: "signup",
			Index: 1,
			Status: "IN_PROGRESS",
			WorkflowID: nil,
		},
		Steps: []WorkflowStep{
			{
				Action: "signup",
				Index: 1,
				Status: "IN_PROGRESS",
				WorkflowID: nil,
			},
			{
				Action: "lead",
				Index: 2,
				Status: "NOT_STARTED",
				WorkflowID: nil,
			},
			{
				Action: "application",
				Index: 2,
				Status: "NOT_STARTED",
				WorkflowID: nil,
			},
		},
	}

	return workflowState
}

func TeacherJourneyWorkflow(ctx workflow.Context) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	logger := workflow.GetLogger(ctx)
	logger.Info("Teacher Onboarding workflow started")
	
	workflowState := createTeacherJourneyState()

	err := workflow.SetQueryHandler(ctx, "state", func(input []byte) (WorkflowState, error) {
		return workflowState, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed: " + err.Error())
	}
	
	// Signup Workflow
	execution := workflow.GetInfo(ctx).WorkflowExecution
	// Parent workflow can choose to specify it's own ID for child execution.  Make sure they are unique for each execution.
	childID := fmt.Sprintf("signup:%v", execution.RunID)
	cwo := workflow.ChildWorkflowOptions{
		// Do not specify WorkflowID if you want cadence to generate a unique ID for child execution
		WorkflowID:                   childID,
		ExecutionStartToCloseTimeout: time.Hour,
	}
	ctx = workflow.WithChildOptions(ctx, cwo)
	var result string
	err = workflow.ExecuteChildWorkflow(ctx, SignupWorkflow).Get(ctx, &result)
	if err != nil {
		logger.Error("Parent execution received child execution failure.", zap.Error(err))
		return "", err
	}

	workflowState.Steps[0].Status = "COMPLETED"
	workflowState.Steps[1].Status = "IN_PROGRESS"
	workflowState.Current = workflowState.Steps[1]


	// Lead Workflow
	childID = fmt.Sprintf("lead:%v", execution.RunID)
	cwo = workflow.ChildWorkflowOptions{
		WorkflowID:                   childID,
		ExecutionStartToCloseTimeout: time.Hour,
	}
	ctx = workflow.WithChildOptions(ctx, cwo)
	err = workflow.ExecuteChildWorkflow(ctx, LeadWorkflow).Get(ctx, &result)
	if err != nil {
		logger.Error("Parent execution received child execution failure.", zap.Error(err))
		return "", err
	}

	workflowState.Steps[1].Status = "COMPLETED"
	workflowState.Steps[2].Status = "IN_PROGRESS"
	workflowState.Current = workflowState.Steps[2]


	// Application Workflow
	childID = fmt.Sprintf("application:%v", execution.RunID)
	err = workflow.ExecuteChildWorkflow(ctx, ApplicationWorkflow).Get(ctx, &result)
	if err != nil {
		logger.Error("Parent execution received child execution failure.", zap.Error(err))
		return "", err
	}

	workflowState.Steps[2].Status = "COMPLETED"
	workflowState.Current = workflowState.Steps[2]

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

	return "Teacher Onboarding Completed", nil
}