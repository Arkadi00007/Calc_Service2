package internal

import (
	"awesomeProject/pkg/calculation"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var jwtKey = []byte("your_secret_key") // ключ для подписи JWT
var db *sql.DB

type contextKey string

const loginContextKey contextKey = "login"

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Expression struct {
	ID         int       `json:"id"`
	Expression *[]string `json:"expression"`
	Status     string    `json:"status"`
	Result     float64   `json:"result,omitempty"`
	SubStatus  string    `json:"sub_status,omitempty"`
	Login      string    `json:"login"`
}

type Task struct {
	ID             int     `json:"id"`
	Arg1           float64 `json:"arg1"`
	Arg2           float64 `json:"arg2"`
	Operation      string  `json:"operation"`
	Operation_time int     `json:"operation_time"`
}

var expressions = make(map[int]*Expression)
var mu sync.Mutex
var mumu sync.Mutex
var idCounter int

func GenerateJWT(login string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["login"] = login
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix() // токен истекает через 24 часа
	return token.SignedString(jwtKey)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Missing token"})
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		login, ok := claims["login"].(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), loginContextKey, login)
		next(w, r.WithContext(ctx))
	}
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		var req map[string]string
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		expression := req["expression"]
		login := r.Context().Value(loginContextKey).(string)

		postfix, err := calculation.Calc(expression)

		if err != nil {
			http.Error(w, "Invalid expression", http.StatusUnprocessableEntity)
			return
		}
		mu.Lock()
		res, errrr := db.Exec("INSERT INTO expressions (expression,login,status,result) VALUES (?,?,?,?)", expression, login, "processing", 0)
		temp, _ := res.LastInsertId()
		idCounter = int(temp)
		expressions[idCounter] = &Expression{
			ID:         idCounter,
			Expression: &postfix,
			Status:     "processing",
			Login:      login,
		}
		if errrr != nil {
			log.Println("ошибка при записи в дб ", errrr)
			return
		}
		mu.Unlock()

		response := map[string]int{"id": idCounter}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetExpressions(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	if r.Method == http.MethodGet {
		var exprList []map[string]interface{}
		login := r.Context().Value(loginContextKey).(string)
		for _, expr := range expressions {
			if expr.Login == login {
				exprList = append(exprList, map[string]interface{}{
					"id":     expr.ID,
					"status": expr.Status,
					"result": expr.Result,
				})
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.MarshalIndent(map[string][]map[string]interface{}{"expressions": exprList}, "", "    ")
		w.Write(response)
	} else {
		http.Error(w, "Method not allowed", http.StatusInternalServerError)
	}

}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodPost {
		type Result struct {
			TaskID int    `json:"taskid"`
			Result string `json:"result"`
		}

		var result Result
		err := json.NewDecoder(r.Body).Decode(&result)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		expr, exists := expressions[result.TaskID]
		if !exists {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		for i := 2; i < len(*expr.Expression); i++ {
			if calculation.IsOperator((*expr.Expression)[i][0]) {

				*expr.Expression = append((*expr.Expression)[:i-2], append([]string{result.Result}, (*expr.Expression)[i+1:]...)...)
				expr.SubStatus = ""
				break
			}
		}

		if len(*expr.Expression) == 1 {
			num, err := strconv.ParseFloat(result.Result, 64)
			if err == nil {
				expr.Status = "completed"
				expr.Result = num

				_, err := db.Exec("UPDATE  expressions SET  status=?, result=? WHERE id=? ", "completed", expr.Result, expr.ID)

				if err != nil {
					http.Error(w, "Error saving result to database", http.StatusInternalServerError)
					return
				}

			} else {
				http.Error(w, "Invalid result format", http.StatusBadRequest)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method == http.MethodGet {

		for _, expr := range expressions {
			if expr.Status == "processing" && expr.SubStatus != "waiting" {
				expr.SubStatus = "waiting"
				tasks := createTasks(expr.ID, (expr.Expression))
				w.WriteHeader(http.StatusOK)

				json.NewEncoder(w).Encode(*tasks)
				return
			}
		}
	}
}

func handleGetExpressionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, expr := range expressions {
		if expr.ID == id {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			response, _ := json.MarshalIndent(map[string]interface{}{
				"id":     expr.ID,
				"status": expr.Status,
				"result": expr.Result,
			}, "", "    ")
			w.Write(response)
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
			ttime := opTime(v[0])

			task = &Task{
				ID:             taskID,
				Arg1:           Ara(strconv.ParseFloat((*expression)[i-2], 64)),
				Arg2:           Ara(strconv.ParseFloat((*expression)[i-1], 64)),
				Operation:      string(v[0]),
				Operation_time: ttime,
			}

			return task

		}
	}
	return task
}

func RunServer() {
	r := mux.NewRouter()
	initDB()
	loadDataFromDB()
	r.HandleFunc("/api/v1/register", handleRegister).Methods("POST")
	r.HandleFunc("/api/v1/login", handleLogin).Methods("POST")
	r.HandleFunc("/api/v1/calculate", AuthMiddleware(handleCalculate))
	r.HandleFunc("/api/v1/expressions", AuthMiddleware(handleGetExpressions))
	r.HandleFunc("/internal/task", handleGetTask)

	r.HandleFunc("/api/v1/expressions/{id}", AuthMiddleware(handleGetExpressionByID))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}

func opTime(op uint8) int {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	switch op {
	case '+':
		ara, _ := strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
		return ara
	case '-':
		ara, _ := strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
		return ara
	case '*':
		ara, _ := strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
		return ara
	case '/':
		ara, _ := strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
		return ara
	}
	panic("invalid operation when attempted to send task")
}
