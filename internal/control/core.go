package control

import (
	"fmt"
	"github.com/codemicro/spacetraders/internal/db"
	"github.com/codemicro/spacetraders/internal/stapi"
	"github.com/codemicro/spacetraders/internal/tool"
	"github.com/logrusorgru/aurora"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

type Core struct {
	user *stapi.User
	allowStartNewFlight bool

	logger zerolog.Logger

	stdoutLock sync.Mutex
}

func NewCore(user *stapi.User) *Core {
	c := new(Core)
	c.user = user
	c.logger = log.With().Str("area", "Core").Str("username", c.user.Username).Logger()
	c.allowStartNewFlight = true

	go c.Start()

	go func() {
		for {
			time.Sleep(time.Second * 10)
			fcont, err := ioutil.ReadFile("killswitch.txt")
			if err != nil {
				c.error(err)
			} else {
				if len(fcont) != 0 {
					c.allowStartNewFlight = false
					c.log("stopping all new flights...")
					return
				}
			}
		}
	}()

	return c
}

func (c *Core) WriteToStdout(format string, a ...interface{}) {
	c.stdoutLock.Lock()
	defer c.stdoutLock.Unlock()
	fmt.Printf(format, a...)
}

func (c *Core) log(format string, a ...interface{}) {
	prefix := c.user.Username + ": "
	x := strings.ReplaceAll(fmt.Sprintf(format, a...), "\n", "\n"+strings.Repeat(" ", len(prefix)))
	c.WriteToStdout("%s%s\n", aurora.Green(prefix), x)
}

func (c *Core) error(err error) {
	c.logger.Error().Err(err).Msg(tool.GetContext(2))
}

func (c *Core) Start() {
	for _, ship := range c.user.Ships {

		dbShip, found, err := db.GetShip(ship.ID)
		if err != nil {
			c.error(err)
			return
		}
		if !found {

			numTraders, err := db.CountShipsOfType(ShipTypeTrader)
			if err != nil {
				c.error(err)
				return
			}
			numProbes, err := db.CountShipsOfType(ShipTypeProbe)
			if err != nil {
				c.error(err)
				return
			}

			var numberOfLocations int
			systemLocations, err := stapi.GetSystemLocations(tool.SystemFromSymbol(ship.Location))
			if err != nil {
				c.error(err)
				return
			}
			numberOfLocations = len(systemLocations)

			targetShipType := ShipTypeTrader
			var targetShipData string
			if numProbes < numTraders && numProbes < numberOfLocations {
				targetShipType = ShipTypeProbe
				{
					currentProbeLocations, err := db.GetShipDataByType(ShipTypeProbe)
					if err !=  nil {
						c.error(err)
						return
					}
					for _, location := range systemLocations {
						if !tool.IsStringInSlice(location.Symbol, currentProbeLocations) {
							targetShipData = location.Symbol
							break
						}
					}
				}
			}

			c.log("found unrecognised ship %s, categorising as %d", ship.ID, targetShipType)

			dbShip = &db.Ship{
				ID:   ship.ID,
				Type: targetShipType,
				Data: targetShipData,
			}
			err = dbShip.Create()
			if err != nil {
				c.error(err)
				return
			}
		}

		_ = NewShipController(ship, c, dbShip.Type, dbShip.Data)
	}
}