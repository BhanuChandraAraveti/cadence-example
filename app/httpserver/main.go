// app/httpserver/main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/BhanuChandraAraveti/cadence-example/app/adapters/cadenceAdapter"
	"github.com/BhanuChandraAraveti/cadence-example/app/config"
	"github.com/BhanuChandraAraveti/cadence-example/app/worker/workflows"

	s "go.uber.org/cadence/.gen/go/shared"
	"go.uber.org/cadence/client"
	"go.uber.org/zap"
)

type Service struct {
	cadenceAdapter *cadenceAdapter.CadenceAdapter
	logger         *zap.Logger
}

func (h *Service) triggerSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		applicantID := r.URL.Query().Get("applicant_id")
		h.logger.Info("####### flow!", zap.String("applicantId", applicantID))

		wo := client.StartWorkflowOptions{
			TaskList:                     workflows.TaskListName,
			ExecutionStartToCloseTimeout: time.Hour * 24,
		}
		execution, err := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.SignupWorkflow, applicantID)
		if err != nil {
			http.Error(w, "Error starting workflow!", http.StatusBadRequest)
			return
		}

		h.logger.Info("Started work flow!", zap.String("WorkflowId", execution.ID), zap.String("RunId", execution.RunID))
		js, _ := json.Marshal(execution)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	} else {
		_, _ = w.Write([]byte("Invalid Method!" + r.Method))
	}
}


func (h *Service) triggerTeacherJourney(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		wo := client.StartWorkflowOptions{
			TaskList:                     workflows.TaskListName,
			ExecutionStartToCloseTimeout: time.Hour * 24,
		}
		execution, err := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.TeacherJourneyWorkflow)
		if err != nil {
			http.Error(w, "Error starting teacher journey workflow!", http.StatusBadRequest)
			return
		}

		h.logger.Info("Started teacher journey flow!", zap.String("WorkflowId", execution.ID), zap.String("RunId", execution.RunID))
		js, _ := json.Marshal(execution)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	} else {
		_, _ = w.Write([]byte("Invalid Method!" + r.Method))
	}
}


func (h *Service) triggerOrientation(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		applicantID := r.URL.Query().Get("applicant_id")
		h.logger.Info("####### flow!", zap.String("applicantId", applicantID))

		wo := client.StartWorkflowOptions{
			TaskList:                     workflows.TaskListName,
			ExecutionStartToCloseTimeout: time.Hour * 24,
		}
		execution, err := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.OrientationWorkflow, applicantID)
		if err != nil {
			http.Error(w, "Error starting orientation workflow!", http.StatusBadRequest)
			return
		}

		resp, err := h.cadenceAdapter.CadenceClient.QueryWorkflowWithOptions(context.Background(), &client.QueryWorkflowWithOptionsRequest{
			WorkflowID:            execution.ID,
			RunID:                 execution.RunID,
			QueryType:             "state",
			QueryConsistencyLevel: s.QueryConsistencyLevelStrong.Ptr(),
		})

		if err != nil {
			http.Error(w, "Error starting orientation workflow!", http.StatusBadRequest)
			return
		}

		queryResult:= workflows.Response{}
		resp.QueryResult.Get(&queryResult.WorkflowState)
		queryResult.WorkflowID = execution.ID
		queryResult.RunID = execution.RunID
		h.logger.Info("Started orientation workflow!", zap.String("WorkflowId", execution.ID), zap.String("RunId", execution.RunID))
		h.logger.Info("Query Result", zap.Any("hasValue", queryResult))
		//execution.AppendObject("query", resp)

		js, _ := json.Marshal(execution)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	} else {
		_, _ = w.Write([]byte("Invalid Method!" + r.Method))
	}
}

