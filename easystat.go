package easystat

import (
	"fmt"
	"io"
	"os"
	"time"
)

// instantiable struct w/ OO style methods
type Stats struct {
	c    chan IncomingStat
	// all access to the data map is wrapped in this mutex
	mutex sync.RWMutex
	data map[string]int64
}

type IncomingStat struct {
	key string
	value int64
}

// return a new initialized Stats struct, and launch goroutine to read from channel and increment 
func New() *Stats {
	stats := &Stats{make(chan IncomingStat), make(map[string]int64)}
	go func() {
		range incoming := range stats.c {
			stats.mutex.Lock()
			stats.data[incoming.key] += incoming.value
			stats.mutex.Unlock()
		}
	}()
	return stats
}

// every interval, print a summary of the statistics to stdout
func (s *Stats) PrintLoop(interval time.Duration) {
	s.WriteLoop(os.Stdout, interval)
}

// every interval, write a summary of statistics to the io.Writer passed
func (s *Stats) WriteLoop(w io.Writer, interval time.Duration) {
	go func() {
		ticker := time.Tick(interval)
		for now := range ticker {
			fmt.Fprintf(w, "%v %v\n", now, s.data)
		}
	}()
}

func (s *Stats) Add(k string, v int64) {
	s.c <- &IncomingStat{k, v}
}
