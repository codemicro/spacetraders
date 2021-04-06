package control

import (
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
