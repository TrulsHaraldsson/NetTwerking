package main

import (
	"fmt"
	"testing"
	"time"

	"../golang"
)

func TestPing(t *testing.T) {
	n1 := d7024e.StartNode(8200, "none", "none")
	fmt.Println(n1)
	time.Sleep(10 * time.Millisecond)
	n2 := d7024e.StartNode(8201, "localhost:8200", "none")
	fmt.Println("Connected")
	n2.Ping("localhost:8200")
}
