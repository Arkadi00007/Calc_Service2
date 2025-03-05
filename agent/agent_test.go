package agent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestGetTask(t *testing.T) {
	mockTask := Task{
		ID:        1,
		Arg1:      10,
		Arg2:      5,
		Operation: "+",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		response := struct {
			Task Task `json:"task"`
		}{Task: mockTask}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	serverURL = server.URL
	task, err := getTask()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if task.ID != mockTask.ID || task.Arg1 != mockTask.Arg1 || task.Arg2 != mockTask.Arg2 || task.Operation != mockTask.Operation {
		t.Errorf("Task data mismatch: got %+v, want %+v", task, mockTask)
	}
}

func TestSendResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var result struct {
			ID     int     `json:"taskid"`
			Result float64 `json:"result"`
		}

		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if result.ID != 1 || result.Result != 15 {
			t.Errorf("Incorrect result: got %+v, want {ID: 1, Result: 15}", result)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	serverURL = server.URL
	sendResult(1, 15)
}

func TestCompute(t *testing.T) {
	os.Setenv("TIME_ADDITION_MS", "10")
	os.Setenv("TIME_SUBTRACTION_MS", "10")
	os.Setenv("TIME_MULTIPLICATIONS_MS", "10")
	os.Setenv("TIME_DIVISIONS_MS", "10")

	tests := []struct {
		name     string
		task     Task
		expected float64
	}{
		{"Addition", Task{ID: 1, Arg1: 10, Arg2: 5, Operation: "+"}, 15},
		{"Subtraction", Task{ID: 2, Arg1: 10, Arg2: 5, Operation: "-"}, 5},
		{"Multiplication", Task{ID: 3, Arg1: 10, Arg2: 5, Operation: "*"}, 50},
		{"Division", Task{ID: 4, Arg1: 10, Arg2: 5, Operation: "/"}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			compute(tt.task)
			duration := time.Since(start)

			if duration < 10*time.Millisecond {
				t.Errorf("Expected at least 10ms delay, got %v", duration)
			}
		})
	}
}

func TestSendError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var errResponse struct {
			ID      int    `json:"taskid"`
			Message string `json:"message"`
		}

		if err := json.NewDecoder(r.Body).Decode(&errResponse); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if errResponse.ID != 1 || errResponse.Message != "Division by zero error" {
			t.Errorf("Incorrect error response: got %+v, want {ID: 1, Message: \"Division by zero error\"}", errResponse)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	serverURL = server.URL // Переопределяем URL для тестов
}
