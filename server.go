package tmbyn

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
)

func Serve(addr string) {
	// Config
	redisAddr := ConfigString("redis", "address")

	// Redis connection
	redisConn, err := redis.Dial("tcp", redisAddr)
	if err != nil {
		panic(err)
	}
	defer redisConn.Close()
	psc := redis.PubSubConn{redisConn}
	psc.Subscribe("tmbyn")

	// Handlers
	http.HandleFunc("/", IndexHandler())
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))
	http.HandleFunc("/ws", WebsocketHandler(psc))

	// Serve
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
