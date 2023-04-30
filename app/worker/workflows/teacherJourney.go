package workflows

import (
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

// This is registration process where you register all your workflows
// and activity function handlers.
func init() {
	workflow.Register(TeacherJourneyWorkflow)
}

type Response2 struct {
	Execution
	WorkflowData WorkflowData `json:"workflow_state"`
}

type Execution struct {
	WorkflowID string `json:"workflow_id"`
	RunID	  string `json:"run_id"`
}
type WorkflowStep2 struct {
    Activity string `json:"activity"`
    Status   string `json:"status"`
}

type WorkflowInfo struct {
    Activity string          `json:"activity"`
    Steps    []WorkflowStep2 `json:"steps"`
}

type WorkflowData struct {
    ParentWorkflowInfo WorkflowInfo  `json:"parent_workflow_info"`
    Activity           string        `json:"activity"`
    Steps              []WorkflowStep2 `json:"steps"`
}


func createTeacherJourneyData() WorkflowData {
	
	workflowData := WorkflowData{
		ParentWorkflowInfo: WorkflowInfo{
			Activity: "lead",
			Steps: []WorkflowStep2{
				{
					Activity: "lead",
					Status: "IN_PROGRESS",
				},
				{
					Activity: "application",
					Status: "NOT_STARTED",
				},
			},
		},
		Activity: "select-degree",
		Steps: []WorkflowStep2{
			{
				Activity: "select-degree",
				Status: "NOT_STARTED",
			},
			{
				Activity: "select-stream",
				Status: "NOT_STARTED",
			},
			{
				Activity: "select-experience",
				Status: "NOT_STARTED",
			},
		},
	}

	return workflowData
}


func TeacherJourneyWorkflow(ctx workflow.Context) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	logger := workflow.GetLogger(ctx)
	logger.Info("Teacher Onboarding workflow started")
	
	workflowData := createTeacherJourneyData()

	err := workflow.SetQueryHandler(ctx, "state", func(input []byte) (WorkflowData, error) {
		return workflowData, nil
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
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))

	workflowData.Activity = "select-stream"
	workflowData.Steps[0].Status = "COMPLETED"


	// SELECT STREAM
	err = workflow.ExecuteActivity(ctx, templateActivity, "Select Stream").Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Select Stream Activity failed.", zap.Error(err))
		return "", err
	}
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))
	workflowData.Activity = "select-experience"
	workflowData.Steps[1].Status = "COMPLETED"

	// SELECT EXPERIENCE
	err = workflow.ExecuteActivity(ctx, templateActivity, "Select Experience").Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Select Experience Activity failed.", zap.Error(err))
		return "", err
	}
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))

	workflowData.ParentWorkflowInfo.Activity = "application"
	workflowData.ParentWorkflowInfo.Steps[0].Status = "COMPLETED"
	workflowData.ParentWorkflowInfo.Steps[1].Status = "IN_PROGRESS"
	workflowData.Activity = "watch-video"
	workflowData.Steps = []WorkflowStep2{
		{
			Activity: "watch-video",
			Status: "NOT_STARTED",
		},
		{
			Activity: "select-grade",
			Status: "NOT_STARTED",
		},
		{
			Activity: "screening",
			Status: "NOT_STARTED",
		},
	}


	// WATCH VIDEO
	err = workflow.ExecuteActivity(ctx, templateActivity, "Watch Video").Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Watch Video Activity failed.", zap.Error(err))
		return "", err
	}


	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))
	workflowData.Activity = "select-grade"
	workflowData.Steps[0].Status = "COMPLETED"


	// SELECT GRADE
	err = workflow.ExecuteActivity(ctx, templateActivity, "Select Grade").Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Select Stream Grade failed.", zap.Error(err))
		return "", err
	}
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))
	workflowData.Activity = "screening"
	workflowData.Steps[1].Status = "COMPLETED"

	// SCREENING
	err = workflow.ExecuteActivity(ctx, templateActivity, "Screening").Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Screening Activity failed.", zap.Error(err))
		return "", err
	}
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))
	workflowData.Steps[2].Status = "IN_PROGRESS"

	// SCREENING IN_PROGRESS
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))
	workflowData.Steps[2].Status = "COMPLETED"
    workflowData.ParentWorkflowInfo.Steps[1].Status = "COMPLETED"
	workflowData.Activity = "success"

	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)

	return "Teacher Journey Completed", nil
}

// func createTeacherJourneyState() WorkflowState {

// 	workflowState := WorkflowState{
// 		Current: WorkflowStep{
// 			Action: "signup",
// 			Index: 1,
// 			Status: "IN_PROGRESS",
// 			WorkflowID: nil,
// 			RunID: nil,
// 		},
// 		Steps: []WorkflowStep{
// 			{
// 				Action: "signup",
// 				Index: 1,
// 				Status: "IN_PROGRESS",
// 				WorkflowID: nil,
// 				RunID: nil,
// 			},
// 			{
// 				Action: "lead",
// 				Index: 2,
// 				Status: "NOT_STARTED",
// 				WorkflowID: nil,
// 				RunID: nil,
// 			},
// 			{
// 				Action: "application",
// 				Index: 2,
// 				Status: "NOT_STARTED",
// 				WorkflowID: nil,
// 				RunID: nil,
// 			},
// 		},
// 	}

