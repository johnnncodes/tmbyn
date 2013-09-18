package kmsp

import (
	"github.com/garyburd/redigo/redis"
	"github.com/pilu/traffic"
	"github.com/robfig/config"
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
	router := traffic.New()
	router.Get("/", IndexHandler(redisConn))
	router.Get("/ws", WebsocketHandler(redisConn))
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))
	http.Handle("/", router)
	http.ListenAndServe(addr, nil)
}
