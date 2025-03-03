package main

import (
	"awesomeProject/agent"
	"awesomeProject/internal"
)

func main() {
	go agent.Man()
	internal.RunServer()
}
