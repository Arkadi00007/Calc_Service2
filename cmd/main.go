package main

import (
	"awesomeProject/agent"
	"awesomeProject/internal"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	go internal.RunServer()

	time.Sleep(time.Second)
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	num, _ := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	for i := 0; i < num; i++ {
		agent.Agents()

	}

}