func (h *Service) triggerSetup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		applicantID := r.URL.Query().Get("applicant_id")
		h.logger.Info("####### flow!", zap.String("applicantId", applicantID))

		wo := client.StartWorkflowOptions{
			TaskList:                     workflows.TaskListName,
			ExecutionStartToCloseTimeout: time.Hour * 24,
		}
		execution, err := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.SetupWorkflow, applicantID)
		if err != nil {
			http.Error(w, "Error starting Setup workflow!", http.StatusBadRequest)
			return
		}

		resp, err := h.cadenceAdapter.CadenceClient.QueryWorkflowWithOptions(context.Background(), &client.QueryWorkflowWithOptionsRequest{
			WorkflowID:            execution.ID,
			RunID:                 execution.RunID,
			QueryType:             "state",
			QueryConsistencyLevel: s.QueryConsistencyLevelStrong.Ptr(),
		})

		if err != nil {
			http.Error(w, "Error starting Setup workflow!", http.StatusBadRequest)
			return
		}

		queryResult:= workflows.Response{}
		resp.QueryResult.Get(&queryResult.WorkflowState)
		queryResult.WorkflowID = execution.ID
		queryResult.RunID = execution.RunID
		h.logger.Info("Started Setup workflow!", zap.String("WorkflowId", execution.ID), zap.String("RunId", execution.RunID))
		h.logger.Info("Query Result", zap.Any("hasValue", queryResult))
		//execution.AppendObject("query", resp)

		js, _ := json.Marshal(queryResult)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	} else {
		_, _ = w.Write([]byte("Invalid Method!" + r.Method))
	}
}

func (h *Service) triggerOnboarding(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		applicantID := r.URL.Query().Get("applicant_id")
		h.logger.Info("####### flow!", zap.String("applicantId", applicantID))

		wo := client.StartWorkflowOptions{
			TaskList:                     workflows.TaskListName,
			ExecutionStartToCloseTimeout: time.Hour * 24,
		}
		execution, err := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.OnboardingWorkflow, applicantID)
		if err != nil {
			http.Error(w, "Error starting Onboarding workflow!", http.StatusBadRequest)
			return
		}

		resp, err := h.cadenceAdapter.CadenceClient.QueryWorkflowWithOptions(context.Background(), &client.QueryWorkflowWithOptionsRequest{
			WorkflowID:            execution.ID,
			RunID:                 execution.RunID,
			QueryType:             "state",
			QueryConsistencyLevel: s.QueryConsistencyLevelStrong.Ptr(),
		})

		if err != nil {
			http.Error(w, "Error starting Onboarding workflow!", http.StatusBadRequest)
			return
		}

		queryResult:= workflows.Response{}
		resp.QueryResult.Get(&queryResult.WorkflowState)
		queryResult.WorkflowID = execution.ID
		queryResult.RunID = execution.RunID
		h.logger.Info("Started Onboarding workflow!", zap.String("WorkflowId", execution.ID), zap.String("RunId", execution.RunID))
		h.logger.Info("Query Result", zap.Any("hasValue", queryResult))
		//execution.AppendObject("query", resp)

		js, _ := json.Marshal(queryResult)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	} else {
		_, _ = w.Write([]byte("Invalid Method!" + r.Method))
	}
}


func (h *Service) getStatusSingle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.logger.Info("$$$$$")

		data := Mystruct{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		workflowID := data.WorkflowId
		runID := data.RunId
		h.logger.Info(workflowID)
		h.logger.Info("payload", zap.Any("data", data))

		resp, err := h.cadenceAdapter.CadenceClient.QueryWorkflowWithOptions(context.Background(), &client.QueryWorkflowWithOptionsRequest{
			WorkflowID:            workflowID,
			RunID:                 runID,
			QueryType:             "state",
			QueryConsistencyLevel: s.QueryConsistencyLevelStrong.Ptr(),
		})

		if err != nil {
			http.Error(w, "Error getting status workflow!", http.StatusBadRequest)
			return
		}

		queryResult:= workflows.Response2{}
		resp.QueryResult.Get(&queryResult.WorkflowData)
		queryResult.WorkflowID = workflowID
		queryResult.RunID = runID
		h.logger.Info("Query Result", zap.Any("hasValue", queryResult))
		//execution.AppendObject("query", resp)

		js, _ := json.Marshal(queryResult)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	} else {
		_, _ = w.Write([]byte("Invalid Method!" + r.Method))
	}
}

