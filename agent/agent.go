package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	serverURL = "http://localhost:8080"
)

var (
	TIME_ADDITION_MS        = os.Getenv("TIME_ADDITION_MS")
	TIME_SUBTRACTION_MS     = os.Getenv("TIME_SUBTRACTION_MS")
	TIME_MULTIPLICATIONS_MS = os.Getenv("TIME_MULTIPLICATIONS_MS")
	TIME_DIVISIONS_MS       = os.Getenv("TIME_DIVISIONS_MS")
)

type Task struct {
	ID            int     `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

func getTask() (*Task, error) {
	resp, err := http.Get(serverURL + "/internal/task")
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("")
	}
	defer resp.Body.Close()

	var task struct {
		Task Task `json:"task"`
	}
	json.NewDecoder(resp.Body).Decode(&task)

	return &task.Task, nil
}

func sendResult(taskID int, result float64) {

	type TaskResult struct {
		ID     int     `json:"taskid"`
		Result float64 `json:"result"`
	}

	// Создаём объект структуры TaskResult
	taskResult := &TaskResult{
		ID:     taskID,
		Result: result,
	}

	// Кодируем структуру в JSON
	data, err := json.Marshal(taskResult)
	if err != nil {
		fmt.Println("Ошибка при кодировании JSON:", err)
		return
	}

	resp, err := http.Post(serverURL+"/internal/task", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Сервер вернул ошибку:", resp.Status)
	}
}

func compute(task Task) {
	var result float64

	switch task.Operation {
	case "+":
		a, _ := strconv.Atoi(TIME_ADDITION_MS)
		time.Sleep(time.Duration(a) * time.Millisecond)
		result = task.Arg1 + task.Arg2
	case "-":
		b, _ := strconv.Atoi(TIME_SUBTRACTION_MS)
		time.Sleep(time.Duration(b) * time.Millisecond)
		result = task.Arg1 - task.Arg2
	case "*":
		c, _ := strconv.Atoi(TIME_MULTIPLICATIONS_MS)
		time.Sleep(time.Duration(c) * time.Millisecond)
		result = task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0 {
			sendError(task.ID)
			return
		}
		d, _ := strconv.Atoi(TIME_DIVISIONS_MS)
		time.Sleep(time.Duration(d) * time.Millisecond)

		result = task.Arg1 / task.Arg2
	}

	sendResult(task.ID, result)
}

func sendError(taskID int) {
	type TaskError struct {
		ID      int    `json:"taskid"`
		Message string `json:"message"`
	}

	// Создаем объект с ошибкой
	taskError := &TaskError{
		ID:      taskID,
		Message: "Division by zero error",
	}

	// Кодируем структуру в JSON
	data, err := json.Marshal(taskError)
	if err != nil {
		fmt.Println("Ошибка при кодировании JSON:", err)
		return
	}

	// Отправляем ошибку на сервер оркестратора
	resp, err := http.Post(serverURL+"/internal/task", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

}

func worker() {
	for {
		task, err := getTask()
		if err != nil {
			time.Sleep(1 * time.Second) // Если задач нет, подождать
			continue
		}
		compute(*task)
	}
}

func Man() {
	a := 4

	for i := 0; i <
		a; i++ {
		go worker()
	}
	select {}
}
