package tmbyn

import (
	"log"
	"net/http"

	_ "github.com/marksteve/tmbyn/statik"
	"github.com/rakyll/statik/fs"
)

func Serve(addr string) {
	// Handlers
	http.HandleFunc("/", IndexHandler())
	http.HandleFunc("/ws", WebsocketHandler())
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(sfs)))
	// Serve
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

var sfs http.FileSystem

func init() {
	sfs, _ = fs.New()
}