func (h *Service) getStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.logger.Info("$$$$$")

		data := Mystruct{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		workflowID := data.WorkflowId
		runID := data.RunId
		h.logger.Info(workflowID)
		h.logger.Info("payload", zap.Any("data", data))

		resp, err := h.cadenceAdapter.CadenceClient.QueryWorkflowWithOptions(context.Background(), &client.QueryWorkflowWithOptionsRequest{
			WorkflowID:            workflowID,
			RunID:                 runID,
			QueryType:             "state",
			QueryConsistencyLevel: s.QueryConsistencyLevelStrong.Ptr(),
		})

		if err != nil {
			http.Error(w, "Error getting status workflow!", http.StatusBadRequest)
			return
		}

		queryResult:= workflows.Response{}
		resp.QueryResult.Get(&queryResult.WorkflowState)
		queryResult.WorkflowID = workflowID
		queryResult.RunID = runID
		h.logger.Info("Query Result", zap.Any("hasValue", queryResult))
		//execution.AppendObject("query", resp)


		resp, err = h.cadenceAdapter.CadenceClient.QueryWorkflowWithOptions(context.Background(), &client.QueryWorkflowWithOptionsRequest{
			WorkflowID:            queryResult.Execution.WorkflowID,
			RunID:                 queryResult.Execution.RunID,
			QueryType:             "state",
			QueryConsistencyLevel: s.QueryConsistencyLevelStrong.Ptr(),
		})

		if err != nil {
			http.Error(w, "Error getting status workflow!", http.StatusBadRequest)
			return
		}

		childQueryResult:= workflows.Response{}
		resp.QueryResult.Get(&childQueryResult.WorkflowState)
		childQueryResult.WorkflowID = workflowID
		childQueryResult.RunID = runID
		h.logger.Info("Child Query Result", zap.Any("hasValue", childQueryResult))

		js, _ := json.Marshal(queryResult)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	} else {
		_, _ = w.Write([]byte("Invalid Method!" + r.Method))
	}
}

func (h *Service) parentStart(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		applicantID := r.URL.Query().Get("applicant_id")
		h.logger.Info("####### flow!", zap.String("applicantId", applicantID))

		wo := client.StartWorkflowOptions{
			TaskList:                     workflows.TaskListName,
			ExecutionStartToCloseTimeout: time.Hour * 24,
		}
		execution, err := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.SampleParentWorkflow)
		if err != nil {
			http.Error(w, "Error starting workflow!", http.StatusBadRequest)
			return
		}

		h.logger.Info("Parent Started work flow!", zap.String("WorkflowId", execution.ID), zap.String("RunId", execution.RunID))
		js, _ := json.Marshal(execution)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	} else {
		_, _ = w.Write([]byte("Invalid Method!" + r.Method))
	}
}


func (h *Service) LastCompletedActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		workflowId := r.URL.Query().Get("workflowId")
		runId := r.URL.Query().Get("runId")
		

		// To iterate all events,
		var isLongPoll bool
		iter := h.cadenceAdapter.CadenceClient.GetWorkflowHistory(context.Background(), workflowId, runId, isLongPoll, s.HistoryEventFilterTypeAllEvent)
		events := []*s.HistoryEvent{}
		h.logger.Info("$$$$$", zap.Any("iter",iter), zap.Any("events", events))
		//var lastActivity string
		for iter.HasNext() {
			event, err := iter.Next()
			if err != nil {
				return
			}
			events = append(events, event)
			eventName := event.GetEventType().String()
			h.logger.Info("Task Ids", zap.String("Event name", eventName))
			if *event.EventType == s.EventTypeActivityTaskCompleted {
				// Store the name of the last completed activity.
				//lastActivity = event.ActivityTaskCompletedEventAttributes.ActivityType.Name
			}
		}
		h.logger.Info("******", zap.String("WorkflowId", workflowId))

		js, _ := json.Marshal("Success")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	} else {
		_, _ = w.Write([]byte("Invalid Method!" + r.Method))
	}
}

func (h *Service) signalHelloWorld(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		workflowId := r.URL.Query().Get("workflowId")
		age, err := strconv.Atoi(r.URL.Query().Get("age"))
		if err != nil {
			h.logger.Error("Failed to parse age from request!")
		}

		err = h.cadenceAdapter.CadenceClient.SignalWorkflow(context.Background(), workflowId, "", workflows.SignalName, age)
		if err != nil {
			http.Error(w, "Error signaling workflow!", http.StatusBadRequest)
			return
		}

		h.logger.Info("Signaled work flow with the following params!", zap.String("WorkflowId", workflowId), zap.Int("Age", age))

		js, _ := json.Marshal("Success")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	} else {
		_, _ = w.Write([]byte("Invalid Method!" + r.Method))
	}
}


type Mystruct struct {
	WorkflowId string `json:"workflowId"`
	RunId string `json:"runId"`
	Payload interface{} `json:"payload"`
	ApplicantId string `json:"applicantId"`
}

