package config

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	Token string
}

var C config

func init() {
	fcont, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(fcont, &C); err != nil {
		panic(err)
	}
}