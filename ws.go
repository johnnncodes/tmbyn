package tmbyn

import (
	"fmt"
	"github.com/fatih/goset"
	"github.com/garyburd/redigo/redis"
	"github.com/trevex/golem"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func id() string {
	i := 100000 + r.Intn(900000)
	return fmt.Sprintf("%x", i)
}

var (
	rooms        = golem.NewRoomManager()
	connUser     = make(map[*golem.Connection]*UserConn)
	roomNames    = make(map[string]*goset.Set)
	userRooms    = make(map[*UserConn]*goset.Set)
	invalidChars = regexp.MustCompile("\\W")
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

func join(uc *UserConn, ru *RoomUser) {
	// Remove invalid user name chars
	ru.User = invalidChars.ReplaceAllString(ru.User, "")

	// Check user name length
	if len(ru.User) == 0 {
		return
	}

	// Get or create room
	if ru.Room == "" {
		ru.Room = id()
	}

	// Init mappings
	if roomNames[ru.Room] == nil {
		roomNames[ru.Room] = goset.New()
	}
	if userRooms[uc] == nil {
		userRooms[uc] = goset.New()
	}

	// Append _ for dupe names
	for roomNames[ru.Room].Has(ru.User) {
		ru.User += "_"
	}

	// Set room
	uc.Name = ru.User

	// Join
	rooms.Join(ru.Room, uc.Connection)

	// Announce
	rooms.Emit(ru.Room, "join_room", ru)
	uc.Emit("join", ru)

	// Update mappings
	roomNames[ru.Room].Add(uc.Name)
	userRooms[uc].Add(ru.Room)

	// Update users
	rooms.Emit(ru.Room, "users", &RoomUsers{roomNames[ru.Room].StringSlice()})

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
			for _, r := range userRooms[uc].StringSlice() {
				roomNames[r].Remove(uc.Name)
				rooms.Emit(r, "leave_room", &RoomUser{r, uc.Name})
				rooms.Emit(r, "users", &RoomUsers{roomNames[r].StringSlice()})
				log.Printf("%s left %s", uc.Name, r)
			}
			delete(userRooms, uc)
		}
		delete(connUser, conn)
	})
	return g.Handler()
}
