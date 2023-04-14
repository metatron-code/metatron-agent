package app

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Task struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type TaskResponse struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"` // Unix timestamp for start of task
	Endtime   int64  `json:"endtime"`   // Unix timestamp for end of task
}

func (app *App) mqttEventTask(_ mqtt.Client, msg mqtt.Message) {
	var task Task

	if err := json.Unmarshal(msg.Payload(), &task); err != nil {
		log.Println("error unmarshal message:", err)
		return
	}

	resp := TaskResponse{
		ID:        task.ID,
		Type:      task.Type,
		Timestamp: time.Now().Unix(),
	}

	switch task.Type {
	case "int-ping":
		if err := app.mqttTaskResponse(resp); err != nil {
			log.Println("error make response:", err)
		}

	default:
		log.Println("unknown task type:", task.Type)
	}
}

func (app *App) mqttTaskResponse(resp TaskResponse) error {
	resp.Endtime = time.Now().Unix()

	dataBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	responseTopic := fmt.Sprintf("metatron-agent/%s/response", app.mqttAuthConf.ThingName)

	if token := app.mqtt.Publish(responseTopic, 1, false, dataBytes); token.Wait() && token.Error() != nil {
		log.Println("error publish:", err)
		app.mqttErrors++
	}

	return nil
}
