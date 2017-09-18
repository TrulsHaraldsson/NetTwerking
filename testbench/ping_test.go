package main

import (
	"testing"
	"time"

	"../golang"
)

func TestPing(t *testing.T) {
	// Node A
	A := d7024e.NewKademlia(8200, "none")
	A.Start()
	time.Sleep(10 * time.Millisecond)

	// Node B
	B := d7024e.NewKademlia(8201, "none")
	B.Start()
	time.Sleep(10 * time.Millisecond)

	A.Ping("localhost:8201")
}
