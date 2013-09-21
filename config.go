package tmbyn

import (
	"github.com/robfig/config"
)

var conf *config.Config

func ReadConfig(confFile string) {
	var err error
	conf, err = config.ReadDefault(confFile)
	if err != nil {
		panic(err)
	}
}

func ConfigString(s string, n string) string {
	v, err := conf.String(s, n)
	if err != nil {
		panic(err)
	}
	return v
}

func ConfigInt(s string, n string) int {
	v, err := conf.Int(s, n)
	if err != nil {
		panic(err)
	}
	return v
}
