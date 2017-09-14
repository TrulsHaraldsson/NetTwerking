package main

import (
	"testing"

	"../golang"
)

func TestPing(t *testing.T) {
	d7024e.StartNode(8000, "none")
	n2 := d7024e.StartNode(8001, "localhost:8000")
	n2.SendPingMessage("localhost:8000")
}
