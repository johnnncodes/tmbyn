package tmbyn

import (
	"github.com/garyburd/redigo/redis"
	"github.com/robfig/config"
	"log"
	"net/http"
)

func Serve(addr string, conf *config.Config) {

	redisAddr, err := conf.String("redis", "address")
	if err != nil {
		panic(err)
	}
	redisConn, err := redis.Dial("tcp", redisAddr)
	if err != nil {
		panic(err)
	}
	defer redisConn.Close()
	psc := redis.PubSubConn{redisConn}
	psc.Subscribe("tmbyn")

	http.HandleFunc("/", IndexHandler())
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))
	http.HandleFunc("/ws", WebsocketHandler(psc))

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
