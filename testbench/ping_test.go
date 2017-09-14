package main

import (
	"testing"
	"time"

	"../golang"
)

func TestPing(t *testing.T) {
	go d7024e.StartNode(8000, "none")

	go d7024e.StartNode(8001, "localhost:8000")
	time.Sleep(500 * time.Millisecond)
}
