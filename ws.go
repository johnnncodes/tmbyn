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

var (
	rooms     = golem.NewRoomManager()
	connUser  = make(map[*golem.Connection]*UserConn)
	roomUsers = make(map[string]map[*UserConn]bool)
	userRooms = make(map[*UserConn]map[string]bool)
)

type UserConn struct {
	*golem.Connection
	Name string
}

func NewUserConn(conn *golem.Connection) *UserConn {
	uc := &UserConn{Connection: conn}
	connUser[conn] = uc
	return uc
}

type RoomUser struct {
	Room string `json:"room"`
	User string `json:"user"`
}

type RoomUsers struct {
	Users []string `json:"users"`
}

type Message struct {
	Room string `json:"room"`
	User string `json:"user"`
	Text string `json:"text"`
}

func getRoomUsers(r string) *RoomUsers {
	u := make([]string, 0)
	for _uc, _ := range roomUsers[r] {
		u = append(u, _uc.Name)
	}
	return &RoomUsers{u}
}

func join(uc *UserConn, ru *RoomUser) {
	// Get or create room
	if ru.Room == "" {
		ru.Room = id()
	}

	// Set room
	// TODO: Unique names
	uc.Name = ru.User

	// Join
	rooms.Join(ru.Room, uc.Connection)

	// Announce
	rooms.Emit(ru.Room, "join_room", ru)
	uc.Emit("join", ru)

	// Hackish mapping
	if roomUsers[ru.Room] == nil {
		roomUsers[ru.Room] = make(map[*UserConn]bool)
	}
	roomUsers[ru.Room][uc] = true
	if userRooms[uc] == nil {
		userRooms[uc] = make(map[string]bool)
	}
	userRooms[uc][ru.Room] = true

	// Update users
	rooms.Emit(ru.Room, "users", getRoomUsers(ru.Room))

	log.Printf("%s joined %s", ru.User, ru.Room)
}

func msg(uc *UserConn, md *Message) {
	md.User = uc.Name
	rooms.Emit(md.Room, "msg", &md)
	log.Printf("%s talked at %s", md.User, md.Room)
}

func WebsocketHandler(psc redis.PubSubConn) func(http.ResponseWriter, *http.Request) {
	g := golem.NewRouter()
	g.SetConnectionExtension(NewUserConn)
	g.On("join", join)
	g.On("msg", msg)
	g.OnClose(func(conn *golem.Connection) {
		rooms.LeaveAll(conn)
		uc, ok := connUser[conn]
		if ok {
			for r, _ := range userRooms[uc] {
				delete(roomUsers[r], uc)
				rooms.Emit(r, "leave_room", &RoomUser{r, uc.Name})
				rooms.Emit(r, "users", getRoomUsers(r))
				log.Printf("%s left %s", uc.Name, r)
			}
			delete(userRooms, uc)
		}
	})
	return g.Handler()
}
