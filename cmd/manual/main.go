package main

import (
	"fmt"
	"github.com/codemicro/spacetraders/internal/config"
	"github.com/codemicro/spacetraders/internal/stapi"
	"os"
	"strings"
)

func main() {
	action := strings.ToLower(os.Args[1])

	switch action {
	case "sellall":
		user, err := stapi.GetUserInfo(config.C.Username)
		if err != nil {
			panic(err)
		}
		for _, ship := range user.Ships {
			for _, good := range ship.Cargo {
				qty := good.Quantity
				for qty > 0 {
					toProcess := 300
					if qty < toProcess {
						toProcess = qty
						qty = 0
					} else {
						qty -= toProcess
					}

					_, _, err = user.SubmitSellOrder(ship.ID, good.Good, toProcess)
					if err != nil {
						err = user.JettisonCargo(ship.ID, good.Good, toProcess)
						if err != nil {
							panic(err)
						}
						fmt.Printf("jettisoned %s %s %d\n", ship.ID, good.Good, toProcess)
					} else {
						fmt.Printf("sold %s %s %d\n", ship.ID, good.Good, toProcess)
					}

				}
			}
			fmt.Println()
		}
	}
}
