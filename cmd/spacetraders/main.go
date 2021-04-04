package main

import (
	"fmt"
	"github.com/codemicro/spacetraders/internal/control"
	"os"
	"os/signal"
)

func main() {
	err := control.Start()
	if err != nil {
		panic(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	fmt.Println("bye")
}
