package easystat

import (
	//"fmt"
	"io"
	"time"
)

type StatMap map[string]int64

type WriteFuncType func(io.Writer, time.Duration, time.Duration, StatMap, StatMap, StatMap)

// instantiable struct w/ OO style methods
type Stats struct {
	c          chan IncomingStat
	data       StatMap
	w          io.Writer
	started    time.Time
	last_time  time.Time
	last_data  StatMap
	write_func WriteFuncType
}

type IncomingStat struct {
	key   string
	value int64
}

// return a new initialized Stats struct, and launch gathering goroutine
func NewWriter(w io.Writer, interval time.Duration, write_func WriteFuncType) *Stats {
	initial_now := time.Now()
	stats := &Stats{make(chan IncomingStat, 100), make(StatMap), w, initial_now, initial_now, make(StatMap), write_func}
	go stats.gather_stats(interval)
	return stats
}

// goroutine to:
// - read from channel and increment
// - every interval, write stats to w, which could be a file, os.Stderr, or os.Stdout
func (stats *Stats) gather_stats(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		//fmt.Println("head of loop")
		select {
		case incoming, ok := <-stats.c:
			//fmt.Println("got incoming", incoming)
			if ok {
				stats.data[incoming.key] += incoming.value
			} else {
				ticker.Stop()
				//fmt.Println("writing from closed chan")
				stats.write()
				return
			}
		case <-ticker.C:
			//fmt.Println("got ticker and writing:", now)
			stats.write()
		}

	}
}

func (s *Stats) write() {
	now := time.Now()
	time_delta := now.Sub(s.last_time)
	since_start := now.Sub(s.started)
	//fmt.Fprintf(s.w, "%v %v: %#v %#v %#v\n", since_start, time_delta, s.data, s.deltas(), s.rates(time_delta))
	s.write_func(s.w, since_start, time_delta, s.data, s.deltas(), s.rates(time_delta))

	// copy current data into previous, so we can compare against it next time
	s.last_time = now
	copy_map(s.data, s.last_data)
}

// return new map showing difference from old to new values
func (s *Stats) deltas() StatMap {
	deltas := make(StatMap)
	for k, v := range s.data {
		deltas[k] = v - s.last_data[k]
	}
	return deltas
}

// return new map showing rates of increment from old to new values
func (s *Stats) rates(time_delta time.Duration) StatMap {
	rates := make(StatMap)
	time_delta_seconds := time_delta.Seconds()
	for k, v := range s.data {
		rates[k] = int64(float64(v-s.last_data[k]) / time_delta_seconds)
	}
	return rates
}

// Increment a statistic by a certain amount (specify negative value to decrement)
func (s *Stats) Add(k string, v int) {
	s.c <- IncomingStat{k, int64(v)}
}

// Dump a copy of our data into a new map for external consumption
// probably not concurrency safe
func (s *Stats) Data() map[string]int64 {
	r := make(map[string]int64)
	copy_map(s.data, r)
	return r
}

// Stop using this set of statistics
func (s *Stats) Stop() {
	//fmt.Println("length:", len(s.c))
	close(s.c)
}

// copy all the keys/values of one map to another
// does not remove any keys
// map variables are actually pointers
func copy_map(from, to StatMap) {
	for k, v := range from {
		to[k] = v
	}
}
