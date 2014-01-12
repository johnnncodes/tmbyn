package tmbyn

import (
	"github.com/codegangsta/martini"
	"log"
	"net/http"
)

func Serve(addr string) {
	m := martini.Classic()
	// Handlers
	http.HandleFunc("/", IndexHandler())
	m.Get("/", IndexHandler())
	m.Get("/ws", WebsocketHandler())
	// Serve
	if err := http.ListenAndServe(addr, m); err != nil {
		log.Fatal(err)
	}
}