// 	return workflowState
// }

// type CHILD struct {
// 	ID    string
// 	RunID string
// }

// func TeacherJourneyWorkflow(ctx workflow.Context) (string, error) {
// 	ctx = workflow.WithActivityOptions(ctx, activityOptions)

// 	logger := workflow.GetLogger(ctx)
// 	logger.Info("Teacher Onboarding workflow started")
	
// 	workflowState := createTeacherJourneyState()

// 	err := workflow.SetQueryHandler(ctx, "state", func(input []byte) (WorkflowState, error) {
// 		return workflowState, nil
// 	})
// 	if err != nil {
// 		logger.Info("SetQueryHandler failed: " + err.Error())
// 	}
	
// 	// Signup Workflow
// 	execution := workflow.GetInfo(ctx).WorkflowExecution
// 	// Parent workflow can choose to specify it's own ID for child execution.  Make sure they are unique for each execution.
// 	childID := fmt.Sprintf("signup:%v", execution.RunID)
// 	cwo := workflow.ChildWorkflowOptions{
// 		// Do not specify WorkflowID if you want cadence to generate a unique ID for child execution
// 		WorkflowID:                   childID,
// 		ExecutionStartToCloseTimeout: time.Hour*24*7*1000,
// 	}
// 	ctx = workflow.WithChildOptions(ctx, cwo)
// 	//var result string
// 	//err = workflow.ExecuteChildWorkflow(ctx, SignupWorkflow, &workflowState).Get(ctx, &result)
// 	childWorkflowFuture := workflow.ExecuteChildWorkflow(ctx, SignupWorkflow)
	
// 	var childWE interface{}
// 	err = childWorkflowFuture.GetChildWorkflowExecution().Get(ctx, &childWE);
// 	if err != nil {
// 		logger.Error("Parent execution received child execution failure.", zap.Error(err))
// 		return "", err
// 	}
// 	// var temp interface{}
// 	// temp = childWE
// 	// childData := temp.(CHILD)
	
// 	// // logger.Info("!!!!!!! workflowID and RunID", zap.Any("errr", something))
// 	// fmt.Println(childWE)
// 	// fmt.Println(childData.ID)
// 	logger.Info("!!!!!!! workflowID and RunID", zap.Any("CHILD", childWE))
// 	logger.Info("!!!!!!! workflowID and RunID", zap.Any("CHILDID", reflect.ValueOf(childWE).FieldByName("ID").Interface()))
// 	logger.Info("!!!!!!! workflowID and RunID", zap.Any("RUNID", reflect.ValueOf(childWE).FieldByName("RunID").Interface()))

// 	childWorkflowID := reflect.ValueOf(childWE).FieldByName("ID").Interface()
// 	childRunID := reflect.ValueOf(childWE).FieldByName("RunID").Interface()

	

// 	// logger.Info("!!!!!!! workflowID and RunID", zap.String("workflowID", childWE.ID), zap.String("RunID", childWE.RunID))

	

// 	workflowState.Steps[0].Status = "COMPLETED"
// 	workflowState.Steps[1].Status = "IN_PROGRESS"
// 	workflowState.Current = workflowState.Steps[1]


// 	// Lead Workflow
// 	childID = fmt.Sprintf("lead:%v", execution.RunID)
// 	cwo = workflow.ChildWorkflowOptions{
// 		WorkflowID:                   childID,
// 		ExecutionStartToCloseTimeout: time.Hour,
// 	}
// 	ctx = workflow.WithChildOptions(ctx, cwo)
// 	var result string
// 	err = workflow.ExecuteChildWorkflow(ctx, LeadWorkflow).Get(ctx, &result)
// 	if err != nil {
// 		logger.Error("Parent execution received child execution failure.", zap.Error(err))
// 		return "", err
// 	}

// 	workflowState.Steps[1].Status = "COMPLETED"
// 	workflowState.Steps[2].Status = "IN_PROGRESS"
// 	workflowState.Current = workflowState.Steps[2]


// 	// Application Workflow
// 	childID = fmt.Sprintf("application:%v", execution.RunID)
// 	err = workflow.ExecuteChildWorkflow(ctx, ApplicationWorkflow).Get(ctx, &result)
// 	if err != nil {
// 		logger.Error("Parent execution received child execution failure.", zap.Error(err))
// 		return "", err
// 	}

// 	workflowState.Steps[2].Status = "COMPLETED"
// 	workflowState.Current = workflowState.Steps[2]

// 	signalName := SignalName
//   	selector := workflow.NewSelector(ctx)
//  	var data Mystruct
// 	signalChan := workflow.GetSignalChannel(ctx, signalName)

// 	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
// 		c.Receive(ctx, &data)
// 		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
// 	})
// 	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)

// 	// Wait for signal
// 	selector.Select(ctx)
// 	logger.Info("payload", zap.Any("data", data))

// 	return "Teacher Onboarding Completed", nil
// }