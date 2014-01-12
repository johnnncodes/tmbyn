package main

import (
	"flag"
	"fmt"

	"github.com/marksteve/tmbyn/lib"
)

var (
	addr     string
	confFile string
)

func main() {
	flag.StringVar(&addr, "address", ":9000", "Address to listen")
	flag.Parse()
	confFile = flag.Arg(0)

	fmt.Printf(`
  TMBYN

  version   %s
  addr      %s
  confFile  %s

`, tmbyn.Version, addr, confFile)

	tmbyn.ReadConfig(confFile)
	tmbyn.Serve(addr)
}
