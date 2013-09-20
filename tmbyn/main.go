package main

import (
	"flag"
	"fmt"
	"github.com/marksteve/tmbyn"
	"github.com/robfig/config"
)

var (
	addr      string
	redisAddr string
)

func main() {
	flag.StringVar(&addr, "address", ":9000", "Address to listen")
	flag.Parse()
	confFile := flag.Arg(0)

	fmt.Printf(`
  TMBYN

  version   %s
  addr      %s
  confFile  %s

`, tmbyn.Version, addr, confFile)

	conf, err := config.ReadDefault(confFile)
	if err != nil {
		panic(err)
	}
	tmbyn.Serve(addr, conf)
}
