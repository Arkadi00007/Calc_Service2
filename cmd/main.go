package main

import (
	"awesomeProject/agent"
	"awesomeProject/internal"
	"time"
)

func main() {

	go internal.RunServer()
	time.Sleep(time.Second)
	agent.Agents()

}
