package main

import (
	"testing"
	"time"

	"../golang"
)

func TestPing(t *testing.T) {
	// Node A
	A := d7024e.NewKademlia("localhost:8200", "none")
	A.StartListening()
	time.Sleep(10 * time.Millisecond)

	// Node B
	B := d7024e.NewKademlia("localhost:8201", "none")
	B.StartListening()
	time.Sleep(10 * time.Millisecond)

	A.Ping("localhost:8201")
}