func (h *Service) submit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.logger.Info("$$$$$")

		data := Mystruct{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		workflowId := data.WorkflowId
		h.logger.Info(workflowId)
		h.logger.Info("payload", zap.Any("data", data))
		//runId := payload.RunId

		err = h.cadenceAdapter.CadenceClient.SignalWorkflow(context.Background(), workflowId, "", workflows.SignalName, data)
		if err != nil {
			http.Error(w, "Error signaling workflow!", http.StatusBadRequest)
			return
		}

		h.logger.Info("Signaled work flow with the following params!", zap.String("WorkflowId", workflowId), zap.Int("Age", 1))

		js, _ := json.Marshal("Success")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	} else {
		_, _ = w.Write([]byte("Invalid Method!" + r.Method))
	}
}


func (h *Service) orientationStart(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		accountId := r.URL.Query().Get("accountId")

		wo := client.StartWorkflowOptions{
			TaskList:                     workflows.TaskListName,
			ExecutionStartToCloseTimeout: time.Hour * 24,
		}
		execution, err := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.SampleParentWorkflow, accountId)
		if err != nil {
			http.Error(w, "Error starting workflow!", http.StatusBadRequest)
			return
		}

		h.logger.Info("Started work flow!", zap.String("WorkflowId", execution.ID), zap.String("RunId", execution.RunID))
		js, _ := json.Marshal(execution)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(js)
	} else {
		_, _ = w.Write([]byte("Invalid Method!" + r.Method))
	}
}

// func (h *Service) listChildWorkflowIDs(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "POST" {
// 		parentWorkflowID := r.URL.Query().Get("parentWorkflowID")
// 		runID := r.URL.Query().Get("runId")

// 		workflowStub := h.cadenceAdapter.CadenceClient.GetWorkflow(context.Background(), parentWorkflowID, runID)
// 		listChildExecutionsResponse, err := workflowStub.ListChildExecutions(context.Background(), &s.ListChildExecutionsRequest{})

// 		if err != nil {
// 			return nil, err
// 		}
	
// 		// Extract the child workflow IDs from the response
// 		var childWorkflowIDs []string
// 		for _, childExecutionInfo := range listChildExecutionsResponse.Executions {
// 			childWorkflowIDs = append(childWorkflowIDs, childExecutionInfo.GetExecution().GetWorkflowId())
// 		}

// 		h.logger.Info("Signaled work flow with the following params!")

// 		js, _ := json.Marshal(childWorkflowIDs)

// 		w.Header().Set("Content-Type", "application/json")
// 		_, _ = w.Write(js)
// 	} else {
// 		_, _ = w.Write([]byte("Invalid Method!" + r.Method))
// 	}
// }

// func listChildWorkflowIDs(parentWorkflowID string) ([]string, error) {
//     // Create a Cadence client
//     c, err := client.NewClient(client.Options{})
//     if err != nil {
//         return nil, err
//     }

//     // Get a workflow stub for the parent workflow
//     workflowOptions := client.StartWorkflowOptions{
//         ID:        parentWorkflowID,
//         TaskQueue: "my-task-queue",
//     }
//     workflowClient := c.WorkflowClient
//     workflowStub := workflowClient.GetWorkflow(workflowOptions.ID, "", workflowOptions)

//     // Use the workflow stub to call ListChildExecutions API
//     listChildExecutionsResponse, err := workflowStub.ListChildExecutions(context.Background(), &shared.ListChildExecutionsRequest{})
//     if err != nil {
//         return nil, err
//     }

//     // Extract the child workflow IDs from the response
//     var childWorkflowIDs []string
//     for _, childExecutionInfo := range listChildExecutionsResponse.Executions {
//         childWorkflowIDs = append(childWorkflowIDs, childExecutionInfo.GetExecution().GetWorkflowId())
//     }

//     return childWorkflowIDs, nil
// }


func (h *Service)getWorkflowHistory(c client.Client, domain string, workflowID string, runID string) (*s.History, error) {
	execution := s.WorkflowExecution{
		WorkflowId: &workflowID,
		RunId:      &runID,
	}
	historyRequest := s.GetWorkflowExecutionHistoryRequest{
		Domain:    &domain,
		Execution: &execution,
	}
	ctx := context.Background()
	h.logger.Info("@@@", zap.Any("WorkflowExecution", execution))
	h.logger.Info("@@@", zap.Any("GetWorkflowExecutionHistoryRequest", historyRequest))
	h.logger.Info("@@@", zap.Any("context Background", context.Background()))
	historyResp, err := h.cadenceAdapter.ServiceClient.GetWorkflowExecutionHistory(ctx, &historyRequest)
	if err != nil {
		h.logger.Info("@@@", zap.Any("err", err))
		return nil, err
	}
	return historyResp.History, nil
}

