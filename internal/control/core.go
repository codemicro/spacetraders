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
	"os"
	"strings"
	"sync"
	"time"
)

type Core struct {
	user *stapi.User
	allowStartNewFlight bool
	sessionProfit int
	stopNotifier chan string

	logger zerolog.Logger

	profitLock sync.Mutex
	stdoutLock sync.Mutex
}

func NewCore(user *stapi.User) *Core {
	c := new(Core)
	c.user = user
	c.logger = log.With().Str("area", "Core").Str("username", c.user.Username).Logger()
	c.allowStartNewFlight = true
	c.stopNotifier = make(chan string, 200)

	go c.Start()

	go func() {
		for {
			time.Sleep(time.Second * 10)
			fcont, err := ioutil.ReadFile("killswitch.txt")
			if err != nil {
				c.error(err)
			} else {
				if len(fcont) != 0 {
					c.TriggerStop()
					return
				}
			}
		}
	}()

	go func(){
		for {
			err := db.DeleteMarketDataOlderThan(time.Now().Add(-10 * time.Minute))
			if err != nil {
				c.error(err)
				return
			}
			time.Sleep(time.Minute)
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
	c.WriteToStdout("%s%s\n", aurora.BrightRed(prefix), x)
}

func (c *Core) error(err error) {
	c.logger.Error().Err(err).Msg(tool.GetContext(2))
}

func (c *Core) ReportProfit(amount int) {
	c.profitLock.Lock()
	c.sessionProfit += amount
	c.log("Received profit of %dcr, session profit now at %dcr", amount, c.sessionProfit)
	c.profitLock.Unlock()
}

func (c *Core) determineNewShipType(ship *stapi.Ship) (int, string, error) {
	numTraders, err := db.CountShipsOfType(ShipTypeTrader)
	if err != nil {
		return 0, "", err
	}
	numProbes, err := db.CountShipsOfType(ShipTypeProbe)
	if err != nil {
		return 0, "", err
	}

	var numberOfLocations int
	systemLocations, err := stapi.GetSystemLocations(tool.SystemFromSymbol(ship.Location))
	if err != nil {
		return 0, "", err
	}
	numberOfLocations = len(systemLocations)

	{
		for _, x := range systemLocations {
			if x.Type == stapi.LocationTypeWormhole {
				numberOfLocations -= 1 // this prevents probe ships being sent to wormholes
			}
		}
	}

	targetShipType := ShipTypeTrader
	var targetShipData string
	if numProbes < numTraders && numProbes < numberOfLocations {
		targetShipType = ShipTypeProbe
		{
			currentProbeLocations, err := db.GetShipDataByType(ShipTypeProbe)
			if err !=  nil {
				return 0, "", err
			}
			for _, location := range systemLocations {
				if !tool.IsStringInSlice(location.Symbol, currentProbeLocations) && location.Type != stapi.LocationTypeWormhole {
					targetShipData = location.Symbol
					break
				}
			}
		}
	}

	return targetShipType, targetShipData, nil
}

func (c *Core) TriggerStop() {
	c.allowStartNewFlight = false
	c.log("stopping all new flights...")
}

func (c *Core) Start() {
	var runningShips []string

	for _, ship := range c.user.Ships {

		dbShip, found, err := db.GetShip(ship.ID)
		if err != nil {
			c.error(err)
			c.TriggerStop()
			break
		}
		if !found {

			targetShipType, targetShipData, err := c.determineNewShipType(ship)
			if err != nil {
				c.error(err)
				c.TriggerStop()
				break
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
				c.TriggerStop()
				break
			}
		}

		runningShips = append(runningShips, ship.ID)
		_ = NewShipController(ship, c, dbShip.Type, dbShip.Data)
		time.Sleep(time.Second * 2) // spaces out requests a bit
	}


	if len(runningShips) == len(c.user.Ships) {
		c.log("all ships started")
	}

	for shipID := range c.stopNotifier {

		var n int
		for _, x := range runningShips {
			if x != shipID {
				runningShips[n] = x
				n++
			}
		}
		runningShips = runningShips[:n]

		remaining := len(runningShips)
		c.log("ship %s stopping - %d remaining", shipID, remaining)

		if remaining == 0 {
			c.log("all ships shutdown - bye!")
			os.Exit(0)
		}
	}

}