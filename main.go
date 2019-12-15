package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dirk1492/lms-control/control"
	"github.com/spf13/cobra"
)

var timeTable string
var host string
var port int
var interval time.Duration

func main() {

	var rootCmd = &cobra.Command{Use: "lms-control",
		Short: "says hello",
		Long:  "Application to limit volume of players connected to a Logitech Mediaserver by timetable",
		Run:   run,
	}

	rootCmd.Flags().StringVarP(&timeTable, "timetable", "t", "", "Comma seperated list of timetable entries (e.g. 22:00:00=20,23:00:00=15,00:00:00=0,05:30=)")
	rootCmd.MarkFlagRequired("time-table")
	rootCmd.Flags().StringVarP(&host, "lms", "l", "", "Hostname of the lms server")
	rootCmd.MarkFlagRequired("lms")
	rootCmd.Flags().IntVarP(&port, "port", "p", 9090, "Port of the lms telnet interface")
	rootCmd.Flags().DurationVarP(&interval, "interval", "i", 1*time.Second, "Duration between 2 checks")

	rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) {
	server, _ := control.Connect(host, port)
	defer server.Close()

	if server == nil {
		os.Exit(1)
	}

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	exit := make(chan bool)

	vol, _ := server.Players[1].GetPower()
	fmt.Printf("%s -> %v\n", server.Players[1].Name, vol)

	table := control.ParseTimeTable(timeTable)

	log.Printf(table.String())

	go func() {
		for {
			select {
			case <-sigchan:
				exit <- true
			case <-time.After(interval):
				server.Check(table)
			}
		}
	}()

	<-exit

}
