package tmbyn

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	json "github.com/bitly/go-simplejson"
	"github.com/trevex/golem"
)

const (
	YTApi = "https://www.googleapis.com/youtube/v3/videos?id=%s&part=snippet,contentDetails&key=%s"
)

type VidInfo struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
}

func getVidInfo(vid string) *VidInfo {
	// TODO: Parse youtube url
	// TODO: Error handling
	resp, _ := http.Get(fmt.Sprintf(YTApi, vid, ConfigString("youtube", "apiKey")))
	b, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	j, _ := json.NewJson(b)
	i := j.Get("items").GetIndex(0)
	return &VidInfo{
		Id:       vid,
		Title:    getTitle(i),
		Duration: getDuration(i),
	}
}

func getTitle(i *json.Json) string {
	s, _ := i.Get("snippet").Get("title").String()
	return s
}

func getDuration(i *json.Json) int {
	s, _ := i.Get("contentDetails").Get("duration").String()
	d, _ := time.ParseDuration(strings.ToLower(strings.TrimPrefix(s, "PT")))
	return int(d / time.Second)
}

type Player struct {
	rm   *golem.RoomManager
	play chan *PlayReq
	stop chan bool
}

func NewPlayer(rm *golem.RoomManager) *Player {
	p := Player{
		rm:   rm,
		play: make(chan *PlayReq),
		stop: make(chan bool),
	}
	go p.run()
	return &p
}

type PlayReq struct {
	Room string
	Vid  string
}

func (p *Player) run() {
	for {
		select {
		case pr := <-p.play:
			vi := getVidInfo(pr.Vid)
			log.Printf("playing vid %s at %s", vi.Id, pr.Room)
			p.rm.Emit(pr.Room, "play", vi)
		case <-p.stop:
			return
		}
	}
}

func (p *Player) Play(pr *PlayReq) {
	p.play <- pr
}
