package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/eclipse/paho.golang/paho"
	"github.com/metatron-code/metatron-agent/internal/tasks"
)

type Task struct {
	ID     string          `json:"id"`
	Type   string          `json:"type"`
	Params json.RawMessage `json:"params"`
}

type TaskResponse struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Timestamp int64           `json:"timestamp"` // Unix timestamp for start of task
	Endtime   int64           `json:"endtime"`   // Unix timestamp for end of task
	Params    json.RawMessage `json:"params,omitempty"`
	Response  json.RawMessage `json:"response,omitempty"`
}

func (app *App) mqttEventTask(msg *paho.Publish) {
	var task Task

	if err := json.Unmarshal(msg.Payload, &task); err != nil {
		log.Println("error unmarshal message:", err)
		return
	}

	resp := TaskResponse{
		ID:        task.ID,
		Type:      task.Type,
		Timestamp: time.Now().Unix(),
	}

	if task.Params != nil {
		resp.Params = task.Params
	}

	log.Printf("received task - type: %s, id: %s", task.Type, task.ID)

	switch task.Type {
	case "icmp-ping":
		task, err := tasks.NewIcmpPing(task.Params)
		if err != nil {
			log.Println("error init icmp-ping task:", err)
			return
		}

		taskResp, err := task.Run()
		if err != nil {
			log.Println("error run icmp-ping task:", err)
			return
		}
		resp.Response = taskResp

		if err := app.mqttTaskResponse(resp); err != nil {
			log.Println("error make response:", err)
		}

	case "int-ping":
		if err := app.mqttTaskResponse(resp); err != nil {
			log.Println("error make response:", err)
		}

	default:
		if len(task.Type) > 1 {
			log.Println("unknown task type:", task.Type)
		}
	}
}

func (app *App) mqttTaskResponse(resp TaskResponse) error {
	resp.Endtime = time.Now().Unix()

	dataBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	msg := &paho.Publish{
		Topic:   fmt.Sprintf("metatron-agent/%s/response", app.mqttAuthConf.ThingName),
		QoS:     1,
		Payload: dataBytes,
		Properties: &paho.PublishProperties{
			ContentType: "application/json",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	if _, err := app.mqtt.Publish(ctx, msg); err != nil {
		log.Println("error publish:", err)
		app.mqttErrors++
	}

	return nil
}
