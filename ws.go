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

type RoomUser struct {
	Room string `json:"room"`
	User string `json:"user"`
}

type Message struct {
	Room string `json:"room"`
	User string `json:"user"`
	Text string `json:"text"`
}

func join(conn *UserConn, ru *RoomUser) {
	// Get or create room
	if ru.Room == "" {
		ru.Room = id()
	}

	// Set room
	// TODO: Unique names
	conn.Name = ru.User

	// Join
	rooms.Join(ru.Room, conn.Connection)

	// Announce
	rooms.Emit(ru.Room, "join_room", ru)
	conn.Emit("join", ru)

	log.Printf("%s joined %s", ru.User, ru.Room)
}

func leave(conn *UserConn, ru *RoomUser) {
	ru.User = conn.Name
	rooms.Emit(ru.Room, "leave_room", ru)
	rooms.Leave(ru.Room, conn.Connection)

	log.Printf("%s left %s", ru.User, ru.Room)
}

func msg(conn *UserConn, md *Message) {
	md.User = conn.Name
	rooms.Emit(md.Room, "msg", &md)

	log.Printf("%s talked at %s", md.User, md.Room)
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
