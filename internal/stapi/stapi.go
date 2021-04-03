package stapi

import (
	"github.com/codemicro/spacetraders/internal/config"
	"github.com/parnurzeal/gorequest"
	"time"
)

var request = gorequest.New().SetDebug(true).Timeout(10 * time.Second).AppendHeader("Authorization", "Bearer " + config.C.Token)