// func (h *Service)getWorkflowHistory(c client.Client, domain string, workflowID string, runID string) (*s.History, error) {
// 	execution := s.WorkflowExecution{
// 		WorkflowId: &workflowID,
// 		RunId:      &runID,
// 	}
// 	historyRequest := s.GetWorkflowExecutionHistoryRequest{
// 		Domain:    &domain,
// 		Execution: &execution,
// 	}
// 	ctx := context.Background()
// 	h.logger.Info("@@@", zap.Any("WorkflowExecution", execution))
// 	h.logger.Info("@@@", zap.Any("GetWorkflowExecutionHistoryRequest", historyRequest))
// 	h.logger.Info("@@@", zap.Any("context Background", context.Background()))
// 	iter := c.GetWorkflowHistory(ctx, workflowID, runID, false, s.HistoryEventFilterTypeAllEvent)
// 	for iter.HasNext() {
// 		event, err := iter.Next()
// 		if err != nil {
// 			h.logger.Info("@@@", zap.Any("err", err))
// 			return nil, err
// 		}
// 		if event.GetEventType() == s.EventTypeActivityTaskCompleted {
// 		}
// 		if event.GetEventType() == s.EventTypeActivityTaskFailed {
// 		}
// 	}
	
// 	return historyResp.History, nil
// }

func processWorkflowHistory(history *s.History) ([]*s.HistoryEvent, *s.HistoryEvent) {
	var taskList []*s.HistoryEvent
	var nextTask *s.HistoryEvent

	for _, event := range history.Events {
		switch event.GetEventType() {
		case s.EventTypeActivityTaskScheduled:
			taskList = append(taskList, event)
			if nextTask == nil {
				nextTask = event
			}
		case s.EventTypeActivityTaskCompleted:
			nextTask = nil
		}
	}

	return taskList, nextTask
}

func (h *Service) check(w http.ResponseWriter, r *http.Request) {
	//cadenceClient := createCadenceClient("7833")
	cadenceClient := h.cadenceAdapter.CadenceClient
	history, err := h.getWorkflowHistory(cadenceClient, "simple-domain","df07595c-c61d-4a23-8187-4c742b4641da", "4872f259-f4c7-4036-9477-70d43d54c1e5")
	if err != nil {
		panic("Failed to get workflow history.")
	}

	taskList, nextTask := processWorkflowHistory(history)
	fmt.Println("Task list:", taskList)
	fmt.Println("Next task:", nextTask)
	fmt.Println("######")
}

// func createCadenceClient(hostPort string) client.Client {
// 	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName("cadence-client"), tchannel.ListenAddr(hostPort))
// 	if err != nil {
// 		panic("Failed to create TChannel transport.")
// 	}

// 	dispatcher := yarpc.NewDispatcher(yarpc.Config{
// 		Name: "cadence-client",
// 		Outbounds: yarpc.Outbounds{
// 			"cadence": {Unary: ch.NewSingleOutbound(hostPort)},
// 		},
// 	})
// 	dispatcher.Start()
// 	defer dispatcher.Stop()

// 	runtime := z.NewScope()
// 	cadenceClient := workflowserviceclient.New(dispatcher.ClientConfig("cadence"), client.Options{MetricsScope: runtime})
// 	return client.New(cadenceClient)
// }




func main() {
	var appConfig config.AppConfig
	appConfig.Setup()
	var cadenceClient cadenceAdapter.CadenceAdapter
	cadenceClient.Setup(&appConfig.Cadence)

	service := Service{&cadenceClient, appConfig.Logger}
	http.HandleFunc("/api/start-teacher-onboarding", service.triggerTeacherJourney)
	http.HandleFunc("/api/start-signup-workflow", service.triggerSignup)
	http.HandleFunc("/api/start-orientation-workflow", service.triggerOrientation)
	http.HandleFunc("/api/start-setup-workflow", service.triggerSetup)
	http.HandleFunc("/api/start-onboarding-workflow", service.triggerOnboarding)
	http.HandleFunc("/api/get-current-screen", service.LastCompletedActivity)
	http.HandleFunc("/api/submit", service.submit)
	http.HandleFunc("/api/signal-hello-world", service.signalHelloWorld)
	http.HandleFunc("/api/orientation-start", service.orientationStart)
	http.HandleFunc("/api/start-parent", service.parentStart)
	http.HandleFunc("/api/history", service.check)
	http.HandleFunc("/api/get-status-single", service.getStatusSingle)
	http.HandleFunc("/api/get-status", service.getStatus)

	addr := ":3030"
	log.Println("Starting Server! Listening on:", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
