package kmsp

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/trevex/golem"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func id() string {
	i := 100000 + r.Intn(900000)
	return fmt.Sprintf("%x", i)
}

var rooms = golem.NewRoomManager()

type UserConn struct {
	*golem.Connection
	Name string
}

func NewUserConn(conn *golem.Connection) *UserConn {
	return &UserConn{Connection: conn}
}

type JoinData struct {
	Room string `json:"room"`
	User string `json:"user"`
}

type RoomData struct {
	Name string `json:"name"`
}

type MessageData struct {
	Room string `json:"room"`
	User string `json:"user"`
	Text string `json:"text"`
}

func join(conn *UserConn, jd *JoinData) {
	// Get or create room
	room := jd.Room
	if room == "" {
		room = id()
	}

	// Set room
	// TODO: Unique names
	conn.Name = jd.User

	// Join
	rooms.Join(room, conn.Connection)
	conn.Emit("join", &RoomData{room})

	log.Printf("%s joined %s", jd.User, room)
}

func leave(conn *UserConn, rd *RoomData) {
	rooms.Leave(rd.Name, conn.Connection)

	log.Printf("%s left %s", conn.Name, rd.Name)
}

func msg(conn *UserConn, md *MessageData) {
	md.User = conn.Name
	rooms.Emit(md.Room, "msg", &md)

	log.Printf("%s talked at %s", conn.Name, md.Room)
}

func WebsocketHandler(psc redis.PubSubConn) func(http.ResponseWriter, *http.Request) {
	g := golem.NewRouter()
	g.SetConnectionExtension(NewUserConn)
	g.On("join", join)
	g.On("leave", leave)
	g.On("msg", msg)
	g.OnClose(func(conn *golem.Connection) {
		rooms.LeaveAll(conn)
	})
	return g.Handler()
}
