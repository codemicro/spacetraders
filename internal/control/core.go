package control

import (
	"fmt"
	"github.com/codemicro/spacetraders/internal/db"
	"github.com/codemicro/spacetraders/internal/stapi"
	"github.com/codemicro/spacetraders/internal/tool"
	"github.com/logrusorgru/aurora"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
	"sync"
)

type Core struct {
	user *stapi.User

	logger zerolog.Logger

	stdoutLock sync.Mutex
}

func NewCore(user *stapi.User) *Core {
	c := new(Core)
	c.user = user
	c.logger = log.With().Str("area", "Core").Str("username", c.user.Username).Logger()

	go c.Start()

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

			c.log("found unrecognised ship %s, categorising as trader", ship.ID)

			dbShip = &db.Ship{
				ID:   ship.ID,
				Type: ShipTypeTrader,
				Data: "",
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