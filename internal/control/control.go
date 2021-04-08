package control

import (
	"github.com/codemicro/spacetraders/internal/config"
	"github.com/codemicro/spacetraders/internal/stapi"
)

func Start() (func(), error) {

	// get user and ships
	userInfo, err := stapi.GetUserInfo(config.C.Username)
	if err != nil {
		return nil, err
	}

	core := NewCore(userInfo)

	return core.TriggerStop, nil

}