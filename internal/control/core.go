package control

import (
	"fmt"
	"github.com/codemicro/spacetraders/internal/stapi"
	"sync"
)

type Core struct {
	user *stapi.User

	stdoutLock sync.Mutex
}

func NewCore(user *stapi.User) *Core {
	c := new(Core)
	c.user = user
	return c
}

func (c *Core) Log(format string, a ...interface{}) {
	c.stdoutLock.Lock()
	defer c.stdoutLock.Unlock()
	fmt.Printf(format, a...)
}