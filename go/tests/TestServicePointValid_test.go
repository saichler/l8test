package tests

import "testing"

func TestMain(m *testing.M) {
	setup()
	m.Run()
	tear()
}

func TestServicePointValid(t *testing.T) {
}
