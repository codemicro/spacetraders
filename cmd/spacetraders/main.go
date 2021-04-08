package main

import (
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
}
