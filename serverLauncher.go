package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/lologarithm/survival/server"
)

func main() {
	exit := make(chan int, 1)

	fmt.Println("Starting Server!")
	// Launch server manager
	go server.RunServer(exit)

	fmt.Println("Server started. Press a ctrl+c to exit.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("Goodbye!")
	exit <- 1
	return
}
