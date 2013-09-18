package main

import (
	"flag"
	"fmt"
	"github.com/marksteve/kwentomosapagong"
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

   _  .----.
  (_\/      \_,
    'uu----uu~'

  KWENTO MO SA PAGONG

  version   %s
  addr      %s
  confFile  %s

`, kmsp.Version, addr, confFile)

	conf, err := config.ReadDefault(confFile)
	if err != nil {
		panic(err)
	}
	kmsp.Serve(addr, conf)
}
