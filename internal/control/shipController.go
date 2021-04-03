package control

import (
	"fmt"
	"github.com/codemicro/spacetraders/internal/stapi"
)

func ShipController(ship *stapi.Ship) {
	fmt.Println(ship.ID[:6] + "| hello")
}