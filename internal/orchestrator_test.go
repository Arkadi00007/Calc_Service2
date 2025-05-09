package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
}

func validateCredentials(login, password string) (bool, error) {
	if login == "correctUser" && password == "correctPass" {
		return true, nil
	}
	return false, nil
}

func handleLoginRequest(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	valid, err := validateCredentials(loginRequest.Login, loginRequest.Password)
	if err != nil || !valid {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	response := LoginResponse{Message: "Login successful"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func TestLoginSuccessful(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/login", handleLoginRequest).Methods("POST")

	loginRequestBody := `{"login": "correctUser", "password": "correctPass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(loginRequestBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status 200 OK")

	var response LoginResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.Nil(t, err, "Expected valid JSON response")
	assert.Equal(t, "Login successful", response.Message, "Expected 'Login successful' message")
}

func TestLoginInvalidCredentials(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/login", handleLoginRequest).Methods("POST")

	loginRequestBody := `{"login": "incorrectUser", "password": "wrongPass"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(loginRequestBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected status 401 Unauthorized")

	body := rr.Body.String()
	assert.Contains(t, body, "Invalid credentials", "Expected 'Invalid credentials' in response body")
}
