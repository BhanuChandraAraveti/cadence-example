package workflows

import (
	"bytes"
	"context"
	"encoding/json"
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
const SignalName = "submit"

type State struct {
	CurrentActivity string
}

// This is registration process where you register all your workflows
// and activity function handlers.
func init() {
	workflow.Register(Workflow)
	workflow.Register(SampleParentWorkflow)
	workflow.Register(SampleChildWorkflow)
	activity.Register(overviewActivity)
	activity.Register(evalSOPActivity)
	activity.Register(evalCETActivity)
	activity.Register(degreeDetailsActivity)
	activity.Register(watchVideoActivity)
	activity.Register(gradeActivity)
	activity.Register(streamSelectionActivity)
	activity.Register(teacherCETAndSOPActivity)
	activity.Register(uploadLessonVideoActivity)
	activity.Register(submitDocumentsActivity)
	activity.Register(orientationActivity)
	activity.Register(basicDetailsActivity)
	activity.Register(agreementActivity)
	activity.Register(profileActivity)
	activity.Register(availabilityActivity)
	activity.Register(templateActivity)
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
	//state.CurrentActivity = "overview"
	logger.Info("Overview activity started")
	call_api()
	return "Overview activity completed", nil
}

func degreeDetailsActivity(ctx context.Context, applicationID string, workflowID string, runID string) (string, error) {

	logger := activity.GetLogger(ctx)
	logger.Info("degree details activity started")
	// Ask frontend to show the degreeDetails Screen
	// call_api()

	msg, err :=sendWorkflowId(ctx, applicationID, workflowID, runID)
	if err != nil {
		return msg, err
	}

	logger.Info("degree details activity ended")
	return "degree details activity ended", nil
}

func watchVideoActivity(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("watch video activity started")
	// Ask frontend to show the watchVideo Screen
	call_api()
	logger.Info("watch video activity ended")
	return "watch video activity ended", nil
}

func gradeActivity(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Grade Selection activity started")
	// Ask frontend to show the watchVideo Screen
	call_api()
	logger.Info("Grade Selection activity ended")
	return "Grade Selection activity ended", nil
}

func streamSelectionActivity(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Stream Selection activity started")
	// Ask frontend to show the watchVideo Screen
	call_api()
	logger.Info("Stream Selection activity ended")
	return "Stream Selection activity ended", nil
}

func teacherCETAndSOPActivity(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("CET and SOP activity started")
	// Ask frontend to show the watchVideo Screen
	call_api()
	logger.Info("CET and SOP activity ended")
	return "CET and SOP activity ended", nil
}

func uploadLessonVideoActivity(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Lesson upload activity started")
	// Ask frontend to show the watchVideo Screen
	call_api()
	logger.Info("Lesson upload activity ended")
	return "Lesson upload activity ended", nil
}

func submitDocumentsActivity(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Submit documents activity started")
	// Ask frontend to show the watchVideo Screen
	call_api()
	logger.Info("Submit documents activity ended")
	return "Submit documents activity ended", nil
}

func orientationActivity(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Orientation activity started")
	call_api()
	logger.Info("Orientation activity ended")
	return "Orientation activity ended", nil
}

func basicDetailsActivity(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Basic details activity started")
	call_api()
	logger.Info("Basic details activity ended")
	return "Basic details activity ended", nil
}

func agreementActivity(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Agreement activity started")
	call_api()
	logger.Info("Agreement activity ended")
	return "Agreement activity ended", nil
}

func profileActivity(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Profile activity started")
	call_api()
	logger.Info("Profile activity ended")
	return "Profile activity ended", nil
}

func templateActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info(name + " activity started")
	logger.Info(name + " activity ended")
	return "Profile activity ended", nil
}



func availabilityActivity(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Availability activity started")
	call_api()
	logger.Info("Availability activity ended")
	return "Availability activity ended", nil
}

type RequestBody struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func call(data Mystruct) (string, error) {
	status, err := updateProfile(data)
	if err != nil { 
		return status, err
	}
	return "", nil
}

type Mystruct struct {
	WorkflowId string `json:"workflowId"`
	RunId string `json:"runId"`
	Payload interface{} `json:"payload"`
	ApplicantId string `json:"applicantId"`
}

