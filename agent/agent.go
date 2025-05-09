package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var serverURL = "http://localhost:8080"

type Task struct {
	ID             int     `json:"id"`
	Arg1           float64 `json:"arg1"`
	Arg2           float64 `json:"arg2"`
	Operation      string  `json:"operation"`
	Operation_time int     `json:"operation_time"`
}

func getTask() (*Task, error) {
	resp, err := http.Get(serverURL + "/internal/task")
	if err != nil {
		log.Printf("getTask request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		log.Println("No tasks available, retrying...")
		return nil, fmt.Errorf("no tasks available")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil || len(body) == 0 {
		//log.Println("getTask received empty response")
		return nil, fmt.Errorf("empty response")
	}

	log.Printf("getTask response: %s", string(body))

	var task Task
	if err := json.Unmarshal(body, &task); err != nil {
		log.Printf("error decoding JSON: %v", err)
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return &task, nil
}

func sendResult(taskID int, result float64) error {
	taskResult := struct {
		TaskID int    `json:"taskid"`
		Result string `json:"result"`
	}{
		TaskID: taskID,
		Result: strconv.FormatFloat(result, 'f', -1, 64),
	}

	data, err := json.Marshal(taskResult)
	if err != nil {
		return fmt.Errorf("Error marshaling result: %v", err)
	}

	resp, err := http.Post(serverURL+"/internal/task", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("Error sending result: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error: Received non-OK response: %s", resp.Status)
	}
	return nil
}

func compute(task Task) (float64, error) {
	var result float64
	time.Sleep(time.Duration(task.Operation_time) * time.Millisecond)
	switch task.Operation {
	case "+":
		result = task.Arg1 + task.Arg2
	case "-":
		result = task.Arg1 - task.Arg2
	case "*":
		result = task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		result = task.Arg1 / task.Arg2
	default:
		return 0, fmt.Errorf("unknown operation: %s", task.Operation)
	}

	return result, nil
}

func worker() {
	for {
		task, err := getTask()
		if err != nil {
			time.Sleep(1 * time.Second)

			continue
		}

		log.Printf("Received task: %+v", task)

		result, err := compute(*task)
		if err != nil {
			log.Printf("Error computing task: %v", err)
			continue
		}

		err = sendResult(task.ID, result)
		if err != nil {
			log.Printf("Error sending result: %v", err)
		} else {
			log.Printf("Sent result: %f for task ID: %d", result, task.ID)
		}

		time.Sleep(1 * time.Second)
	}
}

func Agents() {

	go worker()
	select {}
}
