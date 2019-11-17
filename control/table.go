package control

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TimeTableEntry struct {
	start time.Time
	stop  time.Time
	max   int
}

type TimeTable struct {
	table []TimeTableEntry
}

func ParseTimeTable(txt string) *TimeTable {
	table := TimeTable{}

	parts := strings.Split(txt, ",")
	for _, p := range parts {
		kv := strings.Split(p, "=")
		if len(kv) == 2 {
			table.add(kv[0], kv[1])
		} else if len(kv) == 1 {
			table.add(kv[0], "")
		}
	}

	sort.Sort(table)

	return &table
}

func (t *TimeTable) add(txt string, max string) {
	format := "15:04"
	if len(txt) > 5 {
		format = "15:04:00"
	}

	tx, err := time.Parse(format, txt)
	if err != nil {
		log.Printf("Error parsing time '%s': %v", txt, err)
		return
	}

	if len(max) == 0 {
		t.table = append(t.table, TimeTableEntry{start: tx, max: 100})
	} else {
		val, err := strconv.Atoi(max)
		if err != nil {
			log.Printf("Error parsing max value '%s': %v", txt, err)
		}
		t.table = append(t.table, TimeTableEntry{start: tx, max: val})
	}
}

func (t *TimeTable) String() string {
	var sb strings.Builder

	sb.WriteString("Time table:\n")

	for i, v := range t.table {
		sb.WriteString(fmt.Sprintf("%-2v %v => %v\n", i, v.start.Format("15:04:00"), v.max))
	}

	return sb.String()
}

func (t TimeTable) Len() int {
	return len(t.table)
}

func (t TimeTable) Less(i, j int) bool {
	return t.table[i].start.Before(t.table[j].start)
}

func (t TimeTable) Swap(i, j int) {
	tmp := t.table[i]
	t.table[i] = t.table[j]
	t.table[j] = tmp
}
