package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/dirk1492/lms-control/control"
	"github.com/spf13/cobra"
)

var (
	timeTable string
	host      string
	port      int
	interval  time.Duration
)

func main() {

	// manually set time zone
	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			log.Printf("error loading location '%s': %v\n", tz, err)
		}
	}

	var rootCmd = &cobra.Command{Use: "lms-control",
		Short: "says hello",
		Long:  "Application to limit volume of players connected to a Logitech Mediaserver by timetable",
		Run:   run,
	}

	rootCmd.Flags().StringVarP(&timeTable, "timetable", "t", "", "Comma separated list of timetable entries (e.g. 22:00:00=20,23:00:00=15,00:00:00=0,05:30=100)")
	rootCmd.Flags().StringVarP(&host, "lms", "l", "", "Hostname of the lms server")
	rootCmd.Flags().IntVarP(&port, "port", "p", 9090, "Port of the lms telnet interface")
	rootCmd.Flags().DurationVarP(&interval, "interval", "i", 1*time.Second, "Duration between 2 checks")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func run(cmd *cobra.Command, args []string) {

	timeTable = getEnv("TIMETABLE", timeTable)
	host = getEnv("LMS_SERVER", host)
	port = getInt("LMS_PORT", port)
	interval = getDuration("INTERVAL", interval)

	server, _ := control.Connect(host, port)
	defer server.Close()

	if server == nil {
		os.Exit(1)
	}

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	exit := make(chan bool)

	table := control.ParseTimeTable(timeTable)

	log.Printf(table.String())

	go func() {
		for {
			select {
			case <-sigchan:
				exit <- true
			case <-time.After(interval):
				server.Check(table)
			case <-time.After(300 * time.Second):
				server.UpdatePlayerList()
			}
		}
	}()

	<-exit

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		log.Printf("Set env var %v='%v'", key, value)
		return value
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		val, err := strconv.Atoi(value)
		if err == nil {
			log.Printf("Set env var %v='%v'", key, val)
			return val
		} else {
			log.Printf("Parse env var %v failed: %v", key, err)
			return fallback
		}
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		val, err := time.ParseDuration(value)
		if err == nil {
			log.Printf("Set env var %v='%v'", key, val)
			return val
		} else {
			log.Printf("Parse env var %v failed: %v", key, err)
			return fallback
		}
	}
	return fallback
}
