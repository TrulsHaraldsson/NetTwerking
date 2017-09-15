package main

import (
	"fmt"
	"testing"

	"../golang"
	"time"
)

func TestFindNode(t *testing.T) {
	d7024e.StartNode(8100, "none", "none") //W
	time.Sleep(50 * time.Millisecond)
	n1 := d7024e.StartNode(8101, "localhost:8100", "none") // Q -> W
	d7024e.StartNode(8102, "localhost:8100", "none") // S -> W
	d7024e.StartNode(8103, "localhost:8102", "2111111400000000000000000000000000000000") // A -> S
	fmt.Println("All nodes connected", n1)
	contact := n1.SendFindContactMessage(d7024e.NewKademliaID("2111111400000000000000000000000000000000"))

	if !contact.ID.Equals(d7024e.NewKademliaID("2111111400000000000000000000000000000000")) {
		t.Error("Not correct contact returned", contact)
	}
}