func Workflow(ctx workflow.Context, applicantID string) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	logger := workflow.GetLogger(ctx)
	logger.Info("Teacher signup workflow started")
	logger.Info("Applicant ID: " + applicantID)

	queryResult := 0
	err := workflow.SetQueryHandler(ctx, "state", func(input []byte) (int, error) {
		return queryResult, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed: " + err.Error())
	}

	info := workflow.GetInfo(ctx)
  	workflowID := info.WorkflowExecution.ID
	runID := info.WorkflowExecution.RunID
	// sendWorkflowId(workflowID, runID)


	var activityResult string
	err = workflow.ExecuteActivity(ctx, degreeDetailsActivity, applicantID, workflowID, runID).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Degree Details Activity failed.", zap.Error(err))
		return "", err
	}

	//
	signalName := SignalName
  	selector := workflow.NewSelector(ctx)
 	var data Mystruct
	signalChan := workflow.GetSignalChannel(ctx, signalName)
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		queryResult += 1
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
	// Wait for signal
	selector.Select(ctx)
	logger.Info("payload", zap.Any("data", data))

	var msg string

	// call BE API
	msg, err = call(data)
	logger.Info(msg)
	//
	
	
    // STREAM Selection Activity
	selector = workflow.NewSelector(ctx)
	signalChan = workflow.GetSignalChannel(ctx, signalName)
	err = workflow.ExecuteActivity(ctx, streamSelectionActivity).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Watch Video Activity failed.", zap.Error(err))
		return "", err
	}
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		queryResult += 1
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
	// Wait for signal
	selector.Select(ctx)

	// call BE API
	msg, err = call(data)
	logger.Info(msg)



	// Grade Activity
	selector = workflow.NewSelector(ctx)
	signalChan = workflow.GetSignalChannel(ctx, signalName)
	err = workflow.ExecuteActivity(ctx, gradeActivity).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Watch Video Activity failed.", zap.Error(err))
		return "", err
	}
	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		queryResult += 1
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
	// Wait for signal
	selector.Select(ctx)

	// call BE API
	msg, err = call(data)
	logger.Info(msg)



	// WATCH VIDEO
	selector = workflow.NewSelector(ctx)
	signalChan = workflow.GetSignalChannel(ctx, signalName)
	err = workflow.ExecuteActivity(ctx, watchVideoActivity).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Watch Video Activity failed.", zap.Error(err))
		return "", err
	}

	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		queryResult += 1
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
	// Wait for signal
	selector.Select(ctx)



	// CET and SOP
	selector = workflow.NewSelector(ctx)
	signalChan = workflow.GetSignalChannel(ctx, signalName)
	err = workflow.ExecuteActivity(ctx, teacherCETAndSOPActivity).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Watch Video Activity failed.", zap.Error(err))
		return "", err
	}

	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		queryResult += 1
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
	// Wait for signal
	selector.Select(ctx)


	// Upload Lesson Video
	selector = workflow.NewSelector(ctx)
	signalChan = workflow.GetSignalChannel(ctx, signalName)
	err = workflow.ExecuteActivity(ctx, uploadLessonVideoActivity).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Watch Video Activity failed.", zap.Error(err))
		return "", err
	}

	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		queryResult += 1
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
	// Wait for signal
	selector.Select(ctx)


	// Submit Documents
	selector = workflow.NewSelector(ctx)
	signalChan = workflow.GetSignalChannel(ctx, signalName)
	err = workflow.ExecuteActivity(ctx, submitDocumentsActivity).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Watch Video Activity failed.", zap.Error(err))
		return "", err
	}

	selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &data)
		queryResult += 1
		workflow.GetLogger(ctx).Info("Received the signal!", zap.String("signal", signalName))
	})
	workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
	// Wait for signal
	selector.Select(ctx)

	// type Mystruct struct {
	// 	WorkflowId string `json:"workflowId"`
	// 	RunId string `json:"runId"`
	// 	Payload interface{} `json:"payload"`
	// 	ApplicantId string `json:"applicantId"`
	// }

	data.Payload = struct{
		CreateTeacher bool `json:"CREATE_TEACHER"`
	}{
		CreateTeacher: true,
	}
	// call BE API
	msg, err = call(data)
	logger.Info(msg)
	//

	logger.Info("Workflow completed.")
	return "Workflow completed.", nil
}

type BEStruct struct {
	ApplicantID       string `json:"applicant_id"`
	ProfileAttributes interface{} `json:"profile_attributes"`
}


func updateProfile(data Mystruct) (string, error){
	updateProfileRequest := BEStruct{
		ApplicantID: data.ApplicantId,
		ProfileAttributes: data.Payload,
	}
	fmt.Println(updateProfileRequest)

	url := "https://admin.testenv6.cuemath.com/teacher/applicant-profile?=null"
	requestBody, err := json.Marshal(updateProfileRequest)
	if err != nil {
		return "ERROR BE call", err
	}

	request, err := http.NewRequest("PATCH", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "ERROR BE call", err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "ERROR BE call", err
	}
	defer response.Body.Close()

	fmt.Println("Response Status:", response.Status)
	return response.Status, nil
}


func sendWorkflowId(ctx context.Context, applicantID string, workflowID string, runID string) (string, error){
	updateProfileRequest := BEStruct{
		ApplicantID: applicantID,
		ProfileAttributes: struct{
			WorkflowID string `json:"workflow_id"`
			RunID string `json:"run_id"`
		}{
			WorkflowID: workflowID,
			RunID: runID,
		},
	}

	logger := activity.GetLogger(ctx)
	logger.Info("************************")
	logger.Info("payload", zap.Any("send", updateProfileRequest))

	

	url := "https://admin.testenv6.cuemath.com/teacher/applicant-profile?=null"
	requestBody, err := json.Marshal(updateProfileRequest)
	if err != nil {
		return "ERROR BE call", err
	}

	request, err := http.NewRequest("PATCH", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "ERROR BE call", err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "ERROR BE call", err
	}
	defer response.Body.Close()

	fmt.Println("Response Status:", response.Status)
	return response.Status, nil
}



