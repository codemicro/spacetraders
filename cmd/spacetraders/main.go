package main

import (
	"fmt"
	"github.com/codemicro/spacetraders/internal/control"
	"os"
	"os/signal"
)

func main() {
	stopFunc, err := control.Start()
	if err != nil {
		panic(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	stopFunc()

	<-sig // if it's pressed for a second time, actually stop
	fmt.Println("Stopping without proper shutdown. This is a bad idea.")
}
