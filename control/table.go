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
	start   time.Time
	seconds int
	max     int
}

func (t *TimeTableEntry) copy() TimeTableEntry {
	return TimeTableEntry{start: t.start, seconds: t.seconds, max: t.max}
}

type TimeTable struct {
	table []TimeTableEntry
}

func ParseTimeTable(txt string) *TimeTable {
	table := TimeTable{}

	if len(txt) > 0 {
		parts := strings.Split(txt, ",")
		for _, p := range parts {
			kv := strings.Split(p, "=")
			if len(kv) == 2 {
				table.add(kv[0], kv[1])
			} else if len(kv) == 1 {
				table.add(kv[0], "")
			}
		}
	}

	sort.Sort(table)

	if table.length() == 0 {
		table.add("00:00", "100")
	} else if table.first().seconds != 0 {
		last := table.last()
		if last != nil {
			table.add("00:00", strconv.Itoa(last.max))
		} else {
			table.add("00:00", "100")
		}
	}

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

	seconds := 60*60*tx.Hour() + 60*tx.Minute() + tx.Second()

	if len(max) == 0 {
		t.table = append(t.table, TimeTableEntry{start: tx, seconds: seconds, max: 100})
	} else {
		val, err := strconv.Atoi(max)
		if err != nil {
			log.Printf("Error parsing max value '%s': %v", txt, err)
		}
		t.table = append(t.table, TimeTableEntry{start: tx, seconds: seconds, max: val})
	}
}

func (t *TimeTable) length() int {
	return len(t.table)
}

func (t *TimeTable) last() *TimeTableEntry {
	l := len(t.table)

	if l == 0 {
		return nil
	} else {
		return &t.table[0]
	}
}

func (t *TimeTable) first() *TimeTableEntry {
	l := len(t.table)

	if l == 0 {
		return nil
	} else {
		return &t.table[l-1]
	}
}

func (t *TimeTable) now() *TimeTableEntry {
	tt := time.Now()
	seconds := 60*60*tt.Hour() + 60*tt.Minute() + tt.Second()

	for _, v := range t.table {
		if seconds >= v.seconds {
			return &v
		}
	}

	return &t.table[0]
}

func (t *TimeTable) String() string {
	var sb strings.Builder

	sb.WriteString("Time table:\n")

	for i := len(t.table) - 1; i >= 0; i-- {
		v := t.table[i]
		sb.WriteString(fmt.Sprintf("%-2v %v => %v\n", len(t.table)-i, v.start.Format("15:04:00"), v.max))
	}

	return sb.String()
}

func (t TimeTable) Len() int {
	return len(t.table)
}

func (t TimeTable) Less(i, j int) bool {
	return t.table[j].start.Before(t.table[i].start)
}

func (t TimeTable) Swap(i, j int) {
	tmp := t.table[i]
	t.table[i] = t.table[j]
	t.table[j] = tmp
}
