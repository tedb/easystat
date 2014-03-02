package easystat

import (
	"fmt"
	"io"
	"time"
)

// instantiable struct w/ OO style methods
type Stats struct {
	c         chan IncomingStat
	data      map[string]int64
	w         io.Writer
	started   time.Time
	last_time time.Time
	last_data map[string]int64
}

type IncomingStat struct {
	key   string
	value int64
}

// return a new initialized Stats struct, and launch goroutine to:
// - read from channel and increment
// - every interval, write stats to w, which could be a file, os.Stderr, or os.Stdout
func NewWriter(w io.Writer, interval time.Duration) *Stats {
	now := time.Now()
	stats := &Stats{make(chan IncomingStat), make(map[string]int64), w, now, now, make(map[string]int64)}
	go func() {

		ticker := time.NewTicker(interval)
		for {
			select {
			case incoming, ok := <-stats.c:
				if ok {
					stats.data[incoming.key] += incoming.value
				} else {
					stats.write(time.Now())
					ticker.Stop()
					return
				}
			case now := <-ticker.C:
				stats.write(now)
			}

		}
	}()
	return stats
}

func (s *Stats) write(now time.Time) {
	fmt.Fprintf(s.w, "%v %v\n", now, s.data)

	// copy current data into previous, so we can compare against it next time
	s.last_time = now
	for k, v := range s.data {
		s.last_data[k] = v
	}
}

func (s *Stats) Add(k string, v int64) {
	s.c <- IncomingStat{k, v}
}
