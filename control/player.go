package control

import (
	"fmt"
	"log"
	"strconv"
)

type Player struct {
	Index  int
	ID     string
	UUID   string
	Name   string
	server *Server
}

func (p *Player) String() string {
	return fmt.Sprintf("%2v %-20v %v", p.Index, p.Name, p.ID)
}

func (p *Player) SetPower(on bool) error {
	_, err := p.server.set(fmt.Sprintf("%v power ?", p.ID), strconv.Itoa(p.conv(on)))
	return err
}

func (p *Player) GetVolume() (int, error) {
	val, err := p.server.query(fmt.Sprintf("%v mixer volume ?", p.ID))
	if err != nil {
		log.Printf("Error read volume of player %v: %v", p.Name, err)
		return -1, err
	}

	rc, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Error parse volume of player %v: %v", p.Name, err)
		return -1, err
	}

	return rc, nil
}

func (p *Player) SetVolume(vol int) error {
	_, err := p.server.set(fmt.Sprintf("%v mixer volume ?", p.ID), strconv.Itoa(vol))
	return err
}

func (p *Player) conv(b bool) int {
	if b {
		return 1
	}
	return 0
}
