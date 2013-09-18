package kmsp

import (
	"github.com/garyburd/redigo/redis"
	"github.com/trevex/golem"
	"net/http"
)

func WebsocketHandler(redisConn redis.Conn) func(http.ResponseWriter, *http.Request) {
	psc := redis.PubSubConn{redisConn}
	psc.Subscribe("tmsp")
	g := golem.NewRouter()
	return g.Handler()
}
