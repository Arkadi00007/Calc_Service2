package internal

import (
	"awesomeProject/pkg/calculation"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Структуры
type Expression struct {
	ID         int       `json:"id"`
	Expression *[]string `json:"expression"`
	Status     string    `json:"status"`
	Result     float64   `json:"result,omitempty"`
}

type Task struct {
	ID        int     `json:"id"`
	Arg1      float64 `json:"arg1"`
	Arg2      float64 `json:"arg2"`
	Operation string  `json:"operation"`
}

// Состояние оркестратора
var expressions = make(map[int]*Expression) // Здесь используем указатели
var mu sync.Mutex
var idCounter int

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req map[string]string
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		expression := req["expression"]
		idCounter++

		// Преобразование выражения в постфиксную нотацию
		postfix, err := calculation.Calc(expression)

		if err != nil {
			http.Error(w, "Invalid expression", http.StatusBadRequest)
		}
		// Создание задачи для вычислений
		mu.Lock()
		expressions[idCounter] = &Expression{
			ID:         idCounter,
			Expression: &postfix,
			Status:     "processing",
		}
		mu.Unlock()

		// Ответ с ID задачи
		response := map[string]int{"id": idCounter}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetExpressions(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var exprList []*Expression // Используем указатели
	for _, expr := range expressions {
		exprList = append(exprList, expr)
	}

	json.NewEncoder(w).Encode(map[string][]*Expression{"expressions": exprList}) // Возвращаем указатели
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		type Result struct {
			TaskID int    `json:"taskid"`
			Result string `json:"result"`
		}
		var result Result
		err := json.NewDecoder(r.Body).Decode(&result)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
		}
		for _, expr := range expressions {
			if expr.ID == result.TaskID {
				if result.Result == "Division by zero error" {
					http.Error(w, "Division by zero error", http.StatusNotFound)
					return
				}

				for j := 0; j < len(*expr.Expression); j++ {
					v := (*(*expr).Expression)[j]
					if calculation.IsOperator(v[0]) {
						(*(*expr).Expression) = append(append((*(*expr).Expression)[:j-2], result.Result), (*(*expr).Expression)[j+1:]...)
						break
					}

				}
				if len((*expr.Expression)) == 1 {
					num, _ := strconv.ParseFloat((*expr.Expression)[0], 64)
					expr.Status = "completed"
					expr.Result = num
				}

			}
		}

	}
	if r.Method == http.MethodGet {
		mu.Lock()
		defer mu.Unlock()

		// Поиск задачи
		for _, expr := range expressions {
			if expr.Status == "processing" {
				// Преобразуем выражение в задачи
				tasks := createTasks(expr.ID, (expr.Expression))
				json.NewEncoder(w).Encode(tasks)
				return
			}
		}
		http.Error(w, "No tasks available", http.StatusNotFound)
	}
}

func handleGetExpressionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                 // Получаем параметры из URL
	id, err := strconv.Atoi(vars["id"]) // Конвертируем ID в число
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Поиск выражения
	for _, expr := range expressions {
		if expr.ID == id {
			json.NewEncoder(w).Encode(map[string]*Expression{"expression": expr})
			return
		}
	}

	http.Error(w, "Expression not found", http.StatusNotFound)
}

func Ara(num float64, err error) float64 {
	return num
}

func createTasks(id int, expression *[]string) *Task {
	var task *Task

	for i := 0; i < len(*expression); i++ {
		v := (*expression)[i]
		taskID := id
		if calculation.IsOperator(v[0]) {
			task = &Task{
				ID:        taskID,
				Arg1:      Ara(strconv.ParseFloat((*expression)[i-2], 64)),
				Arg2:      Ara(strconv.ParseFloat((*expression)[i-1], 64)),
				Operation: string(v[0]),
			}
		}
	}
	return task
}

func RunServer() {
	http.HandleFunc("/api/v1/calculate", handleCalculate)
	http.HandleFunc("/api/v1/expressions", handleGetExpressions)
	http.HandleFunc("/internal/task", handleGetTask)
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/expressions/{id}", handleGetExpressionByID)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
