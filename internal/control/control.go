package control

import (
	"github.com/codemicro/spacetraders/internal/config"
	"github.com/codemicro/spacetraders/internal/stapi"
)

func Start() error {

	// get user and ships
	userInfo, err := stapi.GetUserInfo(config.C.Username)
	if err != nil {
		return err
	}

	core := NewCore(userInfo)

	for i, ship := range userInfo.Ships {
		var scout bool
		if i == 0 {
			scout = true
		}
		_ = NewShipController(ship, core, scout)
	}

	return nil

}
