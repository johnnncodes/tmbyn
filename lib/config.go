package tmbyn

import (
	"log"

	"github.com/robfig/config"
)

var conf *config.Config

func ReadConfig(confFile string) {
	var err error
	conf, err = config.ReadDefault(confFile)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ConfigString(s string, n string) string {
	v, err := conf.String(s, n)
	if err != nil {
		log.Fatal(err.Error())
	}
	return v
}

func ConfigInt(s string, n string) int {
	v, err := conf.Int(s, n)
	if err != nil {
		log.Fatal(err.Error())
	}
	return v
}
