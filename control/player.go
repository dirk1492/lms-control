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
	power, _ := p.GetPower()
	vol, _ := p.GetVolume()

	return fmt.Sprintf("%2v %-20v %v (vol: %v, power: %v)", p.Index, p.Name, p.ID, vol, power)
}

func (p *Player) GetPower() (bool, error) {
	val, err := p.server.query(fmt.Sprintf("%v power ?", p.ID))
	if err != nil {
		log.Printf("Error read power state of player %v: %v", p.Name, err)
		return false, err
	}

	if val == "1" {
		return true, err
	} else {
		return false, err
	}
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

func (p *Player) Check(entry *TimeTableEntry) {
	power, err := p.GetPower()
	if err == nil && power {
		vol, err := p.GetVolume()
		if err == nil {
			if vol > entry.max {
				log.Printf("Set volume of player %v to %v", p, entry.max)
				p.SetVolume(entry.max)
			} else if entry.max == 0 {
				log.Printf("Switch player %v off", p)
				p.SetPower(false)
			}
		}
	}

}

func (p *Player) conv(b bool) int {
	if b {
		return 1
	}
	return 0
}
