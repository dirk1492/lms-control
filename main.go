package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dirk1492/lms-control/control"
)

func main() {

	server, _ := control.Connect("lms001", 9090)
	defer server.Close()

	vol, _ := server.Players[1].GetVolume()
	fmt.Printf("%v", vol)

	server.Players[1].SetVolume(25)

	vol, _ = server.Players[1].GetVolume()
	fmt.Printf("%v", vol)

	table := control.ParseTimeTable("22:00:00=25,23:00:00=15,00:00:00=0,06:00=")

	log.Printf(table.String())

	select {
	case <-time.After(1 * time.Second):
		server.Check(table)
	}

}
