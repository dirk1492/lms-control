package control

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"net"
)

type Server struct {
	Host     string
	Port     int
	LastScan time.Time
	Version  string
	UUID     string
	Albums   int
	Artists  int
	Genres   int
	Songs    int
	Duration float64
	Players  *PlayerList

	conn net.Conn
}

func Connect(host string, port int) (*Server, error) {

	url := fmt.Sprintf("%s:%d", host, port)

	log.Printf("Connect to lms server %v\n", url)

	conn, err := net.Dial("tcp", url)
	if err != nil {
		log.Printf("Connection failed %v\n", err)
		return nil, err
	}

	log.Printf("Connected to %v\n", host)

	server := &Server{
		Host: host,
		Port: port,
		conn: conn,
	}

	server.init()

	return server, nil
}

func (s *Server) Close() {
	if s != nil && s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
}

func (s *Server) init() {
	log.Printf("Read server status")

	s.send("serverstatus")
	result, err := s.read()
	if err == nil {

		parts := strings.Split(result, " ")

		for _, p := range parts {
			if len(p) > 0 {
				dp, err := url.QueryUnescape(p)
				if err == nil {
					kv := strings.Split(dp, ":")
					if len(kv) == 2 {
						switch kv[0] {
						case "lastscan":
							val, err := strconv.ParseInt(kv[1], 10, 64)
							if err == nil {
								s.LastScan = time.Unix(val, 0)
							} else {
								log.Printf("Error parse lastscan: %v", err)
							}
						case "version":
							s.Version = kv[1]
						case "uuid":
							s.UUID = kv[1]
						case "info total albums":
							val, err := strconv.Atoi(kv[1])
							if err == nil {
								s.Albums = val
							} else {
								log.Printf("Error parse info total albums: %v", err)
							}
						case "info total artists":
							val, err := strconv.Atoi(kv[1])
							if err == nil {
								s.Artists = val
							} else {
								log.Printf("Error parse info total artists: %v", err)
							}
						case "info total genres":
							val, err := strconv.Atoi(kv[1])
							if err == nil {
								s.Genres = val
							} else {
								log.Printf("Error parse info total genres: %v", err)
							}
						case "info total songs":
							val, err := strconv.Atoi(kv[1])
							if err == nil {
								s.Songs = val
							} else {
								log.Printf("Error parse info total songs: %v", err)
							}
						case "info total duration":
							val, err := strconv.ParseFloat(kv[1], 64)
							if err == nil {
								s.Duration = val
							} else {
								log.Printf("Error parse info total duration: %v", err)
							}
						case "player count":
							//val, err := strconv.Atoi(kv[1])
							//if err == nil {
							s.Players = NewPlayerList(s)
							//}
						}
					}
				}
			}
		}
	}

	log.Printf("Server status:\n%v", s.String())

}

func (s *Server) send(cmd string) (int, error) {
	return s.conn.Write([]byte(cmd + "\n"))
}

func (s *Server) read() (string, error) {
	var rc bytes.Buffer
	var buffer [1]byte
	p := buffer[:]

	for {
		n, err := s.conn.Read(p)
		if err != nil {
			if err != io.EOF {
				return "", err
			}
			break
		}

		if n == 0 {
			continue
		}

		if buffer[0] == '\n' {
			break
		} else {
			rc.Write(buffer[:])
		}
	}

	return rc.String(), nil
}

func (s *Server) query(q string) (string, error) {
	idx := strings.Index(q, "?")
	if idx == -1 {
		return "", fmt.Errorf("? not found in query")
	}

	_, err := s.send(q)
	if err != nil {
		log.Printf("Error send request: %v\n", err)
		return "", err
	}

	result, err := s.read()
	if err != nil {
		log.Printf("Error read response: %v\n", err)
		return "", err
	}

	dp, err := url.QueryUnescape(result)
	if err != nil {
		log.Printf("Error unescape response: %v\n", err)
		return "", err
	}

	return dp[idx:], nil
}

func (s *Server) set(q string, val interface{}) (string, error) {
	idx := strings.Index(q, "?")
	if idx == -1 {
		return "", fmt.Errorf("? not found in query")
	}

	q = strings.Replace(q, "?", fmt.Sprintf("%v", val), 1)

	_, err := s.send(q)
	if err != nil {
		log.Printf("Error send request: %v\n", err)
		return "", err
	}

	result, err := s.read()
	if err != nil {
		log.Printf("Error read response: %v\n", err)
		return "", err
	}

	dp, err := url.QueryUnescape(result)
	if err != nil {
		log.Printf("Error unescape response: %v\n", err)
		return "", err
	}

	return dp[idx:], nil
}

func (s *Server) getPlayerCount() int {
	resp, err := s.query("player count ?")
	if err != nil {
		log.Printf("Error read player count: %v\n", err)
		return -1
	}

	cnt, err := strconv.Atoi(resp)
	if err != nil {
		log.Printf("Error parse player count: %v -> %v\n", resp, err)
		return -1
	}

	return cnt
}

func (s *Server) getPlayer(index int) *Player {

	id, err := s.query(fmt.Sprintf("player id %v ?", index))
	if err != nil {
		log.Printf("Error read player%v id: %v\n", index, err)
		return nil
	}

	uuid, err := s.query(fmt.Sprintf("player uuid %v ?", index))
	if err != nil {
		log.Printf("Error read player%v uuid: %v\n", index, err)
		return nil
	}

	name, err := s.query(fmt.Sprintf("player name %v ?", index))
	if err != nil {
		log.Printf("Error read player%v name: %v\n", index, err)
		return nil
	}

	return &Player{
		Index:  index,
		ID:     id,
		UUID:   uuid,
		Name:   name,
		server: s,
	}
}

func (s *Server) Check(timeTable *TimeTable) {
	entry := timeTable.now()

	s.Players.foreach(func(player *Player) {
		if player != nil {
			player.Check(entry)
		}
	})
}

func (s *Server) UpdatePlayerList() {
	s.Players.update()
}

func (s *Server) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Host:		%v\n", s.Host))
	sb.WriteString(fmt.Sprintf("Version:	%v\n", s.Version))
	sb.WriteString(fmt.Sprintf("UUID:		%v\n", s.UUID))
	sb.WriteString(fmt.Sprintf("Lastscan: 	%v\n", s.LastScan))
	sb.WriteString(fmt.Sprintf("Albums:		%v\n", s.Albums))
	sb.WriteString(fmt.Sprintf("Artists:	%v\n", s.Artists))
	sb.WriteString(fmt.Sprintf("Songs:		%v\n", s.Songs))
	sb.WriteString(fmt.Sprintf("Genres:		%v\n", s.Genres))
	sb.WriteString(fmt.Sprintf("Duration:	%v\n", s.Duration))
	sb.WriteString(fmt.Sprintf("Players:	\n"))

	s.Players.foreach(func(player *Player) {
		sb.WriteString(fmt.Sprintf("%s\n", player.String()))
	})

	return sb.String()
}
