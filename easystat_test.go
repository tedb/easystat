package easystat

import (
	"testing"
)

func TestNew(t *testing.T) {
	_ = New()
}

func TestWriteLoop(t *testing.T) {
	stats := New()
	slice := make([]byte, 0, 1000)
	w := bytes.NewBuffer(slice)

}
