package tmbyn

import (
	"log"
	"net/http"
)

func Serve(addr string) {
	// Handlers
	http.HandleFunc("/", IndexHandler())
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/ws", WebsocketHandler())
	// Serve
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
