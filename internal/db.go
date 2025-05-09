package internal

import (
	"awesomeProject/pkg/calculation"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./calculator.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			login TEXT PRIMARY KEY,
			password TEXT NOT NULL
		);

		

		CREATE TABLE IF NOT EXISTS expressions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			expression TEXT NOT NULL,
			login TEXT,
			status TEXT,
			result REAL,
			FOREIGN KEY(login) REFERENCES users(login)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func getExprsFromDB() *sql.Rows {
	mumu.Lock()
	defer mumu.Unlock()

	db, err := sql.Open("sqlite3", "./calculator.db")
	if err != nil {
		log.Printf("Ошибка открытия соединения с базой данных: %v", err)
	}
	defer db.Close()

	query := "SELECT * FROM expressions ORDER BY id ASC"

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return nil
	}
	return rows
}

func loadDataFromDB() {
	rows := getExprsFromDB()
	if rows == nil {
		return
	}

	for rows.Next() {
		var id int
		var expression string
		var login string
		var status string
		var result float64

		err := rows.Scan(&id, &expression, &login, &status, &result)
		if err != nil {
			log.Fatalf("Ошибка сканирования строки: %v", err)
		}
		if status == "processing" {
			postfix, _ := calculation.Calc(expression)

			expressions[id] = &Expression{
				ID:         id,
				Expression: &postfix,
				Status:     "processing",
				Login:      login,
			}
		} else if status == "completed" {
			expressions[id] = &Expression{
				ID:     id,
				Status: "completed",
				Result: result,
				Login:  login,
			}
		}
	}

}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE login = ?)", user.Login).Scan(&exists)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	_, err = db.Exec("INSERT INTO users (login, password) VALUES (?, ?)", user.Login, user.Password)
	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registered"})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var logins User
	if err := json.NewDecoder(r.Body).Decode(&logins); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE login = ?", logins.Login).Scan(&storedPassword)
	if err != nil || storedPassword != logins.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(logins.Login)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
