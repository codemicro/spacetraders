package config

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
)

type config struct {
	Token        string
	Username     string
	DebugMode    bool
	DatabaseFile string
}

var C config

const (
	configFile = "config.json"
	logFile    = "general.log"
)

func init() {
	fcont, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(fcont, &C); err != nil {
		panic(err)
	}

	if C.DebugMode {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		log.Logger = log.Output(f)
	}

}
