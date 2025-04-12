package tests

import (
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	setup()
	m.Run()
	tear()
}

func TestServicePointValid(t *testing.T) {
	time.Sleep(time.Second * 5)
}
