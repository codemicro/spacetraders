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

	_ = NewCore(userInfo)

	return nil

}