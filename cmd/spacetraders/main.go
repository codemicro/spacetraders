package main

import (
	"github.com/codemicro/spacetraders/internal/control"
	"time"
)

func main() {
	err := control.Start()
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)
}
