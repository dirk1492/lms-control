package control

import "sync"

type PlayerList struct {
	server *Server
	mux    sync.RWMutex
	list   []Player
}

type PlayerFunc func(player *Player)

func NewPlayerList(server *Server) *PlayerList {
	rc := PlayerList{server: server}
	rc.update()
	return &rc
}

func (p *PlayerList) update() {
	var list []Player
	count := p.server.getPlayerCount()
	if count != -1 {
		for i := 0; i < count; i++ {
			player := p.server.getPlayer(i)
			if player != nil {
				list = append(list, *player)
			}
		}
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	p.list = list
}

func (p *PlayerList) foreach(f PlayerFunc) {
	p.mux.RLock()
	defer p.mux.RUnlock()

	for _, player := range p.list {
		f(&player)
	}
}
