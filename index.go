package kmsp

import (
	"github.com/eknkc/amber"
	"github.com/garyburd/redigo/redis"
	"net/http"
)

func IndexHandler(redisConn redis.Conn) func(http.ResponseWriter, *http.Request) {
	t, err := amber.CompileFile("index.amber", amber.Options{false, false})
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, struct{}{})
	}
}
