package agent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestAgentWorkflow(t *testing.T) {
	type testCase struct {
		task     Task
		expected float64
	}

	tests := []testCase{
		{Task{1, 10, 5, "+", 10}, 15},
		{Task{2, 10, 5, "-", 20}, 5},
		{Task{3, 10, 5, "*", 30}, 50},
		{Task{4, 10, 5, "/", 40}, 2},
		{Task{5, 100, 0, "+", 5}, 100},
		{Task{6, 7, 3, "*", 0}, 21},
		{Task{7, 8, 2, "/", 1}, 4},
		{Task{8, 0, 5, "+", 3}, 5},
		{Task{9, 5, 5, "-", 2}, 0},
		{Task{10, -5, 5, "+", 1}, 0},
		{Task{11, -5, -5, "*", 1}, 25},
		{Task{12, 20, -4, "/", 1}, -5},
		{Task{13, 1.5, 2.5, "+", 0}, 4},
		{Task{14, 9, 3, "/", 0}, 3},
		{Task{15, 5, 0, "/", 0}, 0}, // division by zero (ожидаем ошибку)

		// Новые 15 примеров:
		{Task{16, 999999, 1, "+", 1}, 1000000},
		{Task{17, 3.14, 2, "*", 0}, 6.28},
		{Task{18, 0, 0, "+", 1}, 0},
		{Task{19, -100, -200, "-", 0}, 100},
		{Task{20, 2, 2, "*", 1}, 4},
		{Task{21, 9, 3, "-", 1}, 6},
		{Task{22, 3, 3, "/", 1}, 1},
		{Task{23, 7, 8, "+", 0}, 15},
		{Task{24, -5, 10, "*", 0}, -50},
		{Task{25, 20, 4, "/", 0}, 5},
		{Task{26, 1000, 1, "/", 0}, 1000},
		{Task{27, 123456, 654321, "+", 0}, 777777},
		{Task{28, 0.1, 0.2, "+", 0}, 0.3},
		{Task{29, 50, 25, "-", 0}, 25},
		{Task{30, 10, 2, "/", 0}, 5},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("TaskID_%d", tc.task.ID), func(t *testing.T) {
			var resultReceived bool
			var receivedResult float64

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					if strings.HasSuffix(r.URL.Path, "/internal/task") {
						data, _ := json.Marshal(tc.task)
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusOK)
						w.Write(data)
						return
					}
				case http.MethodPost:
					if strings.HasSuffix(r.URL.Path, "/internal/task") {
						var payload struct {
							TaskID int    `json:"taskid"`
							Result string `json:"result"`
						}
						err := json.NewDecoder(r.Body).Decode(&payload)
						if err != nil {
							t.Errorf("неправильный json : %v", err)
							w.WriteHeader(http.StatusBadRequest)
							return
						}
						if payload.TaskID != tc.task.ID {
							t.Errorf("ожидается TaskID %d, got %d", tc.task.ID, payload.TaskID)
						}
						resultReceived = true
						receivedResult, _ = strconv.ParseFloat(payload.Result, 64)
						w.WriteHeader(http.StatusOK)
						return
					}
				}
				http.NotFound(w, r)
			}))
			defer server.Close()

			serverURL = server.URL

			taskFetched, err := getTask()
			if err != nil {
				t.Fatalf("ошибка в получении задачи: %v", err)
			}

			result, err := compute(*taskFetched)
			if err != nil {
				if tc.task.Operation == "/" && tc.task.Arg2 == 0 {
					return
				}
				t.Fatalf("ошибка в вычислении: %v", err)
			}

			err = sendResult(taskFetched.ID, result)
			if err != nil {
				t.Fatalf("ошибка в отправке: %v", err)
			}

			if !resultReceived {
				t.Fatalf("ошибка :сервер не получил данные")
			}

			const tolerance = 0.0001
			diff := receivedResult - tc.expected
			if diff < 0 {
				diff = -diff
			}
			if diff > tolerance {
				t.Errorf("Incorrect result. Expected %f, got %f", tc.expected, receivedResult)
			}
		})
	}
}
