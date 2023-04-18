// app/httpserver/main.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/nndd91/cadence-api-example/app/adapters/cadenceAdapter"
	"github.com/nndd91/cadence-api-example/app/config"
	"github.com/nndd91/cadence-api-example/app/worker/workflows"
	s "go.uber.org/cadence/.gen/go/shared"
	"go.uber.org/cadence/client"
	"go.uber.org/zap"
)

type Service struct {
	cadenceAdapter *cadenceAdapter.CadenceAdapter
	logger         *zap.Logger
}

func (h *Service) triggerHelloWorld(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		accountId := r.URL.Query().Get("accountId")

		wo := client.StartWorkflowOptions{
			TaskList:                     workflows.TaskListName,
			ExecutionStartToCloseTimeout: time.Hour * 24,
		}
		execution, err := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.Workflow, accountId)
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
	Age int `json:"age"`
}

func (h *Service) submit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.logger.Info("$$$$$")

		payload := Mystruct{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		workflowId := payload.WorkflowId
		age:= payload.Age
		//runId := payload.RunId

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


func main() {
	var appConfig config.AppConfig
	appConfig.Setup()
	var cadenceClient cadenceAdapter.CadenceAdapter
	cadenceClient.Setup(&appConfig.Cadence)

	service := Service{&cadenceClient, appConfig.Logger}
	http.HandleFunc("/api/start-signup-workflow", service.triggerHelloWorld)
	http.HandleFunc("/api/get-current-screen", service.LastCompletedActivity)
	http.HandleFunc("/api/submit", service.submit)
	http.HandleFunc("/api/signal-hello-world", service.signalHelloWorld)

	addr := ":3030"
	log.Println("Starting Server! Listening on:", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
