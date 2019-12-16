# lms-control

## Usage
``` shell-script
Application to limit volume of players connected to a Logitech Mediaserver by timetable

Usage:
  lms-control [flags]

Flags:
  -h, --help                help for lms-control
  -i, --interval duration   Duration between 2 checks (default 1s)
  -l, --lms string          Hostname of the lms server
  -p, --port int            Port of the lms telnet interface (default 9090)
  -t, --timetable string    Comma separated list of timetable entries (e.g. 22:00:00=20,23:00:00=15,00:00:00=0,05:30=100)
```

## Environment variables
|Variable|Description|Default|
|--------|-----------|-------|
|LMS_SERVER|IP or hostname of the LMS server|localhost|
|LMS_PORT|Port of the lms telnet interface|9090|
|TIMETABLE|Comma separated list of timetable entries||
|INTERVAL|Duration between 2 checks|1s|

## Run in docker
``` shell-script
docker run -d --env LMS_SERVER=lms --env TIMETABLE="22:00=25,06:00=100" dil001/lms-control
```
