package easystat

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	w := bytes.NewBuffer(make([]byte, 0, 10))
	_ = NewWriter(w, 1*time.Second)
}

func TestWriteLoop(t *testing.T) {
	slice := make([]byte, 0, 10000)
	w := bytes.NewBuffer(slice)
	stats := NewWriter(w, 20*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	stats.Add("foo", 1)
	time.Sleep(20 * time.Millisecond)
	stats.Add("foo", 2)
	time.Sleep(20 * time.Millisecond)
	stats.Add("foo", 4)
	time.Sleep(20 * time.Millisecond)
	stats.Add("foo", 8)
	fmt.Println(w)
	stats.Add("foo", 16) // sum is 31 by now
	time.Sleep(20 * time.Millisecond)
	stats.Add("foo", 32)
	stats.Add("foo", 0) // does this line make the 32 show up??
	stats.Stop()

	fmt.Println(w)
	// there should be 4 lines of data + newline
	lines := bytes.Split(w.Bytes(), []byte("\n"))
	if len(lines) != 5 {
		t.Errorf("expected 5 newlines in output, got %d", len(lines))
	}

	data := stats.Data()
	if data["foo"] != 31 {
		t.Errorf("expected data[foo] == 31 but got %d", data["foo"])
	}
}
